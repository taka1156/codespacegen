#!/usr/bin/env bash
set -euo pipefail

# ============================================================
# setup.sh — lefthook セットアップ・gh 認証確認
# コンテナ起動後 (postCreateCommand) に実行する
# ============================================================

# ---------- lefthook セットアップ ----------
echo "[INFO]  lefthook install を実行します..."
lefthook install

# ---------- gh 認証確認 ----------
if ! gh auth status &>/dev/null; then
  echo ""
  echo "========================================="
  echo " GitHub CLIの認証が必要です"
  echo "-----------------------------------------"
  echo " scripts/run-copilot-review.sh を使うには"
  echo " 以下のコマンドで認証してください:"
  echo ""
  echo "   gh auth login"
  echo ""
  echo " Copilot拡張も必要な場合:"
  echo "   gh extension install github/gh-copilot"
  echo "========================================="
  echo ""
fi
