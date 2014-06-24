package csp

import (
	"database/sql"
	"log"
	"net/url"
	"time"
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
	err := db.QueryRow("select Id,Hostname from Sites where Hostname=$1", hostname).Scan(&site.Id, &site.Hostname)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

type Document struct {
	Id      uint64
	Created time.Time
	Uri     string
}

//var documentCache map[string]Document = make(map[string]Document)

func InsertOrGetDocument(tx *sql.Tx, uri string) (*Document, error) {
	var document Document

	// document, ok := documentCache[uri]
	// if ok {
	// 	return &document, nil
	// }

	err := tx.QueryRow("select Id,Created,Uri from Documents where Uri=$1", uri).Scan(&document.Id, &document.Created, &document.Uri)
	if err != nil {
		if err == sql.ErrNoRows {
			err := tx.QueryRow("insert into Documents (Uri) values ($1)  returning Id,Created,Uri", uri).Scan(&document.Id, &document.Created, &document.Uri)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	// documentCache[uri] = document
	return &document, nil
}

type Referrer struct {
	Id      uint64
	Created time.Time
	Uri     string
}

func InsertOrGetReferrer(tx *sql.Tx, uri string) (*Referrer, error) {
	var referrer Referrer
	err := tx.QueryRow("select Id,Created,Uri from Referrers where Uri=$1", uri).Scan(&referrer.Id, &referrer.Created, &referrer.Uri)
	if err != nil {
		if err == sql.ErrNoRows {
			err := tx.QueryRow("insert into Referrers (Uri) values ($1)  returning Id,Created,Uri", uri).Scan(&referrer.Id, &referrer.Created, &referrer.Uri)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &referrer, nil
}

type Blocker struct {
	Id      uint64
	Created time.Time
	Uri     string
}

// var blockerCache map[string]Blocker = make(map[string]Blocker)

func InsertOrGetBlocker(tx *sql.Tx, uri string) (*Blocker, error) {
	var blocker Blocker

	// blocker, ok := blockerCache[uri]
	// if ok {
	// 	return &blocker, nil
	// }

	err := tx.QueryRow("select Id,Created,Uri from Blockers where Uri=$1", uri).Scan(&blocker.Id, &blocker.Created, &blocker.Uri)
	if err != nil {
		if err == sql.ErrNoRows {
			err := tx.QueryRow("insert into Blockers (Uri) values ($1)  returning Id,Created,Uri", uri).Scan(&blocker.Id, &blocker.Created, &blocker.Uri)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	// blockerCache[uri] = blocker
	return &blocker, nil
}

type Policy struct {
	Id      uint64
	Created time.Time
	Policy  string
}

// var policyCache map[string]Policy = make(map[string]Policy)

func InsertOrGetPolicy(tx *sql.Tx, p string) (*Policy, error) {
	var policy Policy

	// policy, ok := policyCache[p]
	// if ok {
	// 	return &policy, nil
	// }

	err := tx.QueryRow("select Id,Created,Policy from Policies where Policy=$1", p).Scan(&policy.Id, &policy.Created, &policy.Policy)
	if err != nil {
		if err == sql.ErrNoRows {
			err := tx.QueryRow("insert into Policies (Policy) values ($1) returning Id,Created,Policy", p).Scan(&policy.Id, &policy.Created, &policy.Policy)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	// policyCache[p] = policy
	return &policy, nil
}

type UserAgent struct {
	Id        uint64
	Created   time.Time
	UserAgent string
}

func InsertOrGetUserAgent(tx *sql.Tx, ua string) (*UserAgent, error) {
	var userAgent UserAgent
	err := tx.QueryRow("select Id,Created,UserAgent from UserAgents where UserAgent=$1", ua).Scan(&userAgent.Id, &userAgent.Created, &userAgent.UserAgent)
	if err != nil {
		if err == sql.ErrNoRows {
			err := tx.QueryRow("insert into UserAgents (UserAgent) values ($1) returning Id,Created,UserAgent", ua).Scan(&userAgent.Id, &userAgent.Created, &userAgent.UserAgent)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &userAgent, nil
}

type Directive struct {
	Id        uint64
	Created   time.Time
	Directive string
}

func InsertOrGetDirective(tx *sql.Tx, d string) (*Directive, error) {
	var directive Directive
	err := tx.QueryRow("select Id,Created,Directive from Directives where Directive=$1", d).Scan(&directive.Id, &directive.Created, &directive.Directive)
	if err != nil {
		if err == sql.ErrNoRows {
			err := tx.QueryRow("insert into Directives (Directive) values ($1) returning Id,Created,Directive", d).Scan(&directive.Id, &directive.Created, &directive.Directive)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &directive, nil
}

//

func (rr *ReportRequest) insert(db *sql.DB, ua string) error {
	url, error := url.Parse(rr.Report.DocumentURI)
	if error != nil {
		return error
	}

	// See if this site is allowed to report
	site, err := LoadSiteByHostname(db, url.Host)
	if err != nil {
		return err
	}
	//log.Printf("We got a site %+v", site)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	{
		document, err := InsertOrGetDocument(tx, rr.Report.DocumentURI)
		if err != nil {
			return err
		}
		//log.Printf("We got a document %+v", document)

		var referrer *Referrer
		if len(rr.Report.Referrer) != 0 {
			referrer, err = InsertOrGetReferrer(tx, rr.Report.Referrer)
			if err != nil {
				return err
			}
			//log.Printf("We got a referrer %+v", referrer)
		}

		var blocker *Blocker
		if len(rr.Report.BlockedURI) != 0 {
			blocker, err = InsertOrGetBlocker(tx, rr.Report.BlockedURI)
			if err != nil {
				return err
			}
			//log.Printf("We got a blocker %+v", blocker)
		}

		policy, err := InsertOrGetPolicy(tx, rr.Report.OriginalPolicy)
		if err != nil {
			return err
		}
		//log.Printf("We got a policy %+v", policy)

		userAgent, err := InsertOrGetUserAgent(tx, ua)
		if err != nil {
			return err
		}

		log.Printf("Inserting directive: %s", rr.Report.ViolatedDirective)
		directive, err := InsertOrGetDirective(tx, rr.Report.ViolatedDirective)
		if err != nil {
			return err
		}

		// Insert the Report

		var referrerId *uint64
		if referrer != nil {
			referrerId = &referrer.Id
		}

		var blockerId *uint64
		if blocker != nil {
			blockerId = &blocker.Id
		}
		_, err = tx.Exec("insert into Reports (SiteId, ClientIp, DocumentId, ReferrerId, BlockerId, ViolatedDirectiveId, OriginalPolicyId, ClientUserAgent) values ($1, $2, $3, $4, $5, $6, $7, $8)", site.Id, "0.0.0.0", document.Id, referrerId, blockerId, directive.Id, policy.Id, userAgent.Id)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
