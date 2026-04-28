package assemble

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

func (acc *AssembleCodespaceConfig) resolveEntry(language string, ClientConfig entity.ClientConfig, overrides map[string]json.RawMessage, defaultImage string) (entity.LangEntry, error) {
	mergedImages, err := acc.CodespaceConfigResolver.MergeLanguageEntries(overrides)
	if err != nil {
		return entity.LangEntry{}, err
	}

	return acc.CodespaceConfigResolver.ResolveBaseImage(language, ClientConfig.BaseImageValue(), ClientConfig.ImageConfigValue(), mergedImages, defaultImage)
}
