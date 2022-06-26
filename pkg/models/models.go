package models

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Route struct {
	id            int64
	uuid          string
	server_id     int64
	domain_id     int64
	endpoint_id   sql.NullInt64
	endpoint_type sql.NullString
	name          string
	spam_mode     string
	created_at    time.Time
	updated_at    time.Time
	token         string
	mode          string
}
