package input

import (
	"codespacegen/internal/domain/entity"
	"slices"
)

var commonModules = []string{
	"bash",
	"bash-completion",
	"ca-certificates",
	"tzdata",
	"git",
	"git-lfs",
	"vim",
	"curl",
}

var defaultAlpineModules = slices.Concat(commonModules, []string{
	"musl-locales",
	"musl-locales-lang",
})

var defaultDebianLikeModules = slices.Concat(commonModules, []string{
	"locales",
})

type DefaultConfig struct {
}

func NewDefaultConfig() *DefaultConfig {
	return &DefaultConfig{}
}

func (dc *DefaultConfig) GetDefaultSetting() entity.DefaultSetting {
	return entity.DefaultSetting{
		Image:     entity.DefaultImage,
		Timezone:  entity.DefaultTimezone,
		Version:   entity.DefaultVersion,
		VscSchema: entity.DefaultVscSchema,
		OsModules: entity.OsModules{
			AlpineModules:     defaultAlpineModules,
			DebianLikeModules: defaultDebianLikeModules,
		},
	}
}
