package config

import (
	"encoding/json"
	"fmt"

	"github.com/taka1156/codespacegen/internal/domain/entity"
)

type ConfigTemplateGenerator struct{}

func NewConfigTemplateGenerator() *ConfigTemplateGenerator {
	return &ConfigTemplateGenerator{}
}

func (g *ConfigTemplateGenerator) Generate(templateJson entity.TemplateJson) (string, error) {
	devcontainerBytes, err := json.MarshalIndent(templateJson, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to render devcontainer.json: %w", err)
	}
	return string(devcontainerBytes), nil
}
