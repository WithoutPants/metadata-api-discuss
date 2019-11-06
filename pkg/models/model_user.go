package models

import (
	"database/sql"
)

type User struct {
	ID                int            `db:"id" json:"id"`
	Name              string         `db:"name" json:"name"`
	Email             sql.NullString `db:"email" json:"email"`
	APIKey            string         `db:"api_key" json:"api_key"`
	SuccessfulEdits   int            `db:"successful_edits" json:"successful_edits"`
	UnsuccessfulEdits int            `db:"unsuccessful_edits" json:"unsuccessful_edits"`
	SuccessfulVotes   int            `db:"successful_votes" json:"successful_votes"`
	// Votes on unsuccessful edits
	UnsuccessfulVotes int `db:"unsuccessful_votes" json:"unsuccessful_votes"`
	// Calls to the API from this user over a configurable time period
	APICalls int `db:"api_calls" json:"api_calls"`
}
