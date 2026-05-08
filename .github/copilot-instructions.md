# Copilot Instructions

## Response Style
- Always state the **basis (root cause / reasoning)** and **conclusion** explicitly.
- Do not present a conclusion without supporting reasoning.
- Structure responses as: **Conclusion → Basis → Options (if any) → Next action**

## Before Starting Work
- If the intent, target files, or impact scope is unclear, **ask before proceeding**.
- Share the **plan and approach** before making any changes.

## When Proposing Changes
- Always include **reason**, **benefit**, and **trade-off** for each proposed change.
- If alternative approaches exist, list them alongside the recommendation.
- Keep changes as **minimal diff** as possible. Avoid rewriting beyond what is necessary.

## When to Stop and Return Control
Return control to the user (ask for confirmation) before:
- Deleting files or directories
- Changing directory structure
- Breaking existing responsibility boundaries (e.g., moving logic between layers)
- Changing API specifications
- Receiving an ambiguous instruction that could be interpreted in multiple ways
- Starting significant changes (modifying multiple files, changing behavior, refactoring) — warn the user to **stage current changes in git** first

## Authoring Guidelines
Custom agents (`.github/agents/`), skills (`.github/skills/`), and this file must be written in **English** for token efficiency.
