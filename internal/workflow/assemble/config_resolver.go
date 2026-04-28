package assemble

import (
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
	MergeLanguageEntries(commonEntry *entity.CommonEntry, langEntries map[string]*entity.LangEntry) (map[string]entity.LangEntry, error)
	ResolveBaseImage(language string, explicitBaseImage string, jsonEntries map[string]entity.LangEntry, defaultImage string) (entity.LangEntry, error)
}
