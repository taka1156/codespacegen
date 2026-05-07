package workflow

import (
	"github.com/taka1156/codespacegen/internal/domain/service"
	"github.com/taka1156/codespacegen/internal/workflow/assemble"
	"github.com/taka1156/codespacegen/internal/workflow/collect"
	"github.com/taka1156/codespacegen/internal/workflow/generate"
	"github.com/taka1156/codespacegen/internal/workflow/initialize"
	"github.com/taka1156/codespacegen/internal/workflow/update"
)

type collectInputs = collect.CollectInputs

func NewCollectInputs(
	cliInput collect.ClientInputProvider,
	jsonInput collect.JsonConfigLoader,
	defaultConfig collect.DefaultSettingProvider,
) *collectInputs {
	return collect.NewCollectInputs(cliInput, jsonInput, defaultConfig)
}

type assembleCodespaceConfig = assemble.AssembleCodespaceConfig

func NewAssembleCodespaceConfig(
	codespacePrompter assemble.CodespacegenPrompter,
) *assembleCodespaceConfig {
	return assemble.NewAssembleCodespaceConfig(codespacePrompter)
}

type FileWriter = service.LocalFileWriter

type generateCodespaceArtifacts = generate.GenerateCodespaceArtifacts

func NewGenerateCodespaceArtifacts(
	codespaceGenerator service.CodespaceGenerator,
	writer service.LocalFileWriter,
) *generateCodespaceArtifacts {
	return generate.NewGenerateCodespaceArtifacts(codespaceGenerator, writer)
}

type initializeSettingJson = initialize.InitializeSettingJson

func NewInitializeSettingJson(
	settingTemplateGenerator service.SettingTemplateGenerator,
	workdirProvider service.WorkdirProvider,
	writer service.LocalFileWriter,
) *initializeSettingJson {
	return initialize.NewInitializeSettingJson(settingTemplateGenerator, workdirProvider, writer)
}

type updateCommandline = update.UpdateCommandline

func NewUpdateCommandline(
	updateCodespacegenCommandline update.CodespacegenUpdater,
) *updateCommandline {
	return update.NewUpdateCommandline(updateCodespacegenCommandline)
}
