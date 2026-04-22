package input

import (
	"codespacegen/internal/domain/entity"
	"flag"
)

type CliInput struct {
}

func NewCliInput() *CliInput {
	return &CliInput{}
}

func (ci *CliInput) GetCliInput() entity.CliConfig {
	cliConfig := entity.CliConfig{}

	cliConfig.OutputDir = flag.String("output", ".devcontainer", "output directory for generated files")
	cliConfig.ContainerName = flag.String("name", "", "project name (required, mapped to devcontainer name)")
	cliConfig.ServiceName = flag.String("service", "", "docker compose service name")
	cliConfig.Language = flag.String("language", "", "programming language (go/python/node/rust or image-config keys)")
	cliConfig.WorkspaceFolder = flag.String("workspace-folder", "/workspace", "workspace folder inside container")
	cliConfig.BaseImage = flag.String("base-image", "", "base Docker image (overrides -language default)")
	cliConfig.Timezone = flag.String("timezone", "", "timezone inside container (default: image-config timezone or UTC)")
	cliConfig.ImageConfig = flag.String("image-config", "codespacegen.json", "local path or https:// URL to base image config JSON")
	cliConfig.Port = flag.String("port", "", "port mapping (e.g. 3000 or 3000:3000)")
	cliConfig.ComposeFile = flag.String("compose-file", "docker-compose.yaml", "docker compose file name")
	cliConfig.EnableOverwriteFile = flag.Bool("force", false, "overwrite existing files")
	cliConfig.Lang = flag.String("lang", "", "language for CLI messages (en/ja, default: auto-detect)")
	cliConfig.ShowVersion = flag.Bool("v", false, "print version and exit")

	flag.Parse()

	return cliConfig
}
