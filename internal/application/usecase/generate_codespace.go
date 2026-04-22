package usecase

import (
	"codespacegen/internal/application/port"
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/domain/service"
	"fmt"
	"path/filepath"
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

func (u *GenerateCodespaceArtifacts) Execute(config entity.CodespaceConfig, overwrite bool, outputDir string) error {
	if err := config.Validate(); err != nil {
		return err
	}

	files, err := u.generator.Generate(config)
	if err != nil {
		return fmt.Errorf("failed to generate templates: %w", err)
	}

	for _, file := range files {
		if err := u.writer.Write(filepath.Join(outputDir, file.RelativePath), file.Content, overwrite); err != nil {
			return fmt.Errorf("failed to write %s: %w", file.RelativePath, err)
		}
	}

	return nil
}
