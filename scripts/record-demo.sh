#!/bin/bash
# record-demo.sh — asciinema 데모 레코딩 자동화
# Usage: ./scripts/record-demo.sh
#
# asciinema가 설치되어 있어야 합니다: brew install asciinema
set -euo pipefail

CAST_FILE="docs/demo.cast"
DEMO_SCRIPT="$(mktemp)"

# 데모 시나리오 (타이핑 시뮬레이션)
cat > "$DEMO_SCRIPT" << 'DEMO'
#!/bin/bash
# 타이핑 효과 함수
type_cmd() {
    echo ""
    echo -n "$ "
    for ((i=0; i<${#1}; i++)); do
        echo -n "${1:$i:1}"
        sleep 0.04
    done
    echo ""
    sleep 0.3
    eval "$1"
    sleep 1.5
}

comment() {
    echo ""
    echo -e "\033[90m# $1\033[0m"
    sleep 1
}

clear

echo -e "\033[1;36m  upbit — AI-native CLI for Upbit Exchange\033[0m"
echo ""
sleep 1

comment "Check current BTC price"
type_cmd "upbit ticker KRW-BTC"

comment "Multiple markets at once"
type_cmd "upbit ticker KRW-BTC KRW-ETH KRW-XRP"

comment "View order book"
type_cmd "upbit orderbook KRW-BTC"

comment "Recent trades"
type_cmd "upbit trades KRW-BTC -c 5"

comment "Daily candles"
type_cmd "upbit candle KRW-BTC -i 1d -c 5"

comment "Candle with --from (auto-pagination + SQLite cache)"
type_cmd "upbit candle KRW-BTC -i 1d --from 2026-03-01 -f"

comment "List KRW markets"
type_cmd "upbit market -q KRW | head -10"

comment "Tick size for price correction"
type_cmd "upbit tick-size KRW-BTC KRW-XRP KRW-DOGE"

comment "JSON output for AI agents (non-TTY auto-detection)"
type_cmd "upbit ticker KRW-BTC -o json | head -5"

comment "Select specific fields"
type_cmd "upbit ticker KRW-BTC --json market,trade_price,signed_change_rate"

comment "Export tool schema for LLM/MCP integration"
type_cmd "upbit tool-schema buy"

comment "Cache management"
type_cmd "upbit cache"

comment "English locale support"
type_cmd "LANG=en_US.UTF-8 upbit ticker KRW-BTC"

echo ""
echo -e "\033[1;32m✓ Done! Visit https://github.com/kyungw00k/upbit\033[0m"
sleep 2
DEMO

chmod +x "$DEMO_SCRIPT"

echo "Recording demo to $CAST_FILE ..."
mkdir -p docs

# asciinema로 레코딩
asciinema rec "$CAST_FILE" \
    --cols 100 \
    --rows 30 \
    --title "upbit — AI-native CLI for Upbit Exchange" \
    --command "$DEMO_SCRIPT" \
    --overwrite

rm -f "$DEMO_SCRIPT"

echo ""
echo "Done! Cast file: $CAST_FILE"
echo ""
echo "To upload: asciinema upload $CAST_FILE"
echo "To play locally: asciinema play $CAST_FILE"
echo ""
echo "For README embed, add:"
echo "[![demo](https://asciinema.org/a/XXXX.svg)](https://asciinema.org/a/XXXX)"
