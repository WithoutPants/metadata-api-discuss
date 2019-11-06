package api

import (
	"database/sql"

	"github.com/stashapp/stashdb/pkg/models"
	"github.com/stashapp/stashdb/pkg/utils"
)

func resolveNullString(value sql.NullString) (*string, error) {
	if value.Valid {
		return &value.String, nil
	}
	return nil, nil
}

func resolveSQLiteDate(value models.SQLiteDate) (*string, error) {
	if value.Valid {
		result := utils.GetYMDFromDatabaseDate(value.String)
		return &result, nil
	}
	return nil, nil
}

func resolveNullInt64(value sql.NullInt64) (*int, error) {
	if value.Valid {
		result := int(value.Int64)
		return &result, nil
	}
	return nil, nil
}
