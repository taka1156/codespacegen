package assemble

import "github.com/taka1156/codespacegen/internal/domain/entity"

type resolvedCoreValues struct {
	ProjectName     string
	Language        string
	WorkspaceFolder string
	ServiceName     string
	Port            string
	Timezone        string
}

func (acc *AssembleCodespaceConfig) resolvePromptValues(ClientConfig *entity.ClientConfig, defaultSetting entity.DefaultSetting) (resolvedCoreValues, error) {
	resolvedProjectName, err := acc.CodespacegenPrompter.PromptProjectName(ClientConfig.ContainerNameValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedLanguage, err := acc.CodespacegenPrompter.PromptLanguage(ClientConfig.LanguageValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedWorkspaceFolder, err := acc.CodespacegenPrompter.PromptWorkspaceFolder(ClientConfig.WorkspaceFolderValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedServiceName, err := acc.CodespacegenPrompter.PromptServiceName(ClientConfig.ServiceNameValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedPort, err := acc.CodespacegenPrompter.PromptPortMapping(ClientConfig.PortValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	timezoneDefault := ClientConfig.TimezoneValue()
	if timezoneDefault == "" {
		timezoneDefault = defaultSetting.Timezone
	}
	resolvedTimezone, err := acc.CodespacegenPrompter.PromptTimezone(timezoneDefault)
	if err != nil {
		return resolvedCoreValues{}, err
	}

	return resolvedCoreValues{
		ProjectName:     resolvedProjectName,
		Language:        resolvedLanguage,
		WorkspaceFolder: resolvedWorkspaceFolder,
		ServiceName:     resolvedServiceName,
		Port:            resolvedPort,
		Timezone:        resolvedTimezone,
	}, nil
}
