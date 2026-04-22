package collect

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
)

type CollectInputs struct {
	cliInput      CLIInputProvider
	jsonInput     ImageConfigLoader
	defaultConfig DefaultSettingProvider
}

type CollectedInputs struct {
	CliConfig     entity.CliConfig
	JsonConfig    map[string]json.RawMessage
	DefaultConfig entity.DefaultSetting
}

func NewCollectInputs(
	cliInput CLIInputProvider,
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
	cliConfig := ri.cliInput.GetCliInput()
	jsonConfig, err := ri.jsonInput.LoadLanguageImages(cliConfig.ImageConfigValue())
	if err != nil {
		return nil, err
	}
	ds := ri.defaultConfig.GetDefaultSetting()

	return &CollectedInputs{
		CliConfig:     cliConfig,
		JsonConfig:    jsonConfig,
		DefaultConfig: ds,
	}, nil
}
