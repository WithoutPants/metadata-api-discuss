package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stashdb/pkg/database"
)

const scenesForPerformerQuery = `
SELECT scenes.* FROM scenes
LEFT JOIN performers_scenes as performers_join on performers_join.scene_id = scenes.id
LEFT JOIN performers on performers_join.performer_id = performers.id
WHERE performers.id = ?
GROUP BY scenes.id
`

const scenesForStudioQuery = `
SELECT scenes.* FROM scenes
JOIN studios ON studios.id = scenes.studio_id
WHERE studios.id = ?
GROUP BY scenes.id
`

const scenesForTagQuery = `
SELECT scenes.* FROM scenes
LEFT JOIN scenes_tags as tags_join on tags_join.scene_id = scenes.id
LEFT JOIN tags on tags_join.tag_id = tags.id
WHERE tags.id = ?
GROUP BY scenes.id
`

type SceneQueryBuilder struct{}

func NewSceneQueryBuilder() SceneQueryBuilder {
	return SceneQueryBuilder{}
}

func (qb *SceneQueryBuilder) Create(newScene Scene, tx *sqlx.Tx) (*Scene, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO scenes (checksum, path, title, details, url, date, rating, size, duration, video_codec,
                    			    audio_codec, width, height, framerate, bitrate, studio_id, created_at, updated_at)
				VALUES (:checksum, :path, :title, :details, :url, :date, :rating, :size, :duration, :video_codec,
				        :audio_codec, :width, :height, :framerate, :bitrate, :studio_id, :created_at, :updated_at)
		`,
		newScene,
	)
	if err != nil {
		return nil, err
	}
	sceneID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err := tx.Get(&newScene, `SELECT * FROM scenes WHERE id = ? LIMIT 1`, sceneID); err != nil {
		return nil, err
	}
	return &newScene, nil
}

func (qb *SceneQueryBuilder) Update(updatedScene ScenePartial, tx *sqlx.Tx) (*Scene, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE scenes SET `+SQLGenKeysPartial(updatedScene)+` WHERE scenes.id = :id`,
		updatedScene,
	)
	if err != nil {
		return nil, err
	}

	return qb.find(updatedScene.ID, tx)
}

func (qb *SceneQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	return executeDeleteQuery("scenes", id, tx)
}
func (qb *SceneQueryBuilder) Find(id int) (*Scene, error) {
	return qb.find(id, nil)
}

func (qb *SceneQueryBuilder) find(id int, tx *sqlx.Tx) (*Scene, error) {
	query := "SELECT * FROM scenes WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryScene(query, args, tx)
}

func (qb *SceneQueryBuilder) FindByChecksum(checksum string) (*Scene, error) {
	query := "SELECT * FROM scenes WHERE checksum = ? LIMIT 1"
	args := []interface{}{checksum}
	return qb.queryScene(query, args, nil)
}

func (qb *SceneQueryBuilder) FindByPath(path string) (*Scene, error) {
	query := "SELECT * FROM scenes WHERE path = ? LIMIT 1"
	args := []interface{}{path}
	return qb.queryScene(query, args, nil)
}

func (qb *SceneQueryBuilder) FindByPerformerID(performerID int) ([]*Scene, error) {
	args := []interface{}{performerID}
	return qb.queryScenes(scenesForPerformerQuery, args, nil)
}

func (qb *SceneQueryBuilder) CountByPerformerID(performerID int) (int, error) {
	args := []interface{}{performerID}
	return runCountQuery(buildCountQuery(scenesForPerformerQuery), args)
}

func (qb *SceneQueryBuilder) FindByStudioID(studioID int) ([]*Scene, error) {
	args := []interface{}{studioID}
	return qb.queryScenes(scenesForStudioQuery, args, nil)
}

func (qb *SceneQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT scenes.id FROM scenes"), nil)
}

func (qb *SceneQueryBuilder) CountByStudioID(studioID int) (int, error) {
	args := []interface{}{studioID}
	return runCountQuery(buildCountQuery(scenesForStudioQuery), args)
}

func (qb *SceneQueryBuilder) CountByTagID(tagID int) (int, error) {
	args := []interface{}{tagID}
	return runCountQuery(buildCountQuery(scenesForTagQuery), args)
}

func (qb *SceneQueryBuilder) Wall(q *string) ([]*Scene, error) {
	s := ""
	if q != nil {
		s = *q
	}
	query := "SELECT scenes.* FROM scenes WHERE scenes.details LIKE '%" + s + "%' ORDER BY RANDOM() LIMIT 80"
	return qb.queryScenes(query, nil, nil)
}

func appendClause(clauses []string, clause string) []string {
	if clause != "" {
		return append(clauses, clause)
	}

	return clauses
}

func (qb *SceneQueryBuilder) queryScene(query string, args []interface{}, tx *sqlx.Tx) (*Scene, error) {
	results, err := qb.queryScenes(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *SceneQueryBuilder) queryScenes(query string, args []interface{}, tx *sqlx.Tx) ([]*Scene, error) {
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

	scenes := make([]*Scene, 0)
	for rows.Next() {
		scene := Scene{}
		if err := rows.StructScan(&scene); err != nil {
			return nil, err
		}
		scenes = append(scenes, &scene)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scenes, nil
}
