package usecase

import (
	"fmt"

	"codespacegen/internal/application/port"
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/domain/service"
)

type GenerateCodespaceArtifacts struct {
	generator service.TemplateGenerator
	writer    port.FileWriter
}

func NewGenerateCodespaceArtifacts(
	generator service.TemplateGenerator,
	writer port.FileWriter,
) *GenerateCodespaceArtifacts {
	return &GenerateCodespaceArtifacts{
		generator: generator,
		writer:    writer,
	}
}

func (u *GenerateCodespaceArtifacts) Execute(config entity.CodespaceConfig, overwrite bool) error {
	if err := config.Validate(); err != nil {
		return err
	}

	files, err := u.generator.Generate(config)
	if err != nil {
		return fmt.Errorf("failed to generate templates: %w", err)
	}

	for _, file := range files {
		if err := u.writer.Write(file.RelativePath, file.Content, overwrite); err != nil {
			return fmt.Errorf("failed to write %s: %w", file.RelativePath, err)
		}
	}

	return nil
}
