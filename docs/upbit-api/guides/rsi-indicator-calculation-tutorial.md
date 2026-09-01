---
updatedAt: 2026-07-14T09:05:34.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# RSI 지표 산출

RSI 지표에 대해 알아보고 업비트 시세 API 응답을 기반으로 RSI 지표를 산출하는 방법을 알아봅니다.

<Recipe slug="python-rsi-산출" title="[Python] RSI 산출" />

<meta name="robots" content="noindex" />
<div className="callout-section">
  <div className="callout-title">
    <i className="fa-solid fa-graduation-cap"></i> 튜토리얼 예제 코드 안내
    </div>
  	본 튜토리얼은 사용자의 이해를 돕기 위해 단계별 부분 코드 및 API 호출/응답 예시를 포함하고 있습니다. 
		<br />전체 코드 예제는 Recipes 메뉴에서 확인하실 수 있습니다. 
  	<br />위 버튼을 클릭하면 본 튜토리얼의 전체 코드 Recipe 페이지로 이동합니다.
</div>

<br />

## 시작하기

RSI(Relative Strength Index, 상대강도지수)는 주식과 선물 등 전통 금융 시장에서 널리 사용되는 기술적 분석 지표로서 디지털 자산 거래에서도 유의미한 지표로 사용됩니다. RSI의 개념과 계산 방법을 알아보고, 업비트 API를 활용하여 Python으로 디지털 자산의 RSI를 산출하는 방법을 단계별로 안내합니다.

<br />

## RSI(상대강도지수)란?

RSI (Relative Strength Index, 상대강도지수)란 자산의 과매수, 과매도 상태를 판단하는 지표입니다. 0 \~ 100 사이의 값을 가질 수 있으며, 70 이상을 과매수 상태, 30 이하를 과매도 상태로 해석합니다. 과매수 또는 과매도 상태의 정도는 자산의 추세 전환 가능성 예측에 참고할 수 있습니다. RSI 계산 공식은 다음과 같습니다. RS는 상대강도(Relative Strength)를 의미합니다.

```
RSI = 100 - (100 / (1 + RS))
```

  <div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> RSI는 절대적인 투자 지표가 아닙니다.
      </div>
      RSI는 데이터를 기술적으로 분석한 여러 보조 지표 중 하나이며, 매매 판단의 절대 근거로 사용할 수 없습니다.
  </div>

<br />

## RS(상대강도) 계산 공식

RSI를 계산하기 위해서는 먼저 RS(상대강도)와 평균 상승폭, 평균 하락폭을 계산해야 합니다. RS를 계산하는 공식은 다음과 같습니다.

* **RS 계산 공식**
  * RS = 누적 평균 상승폭 / 누적 평균 하락폭

<br />

* **누적 평균 상승폭, 누적 평균 하락폭 계산 공식**
  * 누적 평균 상승폭 = 누적 평균 상승폭 = (누적 평균 상승폭 \* (기간 - 1) + 전일 종가 대비 변화 금액(양수)) / 기간
  * 누적 평균 하락폭 = 누적 평균 하락폭 = (누적 평균 하락폭 \* (기간 - 1) + 전일 종가 대비 변화 금액(음수)) / 기간<br /><br />

이제 실제 디지털 자산의 누적 평균 상승폭, 누적 평균 하락폭을 계산해 RSI를 산출해 보겠습니다.

<br />

## RSI 계산하기

디지털 자산의 전일 종가는 [일(Day) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-days) API를 호출하여 쉽게 확인할 수 있습니다. 이 가이드에서는 비트코인의 캔들 데이터를 활용해 14일 동안의 평균 상승폭과 평균 하락폭을 구하고 RS와 RSI를 계산합니다.

```python
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
```

이 API는 `market`과 `to`, `count`를 파라미터로 받습니다. 각 인자에 대한 설명은 다음과 같습니다.

* `market`: 조회할 마켓의 디지털 자산을 특정합니다. 이 가이드에서는 원화 마켓의 비트코인으로 특정합니다.
* `to`: 캔들 데이터 조회 기간의 종료 시각을 설정합니다. count와 조합하여 조회할 캔들의 기간과 개수를 특정할 수 있습니다. ISO 8601 포맷(UTC)으로 입력합니다. 입력하지 않은 경우, 가장 최근 캔들 데이터 생성 시간을 기준으로 설정합니다. 이 가이드에서는 이 파라미터를 입력하지 않고 가장 최근 캔들 데이터를 기준으로 조회합니다.
* `count`: 응답으로 받을 캔들 데이터의 개수를 설정합니다. 일반적으로 RSI는 14일의 캔들 데이터를 기준으로 계산합니다.

