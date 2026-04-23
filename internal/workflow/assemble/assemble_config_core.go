package assemble

import "codespacegen/internal/domain/entity"

type resolvedCoreValues struct {
	ProjectName     string
	Language        string
	WorkspaceFolder string
	ServiceName     string
	Port            string
}

func (acc *AssembleCodespaceConfig) resolveCoreValues(ClientConfig *entity.ClientConfig) (resolvedCoreValues, error) {
	resolvedProjectName, err := acc.CodespaceConfigResolver.ResolveProjectName(ClientConfig.ContainerNameValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedLanguage, err := acc.CodespaceConfigResolver.ResolveLanguage(ClientConfig.LanguageValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedWorkspaceFolder, err := acc.CodespaceConfigResolver.ResolveWorkspaceFolder(ClientConfig.WorkspaceFolderValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedServiceName, err := acc.CodespaceConfigResolver.ResolveServiceName(ClientConfig.ServiceNameValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedPort, err := acc.CodespaceConfigResolver.ResolvePortMapping(ClientConfig.PortValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	return resolvedCoreValues{
		ProjectName:     resolvedProjectName,
		Language:        resolvedLanguage,
		WorkspaceFolder: resolvedWorkspaceFolder,
		ServiceName:     resolvedServiceName,
		Port:            resolvedPort,
	}, nil
}
