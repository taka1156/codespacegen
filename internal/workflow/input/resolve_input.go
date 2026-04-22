package input

import (
	"encoding/json"

	"codespacegen/internal/config"
	"codespacegen/internal/domain/entity"
)

type ResolveInput struct {
	cliInput      config.CliInput
	jsonInput     config.JsonInput
	defaultConfig config.DefaultConfig
}

func NewResolveInput(
	cliInput config.CliInput,
	jsonInput config.JsonInput,
	defaultConfig config.DefaultConfig,
) *ResolveInput {
	return &ResolveInput{
		cliInput:      cliInput,
		jsonInput:     jsonInput,
		defaultConfig: defaultConfig,
	}
}

func (ri *ResolveInput) Input() (*entity.CliConfig, map[string]json.RawMessage, config.DefaultSetting, error) {
	cliConfig := ri.cliInput.GetCliInput()
	jsonConfig, err := ri.jsonInput.LoadLanguageImages(*cliConfig.ImageConfig)
	if err != nil {
		return nil, nil, config.DefaultSetting{}, err
	}
	ds := ri.defaultConfig.GetDefaultSetting()

	return &cliConfig, jsonConfig, ds, nil
}
