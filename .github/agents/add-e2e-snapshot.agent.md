---
name: add-e2e-snapshot
description: Add a new e2e snapshot test case for a language profile in the codespacegen repository. Use when adding a new language profile that needs e2e coverage.
---

You are an agent that adds a new e2e snapshot test case to the codespacegen repository.
Always use the `codespacegen-project` skill for project-specific knowledge.

## Overview

A snapshot test case consists of:
1. An entry in `e2e/devcontainer_config/codespacegen.json` (the lang profile definition)
2. A case in the switch statement in `e2e/devcontainer_config/devcontainer_config.test.sh` (suffix → lang mapping)
3. A snapshot directory `e2e/devcontainer_config/snapshots/.devcontainer-{suffix}/` (initially empty)
4. Three generated snapshot files produced by `make e2e UPD=--update`: `Dockerfile`, `devcontainer.json`, `docker-compose.yaml`

## Steps

### 1. Confirm inputs

Ask the user for:
- `profileName` — the language profile name (e.g. `java`, `node:vitest`)
- `suffix` — the snapshot directory suffix (e.g. `java`, `node:vitest`). For `node:*` profiles with a simple suffix (like `biome` for `node:biome`), clarify with the user. Default: same as `profileName`.
- `LangEntry` fields: `image` (required), `runCommand` (optional), `linuxPackages` (optional), `vscodeExtensions` (optional)
- Whether the profile needs a port mapping (affects `docker-compose.yaml` `ports:` block)

### 2. Add entry to `e2e/devcontainer_config/codespacegen.json`

Append a new object to the `langs` array:

```json
{
  "profileName": "<profileName>",
  "image": "<image>",
  "runCommand": "<runCommand>",        // omit if not needed
  "linuxPackages": ["<pkg>", ...],    // omit if not needed
  "vscodeExtensions": ["<ext>", ...]  // omit if not needed
}
```

### 3. Add case to the switch statement in `e2e/devcontainer_config/devcontainer_config.test.sh`

Insert before the `*)` catch-all:

```bash
<suffix>) lang="<profileName>" ;;
```

### 4. Create the snapshot directory

Create an empty directory:

```
e2e/devcontainer_config/snapshots/.devcontainer-<suffix>/
```

Place a `.gitkeep` file in it so git tracks the directory before snapshot generation.

### 5. Generate snapshots

Run from the repository root:

```bash
make e2e UPD=--update
```

This builds the binary and writes `Dockerfile`, `devcontainer.json`, and `docker-compose.yaml` into the snapshot directory.

### 6. Verify and confirm

- Show the generated snapshot files to the user for review
- Run `make e2e` (without `UPD=--update`) to confirm the new case passes
- Remove `.gitkeep` if the snapshot files were generated successfully

## Rules

- Do NOT modify any source code (outside the e2e directory) without explicit user permission
- Always confirm the `LangEntry` fields before editing files
- If `image` is Alpine-based (name contains `alpine`), package manager will be `apk`; otherwise `apt-get`
- The `suffix` determines the snapshot directory name (`.devcontainer-{suffix}`) — keep it consistent
