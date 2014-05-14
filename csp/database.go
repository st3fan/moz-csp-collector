package csp

import (
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
)

type DatabaseSession struct {
	url string
	db  *sql.DB
}

func NewSession(url string) (*DatabaseSession, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DatabaseSession{url: url, db: db}, nil
}

func (session *DatabaseSession) Close() {
	session.db.Close()
}

func (session *DatabaseSession) Database() martini.Handler {
	return func(context martini.Context) {
		context.Map(session.db)
		context.Next()
	}
}
