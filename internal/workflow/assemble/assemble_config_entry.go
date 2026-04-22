package assemble

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

func (rcc *ResolveCodespaceConfig) resolveEntry(language string, cliConfig *entity.CliConfig, overrides map[string]json.RawMessage, defaultImage string) (entity.JsonEntry, error) {
	mergedImages, err := rcc.codeSpaceConfigResolver.MergeLanguageEntries(overrides)
	if err != nil {
		return entity.JsonEntry{}, err
	}

	return rcc.codeSpaceConfigResolver.ResolveBaseImage(language, cliConfig.BaseImageValue(), cliConfig.ImageConfigValue(), mergedImages, defaultImage)
}
