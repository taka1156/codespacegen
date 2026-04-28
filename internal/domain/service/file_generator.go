package service

import "github.com/taka1156/codespacegen/internal/domain/entity"

type SettingGenerator interface {
	Generate(templateConfig entity.TemplateJson) ([]entity.GeneratedFile, error)
}

type CodespaceGenerator interface {
	Generate(config entity.CodespaceConfig) ([]entity.GeneratedFile, error)
}

type LocalFileWriter interface {
	Write(path string, content string, overwrite bool) error
}

type SettingTemplateGenerator interface {
	Generate(templateConfig entity.TemplateJson) (string, error)
}

type WorkdirProvider interface {
	GetConfigOutputPath() (string, error)
}
