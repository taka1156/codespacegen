package assemble

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

type AssembleCodespaceConfig struct {
	codeSpaceConfigResolver ConfigResolver
}

func NewAssembleCodespaceConfig(
	codeSpaceConfigResolver ConfigResolver,
) *AssembleCodespaceConfig {
	return &AssembleCodespaceConfig{
		codeSpaceConfigResolver: codeSpaceConfigResolver,
	}
}

func (acc *AssembleCodespaceConfig) Resolve(cliConfig entity.CliConfig, defaultSetting entity.DefaultSetting, overrides map[string]json.RawMessage, defaultTimezone string, defaultImage string) (*entity.CodespaceConfig, error) {
	resolvedValues, err := acc.resolveCoreValues(&cliConfig)
	if err != nil {
		return nil, err
	}

	resolvedEntry, err := acc.resolveEntry(resolvedValues.Language, cliConfig, overrides, defaultImage)
	if err != nil {
		return nil, err
	}

	resolvedTimezone, err := acc.codeSpaceConfigResolver.ResolveTimezone(cliConfig.TimezoneValue(), resolvedEntry.Timezone, defaultTimezone)
	if err != nil {
		return nil, err
	}

	return acc.buildCodespaceConfig(cliConfig, defaultSetting, resolvedValues, resolvedEntry, resolvedTimezone), nil
}
