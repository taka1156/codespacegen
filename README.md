 <picture>
 	  <source media="(prefers-color-scheme: dark)" srcset="./logo-dark.svg">
	  <source media="(prefers-color-scheme: light)" srcset="./logo-light.svg">
	  <img alt="codespacegen logo" src="./logo-light.svg" width="100%" height="100%">
  </picture>

  ![GitHub Release](https://img.shields.io/github/v/release/taka1156/codespacegen?sort=semver&display_name=release&color=60a5fa&link=https%3A%2F%2Fgithub.com%2Ftaka1156%2Fcodespacegen%2Freleases%2F)
  ![GitHub Release Date](https://img.shields.io/github/release-date/taka1156/codespacegen?color=60a5fa)
  ![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/taka1156/codespacegen/release.yml?logo=github&color=60a5fa)
  ![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/taka1156/codespacegen/main.yaml?event=push&logo=github&label=test&color=60a5fa)
	![GitHub License](https://img.shields.io/github/license/taka1156/codespacegen?color=60a5fa)
  ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/taka1156/codespacegen?color=60a5fa&logo=go&logoColor=white)


[日本語版はこちら](README.ja.md)

codespacegen is a CLI that generates the following three files for Codespaces and Dev Container.

- Dockerfile
- devcontainer.json
- docker-compose.yaml

## Install with curl

The latest release is downloaded automatically and installed into `/usr/local/bin`.

```bash
curl -fsSL https://raw.githubusercontent.com/taka1156/codespacegen/master/scripts/install.sh | bash
```

To change the install destination:

```bash
curl -fsSL https://raw.githubusercontent.com/taka1156/codespacegen/master/scripts/install.sh | INSTALL_DIR=$HOME/.local/bin bash
```

## Release with GitHub Actions

Main generated assets:

- `codespacegen_linux_amd64.tar.gz`
- `codespacegen_linux_arm64.tar.gz`
- `codespacegen_darwin_amd64.tar.gz`
- `codespacegen_darwin_arm64.tar.gz`
- `codespacegen_windows_amd64.exe`
- `checksums.txt`

## Architecture

- Domain: rules and models
	- internal/domain/entity
	- internal/domain/service
- App: composition and orchestration
	- internal/app
- Input adapters (CLI/JSON/defaults)
	- internal/input
- Infra (interactive prompt I/O)
	- internal/infra
- Workflows (use cases)
	- internal/workflow/collect
	- internal/workflow/assemble
	- internal/workflow/generate
	- internal/workflow/initialize
- Artifact generation and file writing
	- internal/generator
	- internal/generator/filewriter
	- internal/generator/workdirprovider
- i18n resources
	- internal/i18n
- Entry point: CLI
	- cmd/codespacegen

Dependencies only point inward.

## Usage

### Run

```bash
go run ./cmd/codespacegen
```

By default, files are generated under .devcontainer.

### Initialize codespacegen.json

Run the `init` subcommand to generate a `codespacegen.json` template in the current directory.

```bash
codespacegen init
```

The generated file serves as a starting point for customizing base images and VS Code extensions.

### Main options

| Option | Default | Description |
|---|---|---|
| `-output` | `.devcontainer` | Output directory |
| `-name` | *(interactive, required)* | Project name. Prompted every time and mapped to the `name` field in `devcontainer.json` |
| `-language` | *(interactive, empty on Enter)* | Programming language key. Prompted every time. Any `profileName` defined in the `langs` array of `codespacegen.json` (or the file specified by `-image-config`) can be used. If empty, no language-specific setting is used and `alpine:latest` is selected |
| `-service` | *(interactive, `app` on Enter)* | Docker Compose service name. Prompted every time and reflected in both `devcontainer.json` and `docker-compose.yaml` |
| `-workspace-folder` | *(interactive, `/workspace` on Enter)* | Workspace path inside the container. Prompted every time |
| `-timezone` | *(interactive, default from `common.timezone` or `UTC`)* | Timezone inside the container. Prompted every time and reflected in `ENV TZ` and timezone setup in the Dockerfile |
| `-image-config` | `codespacegen.json` | Local path or `https://` URL for base image definitions. Supports top-level `common` defaults plus a `langs` array for per-language entries. `image` is required when `runCommand` or `linuxPackages` is specified |
| `-port` | *(interactive, no ports on Enter)* | Port mapping. For example, `3000` is normalized to `3000:3000`, and `8080:3000` is also accepted. Prompted every time |
| `-compose-file` | `docker-compose.yaml` | Compose file name |
| `-force` | `false` | Overwrite existing files |
| `-lang` | *(auto-detect)* | Language for CLI messages (`en` or `ja`). Defaults to system locale |
| `-headless` | `false` | Skip all interactive prompts. All required values must be supplied via flags |
| `-v` | — | Print version and exit |

Base image definitions are separated into [codespacegen.json](codespacegen.json) at the repository root.

- If the JSON file exists: values are loaded from the file

In addition, extension IDs from `codespacegen.json` (`common.vscodeExtensions` and per-language `vscodeExtensions`) are appended.

### codespacegen.json format

You can attach the JSON Schema in editors that support JSON Schema validation and completion.

```json
{
	"$schema": "./codespacegen.schema.json",
	"langs": [
		{
			"profileName": "go",
			"image": "golang:1.24-alpine"
		}
	]
}
```

If `codespacegen.json` is at the repository root, `./codespacegen.schema.json` points to the bundled schema file in this repository.

**Language entries (`langs` array)**

Each entry in `langs` requires `profileName` and supports the following fields (`image` is required when `runCommand` or `linuxPackages` is specified):

```json
{
	"langs": [
		{
			"profileName": "go",
			"image": "golang:1.24-alpine",
			"vscodeExtensions": ["golang.Go"]
		},
		{
			"profileName": "moonbit",
			"image": "ubuntu:24.04",
			"runCommand": "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash",
			"vscodeExtensions": ["moonbit.moonbit-lang"]
		},
		{
			"profileName": "gcc",
			"image": "ubuntu:24.04",
			"linuxPackages": ["gcc", "make", "git", "binutils", "libc6-dev"],
			"vscodeExtensions": ["ms-vscode.cpptools"]
		}
	]
}
```

The `runCommand` value is injected as a `RUN` step in the Dockerfile:

```dockerfile
RUN curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash
```

`linuxPackages` specifies Linux system packages. They are merged with the default package list and installed by the package manager appropriate for the base image (e.g. `apt` for Debian/Ubuntu, `apk` for Alpine).

**Shared defaults with `common`**

```json
{
	"common": {
		"timezone": "Asia/Tokyo",
		"locale": {
			"lang": "ja_JP.UTF-8",
			"language": "ja_JP:ja",
			"lcAll": "ja_JP.UTF-8"
		},
		"vscodeExtensions": [
			"MS-CEINTL.vscode-language-pack-ja",
			"streetsidesoftware.code-spell-checker"
		]
	},
	"langs": [
		{
			"profileName": "go",
			"image": "golang:1.24-alpine",
			"vscodeExtensions": ["golang.Go"]
		}
	]
}
```

Merge behavior:

- `common` is applied first, then language-specific values override/append
- `vscodeExtensions` are merged in order and de-duplicated
- `timezone` and `locale` can only be set in `common`, not per-language
- If timezone is not set in flags or `common`, `UTC` is used

Example:

```bash
go run ./cmd/codespacegen \
	-output .devcontainer \
	-name "My Codespace" \
	-language go \
	-service app \
	-workspace-folder /workspace \
	-timezone Asia/Tokyo \
	-compose-file docker-compose.yaml \
	-force
```

Example using a remote JSON URL:

```bash
go run ./cmd/codespacegen -image-config https://example.com/my-base-images.json -language go -force
```

- Only `https://` URLs are supported. `http://` is rejected
- If the JSON is missing or not specified, built-in Alpine defaults are used

Example exposing a port:

```bash
go run ./cmd/codespacegen -language go -port 3000 -force
```

If `-port` is not specified, the CLI prompts for a port interactively during execution.

The generated `docker-compose.yaml` looks like this, with `ports` added only when a port is provided.

```yaml
services:
		app:
			build: .
			tty: true
			volumes:
				- ../:/workspace
```

## Tests

```bash
go test ./...
```
