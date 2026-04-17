# codespacegen

[English version](README.md)

Codespace 向けに以下 3 ファイルを生成する CLI です。

- Dockerfile
- devcontainer.json
- docker-compose.yaml

## アーキテクチャ

- Domain: ルールとモデル
  - internal/domain/entity
  - internal/domain/service
- Application: ユースケース
  - internal/application/usecase
  - internal/application/port
- Infrastructure: 外部 I/O 実装
  - internal/infrastructure/generator
  - internal/infrastructure/persistence
- Entry Point: CLI
  - cmd/codespacegen

依存方向は外側から内側のみです。

## 使い方

### 実行

```bash
go run ./cmd/codespacegen
```

デフォルトでは .devcontainer 配下にファイルを生成します。

### 主なオプション

| オプション | 既定値 | 説明 |
|---|---|---|
| `-output` | `.devcontainer` | 出力先ディレクトリ |
| `-name` | *(対話入力必須)* | プロジェクト名。毎回確認され、`devcontainer.json` の `name` に反映 |
| `-language` | *(対話入力、Enter で空)* | プログラミング言語。毎回確認される。空の場合は言語固有設定を使わず `alpine:latest` を採用 |
| `-service` | *(対話入力、Enter で `app`)* | docker compose のサービス名。毎回確認され、`devcontainer.json` と `docker-compose.yaml` 両方に反映 |
| `-workspace-folder` | *(対話入力、Enter で `/workspace`)* | コンテナ内ワークスペースパス。毎回確認される |
| `-timezone` | `Asia/Tokyo` | コンテナ内のタイムゾーン。Dockerfile の `ENV TZ` と timezone 設定に反映 |
| `-base-image` | *(言語デフォルト)* | Dockerベースイメージを直接指定。`-language` のデフォルトより優先 |
| `-image-config` | `codespacegen.base-images.json` | ベースイメージ定義のローカルパスまたは `https://` URL。`install` のみ指定してイメージを省略した場合は `alpine:latest` を自動採用 |
| `-port` | *(対話入力、Enter で ports なし)* | ポート指定。たとえば `3000` は `3000:3000` に自動正規化され、`8080:3000` も利用できます。毎回確認されます |
| `-compose-file` | `docker-compose.yaml` | Composeファイル名 |
| `-force` | `false` | 既存ファイルを上書き |

言語ごとのデフォルトベースイメージ:

- go: golang:1.24-alpine
- python: python:3.12-alpine
- node: node:22-alpine
- rust: rust:1-alpine

ベースイメージ定義はリポジトリルートの [codespacegen.base-images.json](codespacegen.base-images.json) に分離しています。

- JSON が存在する場合: ファイルの値を読み込み、同名キーでデフォルトを上書き
- JSON が存在しない場合: CLI 内蔵のデフォルト値で動作
- `-base-image` を指定した場合: JSON と内蔵デフォルトの両方より優先

### codespacegen.base-images.json の書き方

**形式 1: 文字列でイメージ名を直接指定**

```json
{
  "go": "golang:1.24-alpine"
}
```

**形式 2: オブジェクトで install コマンドを指定し、`image` 省略時は `alpine:latest` を自動採用**

```json
{
  "moonbit": {
    "install": "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash"
  }
}
```

生成される Dockerfile には以下の `RUN` ステップが追加されます。

```dockerfile
RUN curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash
```

**形式 2 の派生: イメージも明示する場合**

```json
{
  "moonbit": {
    "image": "ubuntu:24.04",
    "install": "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash"
  }
}
```

形式 1 と形式 2 は同一ファイルに混在できます。

例:

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

リモート JSON を URL 指定する例:

```bash
go run ./cmd/codespacegen -image-config https://example.com/my-base-images.json -language go -force
```

- `https://` URL のみ対応し、`http://` は拒否します
- JSON が存在しない、または未指定の場合は内蔵の Alpine デフォルトで動作します

明示イメージで上書きする例:

```bash
go run ./cmd/codespacegen -language python -base-image python:3.12-alpine -force
```

ポートを明示する例:

```bash
go run ./cmd/codespacegen -language go -port 3000 -force
```

`-port` 未指定時は、実行中に対話形式でポート入力を促します。

生成される `docker-compose.yaml` は以下の構成です。ポートを指定した場合のみ `ports` が追加されます。

```yaml
services:
    app:
      build: .
      tty: true
      volumes:
        - ../:/workspace
```

## テスト

```bash
go test ./...
```

## リリース（GitHub Actions）

タグを push すると、GitHub Actions がクロスビルドを実行し、成果物を GitHub Releases に添付します。

```bash
git tag v0.1.0
git push origin v0.1.0---
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
```

生成される主なアセット:

- `codespacegen_linux_amd64.tar.gz`
- `codespacegen_linux_arm64.tar.gz`
- `codespacegen_darwin_amd64.tar.gz`
- `codespacegen_darwin_arm64.tar.gz`
- `codespacegen_windows_amd64.exe`
- `checksums.txt`

## インストール（curl）

最新リリースを自動でダウンロードし、`/usr/local/bin` に配置します。

```bash
curl -fsSL https://raw.githubusercontent.com/taka1156/codespacegen/master/install.sh | bash
```

インストール先を変更する場合:

```bash
curl -fsSL https://raw.githubusercontent.com/taka1156/codespacegen/master/install.sh | INSTALL_DIR=$HOME/.local/bin bash
```
