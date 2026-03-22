#!/usr/bin/env python3
"""
gen-demo-cast.py — demo.cast / demo.ko.cast 생성 스크립트
asciinema v2 형식으로 직접 cast 파일을 작성합니다.
"""
import json
import sys

CHAR_DELAY = 0.04   # 타이핑 간격
PAUSE_AFTER_TYPE = 0.3
PAUSE_AFTER_OUTPUT = 1.5
COMMENT_PAUSE = 1.0
CMD_PROMPT_PAUSE = 0.8  # 프롬프트 전 공백

WIDTH = 110
HEIGHT = 35


def make_cast(title, events):
    header = {
        "version": 2,
        "width": WIDTH,
        "height": HEIGHT,
        "timestamp": 1774200000,
        "title": title,
        "env": {"SHELL": "/bin/zsh", "TERM": "xterm-256color"},
    }
    lines = [json.dumps(header, ensure_ascii=False)]
    for ts, typ, data in events:
        lines.append(json.dumps([round(ts, 6), typ, data], ensure_ascii=False))
    return "\n".join(lines) + "\n"


def build_events(scenario):
    """
    scenario: list of ("comment", text) | ("cmd", cmd_str, output_str)
    returns list of [timestamp, "o", data]
    """
    events = []
    t = 0.0

    def emit(data, delay=0.0):
        nonlocal t
        events.append([t, "o", data])
        t += delay

    # 헤더 타이틀
    emit(scenario["title"])
    emit("\r\n")
    emit("\r\n")
    t += 1.5

    for item in scenario["items"]:
        kind = item[0]

        if kind == "comment":
            text = item[1]
            emit("\r\n")
            emit(f"\033[90m# {text}\033[0m")
            emit("\r\n")
            t += COMMENT_PAUSE

            cmd_str = item[2]
            emit(f"\033[1;32m$\033[0m ")
            for ch in cmd_str:
                emit(ch, CHAR_DELAY)
            emit("\r\n")
            t += PAUSE_AFTER_TYPE

            output = item[3]
            if output:
                emit(output)
                emit("\r\n")
            t += PAUSE_AFTER_OUTPUT

        elif kind == "comment_only":
            # 실행 없이 comment + typing만 (TUI 등)
            text = item[1]
            emit("\r\n")
            emit(f"\033[90m# {text}\033[0m")
            emit("\r\n")
            t += COMMENT_PAUSE

            cmd_str = item[2]
            emit(f"\033[1;32m$\033[0m ")
            for ch in cmd_str:
                emit(ch, CHAR_DELAY)
            emit("\r\n")
            t += PAUSE_AFTER_OUTPUT

    emit("\r\n")
    emit(f"\033[1;32m\u2713 https://github.com/kyungw00k/upbit\033[0m")
    emit("\r\n")

    return events


