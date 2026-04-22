package assemble

import "codespacegen/internal/domain/entity"

type resolvedCoreValues struct {
	ProjectName     string
	Language        string
	WorkspaceFolder string
	ServiceName     string
	Port            string
}

func (acc *AssembleCodespaceConfig) resolveCoreValues(cliConfig *entity.CliConfig) (resolvedCoreValues, error) {
	resolvedProjectName, err := acc.CodespaceConfigResolver.ResolveProjectName(cliConfig.ContainerNameValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedLanguage, err := acc.CodespaceConfigResolver.ResolveLanguage(cliConfig.LanguageValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedWorkspaceFolder, err := acc.CodespaceConfigResolver.ResolveWorkspaceFolder(cliConfig.WorkspaceFolderValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedServiceName, err := acc.CodespaceConfigResolver.ResolveServiceName(cliConfig.ServiceNameValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedPort, err := acc.CodespaceConfigResolver.ResolvePortMapping(cliConfig.PortValue())
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
