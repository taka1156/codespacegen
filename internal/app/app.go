package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/generator"
	"github.com/taka1156/codespacegen/internal/generator/filewriter"
	"github.com/taka1156/codespacegen/internal/generator/workdirprovider"
	"github.com/taka1156/codespacegen/internal/i18n"
	"github.com/taka1156/codespacegen/internal/infra"
	"github.com/taka1156/codespacegen/internal/input"
	"github.com/taka1156/codespacegen/internal/workflow"
)

type InputConfig struct {
	clientInput   *input.ClientInput
	jsonInput     *input.JsonInput
	defaultConfig *input.DefaultConfig
}

type Infra struct {
	CodespacegenUpdater     *infra.CodespacegenUpdater
	CodespacePromptResolver *infra.CodespacegenPrompter
}

type WorkflowCases struct {
	inputCollector             inputCollector
	assembleConfigResolver     assembleConfigResolver
	generateCodespaceArtifacts generateCodespaceArtifacts
	initializeSettingJson      initializeSettingJson
	updateCommandline          updateCommandline
}

type App struct {
	flows WorkflowCases
}

var Version = "dev"

func NewApp() *App {
	ic := InputConfig{
		clientInput:   input.NewClientInput(),
		jsonInput:     input.NewJsonInput(),
		defaultConfig: input.NewDefaultConfig(),
	}

	rs := Infra{
		CodespacegenUpdater:     infra.NewCodespacegenUpdater(),
		CodespacePromptResolver: infra.NewCodespacegenPrompter(os.Stdin),
	}

	codespaceGenerator := generator.NewCodespaceGenerator()
	settingTemplateGenerator := generator.NewSettingTemplateGenerator()
	workdir := workdirprovider.NewWorkdirProvider()
	writer := filewriter.NewLocalFileWriter()

	flows := WorkflowCases{
		inputCollector:             workflow.NewCollectInputs(ic.clientInput, ic.jsonInput, ic.defaultConfig),
		assembleConfigResolver:     workflow.NewAssembleCodespaceConfig(rs.CodespacePromptResolver),
		generateCodespaceArtifacts: workflow.NewGenerateCodespaceArtifacts(codespaceGenerator, writer),
		initializeSettingJson:      workflow.NewInitializeSettingJson(settingTemplateGenerator, workdir, writer),
		updateCommandline:          workflow.NewUpdateCommandline(rs.CodespacegenUpdater),
	}

	return &App{flows: flows}
}

func (a *App) Run() error {

	var args = os.Args

	inputs, err := a.flows.inputCollector.CollectConfig(args)
	if err != nil {
		return err
	}

	commandlineMode := inputs.ClientConfig.Mode.CommandlineModeValue()

	if inputs.ClientConfig.LangValue() != "" {
		i18n.SetLang(inputs.ClientConfig.LangValue())
	}

	switch commandlineMode {
	case "init":
		err = a.flows.initializeSettingJson.Execute(entity.DefaultTemplateJson, inputs.DefaultConfig.SettingJsonFileName)
		if err != nil {
			return err
		}
		return nil
	case "version":
		fmt.Println(Version)
		return nil
	case "update":
		err = a.flows.updateCommandline.Update(Version)
		if err != nil {
			return err
		}
		return nil
	}

	codespaceConfig, err := a.flows.assembleConfigResolver.Resolve(inputs.ClientConfig, inputs.DefaultConfig, inputs.JsonConfig)
	if err != nil {
		return err
	}

	err = a.flows.generateCodespaceArtifacts.Execute(*codespaceConfig, inputs.ClientConfig.EnableOverwriteFileValue(), inputs.ClientConfig.OutputDirValue())
	if err != nil {
		return err
	}

	resolvedOutput, err := filepath.Abs(inputs.ClientConfig.OutputDirValue())
	if err != nil {
		resolvedOutput = inputs.ClientConfig.OutputDirValue()
	}

	fmt.Println(i18n.T("msg_generated_files", map[string]interface{}{"OutputDir": resolvedOutput}))

	return nil
}
