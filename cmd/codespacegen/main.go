package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"codespacegen/internal/application/usecase"
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/i18n"
	"codespacegen/internal/infrastructure/generator"
	"codespacegen/internal/infrastructure/persistence"
)

var (
	version            = "dev"
	portOnlyPattern    = regexp.MustCompile(`^\d+$`)
	portMappingPattern = regexp.MustCompile(`^\d+:\d+$`)
)

func main() {
	var (
		outputDir       = flag.String("output", ".devcontainer", "output directory for generated files")
		containerName   = flag.String("name", "", "project name (required, mapped to devcontainer name)")
		serviceName     = flag.String("service", "", "docker compose service name")
		language        = flag.String("language", "", "programming language (go/python/node/rust or image-config keys)")
		workspaceFolder = flag.String("workspace-folder", "/workspace", "workspace folder inside container")
		baseImage       = flag.String("base-image", "", "base Docker image (overrides -language default)")
		timezone        = flag.String("timezone", "", "timezone inside container (default: image-config timezone or UTC)")
		imageConfig     = flag.String("image-config", "codespacegen.json", "local path or https:// URL to base image config JSON")
		port            = flag.String("port", "", "port mapping (e.g. 3000 or 3000:3000)")
		composeFile     = flag.String("compose-file", "docker-compose.yaml", "docker compose file name")
		overwrite       = flag.Bool("force", false, "overwrite existing files")
		lang            = flag.String("lang", "", "language for CLI messages (en/ja, default: auto-detect)")
		showVersion     = flag.Bool("v", false, "print version and exit")
	)

	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	if *lang != "" {
		i18n.SetLang(*lang)
	}

	resolvedProjectName, err := resolveProjectName(*containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedLanguage, err := resolveLanguage(*language)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedWorkspaceFolder, err := resolveWorkspaceFolder(*workspaceFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedServiceName, err := resolveServiceName(*serviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedPort, err := resolvePortMapping(*port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedEntry, err := resolveBaseImage(resolvedLanguage, *baseImage, *imageConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedTimezone, err := resolveTimezone(*timezone, resolvedEntry.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	config := entity.CodespaceConfig{
		ContainerName:    resolvedProjectName,
		ServiceName:      resolvedServiceName,
		WorkspaceFolder:  resolvedWorkspaceFolder,
		BaseImage:        resolvedEntry.Image,
		Locale:           resolvedEntry.Locale,
		Timezone:         resolvedTimezone,
		ComposeFileName:  *composeFile,
		PortMapping:      resolvedPort,
		InstallCommand:   resolvedEntry.Install,
		VSCodeExtensions: resolvedEntry.VSCodeExtensions,
	}

	generatorImpl := generator.NewDefaultTemplateGenerator()
	writer := persistence.NewLocalFileWriter(*outputDir)
	uc := usecase.NewGenerateCodespaceArtifacts(generatorImpl, writer)

	if err := uc.Execute(config, *overwrite); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedOutput, err := filepath.Abs(*outputDir)
	if err != nil {
		resolvedOutput = *outputDir
	}

	fmt.Println(i18n.T("msg_generated_files", map[string]interface{}{"OutputDir": resolvedOutput}))
}

func resolveProjectName(explicitProjectName string) (string, error) {
	defaultProjectName := strings.TrimSpace(explicitProjectName)
	reader := bufio.NewReader(os.Stdin)

	for {
		if defaultProjectName == "" {
			fmt.Print(i18n.T("prompt_project_name_required"))
		} else {
			fmt.Print(i18n.T("prompt_project_name_with_default", map[string]interface{}{"Default": defaultProjectName}))
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				line = strings.TrimSpace(line)
				if line == "" {
					if defaultProjectName != "" {
						return defaultProjectName, nil
					}
					return "", fmt.Errorf("%s", i18n.T("error_project_name_required"))
				}
				return line, nil
			}
			return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_project_name"), err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			if defaultProjectName != "" {
				return defaultProjectName, nil
			}
			fmt.Println(i18n.T("msg_project_name_mandatory"))
			continue
		}

		return line, nil
	}
}

func promptWithDefault(prompt string, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			line = strings.TrimSpace(line)
			if line == "" {
				return defaultValue, nil
			}
			return line, nil
		}
		return "", err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return defaultValue, nil
	}

	return line, nil
}

func resolveLanguage(explicitLanguage string) (string, error) {
	defaultLanguage := strings.TrimSpace(explicitLanguage)
	value, err := promptWithDefault(i18n.T("prompt_language"), defaultLanguage)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_language"), err)
	}
	return strings.ToLower(strings.TrimSpace(value)), nil
}

func resolveWorkspaceFolder(explicitWorkspaceFolder string) (string, error) {
	defaultWorkspaceFolder := strings.TrimSpace(explicitWorkspaceFolder)
	if defaultWorkspaceFolder == "" {
		defaultWorkspaceFolder = "/workspace"
	}
	value, err := promptWithDefault(i18n.T("prompt_workspace_folder", map[string]interface{}{"Default": defaultWorkspaceFolder}), defaultWorkspaceFolder)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_workspace_folder"), err)
	}
	return strings.TrimSpace(value), nil
}