# ─────────────────────────────────────────────
# EN scenario
# ─────────────────────────────────────────────
EN_TICKER_1 = (
    "Market   Price        Change    Change%  Volume(24h)     Trade Amt(24h)         High         Low        \r\n"
    "KRW-BTC  102,965,000  -817,000  -0.79%   1,314.48474815  136,563,852,318.50903  104,488,000  102,274,000"
)
EN_TICKER_MULTI = (
    "Market   Price        Change    Change%  Volume(24h)          Trade Amt(24h)         High         Low        \r\n"
    "KRW-BTC  102,965,000  -817,000  -0.79%   1,314.48474815       136,563,852,318.50903  104,488,000  102,274,000\r\n"
    "KRW-ETH  3,125,000    -17,000   -0.54%   39,157.55640281      123,358,523,747.61168  3,191,000    3,096,000  \r\n"
    "KRW-XRP  2,100        -22       -1.04%   83,646,151.04691099  177,040,864,454.3182   2,140        2,075      "
)
EN_TRADES = (
    "Trade Time  Trade Price  Trade Vol   Ask/Bid  Change  \r\n"
    "21:57:02    102,965,000  0.00011611  BID      -817,000\r\n"
    "21:57:01    102,965,000  0.00009712  BID      -817,000\r\n"
    "21:57:00    102,961,000  0.00476482  ASK      -821,000\r\n"
    "21:57:00    102,965,000  0.00004856  BID      -817,000\r\n"
    "21:57:00    102,965,000  0.0001214   BID      -817,000"
)
EN_CANDLES = (
    "Time                 Open         High         Low          Close        Volume        \r\n"
    "2026-03-18 18:00:00  109,015,000  110,022,000  104,935,000  106,000,000  2,441.75375374\r\n"
    "2026-03-19 18:00:00  106,000,000  106,392,000  102,580,000  104,078,000  2,022.57056593\r\n"
    "2026-03-20 18:00:00  104,078,000  106,200,000  103,865,000  105,314,000  1,361.28731847\r\n"
    "2026-03-21 18:00:00  105,399,000  106,446,000  103,452,000  103,782,000  726.2663963   \r\n"
    "2026-03-22 18:00:00  103,831,000  104,488,000  102,274,000  102,965,000  935.95281897  "
)
EN_MARKET = (
    "Market       Korean Name                 English Name               \r\n"
    "KRW-WAXP     \uc641\uc2a4                        WAX                        \r\n"
    "KRW-CARV     \uce74\ube0c                        CARV                       \r\n"
    "KRW-LSK      \ub9ac\uc2a4\ud06c                      Lisk                       \r\n"
    "KRW-0G       \uc81c\ub85c\uc9c0                      0G                         \r\n"
    "KRW-BORA     \ubcf4\ub77c                        BORA                       \r\n"
    "KRW-PUNDIX   \ud380\ub514\uc5d1\uc2a4                    Pundi X                    \r\n"
    "KRW-USD1     \uc6d4\ub4dc\ub9ac\ubc84\ud2f0\ud30c\uc774\ub0b8\uc15c\uc720\uc5d0\uc2a4\ub514  World Liberty Financial USD\r\n"
    "KRW-BAT      \ubca0\uc774\uc9c1\uc5b4\ud150\uc158\ud1a0\ud070            Basic Attention Token      "
)
EN_TICKSIZE = (
    "Market    Quote Currency  Tick Size\r\n"
    "KRW-BTC   KRW             1000     \r\n"
    "KRW-XRP   KRW             1        \r\n"
    "KRW-DOGE  KRW             1        "
)
EN_JSON = '[{"market":"KRW-BTC","signed_change_rate":-0.0079108131,"trade_price":102961000}]'
EN_CACHE = (
    "Path   /Users/humphrey.park/.upbit/cache\r\n"
    "Files  2\r\n"
    "Size   27.6 KB"
)
EN_KO_LOCALE = (
    "\ub9c8\ucf13     \ud604\uc7ac\uac00       \uc804\uc77c \ub300\ube44  \ubcc0\ub3d9\ub960  \uac70\ub798\ub7c9(24h)     \uac70\ub798\ub300\uae08(24h)          \uace0\uac00         \uc800\uac00       \r\n"
    "KRW-BTC  102,965,000  -817,000   -0.79%  1,314.48474815  136,563,852,318.50903  104,488,000  102,274,000"
)

en_scenario = {
    "title": "\033[1;36m  upbit \u2014 AI-native CLI for Upbit Exchange\033[0m\r\n",
    "items": [
        ("comment", "Check current BTC price",
         "upbit ticker KRW-BTC",
         EN_TICKER_1),
        ("comment", "Multiple markets at once",
         "upbit ticker KRW-BTC KRW-ETH KRW-XRP",
         EN_TICKER_MULTI),
        ("comment", "Recent trades",
         "upbit trades KRW-BTC -c 5",
         EN_TRADES),
        ("comment", "Daily candles",
         "upbit candle KRW-BTC -i 1d -c 5",
         EN_CANDLES),
        ("comment", "List KRW markets",
         "upbit market -q KRW | head -9",
         EN_MARKET),
        ("comment", "Tick size for price correction",
         "upbit tick-size KRW-BTC KRW-XRP KRW-DOGE",
         EN_TICKSIZE),
        ("comment", "JSON output for AI agents",
         "upbit ticker KRW-BTC --json market,trade_price,signed_change_rate",
         EN_JSON),
        ("comment", "Cache info",
         "upbit cache",
         EN_CACHE),
        ("comment", "Korean locale",
         "LANG=ko_KR.UTF-8 upbit ticker KRW-BTC",
         EN_KO_LOCALE),
        # Watch TUI (AltScreen — comment + typing only)
        ("comment_only", "Watch TUI: real-time price table (Bubble Tea, full-screen)",
         "upbit watch ticker KRW-BTC KRW-ETH KRW-XRP"),
        ("comment_only", "Watch TUI: bid/ask depth chart",
         "upbit watch orderbook KRW-BTC"),
        ("comment_only", "Watch TUI: ASCII candlestick chart",
         "upbit watch candle KRW-BTC -i 1m"),
    ],
}

