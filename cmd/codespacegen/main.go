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
	)

	flag.Parse()

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

	resolvedBaseImage, resolvedInstall, resolvedConfigTimezone, resolvedExtensions, err := resolveBaseImage(resolvedLanguage, *baseImage, *imageConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resolvedTimezone, err := resolveTimezone(*timezone, resolvedConfigTimezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	config := entity.CodespaceConfig{
		ContainerName:    resolvedProjectName,
		ServiceName:      resolvedServiceName,
		WorkspaceFolder:  resolvedWorkspaceFolder,
		BaseImage:        resolvedBaseImage,
		Timezone:         resolvedTimezone,
		ComposeFileName:  *composeFile,
		PortMapping:      resolvedPort,
		InstallCommand:   resolvedInstall,
		VSCodeExtensions: resolvedExtensions,
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
			return "", fmt.Errorf("failed to read project name: %w", err)
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
		return "", fmt.Errorf("failed to read language: %w", err)
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
		return "", fmt.Errorf("failed to read workspace folder: %w", err)
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
		return "", fmt.Errorf("failed to read timezone: %w", err)
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
		return "", fmt.Errorf("failed to read service name: %w", err)
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
			return "", fmt.Errorf("failed to read port input: %w", err)
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

	return "", fmt.Errorf("invalid port mapping: %s", value)
}

type languageEntry struct {
	Image            string
	Install          string
	Timezone         string
	VSCodeExtensions []string
}

func resolveBaseImage(language string, explicitBaseImage string, imageConfig string) (string, string, string, []string, error) {
	if explicitBaseImage != "" {
		return explicitBaseImage, "", "", nil, nil
	}

	if strings.TrimSpace(language) == "" {
		return entity.DefaultImage, "", "", nil, nil
	}

	entries, err := loadLanguageImages(imageConfig)
	if err != nil {
		return "", "", "", nil, err
	}

	key := strings.ToLower(strings.TrimSpace(language))
	entry, ok := entries[key]
	if !ok {
		return "", "", "", nil, fmt.Errorf("unsupported language: %s", language)
	}

	if entry.Image == "" {
		return "", "", "", nil, fmt.Errorf("image is required for language %q: set \"image\" in the config", language)
	}

	return entry.Image, entry.Install, entry.Timezone, entry.VSCodeExtensions, nil
}

func loadLanguageImages(source string) (map[string]languageEntry, error) {
	images := make(map[string]languageEntry)

	raw, err := fetchBaseImageConfig(source)
	if err != nil {
		return nil, err
	}

	var overrides map[string]json.RawMessage
	if err := json.Unmarshal(raw, &overrides); err != nil {
		return nil, fmt.Errorf("failed to parse base image config: %w", err)
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
			return nil, fmt.Errorf("invalid entry for %q: %w", k, err)
		}

		images[normalizedKey] = mergeLanguageEntries(common, entry)
	}

	return images, nil
}

func parseCommonEntry(overrides map[string]json.RawMessage) (languageEntry, error) {
	for k, v := range overrides {
		if strings.ToLower(strings.TrimSpace(k)) != "common" {
			continue
		}

		entry, err := parseLanguageEntry(v)
		if err != nil {
			return languageEntry{}, fmt.Errorf("invalid entry for %q: %w", k, err)
		}
		return entry, nil
	}

	return languageEntry{}, nil
}

func mergeLanguageEntries(base languageEntry, override languageEntry) languageEntry {
	merged := languageEntry{
		Image:    firstNonEmpty(override.Image, base.Image),
		Install:  firstNonEmpty(override.Install, base.Install),
		Timezone: firstNonEmpty(override.Timezone, base.Timezone),
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

func parseLanguageEntry(raw json.RawMessage) (languageEntry, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return languageEntry{Image: strings.TrimSpace(s)}, nil
	}

	var obj struct {
		Image            string   `json:"image"`
		Install          string   `json:"install"`
		Timezone         string   `json:"timezone"`
		VSCodeExtensions []string `json:"vscodeExtensions"`
	}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return languageEntry{}, fmt.Errorf("must be a string or {\"image\",\"install\",\"timezone\",\"vscodeExtensions\"} object: %w", err)
	}

	image := strings.TrimSpace(obj.Image)
	install := strings.TrimSpace(obj.Install)
	timezone := strings.TrimSpace(obj.Timezone)
	vscodeExtensions := make([]string, 0, len(obj.VSCodeExtensions))
	for _, ext := range obj.VSCodeExtensions {
		trimmed := strings.TrimSpace(ext)
		if trimmed != "" {
			vscodeExtensions = append(vscodeExtensions, trimmed)
		}
	}
	if image == "" && install != "" {
		return languageEntry{}, fmt.Errorf("image is required when install command is provided\n(This is because the installation command is highly dependent on the container image.)")
	}

	return languageEntry{Image: image, Install: install, Timezone: timezone, VSCodeExtensions: vscodeExtensions}, nil
}

func fetchBaseImageConfig(source string) ([]byte, error) {
	if strings.HasPrefix(source, "https://") {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(source) //nolint:noctx
		if err != nil {
			return nil, fmt.Errorf("failed to fetch base image config from URL: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("base image config URL returned status %d", resp.StatusCode)
		}
		raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			return nil, fmt.Errorf("failed to read base image config response: %w", err)
		}
		return raw, nil
	}

	if strings.HasPrefix(source, "http://") {
		return nil, fmt.Errorf("http:// is not allowed for -image-config; use https://")
	}

	raw, err := os.ReadFile(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read base image config: %w", err)
	}
	return raw, nil
}
