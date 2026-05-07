package entity

import "slices"

const DefaultRepositoryName = "taka1156/codespacegen"

const DefaultImage = "alpine:latest"

const DefaultTimezone = "UTC"

//nolint:lll // It's a URL, so breaking it into a new line will break it
const DefaultVscSchema = "https://raw.githubusercontent.com/microsoft/vscode/main/extensions/configuration-editing/schemas/devContainer.vscode.schema.json"

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

var DefaultAlpineModules = slices.Concat(commonModules, []string{
	"musl-locales",
	"musl-locales-lang",
})

var DefaultDebianLikeModules = slices.Concat(commonModules, []string{
	"locales",
})

var DefaultLocale = LocaleConfig{
	Lang:     "ja_JP.UTF-8",
	Language: "ja_JP:ja",
	LcAll:    "ja_JP.UTF-8",
}