과거 시점으로 부터 변화된 종가를 계산해야 하므로 최신 순(내림차순)으로 반환된 캔들 데이터를 오름차순으로 재정렬합니다.

  <div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> 캔들은 해당 시간대에 체결이 발생한 경우에만 생성됩니다.
      </div>
      해당 캔들의 시작시각부터 종료시각까지 체결이 발생하지 않은 경우 캔들이 생성되지 않으며, 응답에도 포함되지 않습니다. 예를 들어, candle_date_time이 \`2024-08-31T00:00:00\`인 일 캔들의 경우 \`2024-08-31T00:00:00\`(이상)부터 \`2024-09-01T00:00:00\`(미만)까지 체결이 발생하지 않은 경우 생성되지 않습니다.
  </div>

```json
[
  {
    market: 'KRW-BTC',
    candle_date_time_utc: '2025-06-23T00:00:00',
    candle_date_time_kst: '2025-06-23T09:00:00',
    opening_price: 140908000,
    high_price: 147101000,
    low_price: 139698000,
    trade_price: 146182000,
    timestamp: 1750723199727,
    candle_acc_trade_price: 346783036638.01794,
    candle_acc_trade_volume: 2422.95238077,
    prev_closing_price: 140908000,
    change_price: 5274000,
    change_rate: 0.0374286769
  },
...
  {
    market: 'KRW-BTC',
    candle_date_time_utc: '2025-07-07T00:00:00',
    candle_date_time_kst: '2025-07-07T09:00:00',
    opening_price: 148500000,
    high_price: 149100000,
    low_price: 147196000,
    trade_price: 148100000,
    timestamp: 1751932799103,
    candle_acc_trade_price: 123245812310.68657,
    candle_acc_trade_volume: 832.29461691,
    prev_closing_price: 148500000,
    change_price: -400000,
    change_rate: -0.0026936027
  }
]
```

<br />

### 평균 상승폭, 평균 하락폭 계산하기

응답의 `change_price`필드는 전일 종가 대비 가격의 변화량을 의미합니다. 이 값으로 아래와 같이 평균 상승폭과 평균 하락폭을 계산할 수 있습니다.

```python
def calculate(candle_data: Sequence, period: int = 14) -> Mapping:
    if len(candle_data) < period:
        raise ValueError("At least {period} candle data are required.".format(period=period))
    
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
```

<br />

RSI 계산에 앞서 먼저 충분한 개수의 일봉 캔들 데이터가 존재하는지 확인합니다. 이후 각 캔들 데이터의 `change_price` 값을 확인하여 두 개의 배열(gains, losses)에 값을 추가합니다.<br />•	`gains` 배열: `change_price`가 양수인 경우 해당 값을 추가하고, 음수이거나 0인 경우 0을 추가합니다.<br />•	`losses` 배열: `change_price`가 음수인 경우 해당 값을 절대값으로 변환해 추가하고, 양수이거나 0인 경우 0을 추가합니다.

이 과정을 거치면 `gains`와 `losses`배열은 캔들 데이터 개수와 동일한 길이를 가지게 됩니다.

```python
if len(candle_data) < period:
    raise ValueError("At least {period} candle data are required.".format(period=period))
    
    gains = []
    losses = []
    
    for item in candle_data:
        change = item.get('change_price')
        gains.append(change if change > 0 else 0)
        losses.append(abs(change) if change < 0 else 0)
```

<br />

`gains` 배열과 `losses` 배열의 첫 값과 period 기간의 값 까지의 합을 구하고 이를 period 기간의 수로 나누어 초기 평균 상승폭, 평균 하락폭을 계산합니다. 이후, gains 배열과 losses 배열의 값을 사용해 전체 캔들 데이터의 누적 평균 상승폭, 누적 평균 하락폭을 계산합니다.

```python
initial_au = sum(gains[:period]) / period
initial_ad = sum(losses[:period]) / period

