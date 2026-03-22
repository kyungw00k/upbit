# upbit — AI-native CLI for Upbit Exchange

[한국어](README.ko.md)

> The only Upbit CLI designed for both humans and AI agents.

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/kyungw00k/upbit)](https://github.com/kyungw00k/upbit/releases)

## Demo

![demo](docs/demo.gif)

## Installation

```bash
brew install kyungw00k/upbit/upbit    # Homebrew
# or
curl -sSL https://github.com/kyungw00k/upbit/releases/latest/download/upbit_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m).tar.gz | tar xz && mv upbit ~/.local/bin/
```

## Quick Start

```bash
upbit ticker KRW-BTC                           # Current price
upbit candle KRW-BTC -i 1d --from 2025-01-01   # Historical candles (cached)
upbit watch candle KRW-BTC -i 1m               # Live candlestick chart (TUI)

export UPBIT_ACCESS_KEY=... UPBIT_SECRET_KEY=...
upbit buy KRW-BTC -p 100000000 -V 0.001        # Limit buy
upbit balance                                   # Portfolio with KRW evaluation
```

## Features

**AI-First** — `upbit tool-schema` exports JSON Schema for LLM/MCP tool calling. Non-TTY auto-outputs JSON for seamless AI agent pipelines.

**Real-time TUI** — `watch ticker/orderbook/trade/candle` powered by Bubble Tea. Multi-market Tab switching. ASCII candlestick chart with volume pane.

**Smart Trading** — Tick-size auto-correction prevents order failures across KRW/BTC/USDT markets. Confirmation prompts (`--force` to skip). Test orders (`--test`).

**Candle Cache** — SQLite cache with `--from` auto-pagination fetches full history. `upbit cache` to inspect or `--clear`.

**i18n** — `LANG=ko_KR` → Korean, default → English. POSIX locale standard (LC_ALL > LC_MESSAGES > LANG).

**Self-Update** — `upbit update` downloads latest from GitHub Releases with SHA256 verification. `--check` to preview.

## Commands

Run `upbit --help` for full list. Key commands:

| Category | Commands |
|----------|----------|
| Market Data | `ticker`, `candle`, `orderbook`, `trades`, `market`, `tick-size` |
| Trading | `buy`, `sell`, `balance`, `order list/show/cancel/replace` |
| Wallet | `wallet`, `deposit list/show/address`, `withdraw list/show/request` |
| Real-time | `watch ticker/orderbook/trade/candle/my-order/my-asset` |
| Utilities | `tool-schema`, `api-keys`, `cache`, `update` |

## Output

```bash
upbit ticker KRW-BTC              # Table (terminal)
upbit ticker KRW-BTC | jq .       # JSON (pipe, auto)
upbit ticker KRW-BTC -o csv       # CSV
upbit ticker KRW-BTC --json price # Field selection
```

| Context | Default | Override |
|---------|---------|---------|
| Terminal (TTY) | Aligned table | `-o json`, `-o csv` |
| Pipe (non-TTY) | Compact JSON | `-o table`, `-o jsonl` |

## Authentication

API keys are read **only from environment variables** — never stored on disk.

```bash
export UPBIT_ACCESS_KEY=your_access_key
export UPBIT_SECRET_KEY=your_secret_key
```

Market data commands work without authentication.
Trading, deposits/withdrawals, and personal WebSocket streams (`watch my-order`, `watch my-asset`) require authentication.

## Real-time TUI

Watch commands feature full-screen terminal UI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea):

- **watch ticker** — Real-time price table with color (rise=green, fall=red)
- **watch orderbook** — Bid/ask depth chart centered on spread
- **watch trade** — Scrolling trade stream with color coding
- **watch candle** — ASCII candlestick chart with volume pane

Multi-market support: use Tab or ←/→ to switch between markets.

```bash
upbit watch ticker KRW-BTC KRW-ETH KRW-XRP
upbit watch orderbook KRW-BTC
upbit watch candle KRW-BTC -i 1m
```

## Candle Cache

```bash
upbit candle KRW-BTC --from 2025-01-01   # Auto-paginate (SQLite cached)
upbit candle KRW-BTC --from 2025-01-01 --no-cache
upbit cache                               # Cache info
upbit cache --clear                       # Clear cache
```

## AI Integration

```bash
# Export all command schemas for LLM function calling
upbit tool-schema

# AI agents get JSON automatically (non-TTY)
upbit ticker KRW-BTC | jq '.trade_price'

# Select specific fields
upbit ticker KRW-BTC --json market,trade_price,signed_change_rate
```

## Claude Code Integration

```bash
claude skill add --url https://kyungw00k.github.io/upbit/skill.md
```

## Self-Update

```bash
upbit update          # Download and install latest version
upbit update --check  # Check only, don't download
```

Single static binary, zero runtime dependencies. Cross-platform: Linux, macOS, Windows (amd64, arm64).

## Acknowledgements

- Candlestick chart inspired by [cli-candlestick-chart](https://github.com/Julien-R44/cli-candlestick-chart)

## License

MIT
