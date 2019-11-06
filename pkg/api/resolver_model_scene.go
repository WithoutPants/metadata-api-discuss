package api

import (
	"context"

	"github.com/stashapp/stashdb/pkg/models"
)

type sceneResolver struct{ *Resolver }

func (r *sceneResolver) Title(ctx context.Context, obj *models.Scene) (*string, error) {
	return resolveNullString(obj.Title)
}

func (r *sceneResolver) Details(ctx context.Context, obj *models.Scene) (*string, error) {
	return resolveNullString(obj.Details)
}

func (r *sceneResolver) URL(ctx context.Context, obj *models.Scene) (*string, error) {
	return resolveNullString(obj.URL)
}

func (r *sceneResolver) Date(ctx context.Context, obj *models.Scene) (*string, error) {
	return resolveSQLiteDate(obj.Date)
}

func (r *sceneResolver) Studio(ctx context.Context, obj *models.Scene) (*models.Studio, error) {
	qb := models.NewStudioQueryBuilder()
	return qb.FindBySceneID(obj.ID)
}

func (r *sceneResolver) Tags(ctx context.Context, obj *models.Scene) ([]*models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	return qb.FindBySceneID(obj.ID, nil)
}

func (r *sceneResolver) Performers(ctx context.Context, obj *models.Scene) ([]*models.Performer, error) {
	qb := models.NewPerformerQueryBuilder()
	return qb.FindBySceneID(obj.ID, nil)
}

func (r *sceneResolver) Checksums(ctx context.Context, obj *models.Scene) ([]string, error) {
	panic("not implemented")
}
