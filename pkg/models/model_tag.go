package models

import "database/sql"

type Tag struct {
	ID          int             `db:"id" json:"id"`
	Name        string          `db:"name" json:"name"`
	Description sql.NullString  `db:"description" json:"description"`
	CreatedAt   SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt   SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

func (p *Tag) IsEditTarget() {
}