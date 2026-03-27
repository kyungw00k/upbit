# Trading Guide

[한국어](trading.ko.md)

Detailed guide for the Upbit CLI `buy` and `sell` commands.

## Order Types

### Limit Order

Places an order at a specified price. Fills only when the price is reached.

```bash
upbit buy KRW-BTC -p 50000000 -V 0.001    # Buy 0.001 BTC at 50M KRW
upbit sell KRW-BTC -p 55000000 -V 0.001   # Sell 0.001 BTC at 55M KRW
```

### Market Order

Executes immediately at the current best available price.

```bash
upbit buy KRW-BTC -t 100000               # Market buy 100K KRW worth
upbit sell KRW-BTC -V 0.001               # Market sell 0.001 BTC
```

- Buy: specify total amount with `--total` (-t)
- Sell: specify volume with `--volume` (-V)

### Best Price Limit Order

Places a limit order at the counterparty's best available price. Use the `--best` flag.

```bash
upbit buy KRW-BTC -V 0.001 --best         # Best price limit buy
upbit sell KRW-BTC -V 0.001 --best        # Best price limit sell
```

### Reserved Limit Order

Triggers a limit order when the market reaches the watch price (`--watch`).

```bash
upbit buy KRW-BTC --watch 49000000 -p 49500000 -V 0.001   # Buy at 49.5M when price hits 49M
upbit sell KRW-BTC --watch 55000000 -p 54500000 -V 0.001  # Sell at 54.5M when price hits 55M
```

## Time in Force (TIF)

Use `--tif` to specify order execution conditions.

| TIF | Description | Behavior |
|-----|-------------|----------|
| *(default)* | GTC (Good Till Cancel) | Remains until filled or cancelled |
| `ioc` | Immediate or Cancel | Fills what's available immediately, cancels the rest |
| `fok` | Fill or Kill | Fills entirely or cancels entirely |
| `post_only` | Maker Only | Only places maker orders; cancels if it would take |

### TIF Compatibility by Order Type

| Order Type | GTC | IOC | FOK | post_only |
|------------|:---:|:---:|:---:|:---------:|
| Limit | O | O | O | O |
| Market (price/market) | O | - | - | - |
| Best price limit | O | O | O | - |
| Reserved limit (watch) | O | - | - | - |

### Examples

```bash
upbit buy KRW-BTC -p 50000000 -V 0.001 --tif ioc       # Limit IOC
upbit buy KRW-BTC -p 50000000 -V 0.001 --tif fok       # Limit FOK
upbit buy KRW-BTC -p 50000000 -V 0.001 --tif post_only # Maker only
upbit buy KRW-BTC -V 0.001 --best --tif ioc             # Best price IOC
upbit sell KRW-BTC -V 0.001 --best --tif fok            # Best price FOK
```

## Percentage Orders

Specify percentages for `--volume` (-V) and `--total` (-t).

### Buy

```bash
upbit buy KRW-BTC -p 50000000 -V 50%      # Limit buy with 50% of KRW balance
upbit buy KRW-BTC -t 100%                 # Market buy with full KRW balance
upbit buy KRW-BTC -t 25%                  # Market buy with 25% of KRW balance
```

- **Fees are automatically deducted** for buy orders
- Formula: `available_balance × ratio ÷ (1 + fee_rate)`
- Balance and fee rate are fetched via `GetOrderChance` API

### Sell

```bash
upbit sell KRW-BTC -V 100%                # Market sell entire BTC holding
upbit sell KRW-BTC -p 55000000 -V 50%     # Limit sell 50% of BTC holding
```

- Sell fees are deducted after execution, no pre-calculation needed
- Formula: `holding × ratio`

## Price Keywords

Use keywords with `--price` (-p) to auto-resolve prices from current market data.

| Keyword | Description | Ticker Field |
|---------|-------------|--------------|
| `now` | Current price (last trade) | `trade_price` |
| `open` | Today's opening price | `opening_price` |
| `low` | Today's low | `low_price` |
| `high` | Today's high | `high_price` |

### Examples

```bash
upbit buy KRW-BTC -p now -V 50%           # Limit buy at current price, 50% balance
upbit buy KRW-BTC -p low -V 0.001         # Limit buy at today's low
upbit sell KRW-BTC -p high -V 100%        # Limit sell all at today's high
upbit sell KRW-BTC -p open -V 0.5         # Limit sell 0.5 BTC at opening price
```

Price keywords are resolved before percentage calculation. Keywords and percentages can be combined.

## Fees

| Order Type | Fee Rate |
|------------|----------|
| Limit / Market / Best price limit | 0.05% |
| Reserved limit | 0.139% (incl. VAT) |

- Percentage buy orders automatically deduct fees to prevent insufficient balance errors
- Sell fees are deducted post-execution

## Additional Features

### Tick Size Auto-Adjustment

Limit order prices are automatically adjusted to the nearest valid tick size.

- Buy: rounds down (favorable to buyer)
- Sell: rounds up (favorable to seller)
- Watch prices for reserved orders are also auto-adjusted

### Test Orders

Use `--test` to simulate orders without actual execution.

```bash
upbit buy KRW-BTC -p 50000000 -V 0.001 --test
```

### Skip Confirmation

Use `--force` (-f) to skip the confirmation prompt.

```bash
upbit buy KRW-BTC -t 100000 --force
```

## Order Type Summary

| Flags | Order Type | Fee |
|-------|-----------|-----|
| `-p` + `-V` | Limit | 0.05% |
| `-t` | Market buy | 0.05% |
| `-V` (sell) | Market sell | 0.05% |
| `--best` + `-V` | Best price limit | 0.05% |
| `--watch` + `-p` + `-V` | Reserved limit | 0.139% |
