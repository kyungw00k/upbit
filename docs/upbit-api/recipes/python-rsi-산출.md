---
updatedAt: 2026-01-07T10:53:33.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [Python] RSI 산출

```python Python
from typing import Any
from collections.abc import Mapping, Sequence
import requests
from decimal import Decimal, ROUND_DOWN

def list_days_candles(trading_pair: str, count: int) -> Sequence:
    params = {
        "market": trading_pair,
        "count": count,
    }
    url = "https://api.upbit.com/v1/candles/days"
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    reversed_candle_data = response[::-1]
    return reversed_candle_data


def calculate(candle_data: Sequence, period: int = 14) -> Mapping:
    if len(candle_data) < period:
        raise ValueError(
            "At least {period} candle data are required.".format(period=period))

    gains = []
    losses = []

    for item in candle_data:
        change = item.get('change_price')
        gains.append(change if change > 0 else 0)
        losses.append(abs(change) if change < 0 else 0)

    initial_au = sum(gains[:period]) / period
    initial_ad = sum(losses[:period]) / period

    au = initial_au
    ad = initial_ad

    for i in range(period, len(gains)):
        au = (au * (period - 1) + gains[i]) / period
        ad = (ad * (period - 1) + losses[i]) / period

    rs = float('inf') if ad == 0 else au / ad

    rsi = (100 - 100 / (1 + rs))

    return {
        "AU": str(Decimal(au).quantize(Decimal("1e-4"), rounding=ROUND_DOWN)),
        "AD": str(Decimal(ad).quantize(Decimal("1e-4"), rounding=ROUND_DOWN)),
        "RS": str(Decimal(rs).quantize(Decimal("1e-4"), rounding=ROUND_DOWN)),
        "RSI": str(Decimal(rsi).quantize(Decimal("1e-4"), rounding=ROUND_DOWN))
    }


if __name__ == "__main__":
    candle_data = list_days_candles("KRW-BTC", 200)
    rsi_data = calculate(candle_data, 14)
    print(rsi_data)
```

```json Response Example
{'AU': '981369.4328', 'AD': '458194.5124', 'RS': '2.1418', 'RSI': '68.1712'}
```

# 유틸 라이브러리 Import

<!-- python@1-4 -->

기능 구현을 위해 필요한 모듈을 import 합니다. 별도의 설치가 필요한 모듈의 경우  `pip install <module name>` 명령어를 실행해 설치할 수 있습니다.

# 캔들 데이터 조회

<!-- python@6-17 -->

RSI 계산에 필요한 캔들 데이터를 조회하는 함수입니다. 캔들 데이터 조회 시 조회할 페어(`market`), 마지막 캔들 조회 시간(`to`), 조회할 캔들 데이터의 개수(`count`)를 파라미터로 입력할 수 있습니다. 

캔들 데이터 조회 API의 응답은 캔들 생성 시점을 기준으로 내림차순으로 반환되기 때문에 RSI 계산을 위해 응답의 순서를 오름차순으로 변환해 반환합니다.

# RSI 산출

<!-- python@20-52 -->

RSI를 계산하는 함수입니다. 사용자가 입력하는 기간 만큼의 캔들 데이터에서 전일 대비 변화한 가격(change_price)을 추출해 누적 평균 상승/하락폭을 계산합니다. 이를 바탕으로 RS와 RSI를 계산해 반환합니다.

# 실제 RSI 산출 확인

<!-- python@55-58 -->

앞서 정의한 함수를 실행해 KRW-BTC 페어의 200일 간의 캔들 데이터를 사용해 최근 14일 간의 RSI를 계산합니다. 코드를 실행해 RSI를 확인할 수 있습니다.