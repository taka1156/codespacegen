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

type devcontainerJSON struct {
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
	locale := config.Locale
	if locale.Lang == "" {
		locale = entity.DefaultLocale
	}

	baseSetup := renderBaseSetupBlock(config.BaseImage, locale)
	timezone := strings.TrimSpace(config.Timezone)
	if timezone == "" {
		timezone = entity.DefaultTimezone
	}
	timezoneSetup := renderTimezoneSetupBlock(config.BaseImage, timezone)

	installBlock := ""
	if config.InstallCommand != "" {
		installBlock = fmt.Sprintf("RUN %s\n\n", config.InstallCommand)
	}

	// Dockerfile
	dockerfileTmpl, err := template.ParseFS(templateFiles, "template/Dockerfile.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse Dockerfile template: %w", err)
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
		return nil, fmt.Errorf("failed to render Dockerfile: %w", err)
	}
	dockerfile := dockerfileBuf.String()

	// docker-compose
	composeTmpl, err := template.ParseFS(templateFiles, "template/docker-compose.yaml.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse docker-compose template: %w", err)
	}
	var composeBuf bytes.Buffer
	if err := composeTmpl.Execute(&composeBuf, composeData{
		ServiceName:     config.ServiceName,
		WorkspaceFolder: config.WorkspaceFolder,
		PortMapping:     config.PortMapping,
	}); err != nil {
		return nil, fmt.Errorf("failed to render docker-compose: %w", err)
	}
	compose := composeBuf.String()

	extensions := []string{"GitHub.copilot", "GitHub.copilot-chat"}
	extensions = append(extensions, config.VSCodeExtensions...)
	extensions = uniqueStringsPreserveOrder(extensions)

	devcontainerObj := devcontainerJSON{
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
		return nil, fmt.Errorf("failed to render devcontainer.json: %w", err)
	}
	devcontainer := string(devcontainerBytes) + "\n"

	return []entity.GeneratedFile{
		{RelativePath: "Dockerfile", Content: dockerfile},
		{RelativePath: "devcontainer.json", Content: devcontainer},
		{RelativePath: config.ComposeFileName, Content: compose},
	}, nil
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

func renderBaseSetupBlock(baseImage string, locale entity.LocaleConfig) string {
	if isAlpineImage(baseImage) {
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

func renderTimezoneSetupBlock(baseImage string, timezone string) string {
	if isAlpineImage(baseImage) {
		return `RUN <<-EOF
ln -sf /usr/share/zoneinfo/` + timezone + ` /etc/localtime
echo "` + timezone + `" > /etc/timezone
EOF`
	}

	return `RUN <<-EOF
ln -fs /usr/share/zoneinfo/` + timezone + ` /etc/localtime
dpkg-reconfigure -f noninteractive tzdata
EOF`
}
