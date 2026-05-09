package generate

import (
	"fmt"

	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/domain/service"
	"github.com/taka1156/codespacegen/internal/utils"
)

type GenerateCodespaceArtifacts struct {
	codespaceGenerator service.CodespaceGenerator
	writer             service.LocalFileWriter
}

func NewGenerateCodespaceArtifacts(
	codespaceGenerator service.CodespaceGenerator,
	writer service.LocalFileWriter,
) *GenerateCodespaceArtifacts {
	return &GenerateCodespaceArtifacts{
		codespaceGenerator: codespaceGenerator,
		writer:             writer,
	}
}

func (u *GenerateCodespaceArtifacts) Execute(
	config entity.CodespaceConfig,
	enableOverwriteFile bool,
	outputDir string,
) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	files, err := u.codespaceGenerator.Generate(config)
	if err != nil {
		return fmt.Errorf("failed to generate templates: %w", err)
	}

	for _, file := range files {
		outputPath, err := utils.ResolveOutputPath(outputDir, file.RelativePath)
		if err != nil {
			return fmt.Errorf("failed to resolve output path for %s: %w", file.RelativePath, err)
		}

		err = u.writer.Write(outputPath, file.Content, enableOverwriteFile)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", file.RelativePath, err)
		}
	}

	return nil
}
