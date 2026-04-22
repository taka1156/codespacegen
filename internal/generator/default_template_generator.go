package generator

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"codespacegen/internal/domain/entity"
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
	InstallBlock    string
}

type composeData struct {
	ServiceName     string
	WorkspaceFolder string
	PortMapping     string
}

type DefaultTemplateGenerator struct{}

type baseImageStrategy interface {
	renderBaseSetup(locale entity.LocaleConfig) string
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

func NewDefaultTemplateGenerator() *DefaultTemplateGenerator {
	return &DefaultTemplateGenerator{}
}

func (g *DefaultTemplateGenerator) Generate(config entity.CodespaceConfig) ([]entity.GeneratedFile, error) {
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

func (g *DefaultTemplateGenerator) renderDockerfile(config entity.CodespaceConfig) (string, error) {
	locale := config.Locale
	if locale.Lang == "" {
		locale = entity.DefaultLocale
	}
	strategy := resolveBaseImageStrategy(config.BaseImage)
	baseSetup := strategy.renderBaseSetup(locale)

	timezone := strings.TrimSpace(config.Timezone)
	if timezone == "" {
		timezone = entity.DefaultTimezone
	}
	timezoneSetup := strategy.renderTimezoneSetup(timezone)

	installBlock := ""
	if config.InstallCommand != "" {
		installBlock = fmt.Sprintf("RUN %s\n\n", config.InstallCommand)
	}

	var dockerfileBuf bytes.Buffer
	if err := dockerfileTmpl.Execute(&dockerfileBuf, dockerfileData{
		BaseImage:       config.BaseImage,
		WorkspaceFolder: config.WorkspaceFolder,
		Timezone:        timezone,
		Locale:          locale,
		BaseSetup:       baseSetup,
		TimezoneSetup:   timezoneSetup,
		InstallBlock:    installBlock,
	}); err != nil {
		return "", fmt.Errorf("failed to render Dockerfile: %w", err)
	}

	return dockerfileBuf.String(), nil
}

func (g *DefaultTemplateGenerator) renderCompose(config entity.CodespaceConfig) (string, error) {
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

func (g *DefaultTemplateGenerator) renderDevcontainer(config entity.CodespaceConfig) (string, error) {
	extensions := []string{"GitHub.copilot", "GitHub.copilot-chat"}
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

func (alpineStrategy) renderBaseSetup(_ entity.LocaleConfig) string {
	return `RUN <<-EOF
apk add --no-cache \
  bash \
  bash-completion \
  ca-certificates \
  tzdata \
  git \
  git-lfs \
  vim \
  curl \
  musl-locales \
  musl-locales-lang
EOF`
}

func (debianLikeStrategy) renderBaseSetup(locale entity.LocaleConfig) string {
	return `RUN <<-EOF
apt-get update
apt-get install -y --no-install-recommends \
  bash \
  bash-completion \
  ca-certificates \
  tzdata \
  git \
  git-lfs \
  vim \
  curl \
  locales
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
