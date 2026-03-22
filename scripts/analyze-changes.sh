#!/bin/bash
# Upbit API 문서 변경 사항을 LLM으로 분석하여 CLI 영향도 평가
set -euo pipefail

DOCS_DIR="docs/upbit-api"
CHANGES_LOG="$DOCS_DIR/CHANGES.md"
ANALYSIS_FILE="$DOCS_DIR/ANALYSIS.md"

# ZAI_API_KEY 확인
if [ -z "${ZAI_API_KEY:-}" ]; then
    echo "ZAI_API_KEY not set, skipping LLM analysis"
    # CHANGES.md를 그대로 ANALYSIS.md로 복사
    cp "$CHANGES_LOG" "$ANALYSIS_FILE"
    exit 0
fi

# 변경된 파일 목록 추출
changed_files=$(grep -E '^\- \*\*(추가|변경)\*\*' "$CHANGES_LOG" | sed 's/.*`\(.*\)`.*/\1/' || true)

if [ -z "$changed_files" ]; then
    echo "No changes to analyze"
    cp "$CHANGES_LOG" "$ANALYSIS_FILE"
    exit 0
fi

# 변경된 파일의 내용 수집 (최대 20개)
file_contents=""
count=0
for f in $changed_files; do
    filepath=$(find "$DOCS_DIR" -name "$f" -type f 2>/dev/null | head -1)
    if [ -n "$filepath" ] && [ $count -lt 20 ]; then
        content=$(head -100 "$filepath" 2>/dev/null || true)
        file_contents="${file_contents}

--- FILE: ${f} ---
${content}
--- END ---"
        count=$((count + 1))
    fi
done

# 현재 구현된 API 경로 목록
implemented_apis=$(grep -rh 'path.*"/' internal/api/quotation/*.go internal/api/exchange/*.go internal/api/wallet/*.go 2>/dev/null \
    | grep -oE '"/.+?"' | sort -u | tr '\n' ', ' || true)

# LLM 프롬프트 생성
prompt=$(cat <<PROMPT
You are an API change analyst for the Upbit CLI (Go).

The following Upbit API documentation files have been added or changed:

${file_contents}

Currently implemented API endpoints in the CLI:
${implemented_apis}

Please analyze and output in Markdown:

## Impact Analysis

### 1. New API Endpoints
List any new API endpoints found in the changed docs that are NOT in the implemented list.

### 2. Changed Parameters
List any parameter changes (added/removed/renamed) for existing endpoints.

### 3. Rate Limit Changes
Note any rate limit policy changes.

### 4. Breaking Changes
Flag anything that would break the current CLI implementation.

### 5. Deprecated Fields
List any newly deprecated fields or endpoints.

### 6. Recommended Actions
For each finding, suggest:
- [ ] Action item (e.g., "Add new endpoint GET /v1/xxx to exchange/xxx.go")
- Priority: HIGH / MEDIUM / LOW

Be concise. Only report actual changes, not unchanged items.
PROMPT
)

# Z.AI API 호출
response=$(curl -s -X POST "https://api.z.ai/v1/chat/completions" \
    -H "Authorization: Bearer ${ZAI_API_KEY}" \
    -H "Content-Type: application/json" \
    -d "$(jq -n \
        --arg prompt "$prompt" \
        '{
            "model": "claude-sonnet-4-20250514",
            "max_tokens": 4096,
            "messages": [{"role": "user", "content": $prompt}]
        }'
    )" 2>/dev/null || true)

# 응답 파싱
analysis=$(echo "$response" | jq -r '.choices[0].message.content // empty' 2>/dev/null || true)

if [ -z "$analysis" ]; then
    echo "LLM analysis failed, using changes log only"
    cp "$CHANGES_LOG" "$ANALYSIS_FILE"
    exit 0
fi

# ANALYSIS.md 생성
cat > "$ANALYSIS_FILE" <<EOF
# Upbit API Documentation Change Analysis

$(cat "$CHANGES_LOG")

---

${analysis}
EOF

echo "Analysis complete: $ANALYSIS_FILE"
