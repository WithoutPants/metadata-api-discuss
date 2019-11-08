package paths

import (
	"path/filepath"
)

type Paths struct {
	JSON *jsonPaths
}

func NewPaths() *Paths {
	p := Paths{}
	p.JSON = newJSONPaths()
	return &p
}

func GetConfigDirectory() string {
	return "."
}

func GetDefaultDatabaseFilePath() string {
	return filepath.Join(GetConfigDirectory(), "stash-go.sqlite")
}

func GetDefaultConfigFilePath() string {
	return filepath.Join(GetConfigDirectory(), "config.yml")
}

func GetSSLKey() string {
	return filepath.Join(GetConfigDirectory(), "stash.key")
}

func GetSSLCert() string {
	return filepath.Join(GetConfigDirectory(), "stash.crt")
}
