package collect

import (
	"encoding/json"
	"errors"
	"testing"

	"codespacegen/internal/domain/entity"
)

// --- フェイク実装 ---

type fakeCLIInput struct {
	config entity.CliConfig
}

func (f *fakeCLIInput) GetCliInput() entity.CliConfig {
	return f.config
}

type fakeImageConfigLoader struct {
	result map[string]json.RawMessage
	err    error
}

func (f *fakeImageConfigLoader) LoadLanguageImages(_ string) (map[string]json.RawMessage, error) {
	return f.result, f.err
}

type fakeDefaultSettingProvider struct {
	setting entity.DefaultSetting
}

func (f *fakeDefaultSettingProvider) GetDefaultSetting() entity.DefaultSetting {
	return f.setting
}

// --- テスト ---

func TestCollectInputs_CollectConfig_ReturnsCollectedInputs(t *testing.T) {
	imageConfig := "https://example.com/config.json"
	cliConfig := entity.CliConfig{ImageConfig: &imageConfig}

	jsonResult := map[string]json.RawMessage{
		"python": json.RawMessage(`"python:3.12"`),
	}
	defaultSetting := entity.DefaultSetting{
		Timezone: "UTC",
		Image:    "alpine:latest",
		Version:  "1.0.0",
	}

	ci := NewCollectInputs(
		&fakeCLIInput{config: cliConfig},
		&fakeImageConfigLoader{result: jsonResult},
		&fakeDefaultSettingProvider{setting: defaultSetting},
	)

	got, err := ci.CollectConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.CliConfig.ImageConfigValue() != imageConfig {
		t.Errorf("CliConfig.ImageConfig: got %q, want %q", got.CliConfig.ImageConfigValue(), imageConfig)
	}
	if len(got.JsonConfig) != 1 {
		t.Errorf("JsonConfig length: got %d, want 1", len(got.JsonConfig))
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
	cliConfig := entity.CliConfig{ImageConfig: &imageConfig}

	var capturedSource string
	loader := &captureImageConfigLoader{
		captureSource: func(s string) { capturedSource = s },
		result:        map[string]json.RawMessage{},
	}

	ci := NewCollectInputs(
		&fakeCLIInput{config: cliConfig},
		loader,
		&fakeDefaultSettingProvider{},
	)

	if _, err := ci.CollectConfig(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedSource != imageConfig {
		t.Errorf("LoadLanguageImages called with %q, want %q", capturedSource, imageConfig)
	}
}

func TestCollectInputs_CollectConfig_ReturnsErrorFromImageConfigLoader(t *testing.T) {
	ci := NewCollectInputs(
		&fakeCLIInput{},
		&fakeImageConfigLoader{err: errors.New("load failed")},
		&fakeDefaultSettingProvider{},
	)

	_, err := ci.CollectConfig()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCollectInputs_CollectConfig_EmptyJsonConfigWhenNoImageConfigSet(t *testing.T) {
	ci := NewCollectInputs(
		&fakeCLIInput{},
		&fakeImageConfigLoader{result: map[string]json.RawMessage{}},
		&fakeDefaultSettingProvider{setting: entity.DefaultSetting{Timezone: "UTC"}},
	)

	got, err := ci.CollectConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.JsonConfig) != 0 {
		t.Errorf("expected empty JsonConfig, got %v", got.JsonConfig)
	}
}

// captureImageConfigLoader は LoadLanguageImages に渡された source をキャプチャする。
type captureImageConfigLoader struct {
	captureSource func(string)
	result        map[string]json.RawMessage
	err           error
}

func (c *captureImageConfigLoader) LoadLanguageImages(source string) (map[string]json.RawMessage, error) {
	c.captureSource(source)
	return c.result, c.err
}
