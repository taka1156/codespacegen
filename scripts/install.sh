#!/usr/bin/env bash
set -euo pipefail

# ============================================================
# install.sh — jq / lefthook / github-cli インストールスクリプト
# 対応OS: Alpine Linux / Debian・Ubuntu 系
# ============================================================

# ---------- ヘルパー ----------
info()  { echo "[INFO]  $*"; }
error() { echo "[ERROR] $*" >&2; exit 1; }

has_cmd() { command -v "$1" >/dev/null 2>&1; }

# ---------- OS 判定 ----------
detect_os() {
  if [ -f /etc/alpine-release ]; then
    echo "alpine"
  elif [ -f /etc/debian_version ]; then
    echo "debian"
  else
    error "未対応のOSです。Alpine または Debian/Ubuntu 系が必要です。"
  fi
}

OS=$(detect_os)
info "検出されたOS: $OS"

# ---------- lefthook リポジトリ追加 (Alpine のみ) ----------
add_lefthook_repo_alpine() {
  if ! apk info --installed lefthook >/dev/null 2>&1; then
    info "lefthook リポジトリを追加します..."
    curl -1sLf 'https://dl.cloudsmith.io/public/evilmartians/lefthook/setup.alpine.sh' | sh
  fi
}

# ---------- インストール ----------
install_alpine() {
  add_lefthook_repo_alpine
  info "apk でパッケージをインストールします..."
  apk add --no-cache \
    jq \
    lefthook \
    github-cli
}

install_debian() {
  info "apt でパッケージをインストールします..."

  # jq
  if ! has_cmd jq; then
    apt-get update -qq
    apt-get install -y --no-install-recommends jq
  fi

  # lefthook
  if ! has_cmd lefthook; then
    info "lefthook をインストールします..."
    curl -1sLf 'https://dl.cloudsmith.io/public/evilmartians/lefthook/setup.deb.sh' | bash
    apt-get install -y --no-install-recommends lefthook
  fi

  # github-cli
  if ! has_cmd gh; then
    info "github-cli をインストールします..."
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg \
      | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
    chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] \
https://cli.github.com/packages stable main" \
      > /etc/apt/sources.list.d/github-cli.list
    apt-get update -qq
    apt-get install -y --no-install-recommends gh
  fi
}

# ---------- メイン ----------
case "$OS" in
  alpine) install_alpine ;;
  debian) install_debian ;;
esac

# ---------- バージョン確認 ----------
info "インストール完了。バージョン確認:"
jq       --version
lefthook version
gh       --version
