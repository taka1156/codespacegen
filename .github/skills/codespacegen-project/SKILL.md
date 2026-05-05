---
name: codespacegen-project
description: 'Repository knowledge for the codespacegen project. Use when answering questions about project structure, architecture, CLI behavior, code generation, unit tests, e2e snapshot tests, docs updates, or implementing changes in this repo.'
---

# codespacegen Project Knowledge

## Important Behavioral Rules

- **Do NOT modify any source code without explicit user permission.**
- Before making any code change (editing files, adding files, refactoring, etc.), always ask the user for confirmation first.
- Read-only operations (searching, reading files, running tests to observe output) do not require permission.
- If the user asks a question or requests analysis, provide the answer or analysis only — do not apply changes unless the user explicitly approves them.

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
  - Assemble `entity.CodespaceConfig` — interactive prompts (or headless from flags) and merge logic
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
  - `internal/generator/generator.go` — factory functions for generators
  - `internal/generator/codespace/` — `CodespaceGenerator` (renders `Dockerfile`, `devcontainer.json`, `docker-compose.yaml`)
  - `internal/generator/setting/` — `SettingTemplateGenerator` (renders `codespacegen.json`)
  - `internal/generator/filewriter/` — `LocalFileWriter`
  - `internal/generator/workdirprovider/` — `WorkdirProvider`
- Workflow (use-case layer):
  - `internal/workflow/workflow.go` — facade (type aliases and constructor wrappers only)
  - `internal/workflow/collect/` — `CollectInputs` (gathers CLI, JSON, and default config into `CollectedInputs`)
  - `internal/workflow/assemble/` — `AssembleCodespaceConfig` (interactive prompt resolution, entry merge, config build)
  - `internal/workflow/generate/` — `GenerateCodespaceArtifacts` (validate, render, write files)
  - `internal/workflow/initialize/` — `InitializeSettingJson` (generate and write `codespacegen.json` template)
- Utilities:
  - `internal/utils/` — shared helper functions including `NormalizePortMapping`
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
  - top-level `common` with `locale`, `timezone`, and `vscodeExtensions`
  - top-level `langs` array — each item is a `LangEntry` object with:
    - `profileName` (required) — identifies the language profile (e.g. `"go"`, `"node:eslint"`)
    - `image` — base Docker image
    - `runCommand` — injected as `RUN` in the Dockerfile
    - `linuxPackages` — list of OS packages to install
    - `vscodeExtensions` — language-specific VS Code extension IDs
  - `locale` and `timezone` are **only** in `common`; `LangEntry` does not carry them
- Merge behavior:
  - `common.locale` is applied to all generated output (not per-language)
  - `common.timezone` is used as the fallback default when no flag or prompt value is given
  - `vscodeExtensions`: `common` extensions are prepended to language-specific extensions, then deduplicated
  - `linuxPackages` in a lang entry are appended to the default OS modules for that image type
- Base image resolution priority: `profileName` lookup in `langs` array of JSON config > default image
- Locale resolution: `jsonConfig.Common.Locale` > `defaultSetting.Locale`
- Timezone resolution priority: explicit flag > `jsonConfig.Common.Timezone` > `defaultSetting.Timezone` (UTC)
- `-headless` flag: skips all interactive prompts; all values must be supplied via CLI flags

## Generation Knowledge

- `CodespaceGenerator` in `internal/generator/codespace/` renders all three devcontainer files using embedded templates.
- The generator chooses package setup based on the base image:
  - Alpine-like images use `apk`
  - Non-Alpine images use `apt-get`
- VS Code extensions in `devcontainer.json` are taken directly from `config.VSCodeExtensions` (merged from `common` + lang-specific entries, deduplicated). The generator does not inject any hardcoded extensions.
- `docker-compose.yaml` includes `ports` only when a port mapping is provided.
- Port mapping is normalized at assembly time via `utils.NormalizePortMapping`: a bare port number `N` becomes `N:N`.

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

### E2E snapshot tests

- Run e2e tests from the repository root:
  - `make e2e`
- `make e2e`:
  - builds binary into `bin/codespacegen`
  - copies it into `e2e/devcontainer_config/` and `e2e/codespacegen_config/`
  - runs `e2e/devcontainer_config/devcontainer_config.test.sh` — compares generated devcontainer files with snapshots under `e2e/devcontainer_config/snapshots/.devcontainer-*`
  - runs `e2e/codespacegen_config/codespacegen_config.test.sh` — verifies `codespacegen init` output against `e2e/codespacegen_config/snapshots/codespacegen.json`
- Current snapshot cases for devcontainer_config:
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
