package models

import (
	"database/sql"
)

type Studio struct {
	ID        int             `db:"id" json:"id"`
	Name      string          `db:"name" json:"name"`
	ParentID  sql.NullInt64   `db:"parent_id,omitempty" json:"parent_id"`
	CreatedAt SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

func (p *Studio) IsEditTarget() {
}