package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codespacegen/internal/generator"
	"codespacegen/internal/generator/filewriter"
	"codespacegen/internal/i18n"
	"codespacegen/internal/input"
	"codespacegen/internal/workflow"

	"codespacegen/internal/resolve"
)

type InputConfig struct {
	clientInput   *input.CliInput
	jsonInput     *input.JsonInput
	defaultConfig *input.DefaultConfig
}

type Resolvers struct {
	codeSpaceConfigResolver *resolve.CodeSpaceConfigResolver
}

type WorkflowCases struct {
	inputInputs           *workflow.CollectInputs
	resolveCodespace      *workflow.ResolveCodespaceConfig
	generateCodeArtifacts *workflow.GenerateCodespaceArtifacts
}

type App struct {
	flows WorkflowCases
}

func main() {
	app := newApp()
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}

func newApp() *App {
	ic := InputConfig{
		clientInput:   input.NewCliInput(),
		jsonInput:     input.NewJsonInput(),
		defaultConfig: input.NewDefaultConfig(),
	}

	rs := Resolvers{
		codeSpaceConfigResolver: resolve.NewCodeSpaceConfigResolver(os.Stdin),
	}

	generatorImpl := generator.NewDefaultTemplateGenerator()
	writer := filewriter.NewLocalFileWriter()

	flows := WorkflowCases{
		inputInputs:           workflow.NewCollectInputs(ic.clientInput, ic.jsonInput, ic.defaultConfig),
		resolveCodespace:      workflow.NewResolveCodespaceConfig(rs.codeSpaceConfigResolver),
		generateCodeArtifacts: workflow.NewGenerateCodespaceArtifacts(generatorImpl, writer),
	}

	return &App{flows: flows}
}

func (a *App) Run() error {
	inputs, err := a.flows.inputInputs.CollectConfig()
	if err != nil {
		return err
	}

	cliConfig := inputs.CliConfig

	if cliConfig.ShowVersionValue() {
		fmt.Println(inputs.DefaultConfig.Version)
		return nil
	}

	if cliConfig.LangValue() != "" {
		i18n.SetLang(cliConfig.LangValue())
	}

	codespaceConfig, err := a.flows.resolveCodespace.Resolve(cliConfig, inputs.JsonConfig, inputs.DefaultConfig.Timezone, inputs.DefaultConfig.Image)
	if err != nil {
		return err
	}

	err = a.flows.generateCodeArtifacts.Execute(*codespaceConfig, cliConfig.EnableOverwriteFileValue(), cliConfig.OutputDirValue())
	if err != nil {
		return err
	}

	resolvedOutput, err := filepath.Abs(cliConfig.OutputDirValue())
	if err != nil {
		resolvedOutput = cliConfig.OutputDirValue()
	}

	fmt.Println(i18n.T("msg_generated_files", map[string]interface{}{"OutputDir": resolvedOutput}))

	return nil
}
