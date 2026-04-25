package generator

import (
	"codespacegen/internal/generator/codespace"
	"codespacegen/internal/generator/setting"
)

func NewSettingTemplateGenerator() *setting.SettingTemplateGenerator {
	return setting.NewSettingTemplateGenerator()
}

func NewCodespaceGenerator() *codespace.CodespaceGenerator {
	return codespace.NewCodespaceGenerator()
}
