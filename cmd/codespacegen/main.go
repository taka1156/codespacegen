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
	codeSpaceConfigResolver *resolve.CodeSpaceConfigResolver
}

type WorkflowCases struct {
	collectInputs         *workflow.CollectInputs
	resolveCodespace      *workflow.ResolveCodespaceConfig
	generateCodeArtifacts *workflow.GenerateCodespaceArtifacts
}

func main() {

	ic := InputConfig{
		clientInput:   config.NewCliInput(),
		jsonInput:     config.NewJsonInput(),
		defaultConfig: config.NewDefaultConfig(),
	}

	rs := Resolvers{
		codeSpaceConfigResolver: resolve.NewCodeSpaceConfigResolver(),
	}

	generatorImpl := generator.NewDefaultTemplateGenerator()
	writer := filewriter.NewLocalFileWriter()

	flows := WorkflowCases{
		collectInputs:         workflow.NewCollectInputs(ic.clientInput, ic.jsonInput, ic.defaultConfig),
		resolveCodespace:      workflow.NewResolveCodespaceConfig(*rs.codeSpaceConfigResolver),
		generateCodeArtifacts: workflow.NewGenerateCodespaceArtifacts(generatorImpl, writer),
	}

	cliConfig, jsonConfig, defaultConfig, err := flows.collectInputs.Collect()
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

	codespaceConfig, err := flows.resolveCodespace.Resolve(cliConfig, jsonConfig)
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
