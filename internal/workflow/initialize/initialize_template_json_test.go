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
	writer := &mockLocalFileWriter{
		writeFunc: func(path string, content string, overwrite bool) error {
			if path != "/tmp/setting.json" || content != "test-content" || overwrite != false {
				return errors.New("unexpected args")
			}
			return nil
		},
	}
	isj := NewInitializeSettingJson(gen, writer)
	err := isj.Execute(entity.JsonConfig{}, "setting.json", "/tmp")
	assert.NoError(t, err)
}

func TestInitializeSettingJson_Execute_GenerateError(t *testing.T) {
	gen := &mockSettingTemplateGenerator{
		generateFunc: func(tj entity.JsonConfig) (string, error) {
			return "", errors.New("generate error")
		},
	}
	writer := &mockLocalFileWriter{writeFunc: func(string, string, bool) error { return nil }}
	isj := NewInitializeSettingJson(gen, writer)
	err := isj.Execute(entity.JsonConfig{}, "setting.json", "/tmp")
	assert.ErrorContains(t, err, "failed to generate template JSON")
}

func TestInitializeSettingJson_Execute_ResolveOutputPathError(t *testing.T) {
	gen := &mockSettingTemplateGenerator{generateFunc: func(entity.JsonConfig) (string, error) { return "ok", nil }}
	writer := &mockLocalFileWriter{writeFunc: func(string, string, bool) error { return nil }}
	isj := NewInitializeSettingJson(gen, writer)
	// "../escaping.json" はパス外に出るためエラーになる
	err := isj.Execute(entity.JsonConfig{}, "../escaping.json", "/tmp")
	assert.ErrorContains(t, err, "failed to resolve output path")
}

func TestInitializeSettingJson_Execute_WriteError(t *testing.T) {
	gen := &mockSettingTemplateGenerator{generateFunc: func(entity.JsonConfig) (string, error) { return "ok", nil }}
	writer := &mockLocalFileWriter{writeFunc: func(string, string, bool) error { return errors.New("write error") }}
	isj := NewInitializeSettingJson(gen, writer)
	err := isj.Execute(entity.JsonConfig{}, "setting.json", "/tmp")
	assert.ErrorContains(t, err, "failed to write template JSON")
}
