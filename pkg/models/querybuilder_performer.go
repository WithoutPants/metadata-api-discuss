package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stashdb/pkg/database"
)

type PerformerQueryBuilder struct{}

func NewPerformerQueryBuilder() PerformerQueryBuilder {
	return PerformerQueryBuilder{}
}

func (qb *PerformerQueryBuilder) Create(newPerformer Performer, tx *sqlx.Tx) (*Performer, error) {
	ensureTx(tx)
	fields, values := SQLGenKeysCreate(newPerformer)
	result, err := tx.NamedExec(
		`INSERT INTO performers (`+fields+`)
				VALUES (`+values+`)
		`,
		newPerformer,
	)
	if err != nil {
		return nil, err
	}
	performerID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&newPerformer, `SELECT * FROM performers WHERE id = ? LIMIT 1`, performerID); err != nil {
		return nil, err
	}
	return &newPerformer, nil
}

func (qb *PerformerQueryBuilder) Update(updatedPerformer Performer, tx *sqlx.Tx) (*Performer, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE performers SET `+SQLGenKeys(updatedPerformer)+` WHERE performers.id = :id`,
		updatedPerformer,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&updatedPerformer, `SELECT * FROM performers WHERE id = ? LIMIT 1`, updatedPerformer.ID); err != nil {
		return nil, err
	}
	return &updatedPerformer, nil
}

func (qb *PerformerQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	_, err := tx.Exec("DELETE FROM performers_scenes WHERE performer_id = ?", id)
	if err != nil {
		return err
	}

	return executeDeleteQuery("performers", id, tx)
}

func (qb *PerformerQueryBuilder) Find(id int) (*Performer, error) {
	query := "SELECT * FROM performers WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	results, err := qb.queryPerformers(query, args, nil)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *PerformerQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]*Performer, error) {
	query := `
		SELECT performers.* FROM performers
		LEFT JOIN performers_scenes as scenes_join on scenes_join.performer_id = performers.id
		LEFT JOIN scenes on scenes_join.scene_id = scenes.id
		WHERE scenes.id = ?
		GROUP BY performers.id
	`
	args := []interface{}{sceneID}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByNames(names []string, tx *sqlx.Tx) ([]*Performer, error) {
	query := "SELECT * FROM performers WHERE name IN " + getInBinding(len(names))
	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByAliases(names []string, tx *sqlx.Tx) ([]*Performer, error) {
	query := `SELECT performers.* FROM performers
		left join performer_aliases on performers.id = performer_aliases.performer_id
		WHERE performer_aliases.alias IN ` + getInBinding(len(names))

	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByName(name string, tx *sqlx.Tx) ([]*Performer, error) {
	query := "SELECT * FROM performers WHERE upper(name) = upper(?)"
	var args []interface{}
	args = append(args, name)
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByAlias(name string, tx *sqlx.Tx) ([]*Performer, error) {
	query := `SELECT performers.* FROM performers
		left join performer_aliases on performers.id = performer_aliases.performer_id
		WHERE upper(performer_aliases.alias) = UPPER(?)`

	var args []interface{}
	args = append(args, name)
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT performers.id FROM performers"), nil)
}

func (qb *PerformerQueryBuilder) queryPerformers(query string, args []interface{}, tx *sqlx.Tx) ([]*Performer, error) {
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

	performers := make([]*Performer, 0)
	for rows.Next() {
		performer := Performer{}
		if err := rows.StructScan(&performer); err != nil {
			return nil, err
		}
		performers = append(performers, &performer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return performers, nil
}

func (qb *PerformerQueryBuilder) GetAliases(id int) ([]sql.NullString, error) {
	query := "SELECT alias FROM performer_aliases WHERE performer_id = ?"
	args := []interface{}{id}

	var rows *sqlx.Rows
	var err error
	rows, err = database.DB.Queryx(query, args...)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	aliases := make([]sql.NullString, 0)
	for rows.Next() {
		alias := sql.NullString{}

		if err := rows.Scan(&alias); err != nil {
			return nil, err
		}
		aliases = append(aliases, alias)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return aliases, nil
}
