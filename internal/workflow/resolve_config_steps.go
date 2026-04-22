package workflow

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

type resolvedCoreValues struct {
	ProjectName     string
	Language        string
	WorkspaceFolder string
	ServiceName     string
	Port            string
}

func (rc *ResolveConfig) resolveCoreValues(cliConfig *entity.CliConfig) (resolvedCoreValues, error) {
	resolvedProjectName, err := rc.codeSpaceConfigResolver.ResolveProjectName(*cliConfig.ServiceName)
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

func (rc *ResolveConfig) resolveEntry(language string, cliConfig *entity.CliConfig, jsonEntries map[string]entity.JsonEntry, overrides map[string]json.RawMessage) (entity.JsonEntry, error) {
	mergedImages, err := rc.mergeLanguageResolver.MergeLanguageEntries(jsonEntries["common"], overrides)
	if err != nil {
		return entity.JsonEntry{}, err
	}

	return rc.codeSpaceConfigResolver.ResolveBaseImage(language, *cliConfig.BaseImage, *cliConfig.ImageConfig, mergedImages)
}

func (rc *ResolveConfig) buildCodespaceConfig(cliConfig *entity.CliConfig, coreValues resolvedCoreValues, resolvedEntry entity.JsonEntry, resolvedTimezone string) *entity.CodespaceConfig {
	return &entity.CodespaceConfig{
		ContainerName:    coreValues.ProjectName,
		ServiceName:      coreValues.ServiceName,
		WorkspaceFolder:  coreValues.WorkspaceFolder,
		BaseImage:        resolvedEntry.Image,
		Locale:           resolvedEntry.Locale,
		Timezone:         resolvedTimezone,
		ComposeFileName:  *cliConfig.ComposeFile,
		PortMapping:      coreValues.Port,
		InstallCommand:   resolvedEntry.Install,
		VSCodeExtensions: resolvedEntry.VSCodeExtensions,
	}
}
