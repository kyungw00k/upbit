---
name: upbit
description: "Upbit exchange CLI — market data, trading, deposits/withdrawals, real-time streams. Trigger: crypto price, trading, Upbit, buy/sell, candle, orderbook."
---

# Upbit CLI

AI-native CLI for Upbit cryptocurrency exchange. All commands output JSON in non-TTY mode.

## Authentication

```bash
export UPBIT_ACCESS_KEY=your_key
export UPBIT_SECRET_KEY=your_secret
```

Market data commands work without authentication.

## Commands

### Market Data (no auth)

#### `upbit ticker`
View current price

| Param | Type | Description |
|-------|------|-------------|
| `market` | array | Market code (e.g. KRW-BTC, KRW-ETH) |
| `quote` | string | View all prices by quote currency (e.g. KRW, BTC, USDT) |

Response: `acc_trade_price_24h`, `acc_trade_volume_24h`, `change`, `change_price`, `change_rate`, `high_price`...

#### `upbit candle`
View candles

| Param | Type | Description |
|-------|------|-------------|
| `asc` | boolean | Sort oldest first (default) |
| `count` | integer | Number of results |
| `desc` | boolean | Sort newest first |
| `force` | boolean | Skip confirmation prompt |
| `from` | string | Start time (e.g. 2025-01-01, 2025-01-01T09:00:00+09:00) |
| `interval` | string | Candle interval (1s, 1m, 3m, 5m, 10m, 15m, 30m, 60m, 240m, 1d, 1w, 1M, 1y) |
| `market` | string | Market code (e.g. KRW-BTC) |
| `no-cache` | boolean | Skip cache |

Response: `candle_acc_trade_price`, `candle_acc_trade_volume`, `candle_date_time_kst`, `candle_date_time_utc`, `high_price`, `low_price`...

#### `upbit market`
List markets

| Param | Type | Description |
|-------|------|-------------|
| `quote` | string | Quote currency filter (KRW, BTC, USDT) |

Response: `english_name`, `korean_name`, `market`

#### `upbit orderbook`
View orderbook

| Param | Type | Description |
|-------|------|-------------|
| `market` | array | Market code (e.g. KRW-BTC) |

Response: `market`, `orderbook_units`, `total_ask_size`, `total_bid_size`

#### `upbit trades`
View recent trades

| Param | Type | Description |
|-------|------|-------------|
| `count` | integer | Number of results |
| `market` | string | Market code (e.g. KRW-BTC) |

Response: `ask_bid`, `market`, `trade_price`, `trade_timestamp`, `trade_volume`

#### `upbit tick-size`
View tick size

| Param | Type | Description |
|-------|------|-------------|
| `market` | array |  **(req)** |

Response: `market`, `quote_currency`, `tick_size`

#### `upbit orderbook-levels`
View orderbook level units

| Param | Type | Description |
|-------|------|-------------|
| `market` | array |  **(req)** |

### Trading (auth required)

#### Order Types & Fees

| Flags | Order Type | Fee |
|-------|-----------|-----|
| `-p` + `-V` | Limit order | 0.05% |
| `-t` | Market buy (total amount) | 0.05% |
| `-V` (sell only) | Market sell | 0.05% |
| `--best` + `-V` | Best price limit order | 0.05% |
| `--watch` + `-p` + `-V` | Reserved limit order (triggers at watch price) | 0.139% |

#### Percentage Orders

Volume (`-V`) and total (`-t`) accept percentages (e.g. `50%`, `100%`).
- **Buy**: fee is auto-deducted. Formula: `available × ratio ÷ (1 + fee_rate)`
- **Sell**: uses holding balance directly. Formula: `holding × ratio`

#### Price Keywords

`--price` (`-p`) accepts keywords resolved via Ticker API:

| Keyword | Resolves to |
|---------|-------------|
| `now` | Current trade price |
| `open` | Today's opening price |
| `low` | Today's low price |
| `high` | Today's high price |

Keywords and percentages can be combined: `-p now -V 50%`

#### Time in Force (`--tif`)

| TIF | Description |
|-----|-------------|
| *(default)* | GTC — remains until filled or cancelled |
| `ioc` | Immediate or Cancel — fills available, cancels rest |
| `fok` | Fill or Kill — fills entirely or cancels entirely |
| `post_only` | Maker only — cancels if would take immediately |

