package service

import "codespacegen/internal/domain/entity"

type TemplateGenerator interface {
	Generate(config entity.CodespaceConfig) ([]entity.GeneratedFile, error)
}
