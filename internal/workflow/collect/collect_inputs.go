package collect

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

type CollectInputs struct {
	cliInput      ClientInputProvider
	jsonInput     ImageConfigLoader
	defaultConfig DefaultSettingProvider
}

type CollectedInputs struct {
	ClientConfig     entity.ClientConfig
	JsonConfig    map[string]json.RawMessage
	DefaultConfig entity.DefaultSetting
}

func NewCollectInputs(
	cliInput ClientInputProvider,
	jsonInput ImageConfigLoader,
	defaultConfig DefaultSettingProvider,
) *CollectInputs {
	return &CollectInputs{
		cliInput:      cliInput,
		jsonInput:     jsonInput,
		defaultConfig: defaultConfig,
	}
}

func (ri *CollectInputs) CollectConfig() (*CollectedInputs, error) {
	ClientConfig := ri.cliInput.GetInput()
	jsonConfig, err := ri.jsonInput.LoadLanguageImages(ClientConfig.ImageConfigValue())
	if err != nil {
		return nil, err
	}
	ds := ri.defaultConfig.GetDefaultSetting()

	return &CollectedInputs{
		ClientConfig:     ClientConfig,
		JsonConfig:    jsonConfig,
		DefaultConfig: ds,
	}, nil
}
