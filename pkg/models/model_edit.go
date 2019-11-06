package models

import (
	"database/sql"
)

type Edit struct {
	ID          int            `db:"id" json:"id"`
	UserID      sql.NullInt64  `db:"user_id" json:"user_id"`
	TargetID    sql.NullString `db:"target_id" json:"target_id"`
	TargetType  string         `db:"target_type" json:"target_type"`
	Operation   string         `db:"operation" json:"operation"`
	EditComment sql.NullString `db:"edit_comment" json:"edit_comment"`
	// JSON-encoded operation details
	Details string `db:"details" json:"details"`
	Status  string `db:"status" json:"status"`
	Applied bool   `db:"applied" json:"applied"`
}
