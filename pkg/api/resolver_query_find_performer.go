package api

import (
	"context"

	"github.com/stashapp/stashdb/pkg/models"
)

func (r *queryResolver) FindPerformer(ctx context.Context, id string) (*models.Performer, error) {
	panic("not implemented")
}
func (r *queryResolver) QueryPerformers(ctx context.Context, performerFilter *models.PerformerFilterType, filter *models.QuerySpec) (*models.QueryPerformersResultType, error) {
	panic("not implemented")
}

// func (r *queryResolver) FindPerformer(ctx context.Context, id *string, name *string, includeAliases *bool) (*models.Performer, error) {
// 	if err := validateRead(ctx); err != nil {
// 		return nil, err
// 	}

// 	qb := models.NewPerformerQueryBuilder()

// 	if id != nil {
// 		idInt, _ := strconv.Atoi(*id)
// 		return qb.Find(idInt)
// 	}
// 	if name != nil {
// 		ret, err := qb.FindByName(*name, nil)

// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(ret) > 0 {
// 			return ret[0], err
// 		}

// 		if includeAliases != nil && *includeAliases {
// 			// try to get by alias
// 			ret, err = qb.FindByAlias(*name, nil)

// 			if len(ret) > 0 {
// 				return ret[0], err
// 			}
// 		}

// 		return nil, err
// 	}

// 	return nil, errors.New("Must provide id or name")
// }
