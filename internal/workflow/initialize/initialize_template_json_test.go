package initialize

import (
	"errors"
	"testing"

	"github.com/taka1156/codespacegen/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

type mockSettingTemplateGenerator struct {
	generateFunc func(entity.JsonConfig) (string, error)
}

func (m *mockSettingTemplateGenerator) Generate(tj entity.JsonConfig) (string, error) {
	return m.generateFunc(tj)
}

type mockWorkdirProvider struct {
	getConfigOutputPathFunc func() (string, error)
}

func (m *mockWorkdirProvider) GetConfigOutputPath() (string, error) {
	return m.getConfigOutputPathFunc()
}

type mockLocalFileWriter struct {
	writeFunc func(path string, content string, overwrite bool) error
}

func (m *mockLocalFileWriter) Write(path string, content string, overwrite bool) error {
	return m.writeFunc(path, content, overwrite)
}

func TestInitializeSettingJson_Execute_Success(t *testing.T) {
	gen := &mockSettingTemplateGenerator{
		generateFunc: func(tj entity.JsonConfig) (string, error) {
			return "test-content", nil
		},
	}
	wd := &mockWorkdirProvider{
		getConfigOutputPathFunc: func() (string, error) {
			return "/tmp/", nil
		},
	}
	writer := &mockLocalFileWriter{
		writeFunc: func(path string, content string, overwrite bool) error {
			if path != "/tmp/setting.json" || content != "test-content" || overwrite != false {
				return errors.New("unexpected args")
			}
			return nil
		},
	}
	isj := NewInitializeSettingJson(gen, wd, writer)
	err := isj.Execute(entity.JsonConfig{}, "setting.json")
	assert.NoError(t, err)
}

func TestInitializeSettingJson_Execute_GenerateError(t *testing.T) {
	gen := &mockSettingTemplateGenerator{
		generateFunc: func(tj entity.JsonConfig) (string, error) {
			return "", errors.New("generate error")
		},
	}
	wd := &mockWorkdirProvider{getConfigOutputPathFunc: func() (string, error) { return "/", nil }}
	writer := &mockLocalFileWriter{writeFunc: func(string, string, bool) error { return nil }}
	isj := NewInitializeSettingJson(gen, wd, writer)
	err := isj.Execute(entity.JsonConfig{}, "setting.json")
	assert.ErrorContains(t, err, "failed to generate template JSON")
}

func TestInitializeSettingJson_Execute_GetConfigOutputPathError(t *testing.T) {
	gen := &mockSettingTemplateGenerator{generateFunc: func(entity.JsonConfig) (string, error) { return "ok", nil }}
	wd := &mockWorkdirProvider{getConfigOutputPathFunc: func() (string, error) { return "", errors.New("path error") }}
	writer := &mockLocalFileWriter{writeFunc: func(string, string, bool) error { return nil }}
	isj := NewInitializeSettingJson(gen, wd, writer)
	err := isj.Execute(entity.JsonConfig{}, "setting.json")
	assert.ErrorContains(t, err, "failed to get config output path")
}

func TestInitializeSettingJson_Execute_WriteError(t *testing.T) {
	gen := &mockSettingTemplateGenerator{generateFunc: func(entity.JsonConfig) (string, error) { return "ok", nil }}
	wd := &mockWorkdirProvider{getConfigOutputPathFunc: func() (string, error) { return "/tmp/", nil }}
	writer := &mockLocalFileWriter{writeFunc: func(string, string, bool) error { return errors.New("write error") }}
	isj := NewInitializeSettingJson(gen, wd, writer)
	err := isj.Execute(entity.JsonConfig{}, "setting.json")
	assert.ErrorContains(t, err, "failed to write template JSON")
}
