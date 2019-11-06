package api

import (
	"context"

	"github.com/stashapp/stashdb/pkg/models"
)

type performerResolver struct{ *Resolver }

func (r *performerResolver) Disambiguation(ctx context.Context, obj *models.Performer) (*string, error) {
	return resolveNullString(obj.Disambiguation)
}

func (r *performerResolver) Aliases(ctx context.Context, obj *models.Performer) ([]string, error) {
	qb := models.NewPerformerQueryBuilder()
	aliases, err := qb.GetAliases(obj.ID)

	if err != nil {
		return nil, err
	}

	var ret []string
	for _, alias := range aliases {
		if alias.Valid {
			ret = append(ret, alias.String)
		}
	}

	return ret, nil
}

func (r *performerResolver) Gender(ctx context.Context, obj *models.Performer) (*models.GenderEnum, error) {
	panic("not implemented")
}

func (r *performerResolver) Urls(ctx context.Context, obj *models.Performer) ([]*models.URL, error) {
	panic("not implemented")
}

func (r *performerResolver) Birthdate(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Ethnicity(ctx context.Context, obj *models.Performer) (*string, error) {
	return resolveNullString(obj.Ethnicity)
}

func (r *performerResolver) Country(ctx context.Context, obj *models.Performer) (*string, error) {
	return resolveNullString(obj.Country)
}

func (r *performerResolver) EyeColor(ctx context.Context, obj *models.Performer) (*string, error) {
	return resolveNullString(obj.EyeColor)
}

func (r *performerResolver) Height(ctx context.Context, obj *models.Performer) (*int, error) {
	return resolveNullInt64(obj.Height)
}
func (r *performerResolver) Measurements(ctx context.Context, obj *models.Performer) (*models.Measurements, error) {
	panic("not implemented")
}
func (r *performerResolver) FakeTits(ctx context.Context, obj *models.Performer) (*models.FakeTitsEnum, error) {
	panic("not implemented")
}
func (r *performerResolver) CareerStartYear(ctx context.Context, obj *models.Performer) (*int, error) {
	return resolveNullInt64(obj.CareerStartYear)
}
func (r *performerResolver) CareerEndYear(ctx context.Context, obj *models.Performer) (*int, error) {
	return resolveNullInt64(obj.CareerEndYear)
}
func (r *performerResolver) Tattoos(ctx context.Context, obj *models.Performer) ([]string, error) {
	panic("not implemented")
}
func (r *performerResolver) Piercings(ctx context.Context, obj *models.Performer) ([]string, error) {
	panic("not implemented")
}
