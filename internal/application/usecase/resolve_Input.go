package usecase

import (
	"encoding/json"
	"fmt"

	"codespacegen/internal/config"
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/i18n"
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

func (ri *ResolveInput) Input() (*entity.CliConfig, map[string]entity.JsonEntry, map[string]json.RawMessage, error) {
	cliConfig := ri.cliInput.GetCliInput()
	jsonConfig, overrides, err := ri.jsonInput.LoadLanguageImages(*cliConfig.ImageConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	ds := ri.defaultConfig.GetDefaultSetting()

	if *cliConfig.ShowVersion {
		fmt.Println(ds.Version)
		return nil, nil, nil, nil
	}

	if *cliConfig.Lang != "" {
		i18n.SetLang(*cliConfig.Lang)
	}

	return &cliConfig, jsonConfig, overrides, nil
}
