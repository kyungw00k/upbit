#!/bin/bash
# Upbit API 문서 동기화 스크립트
# llms.txt에서 URL 추출 → 변경/추가된 문서 다운로드
set -euo pipefail

DOCS_DIR="docs/upbit-api"
BASE="https://docs.upbit.com/kr"
LLMS_TXT="$DOCS_DIR/llms.txt"

# llms.txt에서 .md URL 추출
extract_urls() {
    grep -oE '\(https://docs\.upbit\.com/kr/[^)]+\.md\)' "$LLMS_TXT" \
        | tr -d '()' \
        | sort -u
}

# URL → 로컬 경로 매핑
url_to_path() {
    local url="$1"
    local rel="${url#$BASE/}"
    local dir
    case "$rel" in
        docs/*)       dir="guides" ;;
        reference/*)  dir="reference" ;;
        recipes/*)    dir="recipes" ;;
        changelog/*)  dir="changelog" ;;
        *)            dir="pages" ;;
    esac
    local filename
    filename=$(basename "$rel")
    echo "$DOCS_DIR/$dir/$filename"
}

# 메인 실행
main() {
    local added=0
    local updated=0
    local failed=0
    local changes_log="$DOCS_DIR/CHANGES.md"

    echo "# Upbit API 문서 변경 사항" > "$changes_log"
    echo "" >> "$changes_log"
    echo "동기화 시각: $(date -u '+%Y-%m-%d %H:%M:%S UTC')" >> "$changes_log"
    echo "" >> "$changes_log"

    # 디렉토리 확보
    mkdir -p "$DOCS_DIR"/{guides,reference,recipes,changelog}

    # 각 URL 다운로드
    while IFS= read -r url; do
        local path
        path=$(url_to_path "$url")
        local dir
        dir=$(dirname "$path")
        mkdir -p "$dir"

        # 임시 파일에 다운로드
        local tmp="${path}.tmp"
        local http_code
        http_code=$(curl -s -w "%{http_code}" -o "$tmp" "$url" --max-time 15 2>/dev/null) || true

        if [ "$http_code" != "200" ]; then
            rm -f "$tmp"
            failed=$((failed + 1))
            continue
        fi

        if [ ! -f "$path" ]; then
            # 새 파일
            mv "$tmp" "$path"
            added=$((added + 1))
            echo "- **추가**: \`$(basename "$path")\`" >> "$changes_log"
        elif ! diff -q "$path" "$tmp" > /dev/null 2>&1; then
            # 변경된 파일
            mv "$tmp" "$path"
            updated=$((updated + 1))
            echo "- **변경**: \`$(basename "$path")\`" >> "$changes_log"
        else
            # 변경 없음
            rm -f "$tmp"
        fi
    done < <(extract_urls)

    # 요약
    echo "" >> "$changes_log"
    echo "## 요약" >> "$changes_log"
    echo "- 추가: ${added}개" >> "$changes_log"
    echo "- 변경: ${updated}개" >> "$changes_log"
    echo "- 실패: ${failed}개" >> "$changes_log"

    echo "sync-docs: added=$added updated=$updated failed=$failed"

    # 변경 여부를 종료 코드로 전달 (GitHub Actions 용)
    if [ $((added + updated)) -gt 0 ]; then
        echo "HAS_CHANGES=true" >> "${GITHUB_OUTPUT:-/dev/null}"
        return 0
    else
        echo "HAS_CHANGES=false" >> "${GITHUB_OUTPUT:-/dev/null}"
        return 0
    fi
}

main "$@"
