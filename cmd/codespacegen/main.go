package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codespacegen/internal/config"
	"codespacegen/internal/generator"
	"codespacegen/internal/generator/filewriter"
	"codespacegen/internal/i18n"
	"codespacegen/internal/workflow"

	"codespacegen/internal/resolve"
)

type InputConfig struct {
	clientInput   *config.CliInput
	jsonInput     *config.JsonInput
	defaultConfig *config.DefaultConfig
}

type Resolvers struct {
	mergeLanguage           *resolve.MergeLanguageResolver
	codeSpaceConfigResolver *resolve.CodeSpaceConfigResolver
}

type WorkflowCases struct {
	resolveInput          *workflow.ResolveInput
	resolveConfig         *workflow.ResolveConfig
	generateCodeArtifacts *workflow.GenerateCodespaceArtifacts
}

func main() {

	ic := InputConfig{
		clientInput:   config.NewCliInput(),
		jsonInput:     config.NewJsonInput(),
		defaultConfig: config.NewDefaultConfig(),
	}

	rs := Resolvers{
		mergeLanguage:           resolve.NewMergeLanguageResolver(),
		codeSpaceConfigResolver: resolve.NewCodeSpaceConfigResolver(),
	}

	generatorImpl := generator.NewDefaultTemplateGenerator()
	writer := filewriter.NewLocalFileWriter()

	flows := WorkflowCases{
		resolveInput:          workflow.NewResolveInput(*ic.clientInput, *ic.jsonInput, *ic.defaultConfig),
		resolveConfig:         workflow.NewResolveConfig(*rs.mergeLanguage, *rs.codeSpaceConfigResolver),
		generateCodeArtifacts: workflow.NewGenerateCodespaceArtifacts(generatorImpl, writer),
	}

	cliConfig, jsonConfig, overrides, defaultConfig, err := flows.resolveInput.Input()
	if err != nil {
		os.Exit(1)
	}

	if *cliConfig.ShowVersion {
		fmt.Println(defaultConfig.Version)
		os.Exit(0)
	}

	if *cliConfig.Lang != "" {
		i18n.SetLang(*cliConfig.Lang)
	}

	codespaceConfig, err := flows.resolveConfig.Resolve(cliConfig, jsonConfig, overrides)
	if err != nil {
		os.Exit(1)
	}

	err = flows.generateCodeArtifacts.Execute(*codespaceConfig, *cliConfig.EnableOverwriteFile, *cliConfig.OutputDir)
	if err != nil {
		os.Exit(1)
	}

	resolvedOutput, err := filepath.Abs(*cliConfig.OutputDir)
	if err != nil {
		resolvedOutput = *cliConfig.OutputDir
	}

	fmt.Println(i18n.T("msg_generated_files", map[string]interface{}{"OutputDir": resolvedOutput}))
}