au = initial_au
ad = initial_ad

for i in range(period, len(gains)):
    au = (au * (period - 1) + gains[i]) / period
    ad = (ad * (period - 1) + losses[i]) / period
```

<br />

### RS와 RSI 계산하기

누적 평균 상승폭과 누적 평균 하락폭을 통해 상대강도(RS)와 상대강도지수(RSI)를 계산할 수 있습니다. `calculate`함수에 RS와 RSI를 계산하는 공식을 추가합니다.

RS를 계산할 때, 누적 평균 하락폭이 0인 경우 하락이 전혀 없었다는 의미이므로 RSI는 100으로 처리합니다. 누적 하락폭이 0이 아닌 경우, 누적 평균 상승폭을 누적 평균 하락폭으로 나누어 RS를 계산하고 공식에 대입해 RSI를 계산합니다.

```python
def calculate(candle_data: Sequence[Mapping[str, Any]], period: int = 14) -> Mapping[str, Any]:

    ...
    
    rs = float('inf') if ad == 0 else au / ad
    
    rsi = (100 - 100 / (1 + rs))
    
    return {
        "AU": str(Decimal(au).quantize(Decimal("1e-4"), rounding=ROUND_DOWN)),
        "AD": str(Decimal(ad).quantize(Decimal("1e-4"), rounding=ROUND_DOWN)),
        "RS": str(Decimal(rs).quantize(Decimal("1e-4"), rounding=ROUND_DOWN)),
        "RSI": str(Decimal(rsi).quantize(Decimal("1e-4"), rounding=ROUND_DOWN))
    }
```

<br />

## 결과 확인

코드를 실행해 `KRW-BTC` 페어의 RSI를 조회할 수 있습니다.

```python
if __name__ == "__main__":
    candle_data = list_days_candles("KRW-BTC", 200)
    rsi_data = calculate(candle_data, 14)
    print(rsi_data)
```

```
{
  "AU": "626931.7147",
  "AD": "531397.2181",
  "RS": "1.1797",
  "RSI": "54.1238"
}
```

<br />

## 전체 코드 예제

RSI를 산출하는 전체 코드는 다음과 같습니다. 200일 간의 캔들 데이터로 누적 평균 상승,하락폭을 구하고 이를 사용해 14일간의 RSI를 계산할 수 있습니다.

```python
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

<br />

### 다음 단계: 계산한 RSI를 전략으로 검증해보기

지금까지 업비트 캔들 데이터로 RSI를 직접 계산해봤습니다. 그렇다면 "RSI가 특정 값 아래로 내려가면 사고, 다시 올라오면 판다"는 규칙은 과거 시장에서 어떻게 작동했을까요?

이를 코드로 일일이 구현하지 않고 확인해보고 싶다면 <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.3_add-strategy-toolkit-withdrawal-block/docs/upbit-strategy-toolkit">Upbit Strategy Toolkit</Anchor>을 활용할 수 있습니다. AI 에이전트와 대화하듯 매매 규칙을 정리하고, 업비트 과거 캔들 데이터로 백테스트(시뮬레이션)한 결과를 확인하는 도구입니다.

```text
"RSI가 30 이하로 내려가면 매수하고, 70 이상이면 매도하는 규칙을 2024년으로 백테스트해줘"
Upbit Strategy Toolkit은 이용자가 직접 정의한 규칙을 과거 데이터로 계산해보는 백테스트 전용 도구입니다. 특정 종목·매매 시점을 추천하지 않으며, 백테스트 결과는 과거 데이터에 특정 가정을 적용한 예시로 미래 수익을 보장하지 않습니다.
```

수익률·최대 낙폭(MDD)·승률 등 성과 지표를 리포트로 확인하고, RSI 기준값(예: 30 → 25)을 바꿔가며 결과를 비교해볼 수 있습니다.

> Upbit Strategy Toolkit은 이용자가 직접 정의한 규칙을 과거 데이터로 계산해보는 백테스트 전용 도구입니다. 특정 종목·매매 시점을 추천하지 않으며, 백테스트 결과는 미래 수익을 보장하지 않습니다.

<br />

# Sibling pages

* [24시간 누적 거래대금 확인](https://docs.upbit.com/kr/docs/24-hour-accumulated-trade-volume.md)