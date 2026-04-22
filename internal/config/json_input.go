package config

import (
	"codespacegen/internal/i18n"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type JsonInput struct {
	httpsLoader baseImageConfigLoader
	fileLoader  baseImageConfigLoader
}

type baseImageConfigLoader interface {
	Load(source string) ([]byte, error)
}

type httpsConfigLoader struct {
	client *http.Client
}

type fileConfigLoader struct{}

func NewJsonInput() *JsonInput {
	return &JsonInput{
		httpsLoader: httpsConfigLoader{client: &http.Client{Timeout: 10 * time.Second}},
		fileLoader:  fileConfigLoader{},
	}
}

func (ji *JsonInput) LoadLanguageImages(source string) (map[string]json.RawMessage, error) {

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

	return jsonConfig, nil
}

func (ji *JsonInput) fetchBaseImageConfig(source string) ([]byte, error) {
	if strings.HasPrefix(source, "https://") {
		return ji.httpsLoader.Load(source)
	}

	if strings.HasPrefix(source, "http://") {
		return nil, errors.New(i18n.T("error_http_not_allowed_for_image_config"))
	}

	return ji.fileLoader.Load(source)
}

func (l httpsConfigLoader) Load(source string) ([]byte, error) {
	resp, err := l.client.Get(source) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_fetch_base_image_config_url"), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(i18n.T("error_base_image_config_url_status", map[string]interface{}{"StatusCode": resp.StatusCode}))
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_base_image_config_response"), err)
	}

	return raw, nil
}

func (fileConfigLoader) Load(source string) ([]byte, error) {
	raw, err := os.ReadFile(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_base_image_config"), err)
	}

	return raw, nil
}
