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

func (rc *ResolveCodespaceConfig) Resolve(cliConfig *entity.CliConfig, overrides map[string]json.RawMessage) (*entity.CodespaceConfig, error) {
	resolvedValues, err := rc.resolveCoreValues(cliConfig)
	if err != nil {
		return nil, err
	}

	resolvedEntry, err := rc.resolveEntry(resolvedValues.Language, cliConfig, overrides)
	if err != nil {
		return nil, err
	}

	resolvedTimezone, err := rc.codeSpaceConfigResolver.ResolveTimezone(*cliConfig.Timezone, resolvedEntry.Timezone)
	if err != nil {
		return nil, err
	}

	return rc.buildCodespaceConfig(cliConfig, resolvedValues, resolvedEntry, resolvedTimezone), nil
}
