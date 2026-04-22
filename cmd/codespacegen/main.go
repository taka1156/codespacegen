package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codespacegen/internal/adapter/generator"
	"codespacegen/internal/adapter/persistence"
	"codespacegen/internal/application/usecase"
	"codespacegen/internal/config"
	"codespacegen/internal/i18n"

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

type UsecaseCases struct {
	resolveInput          *usecase.ResolveInput
	resolveConfig         *usecase.ResolveConfig
	generateCodeArtifacts *usecase.GenerateCodespaceArtifacts
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
	writer := persistence.NewLocalFileWriter()

	uc := UsecaseCases{
		resolveInput:          usecase.NewResolveInput(*ic.clientInput, *ic.jsonInput, *ic.defaultConfig),
		resolveConfig:         usecase.NewResolveConfig(*rs.mergeLanguage, *rs.codeSpaceConfigResolver),
		generateCodeArtifacts: usecase.NewGenerateCodespaceArtifacts(generatorImpl, writer),
	}

	cliConfig, jsonConfig, overrides, err := uc.resolveInput.Input()
	if err != nil {
		os.Exit(0)
	}

	config, err := uc.resolveConfig.Resolve(cliConfig, jsonConfig, overrides)
	if err != nil {
		os.Exit(0)
	}

	err = uc.generateCodeArtifacts.Execute(config, *cliConfig.Overwrite, *cliConfig.OutputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedOutput, err := filepath.Abs(*cliConfig.OutputDir)
	if err != nil {
		resolvedOutput = *cliConfig.OutputDir
	}

	fmt.Println(i18n.T("msg_generated_files", map[string]interface{}{"OutputDir": resolvedOutput}))
}
