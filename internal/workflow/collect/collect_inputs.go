package collect

import (
	"github.com/taka1156/codespacegen/internal/domain/entity"
)

type CollectInputs struct {
	clientInput   ClientInputProvider
	jsonInput     JsonConfigLoader
	defaultConfig DefaultSettingProvider
}

type CollectedInputs struct {
	ClientConfig  entity.ClientConfig
	JsonConfig    entity.JsonConfig
	DefaultConfig entity.DefaultSetting
}

func NewCollectInputs(
	clientInput ClientInputProvider,
	jsonInput JsonConfigLoader,
	defaultConfig DefaultSettingProvider,
) *CollectInputs {
	return &CollectInputs{
		clientInput:   clientInput,
		jsonInput:     jsonInput,
		defaultConfig: defaultConfig,
	}
}

func (ri *CollectInputs) CollectConfig(args []string) (*CollectedInputs, error) {
	clientConfig := ri.clientInput.GetInput(args)
	jsonConfig, err := ri.jsonInput.LoadLanguageImages(clientConfig.ImageConfigValue())
	if err != nil {
		return nil, err
	}
	if jsonConfig == nil {
		jsonConfig = &entity.JsonConfig{}
	}
	ds := ri.defaultConfig.GetDefaultSetting()

	return &CollectedInputs{
		ClientConfig:  clientConfig,
		JsonConfig:    *jsonConfig,
		DefaultConfig: ds,
	}, nil
}
