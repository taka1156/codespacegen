# codespacegen

Codespace向けの以下3ファイルを生成するCLIです。

- Dockerfile
- devcontainer.json
- docker-compose.yaml

## アーキテクチャ

オニオンアーキテクチャで実装しています。

- Domain: ルールとモデル
	- internal/domain/entity
	- internal/domain/service
- Application: ユースケース
	- internal/application/usecase
	- internal/application/port
- Infrastructure: 外部I/O実装
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
| `-base-image` | *(言語デフォルト)* | Dockerベースイメージを直接指定。`-language` のデフォルトより優先 |
| `-image-config` | `codespacegen.base-images.json` | ベースイメージ定義のローカルパスまたは `https://` URL。`install` のみ指定してイメージを省略した場合は `alpine:latest` を自動採用 |
| `-port` | *(対話入力、Enter で ports なし)* | ポート指定（例: `3000` → `3000:3000` に自動正規化、`8080:3000` も可）。毎回確認される |
| `-compose-file` | `docker-compose.yaml` | Composeファイル名 |
| `-force` | `false` | 既存ファイルを上書き |

言語ごとのデフォルトベースイメージ:

- go: golang:1.24-alpine
- python: python:3.12-alpine
- node: node:22-alpine
- rust: rust:1-alpine

ベースイメージ定義はルートの [codespacegen.base-images.json](codespacegen.base-images.json) に切り出しています。

- JSONが存在する場合: JSONの値を読み込み（同一キーは上書き）
- JSONが存在しない場合: CLI内部のデフォルト値で動作
- -base-image を指定した場合: JSON/デフォルトより優先

### codespacegen.base-images.json の書き方

**形式1: 文字列（イメージ名を直指定）**

```json
{
  "go": "golang:1.24-alpine"
}
```

**形式2: オブジェクト（curl でインストール、`image` 省略時は `alpine:latest` を自動採用）**

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

**形式2 応用: イメージも明示したい場合**

```json
{
  "moonbit": {
    "image": "alpine:3.20",
    "install": "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash"
  }
}
```

形式1・2 は同一ファイルに混在できます。

例:

```bash
go run ./cmd/codespacegen \
	-output .devcontainer \
	-name "My Codespace" \
	-language go \
	-service app \
	-workspace-folder /workspace \
	-compose-file docker-compose.yaml \
	-force
```

リモートJSONをURL指定する例:

```bash
go run ./cmd/codespacegen -image-config https://example.com/my-base-images.json -language go -force
```

- `https://` URLのみ対応（`http://` は拒否）
- JSONが存在しない・未指定の場合は内蔵のAlpineデフォルトで動作

明示イメージで上書きする例:

```bash
go run ./cmd/codespacegen -language python -base-image python:3.12-alpine -force
```

ポートを明示する例:

```bash
go run ./cmd/codespacegen -language go -port 3000 -force
```

`-port` 未指定時は実行中に対話形式でポート入力を促します。

生成される `docker-compose.yaml` は以下の構成です（ポート入力時のみ `ports` を追加）。

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
