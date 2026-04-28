package assemble

import (
	"codespacegen/internal/domain/entity"
)

func (acc *AssembleCodespaceConfig) buildCodespaceConfig(ClientConfig entity.ClientConfig, defaultSetting entity.DefaultSetting, coreValues resolvedCoreValues, resolvedEntry entity.LangEntry, resolvedTimezone string) *entity.CodespaceConfig {
	var locale entity.LocaleConfig
	if resolvedEntry.Locale != nil {
		locale = *resolvedEntry.Locale
	} else {
		locale = entity.DefaultLocale
	}
	return &entity.CodespaceConfig{
		Schema:           defaultSetting.VscSchema,
		ContainerName:    coreValues.ProjectName,
		ServiceName:      coreValues.ServiceName,
		WorkspaceFolder:  coreValues.WorkspaceFolder,
		BaseImage:        resolvedEntry.Image,
		Locale:           locale,
		Timezone:         resolvedTimezone,
		ComposeFileName:  ClientConfig.ComposeFileValue(),
		PortMapping:      coreValues.Port,
		RunCommand:       resolvedEntry.RunCommand,
		VSCodeExtensions: resolvedEntry.VSCodeExtensions,
		OsModules:        defaultSetting.OsModules,
	}
}
