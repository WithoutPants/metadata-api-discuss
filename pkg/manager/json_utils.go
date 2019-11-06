package manager

import (
	"github.com/stashapp/stashdb/pkg/manager/jsonschema"
)

type jsonUtils struct{}

func (jp *jsonUtils) getMappings() (*jsonschema.Mappings, error) {
	return jsonschema.LoadMappingsFile(instance.Paths.JSON.MappingsFile)
}

func (jp *jsonUtils) saveMappings(mappings *jsonschema.Mappings) error {
	return jsonschema.SaveMappingsFile(instance.Paths.JSON.MappingsFile, mappings)
}

func (jp *jsonUtils) getPerformer(checksum string) (*jsonschema.Performer, error) {
	return jsonschema.LoadPerformerFile(instance.Paths.JSON.PerformerJSONPath(checksum))
}

func (jp *jsonUtils) savePerformer(checksum string, performer *jsonschema.Performer) error {
	return jsonschema.SavePerformerFile(instance.Paths.JSON.PerformerJSONPath(checksum), performer)
}

func (jp *jsonUtils) getStudio(checksum string) (*jsonschema.Studio, error) {
	return jsonschema.LoadStudioFile(instance.Paths.JSON.StudioJSONPath(checksum))
}

func (jp *jsonUtils) saveStudio(checksum string, studio *jsonschema.Studio) error {
	return jsonschema.SaveStudioFile(instance.Paths.JSON.StudioJSONPath(checksum), studio)
}

func (jp *jsonUtils) getScene(checksum string) (*jsonschema.Scene, error) {
	return jsonschema.LoadSceneFile(instance.Paths.JSON.SceneJSONPath(checksum))
}

func (jp *jsonUtils) saveScene(checksum string, scene *jsonschema.Scene) error {
	return jsonschema.SaveSceneFile(instance.Paths.JSON.SceneJSONPath(checksum), scene)
}
