package setting

import (
	"encoding/json"
	"fmt"

	"github.com/taka1156/codespacegen/internal/domain/entity"
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
