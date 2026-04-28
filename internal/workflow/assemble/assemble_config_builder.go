package assemble

import (
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/i18n"
	"errors"
	"strings"
)

func (acc *AssembleCodespaceConfig) buildCodespaceConfig(clientConfig entity.ClientConfig, defaultSetting entity.DefaultSetting, promptValues resolvedCoreValues, langEntries map[string]entity.LangEntry) (*entity.CodespaceConfig, error) {
	imageEntry, err := resolveBaseImage(promptValues.Language, clientConfig.BaseImageValue(), langEntries, defaultSetting.Image)
	if err != nil {
		return nil, err
	}

	if imageEntry.Locale == nil {
		imageEntry.Locale = &defaultSetting.Locale
	}

	localeTimezone := resolveTimezone(promptValues.Timezone, clientConfig.TimezoneValue(), imageEntry.Timezone, defaultSetting.Timezone)

	osModules := mergeOsModules(defaultSetting.OsModules, imageEntry.LinuxPackages)

	runCommand := ""
	if imageEntry.RunCommand != nil {
		runCommand = *imageEntry.RunCommand
	}

	vsCodeExtensions := []string{}
	if imageEntry.VSCodeExtensions != nil {
		vsCodeExtensions = *imageEntry.VSCodeExtensions
	}

	portMapping := resolvePort(promptValues.Port, clientConfig.PortValue())

	return &entity.CodespaceConfig{
		Schema:           defaultSetting.VscSchema,
		ContainerName:    promptValues.ProjectName,
		ServiceName:      promptValues.ServiceName,
		WorkspaceFolder:  promptValues.WorkspaceFolder,
		BaseImage:        imageEntry.Image,
		Locale:           *imageEntry.Locale,
		Timezone:         localeTimezone,
		ComposeFileName:  clientConfig.ComposeFileValue(),
		PortMapping:      portMapping,
		RunCommand:       runCommand,
		VSCodeExtensions: vsCodeExtensions,
		OsModules:        osModules,
	}, nil
}

func resolveBaseImage(language string, explicitBaseImage string, jsonEntries map[string]entity.LangEntry, defaultImage string) (entity.LangEntry, error) {
	// priority: explicit(flag) > language(json with selection key) > default
	if explicitBaseImage != "" {
		return entity.LangEntry{Image: explicitBaseImage}, nil
	}

	if strings.TrimSpace(language) == "" {
		image := strings.TrimSpace(defaultImage)
		if image == "" {
			image = entity.DefaultImage
		}
		return entity.LangEntry{Image: image}, nil
	}

	key := strings.ToLower(strings.TrimSpace(language))
	entry, ok := jsonEntries[key]
	if !ok {
		return entity.LangEntry{}, errors.New(i18n.T("error_unsupported_language", map[string]interface{}{"Language": language}))
	}

	if entry.Image == "" {
		return entity.LangEntry{}, errors.New(i18n.T("error_image_required_for_language", map[string]interface{}{"Language": language}))
	}

	return entry, nil
}

func resolveTimezone(promptTimezone string, explicitTimezone string, configTimezone *string, defaultTimezone string) string {
	// priority: prompt > explicit(flag) > config > default
	resolved := strings.TrimSpace(promptTimezone)
	if resolved == "" {
		resolved = strings.TrimSpace(explicitTimezone)
	}
	if resolved == "" && configTimezone != nil {
		resolved = strings.TrimSpace(*configTimezone)
	}
	if resolved == "" {
		resolved = strings.TrimSpace(defaultTimezone)
	}
	if resolved == "" {
		resolved = entity.DefaultTimezone
	}

	return strings.TrimSpace(resolved)
}

func resolvePort(promptPort string, explicitPort string) string {
	// priority: prompt > explicit(flag) > nil
	if strings.TrimSpace(promptPort) != "" {
		return strings.TrimSpace(promptPort)
	}
	return strings.TrimSpace(explicitPort)
}

func mergeOsModules(base entity.OsModules, linuxPackages *[]entity.LinuxPackage) entity.OsModules {
	if linuxPackages == nil {
		return base
	}

	return entity.OsModules{
		AlpineModules:     append(base.AlpineModules, *linuxPackages...),
		DebianLikeModules: append(base.DebianLikeModules, *linuxPackages...),
	}
}
