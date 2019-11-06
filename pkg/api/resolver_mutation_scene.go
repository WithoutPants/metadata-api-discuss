package api

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stashdb/pkg/database"
	"github.com/stashapp/stashdb/pkg/manager"
	"github.com/stashapp/stashdb/pkg/models"
)

func (r *mutationResolver) SceneCreate(ctx context.Context, input models.SceneCreateInput) (*models.Scene, error) {
	if err := validateModify(ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mutationResolver) SceneUpdate(ctx context.Context, input models.SceneUpdateInput) (*models.Scene, error) {
	if err := validateModify(ctx); err != nil {
		return nil, err
	}

	// Start the transaction and save the scene
	tx := database.DB.MustBeginTx(ctx, nil)

	ret, err := r.sceneUpdate(input, tx)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) sceneUpdate(input models.SceneUpdateInput, tx *sqlx.Tx) (*models.Scene, error) {
	// Populate scene from the input
	sceneID, _ := strconv.Atoi(input.ID)
	updatedTime := time.Now()
	updatedScene := models.ScenePartial{
		ID:        sceneID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}
	if input.Title != nil {
		updatedScene.Title = &sql.NullString{String: *input.Title, Valid: true}
	}
	if input.Details != nil {
		updatedScene.Details = &sql.NullString{String: *input.Details, Valid: true}
	}
	if input.URL != nil {
		updatedScene.URL = &sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Date != nil {
		updatedScene.Date = &models.SQLiteDate{String: *input.Date, Valid: true}
	}

	if input.StudioID != nil {
		studioID, _ := strconv.ParseInt(*input.StudioID, 10, 64)
		updatedScene.StudioID = &sql.NullInt64{Int64: studioID, Valid: true}
	} else {
		// studio must be nullable
		updatedScene.StudioID = &sql.NullInt64{Valid: false}
	}

	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()
	scene, err := qb.Update(updatedScene, tx)
	if err != nil {
		return nil, err
	}

	// Save the performers
	var performerJoins []models.PerformersScenes
	for _, pid := range input.PerformerIds {
		performerID, _ := strconv.Atoi(pid)
		performerJoin := models.PerformersScenes{
			PerformerID: performerID,
			SceneID:     sceneID,
		}
		performerJoins = append(performerJoins, performerJoin)
	}
	if err := jqb.UpdatePerformersScenes(sceneID, performerJoins, tx); err != nil {
		return nil, err
	}

	// Save the tags
	var tagJoins []models.ScenesTags
	for _, tid := range input.TagIds {
		tagID, _ := strconv.Atoi(tid)
		tagJoin := models.ScenesTags{
			SceneID: sceneID,
			TagID:   tagID,
		}
		tagJoins = append(tagJoins, tagJoin)
	}
	if err := jqb.UpdateScenesTags(sceneID, tagJoins, tx); err != nil {
		return nil, err
	}

	return scene, nil
}

func (r *mutationResolver) SceneDestroy(ctx context.Context, input models.SceneDestroyInput) (bool, error) {
	if err := validateModify(ctx); err != nil {
		return false, err
	}

	tx := database.DB.MustBeginTx(ctx, nil)

	sceneID, _ := strconv.Atoi(input.ID)
	err := manager.DestroyScene(sceneID, tx)

	if err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
