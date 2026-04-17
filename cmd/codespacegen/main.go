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
	)

	flag.Parse()

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

	resolvedBaseImage, resolvedInstall, resolvedConfigTimezone, err := resolveBaseImage(resolvedLanguage, *baseImage, *imageConfig)
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
		ContainerName:   resolvedProjectName,
		ServiceName:     resolvedServiceName,
		WorkspaceFolder: resolvedWorkspaceFolder,
		BaseImage:       resolvedBaseImage,
		Timezone:        resolvedTimezone,
		ComposeFileName: *composeFile,
		PortMapping:     resolvedPort,
		InstallCommand:  resolvedInstall,
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

	fmt.Printf("Generated Codespace files in %s\n", resolvedOutput)
}

func resolveProjectName(explicitProjectName string) (string, error) {
	defaultProjectName := strings.TrimSpace(explicitProjectName)
	reader := bufio.NewReader(os.Stdin)

	for {
		if defaultProjectName == "" {
			fmt.Print("プロジェクト名を入力してください（必須）: ")
		} else {
			fmt.Printf("プロジェクト名を入力してください（必須、未入力で %s）: ", defaultProjectName)
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				line = strings.TrimSpace(line)
				if line == "" {
					if defaultProjectName != "" {
						return defaultProjectName, nil
					}
					return "", fmt.Errorf("project name is required")
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
			fmt.Println("プロジェクト名は必須です。")
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
	value, err := promptWithDefault("言語を入力してください（未入力で alpine 固定）: ", defaultLanguage)
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
	value, err := promptWithDefault(fmt.Sprintf("ワークスペースを入力してください（未入力で %s）: ", defaultWorkspaceFolder), defaultWorkspaceFolder)
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

	value, err := promptWithDefault(fmt.Sprintf("タイムゾーンを入力してください（未入力で %s）: ", defaultTimezone), defaultTimezone)
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
	value, err := promptWithDefault(fmt.Sprintf("サービス名を入力してください（未入力で %s）: ", defaultServiceName), defaultServiceName)
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
			fmt.Print("公開ポートを入力してください (例: 3000 または 3000:3000、不要ならEnter): ")
		} else {
			fmt.Printf("公開ポートを入力してください (例: 3000 または 3000:3000、未入力で %s): ", defaultPort)
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

		fmt.Println("無効なポート形式です。3000 または 3000:3000 の形式で入力してください。")
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
	Image    string
	Install  string
	Timezone string
}

func resolveBaseImage(language string, explicitBaseImage string, imageConfig string) (string, string, string, error) {
	if explicitBaseImage != "" {
		return explicitBaseImage, "", "", nil
	}

	if strings.TrimSpace(language) == "" {
		return "alpine:latest", "", "", nil
	}

	entries, err := loadLanguageBaseImages(imageConfig)
	if err != nil {
		return "", "", "", err
	}

	key := strings.ToLower(strings.TrimSpace(language))
	entry, ok := entries[key]
	if !ok {
		return "", "", "", fmt.Errorf("unsupported language: %s", language)
	}

	return entry.Image, entry.Install, entry.Timezone, nil
}

func loadLanguageBaseImages(source string) (map[string]languageEntry, error) {
	defaults := map[string]languageEntry{
		"go":     {Image: "golang:1.24-alpine"},
		"python": {Image: "python:3.12-alpine"},
		"node":   {Image: "node:22-alpine"},
		"rust":   {Image: "rust:1-alpine"},
	}

	resolved := map[string]languageEntry{}
	for k, v := range defaults {
		resolved[k] = v
	}

	raw, err := fetchBaseImageConfig(source)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return resolved, nil
	}

	var overrides map[string]json.RawMessage
	if err := json.Unmarshal(raw, &overrides); err != nil {
		return nil, fmt.Errorf("failed to parse base image config: %w", err)
	}

	for k, v := range overrides {
		normalizedKey := strings.ToLower(strings.TrimSpace(k))
		if normalizedKey == "" {
			continue
		}
		entry, err := parseLanguageEntry(v)
		if err != nil {
			return nil, fmt.Errorf("invalid entry for %q: %w", k, err)
		}

		base := resolved[normalizedKey]
		if entry.Image == "" {
			entry.Image = base.Image
		}
		if entry.Install == "" {
			entry.Install = base.Install
		}
		if entry.Timezone == "" {
			entry.Timezone = base.Timezone
		}
		resolved[normalizedKey] = entry
	}

	return resolved, nil
}

func parseLanguageEntry(raw json.RawMessage) (languageEntry, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return languageEntry{Image: strings.TrimSpace(s)}, nil
	}

	var obj struct {
		Image    string `json:"image"`
		Install  string `json:"install"`
		Timezone string `json:"timezone"`
	}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return languageEntry{}, fmt.Errorf("must be a string or {\"image\",\"install\",\"timezone\"} object: %w", err)
	}

	img := strings.TrimSpace(obj.Image)
	install := strings.TrimSpace(obj.Install)
	timezone := strings.TrimSpace(obj.Timezone)
	if img == "" && install != "" {
		img = "alpine:latest"
	}

	return languageEntry{Image: img, Install: install, Timezone: timezone}, nil
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
