---
updatedAt: 2026-01-07T10:53:34.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [Python] 24시간 누적 거래대금이 가장 높은 5개 종목 조회

```python Python
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

```json Response Example
{
  "KRW-XRP": 658205168495.2153,
  "KRW-ETH": 385094335136.0046,
  "KRW-KNC": 244286558989.6853,
  "KRW-ALT": 226701310792.52917,
  "KRW-GLM": 219180445460.17462,
}

```

# 유틸 라이브러리 Import

<!-- python@1-2 -->

기능 구현을 위해 필요한 모듈을 import 합니다. 별도의 설치가 필요한 모듈의 경우  `pip install <module name>` 명령어를 실행해 설치할 수 있습니다.

# 페어 목록 조회

<!-- python@5-16 -->

사용자가 입력한 마켓에서 지원하는 모든 페어의 목록을 조회하는 함수입니다. 거래를 지원하는 모든 페어의 이름을 추출해 업비트 API의 파라미터의 형식에 맞게 쉼표(comma)로 구분하는 문자열로 변환합니다.

# 24시간 누적 거래대금 조회

<!-- python@18-29 -->

사용자가 입력한 페어의 가격 데이터를 조회하는 함수입니다. 다양한 데이터 중 24시간 누적 거래대금을 추출해 “페어 이름” : “24시간 누적 거래대금” 형태로 반환합니다. 이 함수를 실행함으로써 조회하는 모든 페어의 24시간 누적 거래대금을 확인할 수 있습니다.

# 상위 5개 페어 필터링

<!-- python@31-33 -->

조회한 페어의 24시간 누적 거래대금을 확인하고 가장 많은 누적 거래대금을 가진 상위 5개 페어의 이름과 누적 거래대금을 반환하는 함수입니다.

# 24시간 누적 거래대금 상위 5개 페어 조회

<!-- python@35-39 -->

앞서 정의한 함수를 실행해 사용자가 입력한 마켓에서 지원하는 모든 페어 중 24시간 누적 거래대금이 가장 높은 5개 페어를 확인합니다.