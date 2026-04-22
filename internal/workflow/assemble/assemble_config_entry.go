package assemble

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

func (acc *AssembleCodespaceConfig) resolveEntry(language string, cliConfig entity.CliConfig, overrides map[string]json.RawMessage, defaultImage string) (entity.JsonEntry, error) {
	mergedImages, err := acc.codeSpaceConfigResolver.MergeLanguageEntries(overrides)
	if err != nil {
		return entity.JsonEntry{}, err
	}

	return acc.codeSpaceConfigResolver.ResolveBaseImage(language, cliConfig.BaseImageValue(), cliConfig.ImageConfigValue(), mergedImages, defaultImage)
}
