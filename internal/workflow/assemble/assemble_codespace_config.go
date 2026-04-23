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

func (acc *AssembleCodespaceConfig) Resolve(ClientConfig entity.ClientConfig, defaultSetting entity.DefaultSetting, overrides map[string]json.RawMessage, defaultTimezone string, defaultImage string) (*entity.CodespaceConfig, error) {
	resolvedValues, err := acc.resolveCoreValues(&ClientConfig)
	if err != nil {
		return nil, err
	}

	resolvedEntry, err := acc.resolveEntry(resolvedValues.Language, ClientConfig, overrides, defaultImage)
	if err != nil {
		return nil, err
	}

	resolvedTimezone, err := acc.CodespaceConfigResolver.ResolveTimezone(ClientConfig.TimezoneValue(), resolvedEntry.Timezone, defaultTimezone)
	if err != nil {
		return nil, err
	}

	return acc.buildCodespaceConfig(ClientConfig, defaultSetting, resolvedValues, resolvedEntry, resolvedTimezone), nil
}
