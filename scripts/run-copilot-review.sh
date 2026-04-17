#!/usr/bin/env bash
set -euo pipefail

SKILL_FILE=".github/skills/git-hook-code-review/SKILL.md"
REVIEW_SCRIPT="./scripts/review-check.sh"
LANG_OPT="ja"

if [ ! -f "$SKILL_FILE" ]; then
  echo "エラー: スキルファイルが見つかりません: ${SKILL_FILE}" >&2
  exit 1
fi

if [ ! -f "$REVIEW_SCRIPT" ]; then
  echo "エラー: レビュースクリプトが見つかりません: ${REVIEW_SCRIPT}" >&2
  exit 1
fi

# オプション解析
while [[ $# -gt 0 ]]; do
  case "$1" in
    --lang=en) LANG_OPT="en"; shift ;;
    --lang=ja) LANG_OPT="ja"; shift ;;
    *) shift ;;
  esac
done

if [ "$LANG_OPT" = "ja" ]; then
  LANG_INSTRUCTION="title・detail・recommendationはすべて日本語で回答してください。"
else
  LANG_INSTRUCTION="Respond in English for title, detail, and recommendation."
fi

# pushされる差分を取得
# 優先順: origin追跡ブランチ → HEAD~1 → 空ツリー（初回push時）
BRANCH=$(git rev-parse --abbrev-ref HEAD)
EMPTY_TREE="4b825dc642cb6eb9a060e54bf8d69288fbee4904"
DIFF=$(git diff "origin/${BRANCH}" HEAD 2>/dev/null \
  || git diff HEAD~1 HEAD 2>/dev/null \
  || git diff "${EMPTY_TREE}" HEAD 2>/dev/null \
  || true)

if [ -z "$DIFF" ]; then
  echo "差分なし。Copilotレビューをスキップします。"
  exit 0
fi

# スキルのシステムプロンプトを読み込む（YAMLフロントマターを除く）
SYSTEM_PROMPT=$(awk 'BEGIN{p=0} /^---/{p++; next} p>=2{print}' "$SKILL_FILE")

# GitHub Copilot APIトークンを取得
GH_TOKEN=$(gh auth token 2>/dev/null || true)
if [ -z "$GH_TOKEN" ]; then
  echo "エラー: GitHubトークンの取得に失敗しました。'gh auth login' で認証してください。" >&2
  exit 1
fi

COPILOT_TOKEN=$(curl -sf \
  -H "Authorization: token ${GH_TOKEN}" \
  -H "Accept: application/json" \
  "https://api.github.com/copilot_internal/v2/token" | jq -r '.token')

if [ -z "$COPILOT_TOKEN" ] || [ "$COPILOT_TOKEN" = "null" ]; then
  echo "エラー: Copilot APIトークンの取得に失敗しました。GitHub Copilotのサブスクリプションと権限を確認してください。" >&2
  exit 1
fi

# Copilot Chat APIを呼び出し
REQUEST=$(jq -nc \
  --arg system "$SYSTEM_PROMPT" \
  --arg diff "$DIFF" \
  --arg lang_instruction "$LANG_INSTRUCTION" \
  '{
    model: "gpt-4o",
    stream: false,
    messages: [
      {role: "system", content: $system},
      {role: "user", content: ("以下の差分をレビューしてください。" + $lang_instruction + "\n```diff\n" + $diff + "\n```")}
    ]
  }')

RESULT=$(curl -sf \
  -X POST "https://api.githubcopilot.com/chat/completions" \
  -H "Authorization: Bearer ${COPILOT_TOKEN}" \
  -H "Content-Type: application/json" \
  -H "Copilot-Integration-Id: vscode-chat" \
  -H "Editor-Version: vscode/1.85.0" \
  -d "$REQUEST" | jq -r '.choices[0].message.content')

if [ -z "$RESULT" ] || [ "$RESULT" = "null" ]; then
  echo "エラー: Copilotのレビュー結果の取得に失敗しました。APIの応答を確認してください。" >&2
  exit 1
fi

# コードブロック記法が含まれる場合は除去してJSONのみ抽出
RESULT=$(echo "$RESULT" | sed '/^```/d')

echo "=== Copilotレビュー結果 ==="
echo "$RESULT"
echo "=========================="

export RESULT
set +e
bash "$REVIEW_SCRIPT"
EXIT_CODE=$?
set -e

# 人が読めるマークダウンをcopilot-review/review.mdに生成
REVIEW_MD="copilot-review/review.md"
mkdir -p copilot-review

PRIORITY_ORDER="must want imo imho nits info ask"
EMOJI_must="🔴"; EMOJI_want="🟠"; EMOJI_imo="🟡"; EMOJI_imho="🟡"; EMOJI_nits="🔵"; EMOJI_info="💬"; EMOJI_ask="❓"

{
  echo "# Copilotコードレビュー"
  echo ""

  ISSUE_COUNT=$(echo "$RESULT" | jq '.issues | length')
  if [ "$ISSUE_COUNT" -eq 0 ]; then
    echo "問題は見つかりませんでした。"
  else
    for PRIORITY in $PRIORITY_ORDER; do
      EMOJI_VAR="EMOJI_${PRIORITY}"
      EMOJI="${!EMOJI_VAR}"
      echo "$RESULT" | jq -r --arg p "$PRIORITY" --arg e "$EMOJI" '
        .issues[] | select(.priority == $p) |
        "\($e) **\(.priority) — \(.title)**（行: \(.line)）\n\n\(.detail)\n\n> 修正案: \(.recommendation)\n\n---\n"
      '
    done
  fi
} > "$REVIEW_MD"

echo "レビュー結果を $REVIEW_MD に出力しました。"
exit $EXIT_CODE
