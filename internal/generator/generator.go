package generator

import (
	"github.com/taka1156/codespacegen/internal/generator/codespace"
	"github.com/taka1156/codespacegen/internal/generator/setting"
)

func NewSettingTemplateGenerator() *setting.SettingTemplateGenerator {
	return setting.NewSettingTemplateGenerator()
}

func NewCodespaceGenerator() *codespace.CodespaceGenerator {
	return codespace.NewCodespaceGenerator()
}
