package csp

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"net/url"
)

// {
//   "csp-report": {
//     "document-uri": "http://example.com/signup.html",
//     "referrer": "",
//     "blocked-uri": "http://example.com/css/style.css",
//     "violated-directive": "style-src cdn.example.com",
//     "original-policy": "default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports",
//   }
// }

type Report struct {
	DocumentURI       string `json:"document-uri"`
	Referrer          string `json:"referrer"`
	BlockedURI        string `json:"blocked-uri"`
	ViolatedDirective string `json:"violated-directive"`
	OriginalPolicy    string `json:"original-policy"`
	ScriptSource      string `json:"script-source"` // Firefox specific
	ScriptSample      string `json:"script-sample"` // Firefox specific
}

type ReportRequest struct {
	Report Report `json:"csp-report"`
}

func (rr *ReportRequest) valid() bool {
	return true
}

//

type Site struct {
	Id       uint64
	Hostname string
}

func LoadSiteByHostname(db *sql.DB, hostname string) (*Site, error) {
	var site Site
	err := db.QueryRow("select Id,Hostname from Site where Hostname=$1", hostname).Scan(&site.Id, &site.Hostname)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

//

func SelectValue(db *sql.DB, name string, value string) (uint64, error) {
	var id uint64
	err := db.QueryRow(fmt.Sprintf("select Id from %s where Value = $1", name), value).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func InsertValue(db *sql.DB, name string, value string) (uint64, error) {
	log.Printf("InsertValue name=%s value=%s", name, value)
	var id uint64
	err := db.QueryRow(fmt.Sprintf("insert into %s (Value) values ($1) returning Id", name), value).Scan(&id)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			return SelectValue(db, name, value)
		}
		return 0, err
	}
	return id, nil
}

//

func (rr *ReportRequest) insert(db *sql.DB, userAgent string) error {
	url, error := url.Parse(rr.Report.DocumentURI)
	if error != nil {
		return error
	}

	site, err := LoadSiteByHostname(db, url.Host)
	if err != nil {
		return err
	}

	{
		documentUriId, err := InsertValue(db, "DocumentUri", rr.Report.DocumentURI)
		if err != nil {
			return err
		}

		var referrerId *uint64
		if len(rr.Report.Referrer) != 0 {
			id, err := InsertValue(db, "Referrer", rr.Report.Referrer)
			if err != nil {
				return err
			}
			referrerId = &id
		}

		var blockedUriId *uint64
		if len(rr.Report.BlockedURI) != 0 {
			id, err := InsertValue(db, "BlockedUri", rr.Report.BlockedURI)
			if err != nil {
				return err
			}
			blockedUriId = &id
		}

		originalPolicyId, err := InsertValue(db, "OriginalPolicy", rr.Report.OriginalPolicy)
		if err != nil {
			return err
		}
		log.Printf("Got originalPolicyId %d", originalPolicyId)

		userAgentId, err := InsertValue(db, "UserAgent", userAgent)
		if err != nil {
			return err
		}

		violatedDirectiveId, err := InsertValue(db, "ViolatedDirective", rr.Report.ViolatedDirective)
		if err != nil {
			return err
		}

		var scriptSourceId *uint64
		if len(rr.Report.ScriptSource) != 0 {
			id, err := InsertValue(db, "ScriptSource", rr.Report.ScriptSource)
			if err != nil {
				return err
			}
			scriptSourceId = &id
		}

		var scriptSampleId *uint64
		if len(rr.Report.ScriptSample) != 0 {
			id, err := InsertValue(db, "ScriptSample", rr.Report.ScriptSample)
			if err != nil {
				return err
			}
			scriptSampleId = &id
		}

		// Insert the Report

		_, err = db.Exec("insert into Report (SiteId, UserAgentId, RemoteIp, DocumentUriId, ReferrerId, BlockedUriId, ViolatedDirectiveId, OriginalPolicyId, ScriptSourceId, ScriptSampleId) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", site.Id, userAgentId, "0.0.0.0", documentUriId, referrerId, blockedUriId, violatedDirectiveId, originalPolicyId, scriptSourceId, scriptSampleId)
		if err != nil {
			return err
		}
	}

	return nil
}
