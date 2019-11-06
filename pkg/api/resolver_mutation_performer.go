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

func (r *mutationResolver) PerformerCreate(ctx context.Context, input models.PerformerCreateInput) (*models.Performer, error) {
	if err := validateModify(ctx); err != nil {
		return nil, err
	}

	var err error

	if err != nil {
		return nil, err
	}

	// Populate a new performer from the input
	currentTime := time.Now()
	newPerformer := models.Performer{
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}
	newPerformer.Name = input.Name
	
	if input.Gender != nil {
		newPerformer.Gender = sql.NullString{String: input.Gender.String(), Valid: true}
	}
	
	// TODO - set urls
	
	if input.Birthdate != nil {
		newPerformer.Birthdate = models.SQLiteDate{String: *input.Birthdate, Valid: true}
	}

	if input.Ethnicity != nil {
		newPerformer.Ethnicity = sql.NullString{String: *input.Ethnicity, Valid: true}
	}
	if input.Country != nil {
		newPerformer.Country = sql.NullString{String: *input.Country, Valid: true}
	}
	if input.EyeColor != nil {
		newPerformer.EyeColor = sql.NullString{String: *input.EyeColor, Valid: true}
	}
	if input.Height != nil {
		newPerformer.Height = sql.NullInt64{Int64: int64(*input.Height), Valid: true}
	}
	
	// TODO - handle measurements
	// TODO - handle fake tits
	// TODO - handle start/end career
	// TODO - handle tattoos and piercings

	// Start the transaction and save the performer
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewPerformerQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()
	performer, err := qb.Create(newPerformer, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Save the aliases
	var performerJoins []models.PerformerAliases
	for _, alias := range input.Aliases {
		performerJoin := models.PerformerAliases{
			PerformerID: performer.ID,
			Alias:       alias,
		}
		performerJoins = append(performerJoins, performerJoin)
	}
	if err := jqb.CreatePerformerAliases(performerJoins, tx); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return performer, nil
}

func (r *mutationResolver) PerformerUpdate(ctx context.Context, input models.PerformerUpdateInput) (*models.Performer, error) {
	if err := validateModify(ctx); err != nil {
		return nil, err
	}

	// Populate performer from the input
	performerID, _ := strconv.Atoi(input.ID)
	updatedPerformer := models.Performer{
		ID:        performerID,
		UpdatedAt: models.SQLiteTimestamp{Timestamp: time.Now()},
	}
	if input.Name != nil {
		updatedPerformer.Name = input.Name.String()
	}
	if input.URL != nil {
		updatedPerformer.URL = sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Birthdate != nil {
		updatedPerformer.Birthdate = models.SQLiteDate{String: *input.Birthdate, Valid: true}
	}
	if input.Ethnicity != nil {
		updatedPerformer.Ethnicity = sql.NullString{String: *input.Ethnicity, Valid: true}
	}
	if input.Country != nil {
		updatedPerformer.Country = sql.NullString{String: *input.Country, Valid: true}
	}
	if input.EyeColor != nil {
		updatedPerformer.EyeColor = sql.NullString{String: *input.EyeColor, Valid: true}
	}
	if input.Height != nil {
		updatedPerformer.Height = sql.NullString{String: *input.Height, Valid: true}
	}
	if input.Measurements != nil {
		updatedPerformer.Measurements = sql.NullString{String: *input.Measurements, Valid: true}
	}
	if input.FakeTits != nil {
		updatedPerformer.FakeTits = sql.NullString{String: *input.FakeTits, Valid: true}
	}
	if input.CareerLength != nil {
		updatedPerformer.CareerLength = sql.NullString{String: *input.CareerLength, Valid: true}
	}
	if input.Tattoos != nil {
		updatedPerformer.Tattoos = sql.NullString{String: *input.Tattoos, Valid: true}
	}
	if input.Piercings != nil {
		updatedPerformer.Piercings = sql.NullString{String: *input.Piercings, Valid: true}
	}
	if input.Twitter != nil {
		updatedPerformer.Twitter = sql.NullString{String: *input.Twitter, Valid: true}
	}
	if input.Instagram != nil {
		updatedPerformer.Instagram = sql.NullString{String: *input.Instagram, Valid: true}
	}

	// Start the transaction and save the performer
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewPerformerQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()
	performer, err := qb.Update(updatedPerformer, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Save the aliases
	var performerJoins []models.PerformerAliases
	for _, alias := range input.Aliases {
		performerJoin := models.PerformerAliases{
			PerformerID: performer.ID,
			Alias:       alias,
		}
		performerJoins = append(performerJoins, performerJoin)
	}
	if err := jqb.UpdatePerformerAliases(performer.ID, performerJoins, tx); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return performer, nil
}

func (r *mutationResolver) PerformerDestroy(ctx context.Context, input models.PerformerDestroyInput) (bool, error) {
	if err := validateModify(ctx); err != nil {
		return false, err
	}

	qb := models.NewPerformerQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	performerID, _ := strconv.Atoi(input.ID)
	if err := jqb.DestroyPerformerAliases(performerID, tx); err != nil {
		_ = tx.Rollback()
		return false, err
	}

	if err := qb.Destroy(input.ID, tx); err != nil {
		_ = tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}