# ─────────────────────────────────────────────
# KO scenario
# ─────────────────────────────────────────────
KO_TICKER_1 = (
    "\ub9c8\ucf13     \ud604\uc7ac\uac00       \uc804\uc77c \ub300\ube44  \ubcc0\ub3d9\ub960  \uac70\ub798\ub7c9(24h)     \uac70\ub798\ub300\uae08(24h)          \uace0\uac00         \uc800\uac00       \r\n"
    "KRW-BTC  102,965,000  -817,000   -0.79%  1,314.48474815  136,563,852,318.50903  104,488,000  102,274,000"
)
KO_TICKER_MULTI = (
    "\ub9c8\ucf13     \ud604\uc7ac\uac00       \uc804\uc77c \ub300\ube44  \ubcc0\ub3d9\ub960  \uac70\ub798\ub7c9(24h)          \uac70\ub798\ub300\uae08(24h)          \uace0\uac00         \uc800\uac00       \r\n"
    "KRW-BTC  102,965,000  -817,000   -0.79%  1,314.48474815       136,563,852,318.50903  104,488,000  102,274,000\r\n"
    "KRW-ETH  3,125,000    -17,000    -0.54%  39,157.55640281      123,358,523,747.61168  3,191,000    3,096,000  \r\n"
    "KRW-XRP  2,100        -22        -1.04%  83,646,151.04691099  177,040,864,454.3182   2,140        2,075      "
)
KO_TRADES = (
    "\uccb4\uacb0 \uc2dc\uac01  \uccb4\uacb0\uac00       \uccb4\uacb0\ub7c9      \ub9e4\uc218/\ub9e4\ub3c4  \uc804\uc77c \ub300\ube44\r\n"
    "21:57:02   102,965,000  0.00011611  BID        -817,000 \r\n"
    "21:57:01   102,965,000  0.00009712  BID        -817,000 \r\n"
    "21:57:00   102,961,000  0.00476482  ASK        -821,000 \r\n"
    "21:57:00   102,965,000  0.00004856  BID        -817,000 \r\n"
    "21:57:00   102,965,000  0.0001214   BID        -817,000 "
)
KO_CANDLES = (
    "\uc2dc\uac01                 \uc2dc\uac00         \uace0\uac00         \uc800\uac00         \uc885\uac00         \uac70\ub798\ub7c9        \r\n"
    "2026-03-18 18:00:00  109,015,000  110,022,000  104,935,000  106,000,000  2,441.75375374\r\n"
    "2026-03-19 18:00:00  106,000,000  106,392,000  102,580,000  104,078,000  2,022.57056593\r\n"
    "2026-03-20 18:00:00  104,078,000  106,200,000  103,865,000  105,314,000  1,361.28731847\r\n"
    "2026-03-21 18:00:00  105,399,000  106,446,000  103,452,000  103,782,000  726.2663963   \r\n"
    "2026-03-22 18:00:00  103,831,000  104,488,000  102,274,000  102,965,000  935.95281897  "
)
KO_MARKET = (
    "\ub9c8\ucf13         \ud55c\uae00\uba85                      \uc601\ubb38\uba85                     \r\n"
    "KRW-WAXP     \uc641\uc2a4                        WAX                        \r\n"
    "KRW-CARV     \uce74\ube0c                        CARV                       \r\n"
    "KRW-LSK      \ub9ac\uc2a4\ud06c                      Lisk                       \r\n"
    "KRW-0G       \uc81c\ub85c\uc9c0                      0G                         \r\n"
    "KRW-BORA     \ubcf4\ub77c                        BORA                       \r\n"
    "KRW-PUNDIX   \ud380\ub514\uc5d1\uc2a4                    Pundi X                    \r\n"
    "KRW-USD1     \uc6d4\ub4dc\ub9ac\ubc84\ud2f0\ud30c\uc774\ub0b8\uc15c\uc720\uc5d0\uc2a4\ub514  World Liberty Financial USD\r\n"
    "KRW-BAT      \ubca0\uc774\uc9c1\uc5b4\ud150\uc158\ud1a0\ud070            Basic Attention Token      "
)
KO_TICKSIZE = (
    "\ub9c8\ucf13      \ud638\uac00\ud1b5\ud654  \ud638\uac00\ub2e8\uc704\r\n"
    "KRW-BTC   KRW       1000    \r\n"
    "KRW-XRP   KRW       1       \r\n"
    "KRW-DOGE  KRW       1       "
)
KO_JSON = '[{"market":"KRW-BTC","signed_change_rate":-0.0079108131,"trade_price":102961000}]'
KO_CACHE = (
    "\uacbd\ub85c  /Users/humphrey.park/.upbit/cache\r\n"
    "\ud30c\uc77c  2\uac1c\r\n"
    "\ud06c\uae30  27.6 KB"
)
KO_EN_LOCALE = (
    "Market   Price        Change    Change%  Volume(24h)     Trade Amt(24h)         High         Low        \r\n"
    "KRW-BTC  102,965,000  -817,000  -0.79%   1,314.48474815  136,563,852,318.50903  104,488,000  102,274,000"
)

