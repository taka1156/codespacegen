package input

import (
	"flag"
	"fmt"
	"os"

	"github.com/taka1156/codespacegen/internal/domain/entity"
)

type ClientInput struct {
}

func NewClientInput() *ClientInput {
	return &ClientInput{}
}

func (ci *ClientInput) GetInput(args []string) entity.ClientConfig {
	clientConfig := entity.ClientConfig{}

	if len(args) > 1 {
		switch args[1] {
		case "init":
			initCmd := flag.NewFlagSet("init", flag.ExitOnError)
			outputDir := initCmd.String("output", ".devcontainer", "output directory for generated files")
			initCmd.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage: %s init [options]\n\n", os.Args[0])
				fmt.Fprintf(os.Stderr, "Initialize setting JSON\n")
				initCmd.PrintDefaults()
			}
			_ = initCmd.Parse(args[2:])
			clientConfig.OutputDir = outputDir
			clientConfig.Mode = entity.Initialize
			return clientConfig
		case "update":
			updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
			lang := updateCmd.String("lang", "", "language for CLI messages (en/ja, default: auto-detect)")
			updateCmd.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage: %s update [options]\n\n", os.Args[0])
				fmt.Fprintf(os.Stderr, "Update codespacegen to the latest version\n")
				updateCmd.PrintDefaults()
			}
			_ = updateCmd.Parse(args[2:])
			clientConfig.Lang = lang
			clientConfig.Mode = entity.Update
			return clientConfig
		case "version":
			versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
			versionCmd.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage: %s version\n\n", os.Args[0])
				fmt.Fprintf(os.Stderr, "Show current version of codespacegen\n")
				versionCmd.PrintDefaults()
			}
			_ = versionCmd.Parse(args[2:])
			clientConfig.Mode = entity.Version
			return clientConfig
		}
	}

	fs := flag.NewFlagSet("root", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [command]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  version\tShow current version of codespacegen\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [command] [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  init\t\tInitialize setting JSON\n")
		fmt.Fprintf(os.Stderr, "  update\tUpdate codespacegen to the latest version\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	clientConfig.OutputDir = fs.String("output", ".devcontainer", "output directory for generated files")
	clientConfig.ContainerName = fs.String("name", "", "project name (required, mapped to devcontainer name)")
	clientConfig.ServiceName = fs.String("service", "", "docker compose service name")
	clientConfig.Language = fs.String("language", "", "programming language (go/python/node/rust or image-config keys)")
	clientConfig.WorkspaceFolder = fs.String("workspace-folder", "/workspace", "workspace folder inside container")
	clientConfig.Timezone = fs.String("timezone", "", "timezone inside container (default: image-config timezone or UTC)")
	clientConfig.ImageConfig = fs.String("image-config", "codespacegen.json", "local path or https:// URL to base image config JSON")
	clientConfig.Port = fs.String("port", "", "port mapping (e.g. 3000 or 3000:3000)")
	clientConfig.ComposeFile = fs.String("compose-file", "docker-compose.yaml", "docker compose file name")
	clientConfig.EnableOverwriteFile = fs.Bool("force", false, "overwrite existing files")
	clientConfig.Lang = fs.String("lang", "", "language for CLI messages (en/ja, default: auto-detect)")
	clientConfig.Headless = fs.Bool("headless", false, "run in headless mode without interactive prompts")

	if len(args) > 1 {
		_ = fs.Parse(args[1:])
	} else {
		_ = fs.Parse(args)
	}

	return clientConfig
}