func resolveTimezone(explicitTimezone string, configTimezone string) (string, error) {
	defaultTimezone := strings.TrimSpace(explicitTimezone)
	if defaultTimezone == "" {
		defaultTimezone = strings.TrimSpace(configTimezone)
	}
	if defaultTimezone == "" {
		defaultTimezone = entity.DefaultTimezone
	}

	value, err := promptWithDefault(i18n.T("prompt_timezone", map[string]interface{}{"Default": defaultTimezone}), defaultTimezone)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_timezone"), err)
	}

	return strings.TrimSpace(value), nil
}

func resolveServiceName(explicitServiceName string) (string, error) {
	defaultServiceName := strings.TrimSpace(explicitServiceName)
	if defaultServiceName == "" {
		defaultServiceName = "app"
	}
	value, err := promptWithDefault(i18n.T("prompt_service_name", map[string]interface{}{"Default": defaultServiceName}), defaultServiceName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_service_name"), err)
	}
	return strings.TrimSpace(value), nil
}

func resolvePortMapping(explicitPort string) (string, error) {
	defaultPort := strings.TrimSpace(explicitPort)
	reader := bufio.NewReader(os.Stdin)
	for {
		if defaultPort == "" {
			fmt.Print(i18n.T("prompt_port_empty"))
		} else {
			fmt.Print(i18n.T("prompt_port_with_default", map[string]interface{}{"Default": defaultPort}))
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				line = strings.TrimSpace(line)
				if line == "" {
					if defaultPort != "" {
						return normalizePortMapping(defaultPort)
					}
					return "", nil
				}
				return normalizePortMapping(line)
			}
			return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_port_input"), err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			if defaultPort != "" {
				return normalizePortMapping(defaultPort)
			}
			return "", nil
		}

		normalized, err := normalizePortMapping(line)
		if err == nil {
			return normalized, nil
		}

		fmt.Println(i18n.T("error_invalid_port_format"))
	}
}

func normalizePortMapping(value string) (string, error) {
	v := strings.TrimSpace(value)
	if portOnlyPattern.MatchString(v) {
		return fmt.Sprintf("%s:%s", v, v), nil
	}
	if portMappingPattern.MatchString(v) {
		return v, nil
	}

	return "", errors.New(i18n.T("error_invalid_port_mapping", map[string]interface{}{"Value": value}))
}

type jsonEntry struct {
	Image            string
	Install          string
	Locale           entity.LocaleConfig
	Timezone         string
	VSCodeExtensions []string
}

func resolveBaseImage(language string, explicitBaseImage string, imageConfig string) (jsonEntry, error) {
	if explicitBaseImage != "" {
		return jsonEntry{Image: explicitBaseImage}, nil
	}

	if strings.TrimSpace(language) == "" {
		return jsonEntry{Image: entity.DefaultImage}, nil
	}

	entries, err := loadLanguageImages(imageConfig)
	if err != nil {
		return jsonEntry{}, err
	}

	key := strings.ToLower(strings.TrimSpace(language))
	entry, ok := entries[key]
	if !ok {
		return jsonEntry{}, errors.New(i18n.T("error_unsupported_language", map[string]interface{}{"Language": language}))
	}

	if entry.Image == "" {
		return jsonEntry{}, errors.New(i18n.T("error_image_required_for_language", map[string]interface{}{"Language": language}))
	}

	return entry, nil
}

