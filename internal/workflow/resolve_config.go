package workflow

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
	"codespacegen/internal/resolve"
)

type ResolveConfig struct {
	mergeLanguageResolver   resolve.MergeLanguageResolver
	codeSpaceConfigResolver resolve.CodeSpaceConfigResolver
}

func NewResolveConfig(
	mergeLanguageResolver resolve.MergeLanguageResolver,
	codeSpaceConfigResolver resolve.CodeSpaceConfigResolver,
) *ResolveConfig {
	return &ResolveConfig{
		mergeLanguageResolver:   mergeLanguageResolver,
		codeSpaceConfigResolver: codeSpaceConfigResolver,
	}
}

func (rc *ResolveConfig) Resolve(cliConfig *entity.CliConfig, jsonEntries map[string]entity.JsonEntry, overrides map[string]json.RawMessage) (*entity.CodespaceConfig, error) {
	resolvedValues, err := rc.resolveCoreValues(cliConfig)
	if err != nil {
		return nil, err
	}

	resolvedEntry, err := rc.resolveEntry(resolvedValues.Language, cliConfig, jsonEntries, overrides)
	if err != nil {
		return nil, err
	}

	resolvedTimezone, err := rc.codeSpaceConfigResolver.ResolveTimezone(*cliConfig.Timezone, resolvedEntry.Timezone)
	if err != nil {
		return nil, err
	}

	return rc.buildCodespaceConfig(cliConfig, resolvedValues, resolvedEntry, resolvedTimezone), nil
}
