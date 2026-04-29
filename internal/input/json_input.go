package input

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/i18n"
	"github.com/tidwall/jsonc"
)

type JsonInput struct {
	fileLoader baseImageConfigLoader
}

type baseImageConfigLoader interface {
	Load(source string) ([]byte, error)
}

type fileConfigLoader struct{}

func NewJsonInput() *JsonInput {
	return &JsonInput{
		fileLoader: fileConfigLoader{},
	}
}

func (ji *JsonInput) LoadLanguageImages(source string) (*entity.JsonConfig, error) {

	rawJson, err := ji.fetchBaseImageConfig(source)
	if err != nil {
		return nil, err
	}
	if rawJson == nil {
		return nil, nil
	}

	var jsonConfig map[string]json.RawMessage
	if err := json.Unmarshal(rawJson, &jsonConfig); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_base_image_config"), err)
	}

	var jsonConfigEntity = entity.JsonConfig{}

	for key, value := range jsonConfig {
		switch key {
		case "$schema":
			if err := json.Unmarshal(value, &jsonConfigEntity.Schema); err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_base_image_config_schema"), err)
			}
		case "common":
			if err := json.Unmarshal(value, &jsonConfigEntity.Common); err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_base_image_config_common"), err)
			}
		default:
			var langEntry entity.LangEntry
			if err := json.Unmarshal(value, &langEntry); err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_base_image_config_language_entry", map[string]interface{}{"Language": key}), err)
			}
			if jsonConfigEntity.Langs == nil {
				jsonConfigEntity.Langs = make(map[string]*entity.LangEntry)
			}
			jsonConfigEntity.Langs[key] = &langEntry
		}
	}
	if jsonConfigEntity.Langs == nil {
		jsonConfigEntity.Langs = make(map[string]*entity.LangEntry)
	}
	return &jsonConfigEntity, nil
}

func (ji *JsonInput) fetchBaseImageConfig(source string) ([]byte, error) {
	raw, err := ji.fileLoader.Load(source)
	if err != nil {
		return nil, err
	}

	clean := jsonc.ToJSON(raw)

	var vscodeConfig entity.VSCodeConfig
	if strings.Contains(source, ".vscode/settings.json") {
		if err := json.Unmarshal(clean, &vscodeConfig); err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_vscode_settings"), err)
		}

		var m map[string]interface{}
		if err := json.Unmarshal(clean, &m); err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_vscode_settings"), err)
		}
		dev, ok := m["devcontainergen"].(map[string]interface{})
		if !ok || len(dev) == 0 {
			return nil, nil
		}
		jsonBytes, err := json.Marshal(dev)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_vscode_settings"), err)
		}
		return jsonBytes, nil
	}

	return clean, nil
}

func (l fileConfigLoader) Load(source string) ([]byte, error) {
	raw, err := os.ReadFile(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_base_image_config"), err)
	}

	return raw, nil
}
