package csp

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"log"
)

type Server *martini.ClassicMartini

func NewServer(session *DatabaseSession) Server {
	m := Server(martini.Classic())

	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))

	m.Use(session.Database())

	m.Get("/api/v1/status", func(r render.Render, db *sql.DB) {
		r.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	m.Post("/api/v1/report", binding.Json(ReportRequest{}), func(rr ReportRequest, r render.Render, db *sql.DB) {
		if rr.valid() {
			//log.Printf("Received report %+v", rr)
			err := rr.insert(db)
			if err != nil {
				log.Printf("ERROR %s", err)
			}
		} else {
			r.JSON(400, map[string]string{
				"error": "Not a valid report",
			})
		}
	})

	return m
}
