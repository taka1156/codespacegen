package assemble

import (
	"github.com/taka1156/codespacegen/internal/domain/entity"
)

type AssembleCodespaceConfig struct {
	CodespacegenPrompter CodespacegenPrompter
}

func NewAssembleCodespaceConfig(
	codespacegenPrompter CodespacegenPrompter,
) *AssembleCodespaceConfig {
	return &AssembleCodespaceConfig{
		CodespacegenPrompter: codespacegenPrompter,
	}
}

func (acc *AssembleCodespaceConfig) Resolve(clientConfig entity.ClientConfig, defaultSetting entity.DefaultSetting, jsonConfig entity.JsonConfig) (*entity.CodespaceConfig, error) {
	resolvedValues, err := acc.resolvePromptValues(&clientConfig)
	if err != nil {
		return nil, err
	}

	resolvedEntries, err := acc.resolveMergedEntry(jsonConfig)
	if err != nil {
		return nil, err
	}

	return acc.buildCodespaceConfig(clientConfig, defaultSetting, resolvedValues, resolvedEntries)
}
