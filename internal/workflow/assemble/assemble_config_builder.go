package assemble

import "codespacegen/internal/domain/entity"

func (acc *AssembleCodespaceConfig) buildCodespaceConfig(ClientConfig entity.ClientConfig, defaultSetting entity.DefaultSetting, coreValues resolvedCoreValues, resolvedEntry entity.JsonEntry, resolvedTimezone string) *entity.CodespaceConfig {
	return &entity.CodespaceConfig{
		Schema:           defaultSetting.VscSchema,
		ContainerName:    coreValues.ProjectName,
		ServiceName:      coreValues.ServiceName,
		WorkspaceFolder:  coreValues.WorkspaceFolder,
		BaseImage:        resolvedEntry.Image,
		Locale:           resolvedEntry.Locale,
		Timezone:         resolvedTimezone,
		ComposeFileName:  ClientConfig.ComposeFileValue(),
		PortMapping:      coreValues.Port,
		InstallCommand:   resolvedEntry.Install,
		VSCodeExtensions: resolvedEntry.VSCodeExtensions,
		OsModules:        defaultSetting.OsModules,
	}
}
