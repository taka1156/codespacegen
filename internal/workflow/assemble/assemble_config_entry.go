package assemble

import (
	"codespacegen/internal/domain/entity"
)

func (acc *AssembleCodespaceConfig) resolveEntry(language string, ClientConfig entity.ClientConfig, jsonConfig entity.JsonConfig, defaultImage string) (entity.LangEntry, error) {
	mergedImages, err := acc.CodespaceConfigResolver.MergeLanguageEntries(jsonConfig.Common, jsonConfig.Langs)
	if err != nil {
		return entity.LangEntry{}, err
	}

	return acc.CodespaceConfigResolver.ResolveBaseImage(language, ClientConfig.BaseImageValue(), mergedImages, defaultImage)
}
