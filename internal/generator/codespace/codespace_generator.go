package codespace

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/taka1156/codespacegen/internal/domain/entity"
)

//go:embed template/Dockerfile.tmpl template/docker-compose.yaml.tmpl
var templateFiles embed.FS

var (
	dockerfileTmpl = template.Must(template.ParseFS(templateFiles, "template/Dockerfile.tmpl"))
	composeTmpl    = template.Must(template.ParseFS(templateFiles, "template/docker-compose.yaml.tmpl"))
)

type dockerfileData struct {
	BaseImage       string
	WorkspaceFolder string
	Timezone        string
	Locale          entity.LocaleConfig
	BaseSetup       string
	TimezoneSetup   string
	RunCommandBlock string
}

type composeData struct {
	ServiceName     string
	WorkspaceFolder string
	PortMapping     string
}

type CodespaceGenerator struct{}

type baseImageStrategy interface {
	renderBaseSetup(locale entity.LocaleConfig, osModules entity.OsModules) string
	renderTimezoneSetup(timezone string) string
}

type alpineStrategy struct{}

type debianLikeStrategy struct{}

type devcontainerJSON struct {
	Schema          string                     `json:"$schema"`
	Name            string                     `json:"name"`
	Service         string                     `json:"service"`
	WorkspaceFolder string                     `json:"workspaceFolder"`
	DockerCompose   string                     `json:"dockerComposeFile"`
	Customizations  devcontainerCustomizations `json:"customizations"`
}

type devcontainerCustomizations struct {
	VSCode devcontainerVSCode `json:"vscode"`
}

type devcontainerVSCode struct {
	Settings   map[string]string `json:"settings"`
	Extensions []string          `json:"extensions"`
}

func NewCodespaceGenerator() *CodespaceGenerator {
	return &CodespaceGenerator{}
}

func (g *CodespaceGenerator) Generate(config entity.CodespaceConfig) ([]entity.GeneratedFile, error) {
	dockerfile, err := g.renderDockerfile(config)
	if err != nil {
		return nil, err
	}

	compose, err := g.renderCompose(config)
	if err != nil {
		return nil, err
	}

	devcontainer, err := g.renderDevcontainer(config)
	if err != nil {
		return nil, err
	}

	return []entity.GeneratedFile{
		{RelativePath: "Dockerfile", Content: dockerfile},
		{RelativePath: "devcontainer.json", Content: devcontainer},
		{RelativePath: config.ComposeFileName, Content: compose},
	}, nil
}

func (g *CodespaceGenerator) renderDockerfile(config entity.CodespaceConfig) (string, error) {
	locale := config.Locale
	if locale.Lang == "" {
		locale = entity.DefaultLocale
	}
	strategy := resolveBaseImageStrategy(config.BaseImage)
	baseSetup := strategy.renderBaseSetup(locale, config.OsModules)

	timezone := strings.TrimSpace(config.Timezone)
	if timezone == "" {
		timezone = entity.DefaultTimezone
	}
	timezoneSetup := strategy.renderTimezoneSetup(timezone)

	runCommandBlock := ""
	if config.RunCommand != "" {
		runCommandBlock = fmt.Sprintf("RUN %s\n\n", config.RunCommand)
	}

	var dockerfileBuf bytes.Buffer
	if err := dockerfileTmpl.Execute(&dockerfileBuf, dockerfileData{
		BaseImage:       config.BaseImage,
		WorkspaceFolder: config.WorkspaceFolder,
		Timezone:        timezone,
		Locale:          locale,
		BaseSetup:       baseSetup,
		TimezoneSetup:   timezoneSetup,
		RunCommandBlock: runCommandBlock,
	}); err != nil {
		return "", fmt.Errorf("failed to render Dockerfile: %w", err)
	}

	return dockerfileBuf.String(), nil
}

func (g *CodespaceGenerator) renderCompose(config entity.CodespaceConfig) (string, error) {
	var composeBuf bytes.Buffer
	if err := composeTmpl.Execute(&composeBuf, composeData{
		ServiceName:     config.ServiceName,
		WorkspaceFolder: config.WorkspaceFolder,
		PortMapping:     config.PortMapping,
	}); err != nil {
		return "", fmt.Errorf("failed to render docker-compose: %w", err)
	}

	return composeBuf.String(), nil
}

func (g *CodespaceGenerator) renderDevcontainer(config entity.CodespaceConfig) (string, error) {
	extensions := []string{}
	extensions = append(extensions, config.VSCodeExtensions...)
	extensions = uniqueStringsPreserveOrder(extensions)

	devcontainerObj := devcontainerJSON{
		Schema:          config.Schema,
		Name:            config.ContainerName,
		Service:         config.ServiceName,
		WorkspaceFolder: config.WorkspaceFolder,
		DockerCompose:   config.ComposeFileName,
		Customizations: devcontainerCustomizations{
			VSCode: devcontainerVSCode{
				Settings: map[string]string{
					"terminal.integrated.defaultProfile.linux": "bash",
				},
				Extensions: extensions,
			},
		},
	}

	devcontainerBytes, err := json.MarshalIndent(devcontainerObj, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to render devcontainer.json: %w", err)
	}

	return string(devcontainerBytes) + "\n", nil
}

func resolveBaseImageStrategy(baseImage string) baseImageStrategy {
	if isAlpineImage(baseImage) {
		return alpineStrategy{}
	}

	return debianLikeStrategy{}
}

func uniqueStringsPreserveOrder(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}

	return result
}

func isAlpineImage(baseImage string) bool {
	return strings.Contains(strings.ToLower(strings.TrimSpace(baseImage)), "alpine")
}

func (alpineStrategy) renderBaseSetup(_ entity.LocaleConfig, osModules entity.OsModules) string {
	return `RUN <<-EOF
apk add --no-cache \
  ` + strings.Join(osModules.AlpineModules, " \\\n  ") + `
EOF`
}

func (debianLikeStrategy) renderBaseSetup(locale entity.LocaleConfig, osModules entity.OsModules) string {
	return `RUN <<-EOF
apt update
apt install -y --no-install-recommends \
  ` + strings.Join(osModules.DebianLikeModules, " \\\n  ") + `
rm -rf /var/lib/apt/lists/*
locale-gen ` + locale.Lang + `
update-locale LANG=` + locale.Lang + ` LC_ALL=` + locale.LcAll + `
EOF`
}

func (alpineStrategy) renderTimezoneSetup(timezone string) string {
	return `RUN <<-EOF
ln -sf /usr/share/zoneinfo/` + timezone + ` /etc/localtime
echo "` + timezone + `" > /etc/timezone
EOF`
}

func (debianLikeStrategy) renderTimezoneSetup(timezone string) string {
	return `RUN <<-EOF
ln -fs /usr/share/zoneinfo/` + timezone + ` /etc/localtime
dpkg-reconfigure -f noninteractive tzdata
EOF`
}
