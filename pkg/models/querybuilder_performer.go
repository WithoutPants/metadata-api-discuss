package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stashapp/stashdb/pkg/database"
)

type PerformerQueryBuilder struct{}

const performerTable = "performers"
const performerAliasesJoinTable = "performer_aliases"
const performerUrlsJoinTable = "performer_urls"
const performerTattoosJoinTable = "performer_tattoos"
const performerPiercingsJoinTable = "performer_piercings"
const performerJoinKey = "performer_id"

func NewPerformerQueryBuilder() PerformerQueryBuilder {
	return PerformerQueryBuilder{}
}

func (qb *PerformerQueryBuilder) Create(newPerformer Performer, tx *sqlx.Tx) (*Performer, error) {
	performerID, err := insertObject(tx, performerTable, newPerformer)

	if err != nil {
		return nil, errors.Wrap(err, "Error creating performer")
	}

	if err := getByID(tx, performerTable, performerID, &newPerformer); err != nil {
		return nil, errors.Wrap(err, "Error getting performer after create")
	}
	return &newPerformer, nil
}

func (qb *PerformerQueryBuilder) Update(updatedPerformer Performer, tx *sqlx.Tx) (*Performer, error) {
	err := updateObjectByID(tx, performerTable, updatedPerformer)

	if err != nil {
		return nil, errors.Wrap(err, "Error updating performer")
	}

	if err := getByID(tx, performerTable, updatedPerformer.ID, &updatedPerformer); err != nil {
		return nil, errors.Wrap(err, "Error getting performer after update")
	}
	return &updatedPerformer, nil
}

func (qb *PerformerQueryBuilder) Destroy(id int64, tx *sqlx.Tx) error {
	return executeDeleteQuery(performerTable, id, tx)
}

func (qb *PerformerQueryBuilder) CreateAliases(newJoins []PerformerAliases, tx *sqlx.Tx) error {
	return insertJoins(tx, performerAliasesJoinTable, newJoins)
}

func (qb *PerformerQueryBuilder) UpdateAliases(performerID int64, updatedJoins []PerformerAliases, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	err := deleteObjectsByColumn(tx, performerAliasesJoinTable, performerJoinKey, performerID)
	if err != nil {
		return err
	}
	return qb.CreateAliases(updatedJoins, tx)
}

func (qb *PerformerQueryBuilder) CreateUrls(newJoins []PerformerUrls, tx *sqlx.Tx) error {
	return insertJoins(tx, performerUrlsJoinTable, newJoins)
}

func (qb *PerformerQueryBuilder) UpdateUrls(performerID int64, updatedJoins []PerformerUrls, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	err := deleteObjectsByColumn(tx, performerUrlsJoinTable, performerJoinKey, performerID)
	if err != nil {
		return err
	}
	return qb.CreateUrls(updatedJoins, tx)
}

func (qb *PerformerQueryBuilder) CreateTattoos(newJoins []PerformerBodyMods, tx *sqlx.Tx) error {
	return insertJoins(tx, performerTattoosJoinTable, newJoins)
}

func (qb *PerformerQueryBuilder) UpdateTattoos(performerID int64, updatedJoins []PerformerBodyMods, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	err := deleteObjectsByColumn(tx, performerTattoosJoinTable, performerJoinKey, performerID)
	if err != nil {
		return err
	}
	return qb.CreateTattoos(updatedJoins, tx)
}

func (qb *PerformerQueryBuilder) CreatePiercings(newJoins []PerformerBodyMods, tx *sqlx.Tx) error {
	return insertJoins(tx, performerPiercingsJoinTable, newJoins)
}

func (qb *PerformerQueryBuilder) UpdatePiercings(performerID int64, updatedJoins []PerformerBodyMods, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	err := deleteObjectsByColumn(tx, performerPiercingsJoinTable, performerJoinKey, performerID)
	if err != nil {
		return err
	}
	return qb.CreateTattoos(updatedJoins, tx)
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

func (qb *PerformerQueryBuilder) GetAliases(id int64) ([]string, error) {
	query := "SELECT alias FROM performer_aliases WHERE performer_id = ?"
	args := []interface{}{id}

	var rows *sqlx.Rows
	var err error
	rows, err = database.DB.Queryx(query, args...)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	aliases := make([]string, 0)
	for rows.Next() {
		var alias string

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

func (qb *PerformerQueryBuilder) GetUrls(id int64) ([]PerformerUrls, error) {
	query := "SELECT url, type FROM performer_urls WHERE performer_id = ?"
	args := []interface{}{id}

	var rows *sqlx.Rows
	var err error
	rows, err = database.DB.Queryx(query, args...)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	urls := make([]PerformerUrls, 0)
	for rows.Next() {
		var performerUrl PerformerUrls

		if err := rows.Scan(&performerUrl); err != nil {
			return nil, err
		}
		urls = append(urls, performerUrl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func translateBodyMods(rows *sqlx.Rows) ([]PerformerBodyMods, error) {
	ret := make([]PerformerBodyMods, 0)
	for rows.Next() {
		var performerBodyMod PerformerBodyMods

		if err := rows.Scan(&performerBodyMod); err != nil {
			return nil, err
		}
		ret = append(ret, performerBodyMod)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *PerformerQueryBuilder) GetTattoos(id int64) ([]PerformerBodyMods, error) {
	query := "SELECT location, description FROM performer_tattoos WHERE performer_id = ?"
	args := []interface{}{id}

	var rows *sqlx.Rows
	var err error
	rows, err = database.DB.Queryx(query, args...)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	return translateBodyMods(rows)
}

func (qb *PerformerQueryBuilder) GetPiercings(id int64) ([]PerformerBodyMods, error) {
	query := "SELECT location, description FROM performer_piercings WHERE performer_id = ?"
	args := []interface{}{id}

	var rows *sqlx.Rows
	var err error
	rows, err = database.DB.Queryx(query, args...)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	return translateBodyMods(rows)
}
