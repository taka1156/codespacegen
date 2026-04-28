package workflow

import (
	"github.com/taka1156/codespacegen/internal/domain/service"
	"github.com/taka1156/codespacegen/internal/workflow/assemble"
	"github.com/taka1156/codespacegen/internal/workflow/collect"
	"github.com/taka1156/codespacegen/internal/workflow/generate"
	"github.com/taka1156/codespacegen/internal/workflow/initialize"
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
