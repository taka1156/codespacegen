package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"codespacegen/internal/domain/entity"
)

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
	baseSetup := renderBaseSetupBlock(config.BaseImage)
	timezone := strings.TrimSpace(config.Timezone)
	if timezone == "" {
		timezone = entity.DefaultTimezone
	}
	timezoneSetup := renderTimezoneSetupBlock(config.BaseImage, timezone)

	installBlock := ""
	if config.InstallCommand != "" {
		installBlock = fmt.Sprintf("RUN %s\n\n", config.InstallCommand)
	}

	dockerfile := fmt.Sprintf(`FROM %s

WORKDIR %s

ENV LANG=ja_JP.UTF-8 \
    LANGUAGE=ja_JP:ja \
    LC_ALL=ja_JP.UTF-8 \
  TZ=%s

%s

RUN git lfs install

%s

RUN <<-EOF
curl -o ~/.git-prompt.sh https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh
curl -o ~/.git-completion.sh https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash

echo "[ -f ~/.git-prompt.sh ] && source ~/.git-prompt.sh" >> ~/.bashrc
echo "[ -f ~/.git-completion.sh ] && source ~/.git-completion.sh" >> ~/.bashrc
echo "GIT_PS1_SHOWDIRTYSTATE=true" >> ~/.bashrc
echo "GIT_PS1_SHOWUNTRACKEDFILES=true" >> ~/.bashrc
echo "GIT_PS1_SHOWUPSTREAM=auto"  >> ~/.bashrc
git config --system --add safe.directory %s
echo 'export PS1="\[\033[01;32m\]\u@\h\[\033[01;33m\] \w \[\033[01;31m\]\$(__git_ps1 \"(%%s)\") \\n+\[\033[01;34m\]\\$ \[\033[00m\]"' >> ~/.bashrc
EOF

%sCMD ["bash"]
`, config.BaseImage, config.WorkspaceFolder, timezone, baseSetup, timezoneSetup, config.WorkspaceFolder, installBlock)

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

	compose := fmt.Sprintf(`services:
    %s:
      build: .
      tty: true
      volumes:
        - ../:%s
`, config.ServiceName, config.WorkspaceFolder)

	if config.PortMapping != "" {
		compose += fmt.Sprintf("      ports:\n        - \"%s\"\n", config.PortMapping)
	}

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

func renderBaseSetupBlock(baseImage string) string {
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
locale-gen ja_JP.UTF-8
update-locale LANG=ja_JP.UTF-8 LC_ALL=ja_JP.UTF-8
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