func loadLanguageImages(source string) (map[string]jsonEntry, error) {
	images := make(map[string]jsonEntry)

	raw, err := fetchBaseImageConfig(source)
	if err != nil {
		return nil, err
	}

	var overrides map[string]json.RawMessage
	if err := json.Unmarshal(raw, &overrides); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_base_image_config"), err)
	}

	common, err := parseCommonEntry(overrides)
	if err != nil {
		return nil, err
	}

	for k, v := range overrides {
		normalizedKey := strings.ToLower(strings.TrimSpace(k))
		if normalizedKey == "" || normalizedKey == "common" || normalizedKey == "$schema" {
			continue
		}
		entry, err := parseLanguageEntry(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T("error_invalid_entry_for_key", map[string]interface{}{"Key": k}), err)
		}

		images[normalizedKey] = mergeLanguageEntries(common, entry)
	}

	return images, nil
}

func parseCommonEntry(overrides map[string]json.RawMessage) (jsonEntry, error) {
	for k, v := range overrides {
		if strings.ToLower(strings.TrimSpace(k)) != "common" {
			continue
		}

		entry, err := parseLanguageEntry(v)
		if err != nil {
			return jsonEntry{}, fmt.Errorf("%s: %w", i18n.T("error_invalid_entry_for_key", map[string]interface{}{"Key": k}), err)
		}
		return entry, nil
	}

	return jsonEntry{}, nil
}

func mergeLanguageEntries(base jsonEntry, override jsonEntry) jsonEntry {
	locale := override.Locale
	if locale.Lang == "" {
		locale = base.Locale
	}

	merged := jsonEntry{
		Image:    firstNonEmpty(override.Image, base.Image),
		Install:  firstNonEmpty(override.Install, base.Install),
		Timezone: firstNonEmpty(override.Timezone, base.Timezone),
		Locale:   locale,
	}

	merged.VSCodeExtensions = append(merged.VSCodeExtensions, base.VSCodeExtensions...)
	merged.VSCodeExtensions = append(merged.VSCodeExtensions, override.VSCodeExtensions...)

	return merged
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			return trimmed
		}
	}

	return ""
}

func parseLanguageEntry(raw json.RawMessage) (jsonEntry, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return jsonEntry{Image: strings.TrimSpace(s)}, nil
	}

	var setting struct {
		Image    string `json:"image"`
		Install  string `json:"install"`
		Timezone string `json:"timezone"`
		Locale   struct {
			Lang     string `json:"lang"`
			Language string `json:"language"`
			LcAll    string `json:"lcAll"`
		} `json:"locale"`
		VSCodeExtensions []string `json:"vscodeExtensions"`
	}

	if err := json.Unmarshal(raw, &setting); err != nil {
		return jsonEntry{}, fmt.Errorf("%s: %w", i18n.T("error_must_be_string_or_object"), err)
	}

	image := strings.TrimSpace(setting.Image)
	install := strings.TrimSpace(setting.Install)
	timezone := strings.TrimSpace(setting.Timezone)
	locale := entity.LocaleConfig{
		Lang:     strings.TrimSpace(setting.Locale.Lang),
		Language: strings.TrimSpace(setting.Locale.Language),
		LcAll:    strings.TrimSpace(setting.Locale.LcAll),
	}
	vscodeExtensions := make([]string, 0, len(setting.VSCodeExtensions))
	for _, ext := range setting.VSCodeExtensions {
		trimmed := strings.TrimSpace(ext)
		if trimmed != "" {
			vscodeExtensions = append(vscodeExtensions, trimmed)
		}
	}
	if image == "" && install != "" {
		return jsonEntry{}, errors.New(i18n.T("error_image_required_when_install"))
	}

	return jsonEntry{Image: image, Install: install, Locale: locale, Timezone: timezone, VSCodeExtensions: vscodeExtensions}, nil
}

func fetchBaseImageConfig(source string) ([]byte, error) {
	if strings.HasPrefix(source, "https://") {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(source) //nolint:noctx
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_fetch_base_image_config_url"), err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(i18n.T("error_base_image_config_url_status", map[string]interface{}{"StatusCode": resp.StatusCode}))
		}
		raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_base_image_config_response"), err)
		}
		return raw, nil
	}

	if strings.HasPrefix(source, "http://") {
		return nil, errors.New(i18n.T("error_http_not_allowed_for_image_config"))
	}

	raw, err := os.ReadFile(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_base_image_config"), err)
	}
	return raw, nil
}
