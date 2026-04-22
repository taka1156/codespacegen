package collect

import (
	"codespacegen/internal/domain/entity"
	"encoding/json"
)

type CLIInputProvider interface {
	GetCliInput() entity.CliConfig
}

type ImageConfigLoader interface {
	LoadLanguageImages(source string) (map[string]json.RawMessage, error)
}

type DefaultSettingProvider interface {
	GetDefaultSetting() entity.DefaultSetting
}
