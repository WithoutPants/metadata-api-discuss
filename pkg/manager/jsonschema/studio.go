package jsonschema

import (
	"encoding/json"
	"fmt"
	"github.com/stashapp/stashdb/pkg/models"
	"os"
)

type Studio struct {
	Name      string          `json:"name,omitempty"`
	URL       string          `json:"url,omitempty"`
	Image     string          `json:"image,omitempty"`
	CreatedAt models.JSONTime `json:"created_at,omitempty"`
	UpdatedAt models.JSONTime `json:"updated_at,omitempty"`
}

func LoadStudioFile(filePath string) (*Studio, error) {
	var studio Studio
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&studio)
	if err != nil {
		return nil, err
	}
	return &studio, nil
}

func SaveStudioFile(filePath string, studio *Studio) error {
	if studio == nil {
		return fmt.Errorf("studio must not be nil")
	}
	return marshalToFile(filePath, studio)
}