ko_scenario = {
    "title": "\033[1;36m  upbit \u2014 AI-native CLI for Upbit Exchange\033[0m\r\n",
    "items": [
        ("comment", "\ube44\ud2b8\ucf54\uc778 \ud604\uc7ac\uac00 \uc870\ud68c",
         "upbit ticker KRW-BTC",
         KO_TICKER_1),
        ("comment", "\ubcf5\uc218 \ub9c8\ucf13 \ub3d9\uc2dc \uc870\ud68c",
         "upbit ticker KRW-BTC KRW-ETH KRW-XRP",
         KO_TICKER_MULTI),
        ("comment", "\ucd5c\uadfc \uccb4\uacb0 \ub0b4\uc5ed",
         "upbit trades KRW-BTC -c 5",
         KO_TRADES),
        ("comment", "\uc77c\ubd09 \uce94\ub4e4",
         "upbit candle KRW-BTC -i 1d -c 5",
         KO_CANDLES),
        ("comment", "KRW \ub9c8\ucf13 \ubaa9\ub85d",
         "upbit market -q KRW | head -9",
         KO_MARKET),
        ("comment", "\ud638\uac00 \ub2e8\uc704 \uc870\ud68c",
         "upbit tick-size KRW-BTC KRW-XRP KRW-DOGE",
         KO_TICKSIZE),
        ("comment", "AI \uc5d0\uc774\uc804\ud2b8\uc6a9 JSON \ucd9c\ub825",
         "upbit ticker KRW-BTC --json market,trade_price,signed_change_rate",
         KO_JSON),
        ("comment", "\uce90\uc2dc \uc815\ubcf4",
         "upbit cache",
         KO_CACHE),
        ("comment", "\uc601\uc5b4 \ub85c\ucf00\uc77c\ub85c \uc804\ud658",
         "LANG=en_US.UTF-8 upbit ticker KRW-BTC",
         KO_EN_LOCALE),
        # Watch TUI (AltScreen — comment + typing only)
        ("comment_only", "\uc2e4\uc2dc\uac04 TUI: \uc2e4\uc2dc\uac04 \uac00\uaca9 \ud14c\uc774\ube14 (Bubble Tea, \uc804\uccb4 \ud654\uba74)",
         "upbit watch ticker KRW-BTC KRW-ETH KRW-XRP"),
        ("comment_only", "\uc2e4\uc2dc\uac04 TUI: \ub9e4\uc218/\ub9e4\ub3c4 \ud638\uac00 \ucc28\ud2b8",
         "upbit watch orderbook KRW-BTC"),
        ("comment_only", "\uc2e4\uc2dc\uac04 TUI: ASCII \uce94\ub4e4\uc2a4\ud2f1 \ucc28\ud2b8",
         "upbit watch candle KRW-BTC -i 1m"),
    ],
}


def main():
    en_events = build_events(en_scenario)
    ko_events = build_events(ko_scenario)

    en_cast = make_cast("upbit \u2014 AI-native CLI for Upbit Exchange", en_events)
    ko_cast = make_cast("upbit \u2014 AI-native CLI for Upbit Exchange", ko_events)

    with open("docs/demo.cast", "w", encoding="utf-8") as f:
        f.write(en_cast)
    print("Generated: docs/demo.cast")

    with open("docs/demo.ko.cast", "w", encoding="utf-8") as f:
        f.write(ko_cast)
    print("Generated: docs/demo.ko.cast")


if __name__ == "__main__":
    main()
