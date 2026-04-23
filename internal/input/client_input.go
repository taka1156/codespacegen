package input

import (
	"codespacegen/internal/domain/entity"
	"flag"
)

type ClientInput struct {
}

func NewClientInput() *ClientInput {
	return &ClientInput{}
}

func (ci *ClientInput) GetInput() entity.ClientConfig {
	ClientConfig := entity.ClientConfig{}

	ClientConfig.OutputDir = flag.String("output", ".devcontainer", "output directory for generated files")
	ClientConfig.ContainerName = flag.String("name", "", "project name (required, mapped to devcontainer name)")
	ClientConfig.ServiceName = flag.String("service", "", "docker compose service name")
	ClientConfig.Language = flag.String("language", "", "programming language (go/python/node/rust or image-config keys)")
	ClientConfig.WorkspaceFolder = flag.String("workspace-folder", "/workspace", "workspace folder inside container")
	ClientConfig.BaseImage = flag.String("base-image", "", "base Docker image (overrides -language default)")
	ClientConfig.Timezone = flag.String("timezone", "", "timezone inside container (default: image-config timezone or UTC)")
	ClientConfig.ImageConfig = flag.String("image-config", "codespacegen.json", "local path or https:// URL to base image config JSON")
	ClientConfig.Port = flag.String("port", "", "port mapping (e.g. 3000 or 3000:3000)")
	ClientConfig.ComposeFile = flag.String("compose-file", "docker-compose.yaml", "docker compose file name")
	ClientConfig.EnableOverwriteFile = flag.Bool("force", false, "overwrite existing files")
	ClientConfig.Lang = flag.String("lang", "", "language for CLI messages (en/ja, default: auto-detect)")
	ClientConfig.ShowVersion = flag.Bool("v", false, "print version and exit")

	flag.Parse()

	return ClientConfig
}