TIF works with limit and best orders. Market and reserved orders use GTC only.

#### `upbit buy`
Place buy order

| Param | Type | Description |
|-------|------|-------------|
| `market` | string | Market code to buy (e.g. KRW-BTC) **(req)** |
| `price` | string | Order price (number or keyword: now/open/low/high) |
| `volume` | string | Order volume (supports percentage like 50%) |
| `total` | string | Order total for market buy (supports percentage like 50%) |
| `best` | boolean | Best price limit order |
| `watch` | string | Reserved order watch price (triggers when reached) |
| `tif` | string | Time in Force (ioc, fok, post_only) |
| `smp` | string | Self-trade prevention (cancel_maker, cancel_taker, reduce) |
| `id` | string | Client-specified order identifier |
| `test` | boolean | Test order (no actual execution) |
| `force` | boolean | Skip confirmation prompt |

Response: `created_at`, `executed_volume`, `market`, `ord_type`, `price`, `remaining_volume`...

#### `upbit sell`
Place sell order

| Param | Type | Description |
|-------|------|-------------|
| `market` | string | Market code to sell (e.g. KRW-BTC) **(req)** |
| `price` | string | Order price (number or keyword: now/open/low/high) |
| `volume` | string | Order volume (supports percentage like 50%) |
| `best` | boolean | Best price limit order |
| `watch` | string | Reserved order watch price (triggers when reached) |
| `tif` | string | Time in Force (ioc, fok, post_only) |
| `smp` | string | Self-trade prevention (cancel_maker, cancel_taker, reduce) |
| `id` | string | Client-specified order identifier |
| `test` | boolean | Test order (no actual execution) |
| `force` | boolean | Skip confirmation prompt |

Response: `created_at`, `executed_volume`, `market`, `ord_type`, `price`, `remaining_volume`...

#### `upbit balance`
View account balance

| Param | Type | Description |
|-------|------|-------------|
| `currency` | string | Currency code to query (e.g. KRW, BTC) |

Response: `avg_buy_price`, `balance`, `currency`, `eval_krw`, `locked`

#### `upbit order list`
List orders

| Param | Type | Description |
|-------|------|-------------|
| `closed` | boolean | Show closed orders |
| `count` | integer | Number of results |
| `market` | string | Market filter (e.g. KRW-BTC) |
| `page` | integer | Page number |

Response: `created_at`, `executed_volume`, `market`, `ord_type`, `price`, `remaining_volume`...

#### `upbit order show`
View order details

| Param | Type | Description |
|-------|------|-------------|
| `id` | boolean | Query by identifier (instead of UUID) |
| `uuid` | array | Order/deposit/withdrawal UUID **(req)** |

Response: `created_at`, `executed_volume`, `market`, `ord_type`, `price`, `remaining_volume`...

#### `upbit order cancel`
Cancel order

| Param | Type | Description |
|-------|------|-------------|
| `all` | boolean | Cancel all pending orders |
| `force` | boolean | Skip confirmation prompt |
| `market` | string | Market filter (use with --all) |
| `uuid` | string | UUID to cancel |

#### `upbit order replace`
Replace order (cancel and re-order)

| Param | Type | Description |
|-------|------|-------------|
| `force` | boolean | Skip confirmation prompt |
| `ord-type` | string | Order type (limit, best) |
| `price` | string | New order price |
| `uuid` | string | Order UUID to replace **(req)** |
| `volume` | string | New order volume |

#### `upbit order chance`
View order chance by market

| Param | Type | Description |
|-------|------|-------------|
| `market` | string | Market code for order chance (e.g. KRW-BTC) **(req)** |

### Deposits & Withdrawals (auth required)

#### `upbit wallet`
View wallet service status

#### `upbit deposit list`
List deposits

| Param | Type | Description |
|-------|------|-------------|
| `count` | integer | Number of results |
| `currency` | string | Currency code filter (e.g. BTC, KRW) |
| `page` | integer | Page number |
| `state` | string | Deposit state filter (PROCESSING, ACCEPTED, CANCELLED, REJECTED) |

#### `upbit deposit show`
View deposit details

| Param | Type | Description |
|-------|------|-------------|
| `uuid` | string | Order/deposit/withdrawal UUID **(req)** |

#### `upbit deposit address create`
Request deposit address creation

