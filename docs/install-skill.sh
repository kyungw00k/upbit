#!/bin/sh
# upbit AI agent skill installer
# Usage: curl -sSL https://kyungw00k.dev/upbit/install-skill.sh | sh
set -e

SKILL_URL="https://kyungw00k.dev/upbit/skill.md"

echo "Installing upbit skill for AI agents..."

# Claude Code
if command -v claude >/dev/null 2>&1; then
  echo ""
  echo "[Claude Code]"
  mkdir -p .claude/skills/upbit
  curl -sSfL "$SKILL_URL" -o .claude/skills/upbit/skill.md
  echo "  Installed to .claude/skills/upbit/skill.md"
fi

# agentskills.io standard
echo ""
echo "[agentskills.io]"
mkdir -p "$HOME/.agents/skills/upbit"
curl -sSfL "$SKILL_URL" -o "$HOME/.agents/skills/upbit/SKILL.md"
echo "  Installed to ~/.agents/skills/upbit/SKILL.md"

echo ""
echo "Done! The upbit skill is now available to AI agents."
echo "  CLI install: curl -sSL https://kyungw00k.dev/upbit/install.sh | sh"
