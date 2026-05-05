package collect

import (
	"errors"
	"testing"

	"github.com/taka1156/codespacegen/internal/domain/entity"
)

type fakeClientInput struct {
	config entity.ClientConfig
}

func (f *fakeClientInput) GetInput(args []string) entity.ClientConfig {
	return f.config
}

type fakeJsonConfigLoader struct {
	result *entity.JsonConfig
	err    error
}

func (f *fakeJsonConfigLoader) LoadLanguageImages(_ string) (*entity.JsonConfig, error) {
	return f.result, f.err
}

type fakeDefaultSettingProvider struct {
	setting entity.DefaultSetting
}

func (f *fakeDefaultSettingProvider) GetDefaultSetting() entity.DefaultSetting {
	return f.setting
}

func TestCollectInputs_CollectConfig_ReturnsCollectedInputs(t *testing.T) {
	imageConfig := "https://example.com/config.json"
	clientConfig := entity.ClientConfig{ImageConfig: &imageConfig}

	jsonResult := entity.JsonConfig{
		Langs: []*entity.LangEntry{
			{ProfileName: "python", Image: "python:3.12"},
		},
	}
	defaultSetting := entity.DefaultSetting{
		Timezone: "UTC",
		Image:    "alpine:latest",
		Version:  "1.0.0",
	}

	ci := NewCollectInputs(
		&fakeClientInput{config: clientConfig},
		&fakeJsonConfigLoader{result: &jsonResult},
		&fakeDefaultSettingProvider{setting: defaultSetting},
	)

	got, err := ci.CollectConfig([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ClientConfig.ImageConfigValue() != imageConfig {
		t.Errorf("ClientConfig.ImageConfig: got %q, want %q", got.ClientConfig.ImageConfigValue(), imageConfig)
	}
	if len(got.JsonConfig.Langs) != 1 {
		t.Errorf("JsonConfig length: got %d, want 1", len(got.JsonConfig.Langs))
	}
	if got.DefaultConfig.Timezone != "UTC" {
		t.Errorf("DefaultConfig.Timezone: got %q, want %q", got.DefaultConfig.Timezone, "UTC")
	}
	if got.DefaultConfig.Image != "alpine:latest" {
		t.Errorf("DefaultConfig.Image: got %q, want %q", got.DefaultConfig.Image, "alpine:latest")
	}
	if got.DefaultConfig.Version != "1.0.0" {
		t.Errorf("DefaultConfig.Version: got %q, want %q", got.DefaultConfig.Version, "1.0.0")
	}
}

func TestCollectInputs_CollectConfig_PassesImageConfigToLoader(t *testing.T) {
	imageConfig := "https://example.com/my-config.json"
	ClientConfig := entity.ClientConfig{ImageConfig: &imageConfig}

	var capturedSource string
	loader := &captureJsonConfigLoader{
		captureSource: func(s string) { capturedSource = s },
		result:        &entity.JsonConfig{},
	}

	ci := NewCollectInputs(
		&fakeClientInput{config: ClientConfig},
		loader,
		&fakeDefaultSettingProvider{},
	)

	if _, err := ci.CollectConfig([]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedSource != imageConfig {
		t.Errorf("LoadLanguageImages called with %q, want %q", capturedSource, imageConfig)
	}
}

func TestCollectInputs_CollectConfig_ReturnsErrorFromImageConfigLoader(t *testing.T) {
	ci := NewCollectInputs(
		&fakeClientInput{},
		&fakeJsonConfigLoader{err: errors.New("load failed")},
		&fakeDefaultSettingProvider{},
	)

	_, err := ci.CollectConfig([]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCollectInputs_CollectConfig_EmptyJsonConfigWhenNoImageConfigSet(t *testing.T) {
	ci := NewCollectInputs(
		&fakeClientInput{},
		&fakeJsonConfigLoader{result: &entity.JsonConfig{}},
		&fakeDefaultSettingProvider{setting: entity.DefaultSetting{Timezone: "UTC"}},
	)

	got, err := ci.CollectConfig([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.JsonConfig.Langs) != 0 {
		t.Errorf("expected empty JsonConfig, got %v", got.JsonConfig)
	}
}

func TestCollectInputs_CollectConfig_SkipsWhenJsonConfigLoaderReturnsNil(t *testing.T) {
	ci := NewCollectInputs(
		&fakeClientInput{},
		&fakeJsonConfigLoader{result: nil},
		&fakeDefaultSettingProvider{setting: entity.DefaultSetting{Timezone: "UTC"}},
	)

	got, err := ci.CollectConfig([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.JsonConfig.Langs != nil {
		t.Errorf("expected empty JsonConfig, got %v", got.JsonConfig)
	}
}

type captureJsonConfigLoader struct {
	captureSource func(string)
	result        *entity.JsonConfig
	err           error
}

func (c *captureJsonConfigLoader) LoadLanguageImages(source string) (*entity.JsonConfig, error) {
	c.captureSource(source)
	return c.result, c.err
}
