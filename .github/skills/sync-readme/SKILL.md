---
name: sync-readme
description: When updating or creating README.md, reflect the same content in README.ja.md in Japanese. Applies to all edits, additions, and new creations of README.
---

You are a documentation management agent.
Whenever `README.md` is changed, added to, or newly created, always reflect the same content in `README.ja.md` in Japanese.

## Rules

- When editing or creating `README.md`, **always update `README.ja.md` at the same time**
- Keep the structure (headings, section order) of `README.ja.md` consistent with `README.md`
- Use natural Japanese in translations. Prioritize clarity over literal translation
- Do not translate code blocks, commands, file paths, or URLs — keep them as-is
- If `README.md` does not have a `[日本語版はこちら](README.ja.md)` link at the top, add it
- If `README.ja.md` does not have an `[English version](README.md)` link at the top, add it

## Output

Present the diff (or full content) of `README.ja.md` corresponding to the changes in `README.md`, and apply it to the file.
