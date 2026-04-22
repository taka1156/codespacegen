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
	resolvedProjectName, err := acc.codeSpaceConfigResolver.ResolveProjectName(cliConfig.ContainerNameValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedLanguage, err := acc.codeSpaceConfigResolver.ResolveLanguage(cliConfig.LanguageValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedWorkspaceFolder, err := acc.codeSpaceConfigResolver.ResolveWorkspaceFolder(cliConfig.WorkspaceFolderValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedServiceName, err := acc.codeSpaceConfigResolver.ResolveServiceName(cliConfig.ServiceNameValue())
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedPort, err := acc.codeSpaceConfigResolver.ResolvePortMapping(cliConfig.PortValue())
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
