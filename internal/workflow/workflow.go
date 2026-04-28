package workflow

import (
	"codespacegen/internal/domain/service"
	"codespacegen/internal/workflow/assemble"
	"codespacegen/internal/workflow/collect"
	"codespacegen/internal/workflow/generate"
	"codespacegen/internal/workflow/initialize"
)

type CollectInputs = collect.CollectInputs

func NewCollectInputs(
	cliInput collect.ClientInputProvider,
	jsonInput collect.JsonConfigLoader,
	defaultConfig collect.DefaultSettingProvider,
) *CollectInputs {
	return collect.NewCollectInputs(cliInput, jsonInput, defaultConfig)
}

type AssembleCodespaceConfig = assemble.AssembleCodespaceConfig

func NewAssembleCodespaceConfig(
	CodespacePrompter assemble.CodespacegenPrompter,
) *AssembleCodespaceConfig {
	return assemble.NewAssembleCodespaceConfig(CodespacePrompter)
}

type FileWriter = service.LocalFileWriter

type GenerateCodespaceArtifacts = generate.GenerateCodespaceArtifacts

func NewGenerateCodespaceArtifacts(
	codespaceGenerator service.CodespaceGenerator,
	writer service.LocalFileWriter,
) *GenerateCodespaceArtifacts {
	return generate.NewGenerateCodespaceArtifacts(codespaceGenerator, writer)
}

type InitializeSettingJson = initialize.InitializeSettingJson

func NewInitializeSettingJson(settingTemplateGenerator service.SettingTemplateGenerator, workdirProvider service.WorkdirProvider, writer service.LocalFileWriter) *InitializeSettingJson {
	return initialize.NewInitializeSettingJson(settingTemplateGenerator, workdirProvider, writer)
}
