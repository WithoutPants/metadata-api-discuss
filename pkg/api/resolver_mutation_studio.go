package api

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/stashapp/stashdb/pkg/database"
	"github.com/stashapp/stashdb/pkg/models"
	"github.com/stashapp/stashdb/pkg/utils"
)

func (r *mutationResolver) StudioCreate(ctx context.Context, input models.StudioCreateInput) (*models.Studio, error) {
	if err := validateModify(ctx); err != nil {
		return nil, err
	}

	var imageData []byte
	var err error

	// Process the base 64 encoded image string
	if input.Image != nil {
		_, imageData, err = utils.ProcessBase64Image(*input.Image)
		if err != nil {
			return nil, err
		}
	}

	// Populate a new studio from the input
	currentTime := time.Now()
	newStudio := models.Studio{
		Image:     imageData,
		Name:      sql.NullString{String: input.Name, Valid: true},
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}
	if input.URL != nil {
		newStudio.URL = sql.NullString{String: *input.URL, Valid: true}
	}

	// Start the transaction and save the studio
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewStudioQueryBuilder()
	studio, err := qb.Create(newStudio, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return studio, nil
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input models.StudioUpdateInput) (*models.Studio, error) {
	if err := validateModify(ctx); err != nil {
		return nil, err
	}

	// Populate studio from the input
	studioID, _ := strconv.Atoi(input.ID)
	updatedStudio := models.Studio{
		ID:        studioID,
		UpdatedAt: models.SQLiteTimestamp{Timestamp: time.Now()},
	}
	if input.Image != nil {
		_, imageData, err := utils.ProcessBase64Image(*input.Image)
		if err != nil {
			return nil, err
		}
		updatedStudio.Image = imageData
	}
	if input.Name != nil {
		updatedStudio.Name = sql.NullString{String: *input.Name, Valid: true}
	}
	if input.URL != nil {
		updatedStudio.URL = sql.NullString{String: *input.URL, Valid: true}
	}

	// Start the transaction and save the studio
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewStudioQueryBuilder()
	studio, err := qb.Update(updatedStudio, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return studio, nil
}

func (r *mutationResolver) StudioDestroy(ctx context.Context, input models.StudioDestroyInput) (bool, error) {
	if err := validateModify(ctx); err != nil {
		return false, err
	}

	qb := models.NewStudioQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)
	if err := qb.Destroy(input.ID, tx); err != nil {
		_ = tx.Rollback()
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}
