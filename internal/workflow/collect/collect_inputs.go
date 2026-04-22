package collect

import (
	"encoding/json"

	"codespacegen/internal/config"
	"codespacegen/internal/domain/entity"
)


type CollectInputs struct {
	cliInput      CLIInputProvider
	jsonInput     ImageConfigLoader
	defaultConfig DefaultSettingProvider
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

func (ri *CollectInputs) Collect() (*entity.CliConfig, map[string]json.RawMessage, config.DefaultSetting, error) {
	cliConfig := ri.cliInput.GetCliInput()
	jsonConfig, err := ri.jsonInput.LoadLanguageImages(*cliConfig.ImageConfig)
	if err != nil {
		return nil, nil, config.DefaultSetting{}, err
	}
	ds := ri.defaultConfig.GetDefaultSetting()

	return &cliConfig, jsonConfig, ds, nil
}
