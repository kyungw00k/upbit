---
updatedAt: 2026-05-04T05:48:52.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 24시간 누적 거래대금 확인

업비트 API는 사용자의 편리한 거래 분석을 위해 다양한 데이터를 제공합니다. 24시간 누적 거래 대금을 확인하는 예제를 제공합니다.

[block:tutorial-tile]
{
  "backgroundColor": "#ffffff",
  "emoji": "📊",
  "id": "691be0e426e5b5b0e24ced15",
  "link": "https://docs.upbit.com/v1.6.0/recipes/python-24시간-누적-거래대금이-가장-높은-종목-조회",
  "slug": "python-24시간-누적-거래대금이-가장-높은-종목-조회",
  "title": "[Python] 24시간 누적 거래대금이 가장 높은 종목 조회"
}
[/block]

[block:html]
{
  "html": "<meta name=\"robots\" content=\"noindex\">\n<div class=\"callout-section\">\n  <div class=\"callout-title\">\n    <i class=\"fa-solid fa-graduation-cap\"></i> 튜토리얼 예제 코드 안내\n    </div>\n  \t본 튜토리얼은 사용자의 이해를 돕기 위해 단계별 부분 코드 및 API 호출/응답 예시를 포함하고 있습니다. \n\t\t<br>전체 코드 예제는 Recipes 메뉴에서 확인하실 수 있습니다. \n  \t<br>위 버튼을 클릭하면 본 튜토리얼의 전체 코드 Recipe 페이지로 이동합니다.\n</div>"
}
[/block]

<br />

## 시작하기

특정 마켓에서 지원하는 모든 페어 중 24시간 동안 누적 거래대금이 많은 상위 5개의 페어를 확인하는 과정을 안내합니다. Python 예제 코드를 함께 제공합니다.

<br />

## 페어 목록 조회

[페어 목록 조회](https://docs.upbit.com/kr/reference/list-trading-pairs) API를 호출하여 사용자가 입력한 마켓에서 지원하는 모든 페어를 조회하는 함수를 아래와 같이 구현할 수 있습니다. JSON 배열 형식의 API 응답을 가공하여 페어 코드 목록을 쉼표(,)로 구분된 문자열 형식으로 반환합니다.

```python
def list_markets(quote_currencies: str) -> str:
    params = {
        "quote_currencies": quote_currencies
    }
    url = "https://api.upbit.com/v1/ticker/all"
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params)
    market_list = [item.get("market") for item in response.json()]
    string_market_list = ",".join(market_list)
    return string_market_list
```

<br />

## 페어 별 24시간 누적 거래대금 조회

[페어 단위 현재가 조회](https://docs.upbit.com/kr/reference/list-tickers) API에 쉼표(,)로 구분된 페어 목록을 파라미터로 호출하여 모든 페어의 현재가를 조회할 수 있습니다. 응답 중 `acc_trade_price_24h` 필드의 값이 24시간 누적 거래대금을 의미합니다.

```json
[
  {
    "market": "KRW-BTC",
    "trade_date": "20250704",
    "trade_time": "051400",
    "trade_date_kst": "20250704",
    "trade_time_kst": "141400",
    "trade_timestamp": 1751606040365,
    "opening_price": "148737000.00000000,",
    "high_price": "149360000.00000000,",
    "low_price": "148288000.00000000,",
    "trade_price": "148601000.00000000,",
    "prev_closing_price": "148737000.00000000,",
    "change": "FALL",
    "change_price": 136000,
    "change_rate": 0.0009143656,
    "signed_change_price": "-136000.00000000,",
    "signed_change_rate": -0.0009143656,
    "trade_volume": 0.00016823,
    "acc_trade_price": 31615925234.05438,
    "acc_trade_price_24h": 178448329314.96686,
    "acc_trade_volume": 212.38911576,
    "acc_trade_volume_24h": 1198.26954807,
    "highest_52_week_price": 163325000,
    "highest_52_week_date": "2025-01-20",
    "lowest_52_week_price": 72100000,
    "lowest_52_week_date": "2024-08-05",
    "timestamp": 1751606040403
  },
...
]
```

<br />

입력한 모든 페어의 24시간 누적 거래대금을 추출하는 함수를 구현합니다. 쉼표로 구분된 문자열 형식의 페어 목록을 입력하여 각 페어의 이름과 해당 페어의 24시간 누적 거래대금의 Mapping 객체를 반환합니다.

```python
def list_acc_trade_price_24h(trading_pairs: str) -> Mapping:
    params = {
        "markets": trading_pairs
    }
    url = "https://api.upbit.com/v1/ticker"
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    all_acc_trade_price_24h = {item.get("market"): item.get(
        "acc_trade_price_24h") for item in response}
    return all_acc_trade_price_24h
```

<br />

## 누적 거래대금 상위 5개 페어 조회

모든 페어의 24시간 누적 거래대금을 비교한 뒤 상위 5개 페어를 반환하는 함수를 구현합니다.

```python
def list_top_5_high_acc_trade_price_24h(list_acc_trade_price_24h: Mapping) -> Mapping:
    top_5_list = {k: v for k, v in sorted(list_acc_trade_price_24h.items(), key=lambda x: x[1], reverse=True)[:5]}
    return top_5_list
```

<br />

## 결과 확인

구현한 함수들을 실행하는 코드를 실행하여 지난 24시간 동안 누적 거래대금이 가장 높은 상위 5개 페어와 누적 거래대금을 아래와 같이 확인할 수 있습니다.

```python
if __name__ == "__main__":
    markets = list_markets("KRW")
    list_price_24h = list_acc_trade_price_24h(markets)
    top_5_price_24h = list_top_5_high_acc_trade_price_24h(list_price_24h)
    print(top_5_price_24h))
```

```json
{
  "KRW-OMNI": 529730651468.0878,
  "KRW-XRP": 434360960153.3289,
  "KRW-ETH": 283333357272.6667,
  "KRW-BTC": 188002200480.52707,
  "KRW-USDT": 127548558763.43745,
}

```

<br />

## 전체 코드 예제

전체 코드 예제는 아래와 같습니다.

```python
import requests
from collections.abc import Mapping


def list_markets(quote_currencies: str) -> str:
    params = {
        "quote_currencies": quote_currencies
    }
    url = "https://api.upbit.com/v1/ticker/all"
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params)
    market_list = [item.get("market") for item in response.json()]
    string_market_list = ",".join(market_list)
    return string_market_list

def list_acc_trade_price_24h(trading_pairs: str) -> Mapping:
    params = {
        "markets": trading_pairs
    }
    url = "https://api.upbit.com/v1/ticker"
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    all_acc_trade_price_24h = {item.get("market"): item.get(
        "acc_trade_price_24h") for item in response}
    return all_acc_trade_price_24h

def list_top_5_high_acc_trade_price_24h(list_acc_trade_price_24h: Mapping) -> Mapping:
    top_5_list = {k: v for k, v in sorted(list_acc_trade_price_24h.items(), key=lambda x: x[1], reverse=True)[:5]}
    return top_5_list

if __name__ == "__main__":
    markets = list_markets("KRW")
    list_price_24h = list_acc_trade_price_24h(markets)
    top_5_price_24h = list_top_5_high_acc_trade_price_24h(list_price_24h)
    print(top_5_price_24h)
```

# Sibling pages

* [RSI 지표 산출](https://docs.upbit.com/kr/docs/rsi-indicator-calculation-tutorial.md)