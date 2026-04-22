package config

import (
	"codespacegen/internal/domain/entity"
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
}

func NewJsonInput() *JsonInput {
	return &JsonInput{}
}

func (ji *JsonInput) LoadLanguageImages(source string) (map[string]entity.JsonEntry, map[string]json.RawMessage, error) {
	images := make(map[string]entity.JsonEntry)

	rawJson, err := fetchBaseImageConfig(source)
	if err != nil {
		return images, nil, err
	}

	var overrides map[string]json.RawMessage
	if err := json.Unmarshal(rawJson, &overrides); err != nil {
		return images, overrides, fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_base_image_config"), err)
	}

	return images, overrides, nil
}

func fetchBaseImageConfig(source string) ([]byte, error) {
	if strings.HasPrefix(source, "https://") {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(source) //nolint:noctx
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

	if strings.HasPrefix(source, "http://") {
		return nil, errors.New(i18n.T("error_http_not_allowed_for_image_config"))
	}

	raw, err := os.ReadFile(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_base_image_config"), err)
	}
	return raw, nil
}
