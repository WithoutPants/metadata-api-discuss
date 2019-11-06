package models

import (
	"database/sql"
)

type Performer struct {
	ID                int             `db:"id" json:"id"`
	Name              string          `db:"name" json:"name"`
	Disambiguation    sql.NullString  `db:"disambiguation" json:"disambiguation"`
	Gender			  sql.NullString  `db:"gender" json:"gender"`
	Birthdate         SQLiteDate      `db:"birthdate" json:"birthdate"`
	Ethnicity         sql.NullString  `db:"ethnicity" json:"ethnicity"`
	Country           sql.NullString  `db:"country" json:"country"`
	EyeColor          sql.NullString  `db:"eye_color" json:"eye_color"`
	Height            sql.NullInt64   `db:"height" json:"height"`
	CupSize           sql.NullString  `db:"cup_size" json:"cup_size"`
	MeasurementsBust  sql.NullInt64   `db:"measurements_bust" json:"measurements_bust"`
	MeasurementsWaist sql.NullInt64   `db:"measurements_waist" json:"measurements_waist"`
	MeasurementsHip   sql.NullInt64   `db:"measurements_hip" json:"measurements_hip"`
	FakeTits          sql.NullString  `db:"fake_tits" json:"fake_tits"`
	CareerStartYear   sql.NullInt64   `db:"career_start_year" json:"career_start_year"`
	CareerEndYear     sql.NullInt64   `db:"career_end_year" json:"career_end_year"`
	CreatedAt         SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt         SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

func (p *Performer) IsEditTarget() {
}