| Param | Type | Description |
|-------|------|-------------|
| `currency` | string |  **(req)** |
| `net-type` | string | Network type (default: same as currency) |

#### `upbit withdraw list`
List withdrawals

| Param | Type | Description |
|-------|------|-------------|
| `count` | integer | Number of results |
| `currency` | string | Currency code filter (e.g. BTC, KRW) |
| `page` | integer | Page number |
| `state` | string | Withdrawal state filter (WAITING, PROCESSING, DONE, FAILED, CANCELLED, REJECTED) |

#### `upbit withdraw show`
View withdrawal details

| Param | Type | Description |
|-------|------|-------------|
| `uuid` | string | Order/deposit/withdrawal UUID **(req)** |

#### `upbit withdraw request`
Request withdrawal

| Param | Type | Description |
|-------|------|-------------|
| `amount` | string | Withdrawal amount (required) |
| `currency` | string | Currency code to withdraw (e.g. BTC, KRW) **(req)** |
| `force` | boolean | Skip confirmation prompt |
| `net-type` | string | Network type (default: same as currency) |
| `secondary-address` | string | Secondary address (tag/memo) |
| `to` | string | Recipient address (required for digital assets) |
| `two-factor` | string | 2FA method (kakao, naver, hana) — KRW only |
| `tx-type` | string | Transaction type (default, internal) |

#### `upbit withdraw cancel`
Cancel withdrawal

| Param | Type | Description |
|-------|------|-------------|
| `force` | boolean | Skip confirmation prompt |
| `uuid` | string | UUID to cancel **(req)** |

### Real-time WebSocket

#### `upbit watch ticker`
Real-time ticker stream

| Param | Type | Description |
|-------|------|-------------|
| `market` | array | Market code (e.g. KRW-BTC, KRW-ETH) **(req)** |

#### `upbit watch orderbook`
Real-time orderbook stream

| Param | Type | Description |
|-------|------|-------------|
| `market` | array | Market code (e.g. KRW-BTC) **(req)** |

#### `upbit watch trade`
Real-time trade stream

| Param | Type | Description |
|-------|------|-------------|
| `market` | array |  **(req)** |

#### `upbit watch candle`
Real-time candle stream

| Param | Type | Description |
|-------|------|-------------|
| `interval` | string | Candle unit (1s, 1m, 3m, 5m, 10m, 15m, 30m, 60m, 240m) |
| `market` | array | Market code (e.g. KRW-BTC) **(req)** |

#### `upbit watch my-order`
Real-time my orders stream (auth required)

| Param | Type | Description |
|-------|------|-------------|
| `market` | array |  |

#### `upbit watch my-asset`
Real-time my assets stream (auth required)

### Utilities

#### `upbit api-keys`
List API keys

| Param | Type | Description |
|-------|------|-------------|
| `all` | boolean | Show all including expired |

#### `upbit cache`
Cache management

| Param | Type | Description |
|-------|------|-------------|
| `clear` | boolean | Clear cache |
| `force` | boolean | Skip confirmation prompt |

#### `upbit update`
Check and install latest version

| Param | Type | Description |
|-------|------|-------------|
| `check` | boolean | Check only, don't download |

## Output

- Non-TTY: JSON (automatic)
- `-o json/csv/jsonl/table`
- `--json field1,field2` for field selection

## Examples

```bash
# Market data
upbit ticker KRW-BTC -o json              # BTC price
upbit candle KRW-BTC -i 1d --from 2025-01-01 -o csv  # historical candles

# Trading — basic
upbit buy KRW-BTC -p 100000000 -V 0.001 --force      # limit buy
upbit sell KRW-BTC -V 0.001                           # market sell

# Trading — percentage & keywords
upbit buy KRW-BTC -p now -V 50%           # limit buy at current price, 50% of balance
upbit buy KRW-BTC -t 100%                 # market buy with full KRW balance
upbit sell KRW-BTC -V 100%                # market sell entire holding
upbit sell KRW-BTC -p high -V 50%         # limit sell 50% at today's high

# Trading — best price & reserved
upbit buy KRW-BTC -V 0.001 --best --tif ioc           # best price IOC buy
upbit sell KRW-BTC --watch 55000000 -p 54500000 -V 0.001  # reserved sell

# Portfolio
upbit balance -o json                      # portfolio with KRW eval
upbit order list -o json                   # open orders
```