package api

import (
	"context"
	"errors"
	"strconv"

	"github.com/stashapp/stashdb/pkg/models"
)

func (r *queryResolver) FindStudio(ctx context.Context, id *string, name *string) (*models.Studio, error) {
	if err := validateRead(ctx); err != nil {
		return nil, err
	}

	qb := models.NewStudioQueryBuilder()

	if id != nil {
		idInt, _ := strconv.Atoi(*id)
		return qb.Find(idInt, nil)
	}
	if name != nil {
		return qb.FindByName(*name, nil)
	}

	return nil, errors.New("Must provide id or name")
}

func (r *queryResolver) FindChildStudios(ctx context.Context, id string) ([]*models.Studio, error) {
	panic("not implemented")
}

func (r *queryResolver) QueryStudios(ctx context.Context, studioFilter *models.StudioFilterType, filter *models.QuerySpec) (*models.QueryStudiosResultType, error) {
	panic("not implemented")
}
