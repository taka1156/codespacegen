package input

import (
	"codespacegen/internal/domain/entity"
)

type DefaultConfig struct {
}

func NewDefaultConfig() *DefaultConfig {
	return &DefaultConfig{}
}

func (dc *DefaultConfig) GetDefaultSetting() entity.DefaultSetting {
	return entity.DefaultSetting{
		Image:     entity.DefaultImage,
		Timezone:  entity.DefaultTimezone,
		VscSchema: entity.DefaultVscSchema,
		Locale:    entity.DefaultLocale,
		OsModules: entity.OsModules{
			AlpineModules:     entity.DefaultAlpineModules,
			DebianLikeModules: entity.DefaultDebianLikeModules,
		},
		SettingJsonFileName: entity.DefaultTemplateJsonPath,
	}
}
