package workflow

import (
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/domain/service"
	"fmt"
	"path/filepath"
	"strings"
)

type GenerateCodespaceArtifacts struct {
	generator service.TemplateGenerator
	writer    FileWriter
}

func NewGenerateCodespaceArtifacts(
	generator service.TemplateGenerator,
	writer FileWriter,
) *GenerateCodespaceArtifacts {
	return &GenerateCodespaceArtifacts{
		generator: generator,
		writer:    writer,
	}
}

func (u *GenerateCodespaceArtifacts) Execute(config entity.CodespaceConfig, enableOverwriteFile bool, outputDir string) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	files, err := u.generator.Generate(config)
	if err != nil {
		return fmt.Errorf("failed to generate templates: %w", err)
	}

	for _, file := range files {
		outputPath, err := resolveOutputPath(outputDir, file.RelativePath)
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

func resolveOutputPath(outputDir string, relativePath string) (string, error) {
	cleanRelativePath := filepath.Clean(relativePath)
	if cleanRelativePath == "." || cleanRelativePath == "" {
		return "", fmt.Errorf("invalid file path: %s", relativePath)
	}

	if filepath.IsAbs(cleanRelativePath) {
		return "", fmt.Errorf("absolute path is not allowed: %s", relativePath)
	}

	joinedPath := filepath.Join(outputDir, cleanRelativePath)
	cleanOutputDir := filepath.Clean(outputDir)
	relativeToOutputDir, err := filepath.Rel(cleanOutputDir, joinedPath)
	if err != nil {
		return "", fmt.Errorf("failed to calculate relative path: %w", err)
	}

	if relativeToOutputDir == ".." || strings.HasPrefix(relativeToOutputDir, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes output directory: %s", relativePath)
	}

	return joinedPath, nil
}
