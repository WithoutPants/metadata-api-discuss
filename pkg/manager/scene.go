package manager

import (
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stashdb/pkg/models"
)

func DestroyScene(sceneID int, tx *sqlx.Tx) error {
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	_, err := qb.Find(sceneID)
	if err != nil {
		return err
	}

	if err := jqb.DestroyScenesTags(sceneID, tx); err != nil {
		return err
	}

	if err := jqb.DestroyPerformersScenes(sceneID, tx); err != nil {
		return err
	}

	if err := jqb.DestroyScenesMarkers(sceneID, tx); err != nil {
		return err
	}

	if err := jqb.DestroyScenesGalleries(sceneID, tx); err != nil {
		return err
	}

	if err := qb.Destroy(strconv.Itoa(sceneID), tx); err != nil {
		return err
	}

	return nil
}
