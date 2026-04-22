package workflow

import (
	"codespacegen/internal/config"
	"codespacegen/internal/domain/service"
	"codespacegen/internal/resolve"
	"codespacegen/internal/workflow/assemble"
	"codespacegen/internal/workflow/generate"
	"codespacegen/internal/workflow/input"
)

type ResolveInput = input.ResolveInput

func NewResolveInput(
	cliInput config.CliInput,
	jsonInput config.JsonInput,
	defaultConfig config.DefaultConfig,
) *ResolveInput {
	return input.NewResolveInput(cliInput, jsonInput, defaultConfig)
}

type ResolveCodespaceConfig = assemble.ResolveCodespaceConfig

func NewResolveCodespaceConfig(
	codeSpaceConfigResolver resolve.CodeSpaceConfigResolver,
) *ResolveCodespaceConfig {
	return assemble.NewResolveCodespaceConfig(codeSpaceConfigResolver)
}

type FileWriter = generate.FileWriter

type GenerateCodespaceArtifacts = generate.GenerateCodespaceArtifacts

func NewGenerateCodespaceArtifacts(
	generator service.TemplateGenerator,
	writer FileWriter,
) *GenerateCodespaceArtifacts {
	return generate.NewGenerateCodespaceArtifacts(generator, writer)
}
