package assemble

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

func (rc *ResolveCodespaceConfig) resolveEntry(language string, cliConfig *entity.CliConfig, overrides map[string]json.RawMessage) (entity.JsonEntry, error) {
	mergedImages, err := rc.codeSpaceConfigResolver.MergeLanguageEntries(overrides)
	if err != nil {
		return entity.JsonEntry{}, err
	}

	return rc.codeSpaceConfigResolver.ResolveBaseImage(language, *cliConfig.BaseImage, *cliConfig.ImageConfig, mergedImages)
}
