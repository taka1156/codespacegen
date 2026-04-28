---
name: codespacegen-project
description: 'Repository knowledge for the codespacegen project. Use when answering questions about project structure, architecture, CLI behavior, code generation, unit tests, e2e snapshot tests, docs updates, or implementing changes in this repo.'
---

# codespacegen Project Knowledge

## When to Use

- Answering questions about how this repository is organized
- Generating or modifying code in this repository
- Reviewing changes that depend on repository-specific behavior
- Updating README or other documentation
- Adding or fixing unit tests or e2e tests

## Project Overview

- `codespacegen` is a Go CLI that generates three devcontainer artifacts:
  - `Dockerfile`
  - `devcontainer.json`
  - `docker-compose.yaml`
- Entry point:
  - `cmd/codespacegen/main.go` — calls `app.NewApp().Run()`
- Main flow:
  - Parse CLI flags (`input.ClientInput`)
  - Collect inputs: CLI flags, JSON config, default settings
  - If `init` subcommand: generate `codespacegen.json` and exit
  - If `-v` flag: print version and exit
  - Assemble `entity.CodespaceConfig` via interactive prompts and merge logic
  - Execute `GenerateCodespaceArtifacts` — render templates and write files

## Architecture

- App / DI root:
  - `internal/app/app.go` — `App` struct, `NewApp()` wires all dependencies, `Run()` orchestrates the workflow
  - `internal/app/app_interfaces.go` — internal interfaces used by `App`
- Domain:
  - `internal/domain/entity/` — entity types (`CodespaceConfig`, `ClientConfig`, `JsonConfig`, `LangEntry`, `GeneratedFile`, `TemplateJson`, etc.)
  - `internal/domain/service/` — service interfaces (`CodespaceGenerator`, `LocalFileWriter`, `SettingTemplateGenerator`, `WorkdirProvider`)
- Input adapters:
  - `internal/input/` — `ClientInput` (CLI flags), `JsonInput` (JSON config loader via file or HTTPS), `DefaultConfig` (hardcoded defaults)
- Infra (external I/O):
  - `internal/infra/infra.go` — type alias facade; exports `CodespacePrompter`
  - `internal/infra/prompt/` — `CodespacegenPrompter` (stdin-based interactive prompter)
- Generator (template rendering and file writing):
  - `internal/generator/genarator.go` — factory functions for generators
  - `internal/generator/codespace/` — `CodespaceGenerator` (renders `Dockerfile`, `devcontainer.json`, `docker-compose.yaml`)
  - `internal/generator/setting/` — `SettingTemplateGenerator` (renders `codespacegen.json`)
  - `internal/generator/config/` — `ConfigTemplateGenerator` (JSON marshal for template config)
  - `internal/generator/filewriter/` — `LocalFileWriter`
  - `internal/generator/workdirprovider/` — `WorkdirProvider`
- Workflow (use-case layer):
  - `internal/workflow/workflow.go` — facade (type aliases and constructor wrappers only)
  - `internal/workflow/collect/` — `CollectInputs` (gathers CLI, JSON, and default config into `CollectedInputs`)
  - `internal/workflow/assemble/` — `AssembleCodespaceConfig` (interactive prompt resolution, entry merge, config build)
  - `internal/workflow/generate/` — `GenerateCodespaceArtifacts` (validate, render, write files)
  - `internal/workflow/initialize/` — `InitializeSettingJson` (generate and write `codespacegen.json` template)
- Utilities:
  - `internal/utils/` — shared helper functions
  - `internal/i18n/` — locale-based message lookup (`locales/en.yaml`, `locales/ja.yaml`)

Dependencies point inward (domain has no outward dependencies).

## Interface Mapping

| Interface | Defined in | Implemented by |
|---|---|---|
| `service.CodespaceGenerator` | `internal/domain/service` | `codespace.CodespaceGenerator` |
| `service.LocalFileWriter` | `internal/domain/service` | `filewriter.LocalFileWriter` |
| `service.SettingTemplateGenerator` | `internal/domain/service` | `setting.SettingTemplateGenerator` |
| `service.WorkdirProvider` | `internal/domain/service` | `workdirprovider.WorkdirProvider` |
| `collect.ClientInputProvider` | `internal/workflow/collect` | `input.ClientInput` |
| `collect.JsonConfigLoader` | `internal/workflow/collect` | `input.JsonInput` |
| `collect.DefaultSettingProvider` | `internal/workflow/collect` | `input.DefaultConfig` |
| `assemble.CodespacegenPrompter` | `internal/workflow/assemble` | `infra.CodespacePrompter` (= `prompt.CodespacegenPrompter`) |

