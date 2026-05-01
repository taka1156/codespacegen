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
	var resolvedValues resolvedCoreValues
	var err error

	resolvedEntries, err := acc.resolveMergedEntry(jsonConfig)
	if err != nil {
		return nil, err
	}

	if clientConfig.HeadlessValue() {
		resolvedValues = resolvedCoreValues{
			ProjectName:     clientConfig.ContainerNameValue(),
			Language:        clientConfig.LanguageValue(),
			WorkspaceFolder: clientConfig.WorkspaceFolderValue(),
			ServiceName:     clientConfig.ServiceNameValue(),
			Port:            clientConfig.PortValue(),
			Timezone:        clientConfig.TimezoneValue(),
		}
	} else {
		var commonTimezone *string
		if jsonConfig.Common != nil {
			commonTimezone = jsonConfig.Common.Timezone
		}
		resolvedValues, err = acc.resolvePromptValues(&clientConfig, defaultSetting, commonTimezone)
		if err != nil {
			return nil, err
		}
	}

	return acc.buildCodespaceConfig(clientConfig, defaultSetting, resolvedValues, resolvedEntries, jsonConfig)
}
