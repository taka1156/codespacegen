package collect

import (
	"github.com/taka1156/codespacegen/internal/domain/entity"
)

type ClientInputProvider interface {
	GetInput(args []string) entity.ClientConfig
}

type JsonConfigLoader interface {
	LoadLanguageImages(source string) (*entity.JsonConfig, error)
}

type DefaultSettingProvider interface {
	GetDefaultSetting() entity.DefaultSetting
}
