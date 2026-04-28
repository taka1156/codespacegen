package assemble

import (
	"codespacegen/internal/domain/entity"
)

type AssembleCodespaceConfig struct {
	CodespaceConfigResolver ConfigResolver
}

func NewAssembleCodespaceConfig(
	CodespaceConfigResolver ConfigResolver,
) *AssembleCodespaceConfig {
	return &AssembleCodespaceConfig{
		CodespaceConfigResolver: CodespaceConfigResolver,
	}
}

func (acc *AssembleCodespaceConfig) Resolve(clientConfig entity.ClientConfig, defaultSetting entity.DefaultSetting, jsonConfig entity.JsonConfig) (*entity.CodespaceConfig, error) {
	resolvedValues, err := acc.resolveCoreValues(&clientConfig)
	if err != nil {
		return nil, err
	}

	resolvedEntry, err := acc.resolveEntry(resolvedValues.Language, clientConfig, jsonConfig, defaultSetting.Image)
	if err != nil {
		return nil, err
	}

	resolvedTimezone, err := acc.CodespaceConfigResolver.ResolveTimezone(clientConfig.TimezoneValue(), *resolvedEntry.Timezone, defaultSetting.Timezone)
	if err != nil {
		return nil, err
	}

	return acc.buildCodespaceConfig(clientConfig, defaultSetting, resolvedValues, resolvedEntry, resolvedTimezone), nil
}