## Configuration Knowledge

- Base image resolution supports built-in language keys and custom keys from `codespacegen.json`.
- `-image-config` accepts a local path or `https://` URL.
- `codespacegen.json` supports:
  - top-level `common`
  - per-language object entries with `image`, `runCommand`, `timezone`, `locale`, and `vscodeExtensions`
- Merge behavior:
  - `common` is applied first
  - language-specific values override or extend the common values
  - `vscodeExtensions` are appended and later deduplicated in generated output
- Base image resolution priority: explicit `-base-image` flag > language key lookup in JSON config > default image
- If timezone is not provided by prompt, flags, or config, the effective default is `UTC`.

## Generation Knowledge

- `CodespaceGenerator` in `internal/generator/codespace/` renders all three devcontainer files using embedded templates.
- The generator chooses package setup based on the base image:
  - Alpine-like images use `apk`
  - Non-Alpine images use `apt-get`
- Generated `devcontainer.json` always includes:
  - `GitHub.copilot`
  - `GitHub.copilot-chat`
- Additional VS Code extensions from config are merged and deduplicated.
- `docker-compose.yaml` includes `ports` only when a port mapping is provided.

## Testing Knowledge

### Unit tests

- Run all unit tests:
  - `go test ./...`
- Main test files:
  - `internal/workflow/generate/generate_codespace_artifacts_test.go` — file write behavior and error propagation
  - `internal/generator/codespace/codespace_generator_test.go` — template rendering details (package manager selection, timezone, extensions, key order)
  - `internal/workflow/assemble/assemble_codespace_config_test.go` — config assembly with mocked prompter
  - `internal/workflow/assemble/assemble_config_builder_test.go` — base image resolution and config build logic
  - `internal/workflow/assemble/assemble_config_entry_test.go` — JSON entry merge logic
  - `internal/workflow/collect/collect_inputs_test.go` — input collection logic
  - `internal/workflow/initialize/initialize_template_json_test.go` — setting JSON initialization
  - `internal/input/json_input_test.go` — HTTP/file loading and validation
  - `internal/infra/prompt/prompt_test.go` — interactive prompter behavior
  - `internal/generator/setting/setting_template_generator_test.go` — setting template generation
  - `internal/generator/config/config_template_generator_test.go` — config template generation

### E2E snapshot tests

- Run e2e tests from the repository root:
  - `make e2e`
- `make e2e`:
  - builds binary into `bin/codespacegen`
  - copies it into `e2e/devcontainer_config/` and `e2e/codespacegen_config/`
  - runs `e2e/devcontainer_config/devcontainer_config.test.sh` — compares generated devcontainer files with snapshots under `e2e/devcontainer_config/snapshots/.devcontainer-*`
  - runs `e2e/codespacegen_config/codespacegen_config.test.sh` — verifies `codespacegen init` output against `e2e/codespacegen_config/snapshots/codespacegen.json`
- Current snapshot cases for devcontainer_config:
  - `go`
  - `python`
  - `rust`
  - `moonbit`
  - `node:biome`
  - `node:eslint`
  - `node:react` (suffix `react`)
  - `gcc`
  - `node:zenn`
- Important e2e behavior:
  - the script uses `e2e/devcontainer_config/codespacegen.json` as `-image-config`
  - it adds `-port 3000` only when the snapshot `docker-compose.yaml` contains a `ports:` block
  - it fails on any diff in `Dockerfile`, `devcontainer.json`, or `docker-compose.yaml`
  - snapshots can be updated with `make e2e UPD=--update`

## Change Guidance

- Prefer changes that preserve deterministic generated output.
- When changing flag behavior or defaults, update tests and README consistently.
- When changing template output, expect e2e snapshots to need updates (run `make e2e UPD=--update`).
- When changing config merge behavior, check both unit tests and snapshot coverage.
