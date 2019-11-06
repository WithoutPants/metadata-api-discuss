package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stashdb/pkg/database"
)

type TagQueryBuilder struct{}

func NewTagQueryBuilder() TagQueryBuilder {
	return TagQueryBuilder{}
}

func (qb *TagQueryBuilder) Create(newTag Tag, tx *sqlx.Tx) (*Tag, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO tags (name, created_at, updated_at)
				VALUES (:name, :created_at, :updated_at)
		`,
		newTag,
	)
	if err != nil {
		return nil, err
	}
	studioID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&newTag, `SELECT * FROM tags WHERE id = ? LIMIT 1`, studioID); err != nil {
		return nil, err
	}
	return &newTag, nil
}

func (qb *TagQueryBuilder) Update(updatedTag Tag, tx *sqlx.Tx) (*Tag, error) {
	ensureTx(tx)
	query := `UPDATE tags SET ` + SQLGenKeys(updatedTag) + ` WHERE tags.id = :id`
	_, err := tx.NamedExec(
		query,
		updatedTag,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&updatedTag, `SELECT * FROM tags WHERE id = ? LIMIT 1`, updatedTag.ID); err != nil {
		return nil, err
	}
	return &updatedTag, nil
}

func (qb *TagQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	// delete tag from scenes and markers first
	_, err := tx.Exec("DELETE FROM scenes_tags WHERE tag_id = ?", id)
	if err != nil {
		return err
	}

	return executeDeleteQuery("tags", id, tx)
}

func (qb *TagQueryBuilder) Find(id int, tx *sqlx.Tx) (*Tag, error) {
	query := "SELECT * FROM tags WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryTag(query, args, tx)
}

func (qb *TagQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]*Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN scenes_tags as scenes_join on scenes_join.tag_id = tags.id
		LEFT JOIN scenes on scenes_join.scene_id = scenes.id
		WHERE scenes.id = ?
		GROUP BY tags.id
	`
	query += qb.getTagSort()
	args := []interface{}{sceneID}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) FindByName(name string, tx *sqlx.Tx) (*Tag, error) {
	query := "SELECT * FROM tags WHERE name = ? LIMIT 1"
	args := []interface{}{name}
	return qb.queryTag(query, args, tx)
}

func (qb *TagQueryBuilder) FindByNames(names []string, tx *sqlx.Tx) ([]*Tag, error) {
	query := "SELECT * FROM tags WHERE name IN " + getInBinding(len(names))
	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT tags.id FROM tags"), nil)
}

func (qb *TagQueryBuilder) All() ([]*Tag, error) {
	return qb.queryTags(selectAll("tags")+qb.getTagSort(), nil, nil)
}

func (qb *TagQueryBuilder) getTagSort() string {
	var sort string
	var direction string
	sort = "name"
	direction = "ASC"
	return getSort(sort, direction, "tags")
}

func (qb *TagQueryBuilder) queryTag(query string, args []interface{}, tx *sqlx.Tx) (*Tag, error) {
	results, err := qb.queryTags(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *TagQueryBuilder) queryTags(query string, args []interface{}, tx *sqlx.Tx) ([]*Tag, error) {
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

	tags := make([]*Tag, 0)
	for rows.Next() {
		tag := Tag{}
		if err := rows.StructScan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
