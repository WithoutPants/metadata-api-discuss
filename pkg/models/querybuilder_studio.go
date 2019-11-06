package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stashdb/pkg/database"
)

type StudioQueryBuilder struct{}

func NewStudioQueryBuilder() StudioQueryBuilder {
	return StudioQueryBuilder{}
}

func (qb *StudioQueryBuilder) Create(newStudio Studio, tx *sqlx.Tx) (*Studio, error) {
	ensureTx(tx)
	fields, values := SQLGenKeysCreate(newStudio)
	result, err := tx.NamedExec(
		`INSERT INTO studios (`+fields+`)
		 		VALUES (`+values+`)
		`,
		newStudio,
	)
	if err != nil {
		return nil, err
	}
	studioID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&newStudio, `SELECT * FROM studios WHERE id = ? LIMIT 1`, studioID); err != nil {
		return nil, err
	}
	return &newStudio, nil
}

func (qb *StudioQueryBuilder) Update(updatedStudio Studio, tx *sqlx.Tx) (*Studio, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE studios SET `+SQLGenKeys(updatedStudio)+` WHERE studios.id = :id`,
		updatedStudio,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&updatedStudio, `SELECT * FROM studios WHERE id = ? LIMIT 1`, updatedStudio.ID); err != nil {
		return nil, err
	}
	return &updatedStudio, nil
}

func (qb *StudioQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	// remove studio from scenes
	_, err := tx.Exec("UPDATE scenes SET studio_id = null WHERE studio_id = ?", id)
	if err != nil {
		return err
	}

	return executeDeleteQuery("studios", id, tx)
}

func (qb *StudioQueryBuilder) Find(id int, tx *sqlx.Tx) (*Studio, error) {
	query := "SELECT * FROM studios WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryStudio(query, args, tx)
}

func (qb *StudioQueryBuilder) FindBySceneID(sceneID int) (*Studio, error) {
	query := "SELECT studios.* FROM studios JOIN scenes ON studios.id = scenes.studio_id WHERE scenes.id = ? LIMIT 1"
	args := []interface{}{sceneID}
	return qb.queryStudio(query, args, nil)
}

func (qb *StudioQueryBuilder) FindByName(name string, tx *sqlx.Tx) (*Studio, error) {
	query := "SELECT * FROM studios WHERE name = ? LIMIT 1"
	args := []interface{}{name}
	return qb.queryStudio(query, args, tx)
}

func (qb *StudioQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT studios.id FROM studios"), nil)
}

func (qb *StudioQueryBuilder) queryStudio(query string, args []interface{}, tx *sqlx.Tx) (*Studio, error) {
	results, err := qb.queryStudios(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *StudioQueryBuilder) queryStudios(query string, args []interface{}, tx *sqlx.Tx) ([]*Studio, error) {
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, args...)
	} else {
		rows, err = database.DB.Queryx(query, args...)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	studios := make([]*Studio, 0)
	for rows.Next() {
		studio := Studio{}
		if err := rows.StructScan(&studio); err != nil {
			return nil, err
		}
		studios = append(studios, &studio)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return studios, nil
}
