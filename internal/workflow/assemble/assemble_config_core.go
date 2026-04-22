package assemble

import "codespacegen/internal/domain/entity"

type resolvedCoreValues struct {
	ProjectName     string
	Language        string
	WorkspaceFolder string
	ServiceName     string
	Port            string
}

func (rc *ResolveCodespaceConfig) resolveCoreValues(cliConfig *entity.CliConfig) (resolvedCoreValues, error) {
	resolvedProjectName, err := rc.codeSpaceConfigResolver.ResolveProjectName(*cliConfig.ContainerName)
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedLanguage, err := rc.codeSpaceConfigResolver.ResolveLanguage(*cliConfig.Language)
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedWorkspaceFolder, err := rc.codeSpaceConfigResolver.ResolveWorkspaceFolder(*cliConfig.WorkspaceFolder)
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedServiceName, err := rc.codeSpaceConfigResolver.ResolveServiceName(*cliConfig.ServiceName)
	if err != nil {
		return resolvedCoreValues{}, err
	}

	resolvedPort, err := rc.codeSpaceConfigResolver.ResolvePortMapping(*cliConfig.Port)
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
