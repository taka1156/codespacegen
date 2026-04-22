package workflow

import (
	"codespacegen/internal/domain/service"
	"codespacegen/internal/workflow/assemble"
	"codespacegen/internal/workflow/collect"
	"codespacegen/internal/workflow/generate"
)

type CollectInputs = collect.CollectInputs

func NewCollectInputs(
	cliInput collect.CLIInputProvider,
	jsonInput collect.ImageConfigLoader,
	defaultConfig collect.DefaultSettingProvider,
) *CollectInputs {
	return collect.NewCollectInputs(cliInput, jsonInput, defaultConfig)
}

type ResolveCodespaceConfig = assemble.ResolveCodespaceConfig

func NewResolveCodespaceConfig(
	codeSpaceConfigResolver assemble.ConfigResolver,
) *ResolveCodespaceConfig {
	return assemble.NewResolveCodespaceConfig(codeSpaceConfigResolver)
}

type FileWriter = service.FileWriter

type GenerateCodespaceArtifacts = generate.GenerateCodespaceArtifacts

func NewGenerateCodespaceArtifacts(
	generator service.TemplateGenerator,
	writer service.FileWriter,
) *GenerateCodespaceArtifacts {
	return generate.NewGenerateCodespaceArtifacts(generator, writer)
}
