package app

import (
	"fmt"
	"os"
	"path/filepath"

	"codespacegen/internal/generator"
	"codespacegen/internal/generator/filewriter"
	"codespacegen/internal/i18n"
	"codespacegen/internal/input"
	"codespacegen/internal/resolve"
	"codespacegen/internal/workflow"
)

type InputConfig struct {
	clientInput   *input.CliInput
	jsonInput     *input.JsonInput
	defaultConfig *input.DefaultConfig
}

type Resolvers struct {
	CodespaceConfigResolver *resolve.CodespaceConfigResolver
}

type WorkflowCases struct {
	inputInputs           inputCollector
	resolveCodespace      configAssembler
	generateCodeArtifacts artifactExecutor
}

type App struct {
	flows WorkflowCases
}

func NewApp() *App {
	ic := InputConfig{
		clientInput:   input.NewCliInput(),
		jsonInput:     input.NewJsonInput(),
		defaultConfig: input.NewDefaultConfig(),
	}

	rs := Resolvers{
		CodespaceConfigResolver: resolve.NewCodespaceConfigResolver(os.Stdin),
	}

	generatorImpl := generator.NewDefaultTemplateGenerator()
	writer := filewriter.NewLocalFileWriter()

	flows := WorkflowCases{
		inputInputs:           workflow.NewCollectInputs(ic.clientInput, ic.jsonInput, ic.defaultConfig),
		resolveCodespace:      workflow.NewAssembleCodespaceConfig(rs.CodespaceConfigResolver),
		generateCodeArtifacts: workflow.NewGenerateCodespaceArtifacts(generatorImpl, writer),
	}

	return &App{flows: flows}
}

func (a *App) Run() error {
	inputs, err := a.flows.inputInputs.CollectConfig()
	if err != nil {
		return err
	}

	if inputs.CliConfig.ShowVersionValue() {
		fmt.Println(inputs.DefaultConfig.Version)
		return nil
	}

	if inputs.CliConfig.LangValue() != "" {
		i18n.SetLang(inputs.CliConfig.LangValue())
	}

	codespaceConfig, err := a.flows.resolveCodespace.Resolve(inputs.CliConfig, inputs.DefaultConfig, inputs.JsonConfig)
	if err != nil {
		return err
	}

	err = a.flows.generateCodeArtifacts.Execute(*codespaceConfig, inputs.CliConfig.EnableOverwriteFileValue(), inputs.CliConfig.OutputDirValue())
	if err != nil {
		return err
	}

	resolvedOutput, err := filepath.Abs(inputs.CliConfig.OutputDirValue())
	if err != nil {
		resolvedOutput = inputs.CliConfig.OutputDirValue()
	}

	fmt.Println(i18n.T("msg_generated_files", map[string]interface{}{"OutputDir": resolvedOutput}))

	return nil
}
