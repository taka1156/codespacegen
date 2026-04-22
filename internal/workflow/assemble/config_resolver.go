package assemble

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

// ConfigResolver abstracts the interactive resolution and config merge operations
// required to build a CodespaceConfig.
type ConfigResolver interface {
	ResolveProjectName(explicitProjectName string) (string, error)
	ResolveLanguage(explicitLanguage string) (string, error)
	ResolveWorkspaceFolder(explicitWorkspaceFolder string) (string, error)
	ResolveServiceName(explicitServiceName string) (string, error)
	ResolvePortMapping(explicitPort string) (string, error)
	ResolveTimezone(explicitTimezone string, configTimezone string, defaultTimezone string) (string, error)
	MergeLanguageEntries(overrides map[string]json.RawMessage) (map[string]entity.JsonEntry, error)
	ResolveBaseImage(language string, explicitBaseImage string, imageConfig string, jsonEntries map[string]entity.JsonEntry, defaultImage string) (entity.JsonEntry, error)
}
