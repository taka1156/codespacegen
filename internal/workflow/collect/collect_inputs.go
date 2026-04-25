package collect

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

type CollectInputs struct {
	clientInput   ClientInputProvider
	jsonInput     ImageConfigLoader
	defaultConfig DefaultSettingProvider
}

type CollectedInputs struct {
	ClientConfig  entity.ClientConfig
	JsonConfig    map[string]json.RawMessage
	DefaultConfig entity.DefaultSetting
}

func NewCollectInputs(
	clientInput ClientInputProvider,
	jsonInput ImageConfigLoader,
	defaultConfig DefaultSettingProvider,
) *CollectInputs {
	return &CollectInputs{
		clientInput:   clientInput,
		jsonInput:     jsonInput,
		defaultConfig: defaultConfig,
	}
}

func (ri *CollectInputs) CollectConfig(args []string) (*CollectedInputs, error) {
	ClientConfig := ri.clientInput.GetInput(args)
	jsonConfig, err := ri.jsonInput.LoadLanguageImages(ClientConfig.ImageConfigValue())
	if err != nil {
		return nil, err
	}
	ds := ri.defaultConfig.GetDefaultSetting()

	return &CollectedInputs{
		ClientConfig:  ClientConfig,
		JsonConfig:    jsonConfig,
		DefaultConfig: ds,
	}, nil
}
