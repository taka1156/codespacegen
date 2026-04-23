package collect

import (
	"codespacegen/internal/domain/entity"
	"encoding/json"
)

type ClientInputProvider interface {
	GetInput() entity.ClientConfig
}

type ImageConfigLoader interface {
	LoadLanguageImages(source string) (map[string]json.RawMessage, error)
}

type DefaultSettingProvider interface {
	GetDefaultSetting() entity.DefaultSetting
}
