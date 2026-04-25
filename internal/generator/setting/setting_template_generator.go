package setting

import (
	"codespacegen/internal/domain/entity"
	"encoding/json"
	"fmt"
)

type SettingTemplateGenerator struct{}

func NewSettingTemplateGenerator() *SettingTemplateGenerator {
	return &SettingTemplateGenerator{}
}

func (g *SettingTemplateGenerator) Generate(templateJson entity.TemplateJson) (string, error) {
	devcontainerBytes, err := json.MarshalIndent(templateJson, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to render devcontainer.json: %w", err)
	}
	return string(devcontainerBytes), nil
}
