package api

import (
	"context"
	"errors"
	"strconv"

	"github.com/stashapp/stashdb/pkg/models"
)

func (r *queryResolver) FindTag(ctx context.Context, id *string, name *string) (*models.Tag, error) {
	if err := validateRead(ctx); err != nil {
		return nil, err
	}

	qb := models.NewTagQueryBuilder()

	if id != nil {
		idInt, _ := strconv.Atoi(*id)
		return qb.Find(idInt, nil)
	}
	if name != nil {
		return qb.FindByName(*name, nil)
	}

	return nil, errors.New("Must provide id or name")
}

func (r *queryResolver) QueryTags(ctx context.Context, tagFilter *models.TagFilterType, filter *models.QuerySpec) (*models.QueryTagsResultType, error) {
	panic("not implemented")
}
