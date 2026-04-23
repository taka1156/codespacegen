package input

import (
	"codespacegen/internal/domain/entity"
	"slices"
)

const defaultImage = "alpine:latest"
const defaultTimezone = "UTC"
const defaultVersion = "dev"
const defaultVscSchema = "https://raw.githubusercontent.com/microsoft/vscode/main/extensions/configuration-editing/schemas/devContainer.vscode.schema.json"

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
		Image:     defaultImage,
		Timezone:  defaultTimezone,
		Version:   defaultVersion,
		VscSchema: defaultVscSchema,
		OsModules: entity.OsModules{
			AlpineModules:     defaultAlpineModules,
			DebianLikeModules: defaultDebianLikeModules,
		},
	}
}
