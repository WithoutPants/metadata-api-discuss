package models

import (
	"database/sql"
)

type Scene struct {
	ID        int             `db:"id" json:"id"`
	Title     sql.NullString  `db:"title" json:"title"`
	Details   sql.NullString  `db:"details" json:"details"`
	URL       sql.NullString  `db:"url" json:"url"`
	Date      SQLiteDate      `db:"date" json:"date"`
	StudioID  sql.NullInt64   `db:"studio_id,omitempty" json:"studio_id"`
	CreatedAt SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type ScenePartial struct {
	ID         int             `db:"id" json:"id"`
	Title      *sql.NullString  `db:"title" json:"title"`
	Details    *sql.NullString  `db:"details" json:"details"`
	URL        *sql.NullString  `db:"url" json:"url"`
	Date       *SQLiteDate      `db:"date" json:"date"`
	StudioID   *sql.NullInt64   `db:"studio_id,omitempty" json:"studio_id"`
	CreatedAt  *SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt  *SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

func (p *Scene) IsEditTarget() {
}
