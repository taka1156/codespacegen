package assemble

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
	"codespacegen/internal/resolve"
)

type ResolveCodespaceConfig struct {
	codeSpaceConfigResolver resolve.CodeSpaceConfigResolver
}

func NewResolveCodespaceConfig(
	codeSpaceConfigResolver resolve.CodeSpaceConfigResolver,
) *ResolveCodespaceConfig {
	return &ResolveCodespaceConfig{
		codeSpaceConfigResolver: codeSpaceConfigResolver,
	}
}

func (rcc *ResolveCodespaceConfig) Resolve(cliConfig *entity.CliConfig, overrides map[string]json.RawMessage, defaultTimezone string, defaultImage string) (*entity.CodespaceConfig, error) {
	resolvedValues, err := rcc.resolveCoreValues(cliConfig)
	if err != nil {
		return nil, err
	}

	resolvedEntry, err := rcc.resolveEntry(resolvedValues.Language, cliConfig, overrides, defaultImage)
	if err != nil {
		return nil, err
	}

	resolvedTimezone, err := rcc.codeSpaceConfigResolver.ResolveTimezone(cliConfig.TimezoneValue(), resolvedEntry.Timezone, defaultTimezone)
	if err != nil {
		return nil, err
	}

	return rcc.buildCodespaceConfig(cliConfig, resolvedValues, resolvedEntry, resolvedTimezone), nil
}
