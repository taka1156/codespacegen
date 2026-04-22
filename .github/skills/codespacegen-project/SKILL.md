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
  - `cmd/codespacegen/main.go`
- Main flow:
  - Parse CLI flags
  - Resolve interactive/default values
  - Build `entity.CodespaceConfig`
  - Execute `GenerateCodespaceArtifacts`
  - Render templates and write files

## Architecture

- Domain:
  - `internal/domain/entity` — エンティティ型 (`CodespaceConfig`, `CliConfig`, `JsonEntry`, `GeneratedFile` など)
  - `internal/domain/service` — `TemplateGenerator` インターフェース
- Config (入力アダプター):
  - `internal/config/` — CLI (`CliInput`)、JSON (`JsonInput`)、デフォルト値 (`DefaultConfig`)
- Resolve (インタラクティブ解決):
  - `internal/resolve/` — `CodeSpaceConfigResolver` (stdin プロンプト、ベースイメージ解決、マージ)
- Generator (テンプレート生成・書き込み実装):
  - `internal/generator/` — `DefaultTemplateGenerator`
  - `internal/generator/filewriter/` — `LocalFileWriter`
- Workflow (ユースケース層):
  - `internal/workflow/workflow.go` — ファサード（型エイリアスのみ）
  - `internal/workflow/collect/` — 入力収集 (`CollectInputs`)
  - `internal/workflow/assemble/` — 設定解決・構築 (`ResolveCodespaceConfig`)
  - `internal/workflow/generate/` — アーティファクト生成 (`GenerateCodespaceArtifacts`)
- i18n:
  - `internal/i18n/` — ロケール別メッセージ (`locales/ja.yaml`, `locales/en.yaml`)
- Entry point:
  - `cmd/codespacegen/main.go` — DI ルート、`App` 構造体

Dependencies point inward.

## Interface Mapping

| インターフェース | 定義場所 | 実装 |
|---|---|---|
| `service.TemplateGenerator` | `internal/domain/service` | `generator.DefaultTemplateGenerator` |
| `generate.FileWriter` | `internal/workflow/generate` | `filewriter.LocalFileWriter` |
| `collect.CLIInputProvider` | `internal/workflow/collect` | `config.CliInput` |
| `collect.ImageConfigLoader` | `internal/workflow/collect` | `config.JsonInput` |
| `collect.DefaultSettingProvider` | `internal/workflow/collect` | `config.DefaultConfig` |

> **注意**: `resolve.CodeSpaceConfigResolver` は現在インターフェースではなく具体構造体。`assemble` 層が直接依存しており、モック置き換えが不可。

## Configuration Knowledge

- Base image resolution supports built-in language keys and custom keys from `codespacegen.json`.
- `-image-config` accepts a local path or `https://` URL.
- `codespacegen.json` supports:
  - top-level `common`
  - per-language string entries
  - per-language object entries with `image`, `install`, `timezone`, and `vscodeExtensions`
- Merge behavior:
  - `common` is applied first
  - language-specific values override or extend the common values
  - `vscodeExtensions` are appended and later deduplicated in generated output
- If `image` is omitted and `install` is present, `alpine:latest` is used.
- If timezone is not provided by flags or config, the effective default is `UTC`.

## Generation Knowledge

- Template generation happens in `internal/generator/default_template_generator.go`.
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
  - `internal/workflow/generate/generate_artifacts_test.go` — ファイル書き込みの動作・エラー伝播
  - `internal/generator/default_template_generator_test.go` — テンプレートレンダリング詳細
- Main unit test responsibilities:
  - use case write behavior and error propagation
  - template rendering details such as package manager selection, timezone setup, extension merging, and key order
- **テストが存在しないパッケージ（今後追加が必要）**:
  - `internal/workflow/assemble/` — 設定解決ロジック全体
  - `internal/resolve/` — インタラクティブ解決ロジック

### E2E snapshot tests

- Run e2e tests from the repository root:
  - `make e2e`
- `make e2e`:
  - builds `e2e/codespacegen`
  - executes `e2e/e2e.sh`
  - compares generated files with snapshots under `e2e/snapshots/.devcontainer-*`
- Current snapshot cases include:
  - `python`
  - `rust`
  - `moonbit`
  - `node:biome`
  - `node:eslint`
  - `node:react`
- Important e2e behavior:
  - the script uses `e2e/codespacegen.json` as `-image-config`
  - it adds `-port 3000` only when the snapshot `docker-compose.yaml` contains a `ports:` block
  - it fails on any diff in `Dockerfile`, `devcontainer.json`, or `docker-compose.yaml`

## Change Guidance

- Prefer changes that preserve deterministic generated output.
- When changing flag behavior or defaults, update tests and README consistently.
- When changing template output, expect e2e snapshots to need updates.
- When changing config merge behavior, check both unit tests and snapshot coverage.
