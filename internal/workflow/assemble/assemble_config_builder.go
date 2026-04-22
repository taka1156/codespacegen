package assemble

import "codespacegen/internal/domain/entity"

func (rc *ResolveCodespaceConfig) buildCodespaceConfig(cliConfig *entity.CliConfig, coreValues resolvedCoreValues, resolvedEntry entity.JsonEntry, resolvedTimezone string) *entity.CodespaceConfig {
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
