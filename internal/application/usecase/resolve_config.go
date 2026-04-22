package usecase

import (
	"encoding/json"
	"fmt"
	"os"

	"codespacegen/internal/domain/entity"
	"codespacegen/internal/resolve"
)

type ResolveConfig struct {
	mergeLanguageResolver   resolve.MergeLanguageResolver
	codeSpaceConfigResolver resolve.CodeSpaceConfigResolver
}

func NewResolveConfig(
	mergeLanguageResolver resolve.MergeLanguageResolver,
	codeSpaceConfigResolver resolve.CodeSpaceConfigResolver,
) *ResolveConfig {
	return &ResolveConfig{
		mergeLanguageResolver:   mergeLanguageResolver,
		codeSpaceConfigResolver: codeSpaceConfigResolver,
	}
}

func (rc *ResolveConfig) Resolve(cliConfig *entity.CliConfig, jsonEntries map[string]entity.JsonEntry, overrides map[string]json.RawMessage) (entity.CodespaceConfig, error) {

	resolvedProjectName, err := rc.codeSpaceConfigResolver.ResolveProjectName(*cliConfig.ServiceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedLanguage, err := rc.codeSpaceConfigResolver.ResolveLanguage(*cliConfig.Language)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedWorkspaceFolder, err := rc.codeSpaceConfigResolver.ResolveWorkspaceFolder(*cliConfig.WorkspaceFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedServiceName, err := rc.codeSpaceConfigResolver.ResolveServiceName(*cliConfig.ServiceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedPort, err := rc.codeSpaceConfigResolver.ResolvePortMapping(*cliConfig.Port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	mergedImages, err := rc.mergeLanguageResolver.MergeLanguageEntries(jsonEntries["common"], overrides)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedEntry, err := rc.codeSpaceConfigResolver.ResolveBaseImage(resolvedLanguage, *cliConfig.BaseImage, *cliConfig.ImageConfig, mergedImages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedTimezone, err := rc.codeSpaceConfigResolver.ResolveTimezone(*cliConfig.Timezone, resolvedEntry.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	return entity.CodespaceConfig{
		ContainerName:    resolvedProjectName,
		ServiceName:      resolvedServiceName,
		WorkspaceFolder:  resolvedWorkspaceFolder,
		BaseImage:        resolvedEntry.Image,
		Locale:           resolvedEntry.Locale,
		Timezone:         resolvedTimezone,
		ComposeFileName:  *cliConfig.ComposeFile,
		PortMapping:      resolvedPort,
		InstallCommand:   resolvedEntry.Install,
		VSCodeExtensions: resolvedEntry.VSCodeExtensions,
	}, nil
}
