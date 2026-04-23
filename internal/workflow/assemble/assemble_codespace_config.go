package assemble

import (
	"encoding/json"

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

func (acc *AssembleCodespaceConfig) Resolve(cliConfig entity.CliConfig, defaultSetting entity.DefaultSetting, jsonConfig map[string]json.RawMessage) (*entity.CodespaceConfig, error) {
	resolvedValues, err := acc.resolveCoreValues(&cliConfig)
	if err != nil {
		return nil, err
	}

	resolvedEntry, err := acc.resolveEntry(resolvedValues.Language, cliConfig, jsonConfig, defaultSetting.Image)
	if err != nil {
		return nil, err
	}

	resolvedTimezone, err := acc.CodespaceConfigResolver.ResolveTimezone(cliConfig.TimezoneValue(), resolvedEntry.Timezone, defaultSetting.Timezone)
	if err != nil {
		return nil, err
	}

	return acc.buildCodespaceConfig(cliConfig, defaultSetting, resolvedValues, resolvedEntry, resolvedTimezone), nil
}
