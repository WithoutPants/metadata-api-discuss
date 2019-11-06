package api

import (
	"context"
	"strconv"
	"time"

	"github.com/stashapp/stashdb/pkg/database"
	"github.com/stashapp/stashdb/pkg/models"
)

func (r *mutationResolver) TagCreate(ctx context.Context, input models.TagCreateInput) (*models.Tag, error) {
	if err := validateModify(ctx); err != nil {
		return nil, err
	}

	// Populate a new tag from the input
	currentTime := time.Now()
	newTag := models.Tag{
		Name:      input.Name,
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}

	// Start the transaction and save the studio
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewTagQueryBuilder()
	tag, err := qb.Create(newTag, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *mutationResolver) TagUpdate(ctx context.Context, input models.TagUpdateInput) (*models.Tag, error) {
	if err := validateModify(ctx); err != nil {
		return nil, err
	}

	// Populate tag from the input
	tagID, _ := strconv.Atoi(input.ID)
	updatedTag := models.Tag{
		ID:        tagID,
		Name:      input.Name,
		UpdatedAt: models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	// Start the transaction and save the tag
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewTagQueryBuilder()
	tag, err := qb.Update(updatedTag, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *mutationResolver) TagDestroy(ctx context.Context, input models.TagDestroyInput) (bool, error) {
	if err := validateModify(ctx); err != nil {
		return false, err
	}

	qb := models.NewTagQueryBuilder()
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
