
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

[English version](README.md)

codespacegen は、Codespaces および Dev Container 向けに以下の 3 つのファイルを生成するコマンドラインツールです。

- Dockerfile
- devcontainer.json
- docker-compose.yaml


## インストール（curl）

最新リリースを自動でダウンロードし、`/usr/local/bin` に配置します。

```bash
curl -fsSL https://raw.githubusercontent.com/taka1156/codespacegen/master/scripts/install.sh | bash
```

インストール先を変更する場合:

```bash
curl -fsSL https://raw.githubusercontent.com/taka1156/codespacegen/master/scripts/install.sh | INSTALL_DIR=$HOME/.local/bin bash
```

## リリース（GitHub Actions）

生成される主なアセット:

- `codespacegen_linux_amd64.tar.gz`
- `codespacegen_linux_arm64.tar.gz`
- `codespacegen_darwin_amd64.tar.gz`
- `codespacegen_darwin_arm64.tar.gz`
- `codespacegen_windows_amd64.exe`
- `checksums.txt`

## アーキテクチャ

- Domain: ルールとモデル
  - internal/domain/entity
  - internal/domain/service
- App: 構成とオーケストレーション
  - internal/app
- Input adapters: CLI/JSON/デフォルト値の入力
  - internal/input
- Infra: 対話入力（標準入力プロンプト）
  - internal/infra
- Workflow: ユースケース
  - internal/workflow/collect
  - internal/workflow/assemble
  - internal/workflow/generate
  - internal/workflow/initialize
- Generator: テンプレート生成とファイル書き込み
  - internal/generator
  - internal/generator/filewriter
  - internal/generator/workdirprovider
- i18n: ローカライズリソース
  - internal/i18n
- Entry Point: CLI
  - cmd/codespacegen

依存方向は外側から内側のみです。

## 使い方

### 実行

```bash
go run ./cmd/codespacegen
```

デフォルトでは .devcontainer 配下にファイルを生成します。

### codespacegen.json の初期化

`init` サブコマンドを実行すると、カレントディレクトリに `codespacegen.json` のテンプレートを生成します。

```bash
codespacegen init
```

生成されたファイルをベースイメージや VS Code 拡張機能のカスタマイズの出発点として使用できます。

### 主なオプション

| オプション | 既定値 | 説明 |
|---|---|---|
| `-output` | `.devcontainer` | 出力先ディレクトリ |
| `-name` | *(対話入力必須)* | プロジェクト名。毎回確認され、`devcontainer.json` の `name` に反映 |
| `-language` | *(対話入力、Enter で空)* | 言語キー。毎回確認される。`codespacegen.json`（または `-image-config` で指定したファイル）の `langs` 配列内の `profileName` に定義された任意のキーを利用できる。空の場合は言語固有設定を使わず `alpine:latest` を採用 |
| `-service` | *(対話入力、Enter で `app`)* | docker compose のサービス名。毎回確認され、`devcontainer.json` と `docker-compose.yaml` 両方に反映 |
| `-workspace-folder` | *(対話入力、Enter で `/workspace`)* | コンテナ内ワークスペースパス。毎回確認される |
| `-timezone` | *(対話入力、`common.timezone` または `UTC` が既定)* | コンテナ内のタイムゾーン。毎回確認され、Dockerfile の `ENV TZ` と timezone 設定に反映 |
| `-image-config` | `codespacegen.json` | ベースイメージ定義のローカルパスまたは `https://` URL。トップレベル `common` と、言語別設定を持つ `langs` 配列をサポート。`runCommand` または `linuxPackages` を指定する場合は `image` が必須 |
| `-port` | *(対話入力、Enter で ports なし)* | ポート指定。たとえば `3000` は `3000:3000` に自動正規化され、`8080:3000` も利用できます。毎回確認されます |
| `-compose-file` | `docker-compose.yaml` | Composeファイル名 |
| `-force` | `false` | 既存ファイルを上書き |
| `-lang` | *(自動検出)* | CLI メッセージの言語 (`en` or `ja`)。未指定の場合はシステムロケールを使用 |
| `-headless` | `false` | 対話入力を一切スキップします。必要な値はすべてフラグで指定してください |
| `-v` | — | バージョンを出力して終了 |

ベースイメージ定義はリポジトリルートの [codespacegen.json](codespacegen.json) に分離しています。

- JSON が存在する場合: ファイルの値を読み込んで使用


加えて、`codespacegen.json` の `common.vscodeExtensions` と言語別 `vscodeExtensions` が追記されます。

### codespacegen.json の書き方

JSON Schema の検証や補完に対応したエディタでは、次のように schema を関連付けできます。

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

`codespacegen.json` をリポジトリルートに置く場合、`./codespacegen.schema.json` で同梱の schema ファイルを参照できます。

**言語エントリ（`langs` 配列）**

`langs` の各エントリには `profileName` が必須で、以下のフィールドをサポートします（`runCommand` または `linuxPackages` を指定する場合は `image` が必須）:

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

`runCommand` の値は Dockerfile の `RUN` ステップとして挿入されます。

```dockerfile
RUN curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash
```

`linuxPackages` を使うと Linux システムパッケージをリストで指定できます。デフォルトのパッケージリストにマージされ、ベースイメージに応じたパッケージマネージャー（Debian/Ubuntu は `apt`、Alpine は `apk`）でインストールされます。

**`common` で共通設定を定義**

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

マージ挙動:

- `common` を基本設定として先に適用し、言語別エントリでは `vscodeExtensions` のみ追記できる
- `vscodeExtensions` は順序を保って結合し、重複は除去
- `timezone` と `locale` は `common` にのみ設定でき、言語別エントリには指定できない
- フラグ・`common` のどちらにも timezone がない場合は `UTC` を使用

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
