---
updatedAt: 2026-07-14T09:05:34.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 지정가 매수 주문 생성

업비트 API를 사용해 지정가로 디지털 자산을 매수하는 주문을 생성하는 방법을 안내합니다.

<Recipe slug="python-지정가-매수-주문-생성" title="[Python] 지정가 매수 주문 생성" />

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

본 튜토리얼에서는 디지털 자산을 지정가로 매수하기 위한 절차들을 API로 수행하는 간단한 자동 매매 코드를 구현해볼 수 있습니다.

업비트 사용자가 디지털 자산 매수 시 일반적으로 수행하는 절차를 순서대로 정리하면 다음과 같습니다.

1. 원화를 입금합니다.
2. 매수할 디지털 자산 선택합니다.
3. 현재 호가 정보를 확인합니다.
4. 원하는 주문 단가과 주문 수량을 입력하여 매수 주문을 생성합니다.
5. 생성된 주문의 상태를 확인합니다.

이 가이드에서는 위 순서에 따라 매수할 디지털 자산의 최고 매수 호가를 조회하고 해당 호가보다 3% 낮은 가격으로 지정가 매수 주문을 생성하는 시나리오를 진행합니다. Python 예제 코드를 함께 제공합니다.

  <div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> 본 가이드는 실제 자산을 이용한 디지털 자산 매수 과정을 포함하고 있습니다.
      </div>
    따라서 가이드에 작성된 코드를 그대로 작성 및 실행할 경우, 사용자의 실제 자산을 사용해 디지털 자산을 매수할 수 있습니다. 디지털 자산 매수를 원하지 않는 사용자는 업비트 API 사용 방법만 참조하시고, 실제 매수 관련 함수는 호출하지 않도록 주의해 주시기 바랍니다.
  </div>

<br />

## 인증 안내

일부 API는 호출 시 인증을 위한 JWT를 생성해 요청의 헤더에 입력해야 합니다. [인증](https://docs.upbit.com/kr/reference/auth) 문서와 아래 레시피를 참고하여 모든 Exchange API 호출 시 인증 헤더를 추가하시기 바랍니다.

<Recipe slug="python-인증-토큰jwt-생성" title="[Python] 인증 토큰(JWT) 생성" />

<Recipe slug="java-인증-토큰jwt-생성" title="[Java] 인증 토큰(JWT) 생성" />

<Recipe slug="nodejs-인증-토큰jwt-생성" title="[Node.js] 인증 토큰(JWT) 생성" />

<br />

<br />

## 원화 입금 및 확인

### 원화 입금

디지털 자산을 매수하기 위해 사전에 충분한 수량의 원화를 거래할 포켓에 입금해야 합니다. [원화 입금](https://docs.upbit.com/kr/reference/deposit-krw) API를 호출하여 원화를 입금하고 입금 결과를 확인하는 함수를 아래와 같이 구현합니다. 입출금은 메인포켓만 지원하는 기능입니다. 메인포켓 API Key의 입금하기 권한을 받은 Key로 이용 바랍니다.

  <div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i>  원화 입금은 케이뱅크 연동 및 2차 인증을 완료해야 진행이 가능합니다.
    </div>
업비트에 케이뱅크 연동이 되어있지 않거나, 2차 인증을 완료하지 않은 경우 튜토리얼을 진행할 수 없습니다. 다음 가이드를 진행하기 전, 케이뱅크 연동 및 2차 인증 설정을 확인해 주시기 바랍니다.
  </div>

```python
def deposit_krw(amount: int, two_factor_type: str) -> Mapping:
    body = {
        "amount": amount,   
        "two_factor_type": two_factor_type
    }
    query_string = _build_query_string(body)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/deposits/krw"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=body).json()
    uuid = response.get('uuid')
    return uuid, response
```

<br />

케이뱅크 계좌에 보유 중인 원화 중 입금할 총액(원화)과 2차 인증을 진행할 수단을 파라미터로 입력하고 함수를 실행합니다. 다음과 같은 응답을 확인할 수 있습니다.

```json
 {
  "type": "deposit",
  "uuid": "<your_krw_deposit_uuid>",
  "currency": "KRW",
  "net_type": None,
  "txid": "<your_krw_deposit_txid>",
  "state": "PROCESSING",
  "created_at": "2025-07-09T14:10:08+09:00",
  "done_at": None,
  "amount": "5000.0",
  "fee": "0.0",
  "transaction_type": "default"
}

```

<br />

### 원화 입금 확인

사용자가 원화 입금 함수를 실행해 원화 입금을 요청했지만 2차 인증을 완료하지 않았기 때문에 입금의  state 필드가 PROCESSING, 입금 완료 시각인 done\_at 필드가 None인 것을 확인할 수 있습니다. [단일 입금 조회](https://docs.upbit.com/kr/reference/get-deposit) API를 호출하여 입금 정보를 조회하는 함수를 구현하고 2차 인증을 완료한 후, 입금 정보의 UUID로 해당 입금 건의 반영 여부를 확인합니다.

```python
def get_deposit_krw(uuid: str) -> Mapping:
    params = {
        "uuid": uuid
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/deposit"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    return response
```

<br />

2차 인증을 완료한 후, 동일한 입금 건의 state 필드와 done\_at 필드의 값이 변경된 것을 확인할 수 있습니다.

```json
{
  "type": "deposit",
  "uuid": "<your_deposit_krw_uuid>",
  "currency": "KRW",
  "net_type": None,
  "txid": "<your_deposit_krw_txid>",
  "state": "ACCEPTED",
  "created_at": "2025-07-09T14:10:08+09:00",
  "done_at": "2025-07-09T14:10:36+09:00",
  "amount": "5000.0",
  "fee": "0.0",
  "transaction_type": "default"
}
```

<br />

## 지정가 매수 주문 생성

### 거래 가능 여부 확인

사용자가 매수하고 싶은 디지털 자산의 페어를 입력해 업비트가 해당 페어의 거래를 지원하는지 확인하는 함수를 구현합니다. [페어 목록 조회](https://docs.upbit.com/kr/reference/list-trading-pairs) API를 활용합니다.

```python
def get_trading_pair(trading_pair: str) -> str:
    url = "https://api.upbit.com/v1/market/all"
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers).json()
    trading_pair_list = [
        item for item in response if item.get('market') == trading_pair]
    if len(trading_pair_list) == 0:
        raise ValueError("The trading pair list is empty.")
    return trading_pair_list[0].get('market')
```

<br />

함수 실행 시 업비트가 거래를 지원하는 페어 객체의 배열을 반환합니다. 객체 내 market 필드의 값을 사용자가 입력한 페어와 비교하여 일치하는 값이 있는 경우 market필드의 값을 반환하고 일치하는 값이 없는 경우 에러를 발생시킵니다.

```json
[
  {
    "market": "KRW-BTC",
    "korean_name": "비트코인",
    "english_name": "Bitcoin"
  },
  {
    "market": "KRW-ETH",
    "korean_name": "이더리움",
    "english_name": "Ethereum"
  },
  {
    "market": "KRW-XRP",
    "korean_name": "엑스알피(리플)",
    "english_name": "XRP"
  },
  ...
  {
    "market": "KRW-SOL",
    "korean_name": "솔라나",
    "english_name": "Solana"
  }
]
```

<br />

### 현재 호가 조회

[호가 조회](https://docs.upbit.com/kr/reference/list-orderbooks) API를 호출하여 사용자가 선택한 디지털 자산의 현재 호가를 조회하는 함수를 구현합니다. 이 튜토리얼에서는 현재 매수 호가 중 가장 높은 호가를 반환하는 로직을 구현하였습니다.

```python
getcontext().prec = 16

def get_best_bid_price(trading_pair: str) -> Decimal:
    params = {
        "markets": trading_pair
    }
    url = "https://api.upbit.com/v1/orderbook"
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    orderbook_units = response[0].get('orderbook_units')
    highest_bid_price = Decimal(str(orderbook_units[0].get('bid_price')))
    if highest_bid_price is None:
        raise ValueError(
            "Please check the orderbook. {response}".format(response=response))
    else:
        return highest_bid_price
```

<br />

함수 실행 시 다음과 같이 해당 디지털 자산의 가장 높은 매수 호가를 확인할 수 있습니다.

```
148446000
```

<br />

### 주문 가격 단위 확인 및 주문 단가 설정

업비트는 마켓(원화 마켓, BTC 마켓, USDT 마켓)별로 디지털 자산의 주문 가격 단위를 제공합니다. 주문 가격 단위의 경우, 디지털 자산의 가격 구간에 따라 달라지며 자세한 내용은 링크를 통해 확인할 수 있습니다.

* [원화 마켓 주문 가격 단위 알아보기](https://docs.upbit.com/kr/docs/krw-market-info)
* [BTC 마켓 주문 가격 단위 알아보기](https://docs.upbit.com/kr/docs/btc-market-info)
* [USDT 마켓 주문 가격 단위 알아보기](https://docs.upbit.com/kr/docs/usdt-market-info)

이 가이드에서는 원화 마켓의 디지털 자산 주문 가격 단위 조회 기능만 제공합니다. BTC 마켓과 USDT 마켓의 주문 가격 단위가 필요한 경우 별도의 구현이 필요합니다.

<br />

주문 가격 단위를 조회할 수 있는 함수를 구현합니다. 사용자가 입력하는 디지털 자산의 가격으로 해당 가격에 부합하는 주문 가격 단위를 반환합니다.

```python
def get_tick_size(price: Decimal) -> Decimal:
    if price <= 0:
        raise ValueError("price must be > 0")

    if price < Decimal("0.00001"):
        return Decimal("1e-8")

    decade = int(price.log10().to_integral_value(rounding=ROUND_DOWN))

    if decade < 3:
        return Decimal(10) ** (decade - 2)

    if decade >= 6:
        return Decimal("1000")

    base = Decimal(10) ** (decade - 3)
    leading = price / (Decimal(10) ** decade)
    step = Decimal("5") if leading >= Decimal("5") else Decimal("1")
    return min(base * step, Decimal("1000"))
```

<br />

앞서 구현한 가격 단위 조회 함수를 사용하여 주문 단가를 설정하는 함수를 구현합니다. 사용자가 입력한 디지털 자산 가격을 호가 단위에 맞춰 내림 처리합니다. 이를 통해 사용자는 주문 가격을 호가 단위 기준에 맞게 설정할 수 있으며, 주문 금액 초과를 방지할 수 있습니다.

```python
def round_price_by_tick_size(price: Decimal) -> Decimal:
    tick = get_tick_size(price)
    return (price // tick) * tick
```

<br />

### 주문 생성

지정가 매수 주문을 생성하기 위해서는 다음과 같은 값이 필요합니다.

1. 매수하고자 하는 디지털 자산의 페어(예: KRW-BTC)
2. 주문 방향(매수, 매도)
3. 주문 유형(지정가, 시장가 등)
4. 디지털 자산의 매수 단가 (디지털 자산 1개의 가격)
5. 매수할 디지털 자산의 수량

이 튜토리얼은 최고 매수 호가 대비 3% 낮은 가격으로 지정가 매수 주문을 생성하는 튜토리얼입니다. 따라서 주문 방향, 주문 유형은 각각 bid, limit 으로 정의되어 있습니다. 또한 디지털 자산의 주문 가격 단위에 맞는 단가를 계산하는 함수를 구현했기 때문에 매수할 디지털 자산의 페어와 수량을 입력해 지정가 매수 주문을 생성할 수 있습니다.

  <div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i>주문을 생성하기 전 확인해 주세요!
    </div>
실제 주문 생성시, 주문 정책에 따라 사용자 잔고와 무관하게 주문이 거절될 수 있습니다. 주문을 실행하기 전, 주문 조건이 정책에 부합한지 확인 후 실행해 주시기 바랍니다.
    <br />
<br />

<p> • 현재가 기준 ±300% 이상 매수 주문 </p>
<p> • 시장가 기준 30호가 이상 매수/매도 주문 </p>
<p> • 기타 내부 정책 </p>
  </div>

  <div class="callout-section">
    <div class="callout-title">
      <i class="fa-solid fa-circle-exclamation"></i>매수 주문에 실제 지불되는 금액은 매수 총액과 차이가 있을 수 있습니다.
      </div>
위 예제 코드는 간단한 매수 수량 계산 예제로, 매수 총액(amount)과 실제 지불되는 금액은 체결 조건과 수수료에 따라 상이할 수 있습니다.
    <br><br>
    <li>실제 체결 단가(디지털 자산의 단가)는 슬리피지, 유동성에 따라 변동될 수 있습니다.</li>
		<li>매수에 실제 지불되는 금액은 (계산된 수량 X 실제 체결 단가 + 수수료)입니다.</li>
  </div>

<br />

주문을 생성하는 함수를 구현합니다. 사용자가 입력한 페어와 디지털 자산의 매수 단가, 매수 시 지불할 매수 총액(원화)을 입력해 지정가 매수 주문을 생성하고 주문의 UUID를 반환합니다.

실제 매수 주문을 체결하기 위해서는 **\[주문 금액] + \[주문 금액에 대한 수수료] 만큼의 자산이 필요합니다.** 따라서 사용자의 보유 금액이 주문 금액과 정확히 일치하더라도, 수수료에 해당하는 금액이 부족한 경우 주문은 생성되지 않습니다.

```python
def create_order(
    trading_pair: str,
    price: str,
    volume: str
) -> str:
    body = {
        "market": trading_pair,
        "side": "bid",
        "ord_type": "limit",
        "price": price,
        "volume": volume,
    }
    query_string = _build_query_string(body)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/orders"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=body).json()
    uuid = response.get('uuid')
    if uuid is None:
        raise ValueError(
            "Please check the response. {response}".format(response=response))
    else:
        return uuid
```

<br />

### 주문 상태 조회

주문 정보의 UUID를 사용해 주문의 현재 상태를 조회하는 함수를 구현합니다. 이 함수를 사용해 지정가 매수 주문의 상태를 조회합니다.

```python
def get_order(uuid: str) -> Mapping:
    params = {
        "uuid": uuid
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/order"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    return response
```

<br />

함수 실행 시 다음과 같이 주문 정보를 반환합니다. `state` 필드를 통해 주문 체결 대기, 체결 완료 등의 주문 상태를 확인할 수 있습니다.

```json
{
  "uuid": "<your_order_uuid>",
  "side": "bid",
  "ord_type": "limit",
  "price": "150947000",
  "state": "wait",
  "market": "KRW-BTC",
  "created_at": "2025-07-10T13:15:08+09:00",
  "volume": "0.00006625",
  "remaining_volume": "150947000",
  "prevented_volume": "0",
  "reserved_fee": "5.000119375",
  "remaining_fee": "5.000119375",
  "paid_fee": "0",
  "locked": "10005.238869375",
  "prevented_locked": "0",
  "executed_volume": "0",
  "trades_count": 0,
  "trades": []
}
```

<br />

## 전체 코드

원화를 입금하고 디지털 자산의 최고 호가 대비 3% 낮은 가격으로 지정가 매수 주문을 생성하는 전체 코드는 다음과 같습니다. 원화 입금의 중복 실행을 방지하기 위해 각각의 파일로 분류해 실행하는 것을 권장합니다.

### 원화 입금

```python
from urllib.parse import unquote, urlencode
from collections.abc import Mapping
import hashlib
import uuid
import jwt # PyJWT
import requests

access_key = "<YOUR_ACCESS_KEY>"
secret_key = "<YOUR_SECRET_KEY>"

# Auth Token을 생성하는 로직은 여기에 추가해 주시기 바랍니다.
  
def deposit_krw(amount: int, two_factor_type: str) -> Mapping:
    body = {
        "amount": amount,   
        "two_factor_type": two_factor_type
    }
    query_string = _build_query_string(body)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/deposits/krw"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=body).json()
    uuid = response.get('uuid')
    return uuid, response

def get_deposit_krw(uuid: str) -> Mapping:
    params = {
        "uuid": uuid
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/deposit"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    return response

# 주석을 해제하고 코드를 실행할 경우, 실제 원황 입금이 실행될 수 있습니다. 실행하기 전 다시 한 번 확인해 주시기 바랍니다.
# if __name__ == "__main__":
#     deposit_uuid, response = deposit_krw(amount=10000, two_factor_type="kakao")
#     
# if deposit_uuid is None:
#         raise ValueError("Please check the deposit request. {response}".format(response=response))
#     else:
#         get_deposit_krw_response = get_deposit_krw(deposit_uuid)
#         print(get_deposit_krw_response)
```

### 지정가 매수 주문 생성

```python
from urllib.parse import unquote, urlencode
from typing import Any, Union
from collections.abc import Mapping
import hashlib
import uuid
import jwt # PyJWT
import requests
from decimal import Decimal, getcontext, ROUND_DOWN

access_key = "<YOUR_ACCESS_KEY>"
secret_key = "<YOUR_SECRET_KEY>"

# Auth Token을 생성하는 로직은 여기에 추가해 주시기 바랍니다.

def get_trading_pair(trading_pair: str) -> str:
    url = "https://api.upbit.com/v1/market/all"
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers).json()
    trading_pair_list = [
        item for item in response if item.get('market') == trading_pair]
    if len(trading_pair_list) == 0:
        raise ValueError("The trading pair list is empty.")
    return trading_pair_list[0].get('market')

getcontext().prec = 16

def get_best_bid_price(trading_pair: str) -> Decimal:
    params = {
        "markets": trading_pair
    }
    url = "https://api.upbit.com/v1/orderbook"
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    orderbook_units = response[0].get('orderbook_units')
    highest_bid_price = Decimal(str(orderbook_units[0].get('bid_price')))
    if highest_bid_price is None:
        raise ValueError(
            "Please check the orderbook. {response}".format(response=response))
    else:
        return highest_bid_price


def get_tick_size(price: Decimal) -> Decimal:
    if price <= 0:
        raise ValueError("price must be > 0")

    if price < Decimal("0.00001"):
        return Decimal("1e-8")

    decade = int(price.log10().to_integral_value(rounding=ROUND_DOWN))

    if decade < 3:
        return Decimal(10) ** (decade - 2)

    if decade >= 6:
        return Decimal("1000")

    base = Decimal(10) ** (decade - 3)
    leading = price / (Decimal(10) ** decade)
    step = Decimal("5") if leading >= Decimal("5") else Decimal("1")
    return min(base * step, Decimal("1000"))


def round_price_by_tick_size(price: Decimal) -> Decimal:
    tick = get_tick_size(price)
    return (price // tick) * tick


def create_order(
    trading_pair: str,
    price: str,
    volume: str
) -> str:
    body = {
        "market": trading_pair,
        "side": "bid",
        "ord_type": "limit",
        "price": price,
        "volume": volume,
    }
    query_string = _build_query_string(body)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/orders"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=body).json()
    uuid = response.get('uuid')
    if uuid is None:
        raise ValueError(
            "Please check the response. {response}".format(response=response))
    else:
        return uuid


def get_order(uuid: str) -> Mapping:
    params = {
        "uuid": uuid
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/order"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    return response

# 주석을 해제하고 코드를 실행할 경우, 실제 주문이 생성될 수 있습니다! 실행하기 전 다시 한 번 확인해 주시기 바랍니다.
# if __name__ == '__main__':
#     trading_pair = "KRW-BTC"
#     volume = "0.0001"
#     trading_pair = get_trading_pair(trading_pair)
#     orderbook_unit = get_best_bid_price(trading_pair)
#     price_3percent_rounded = str(
#         round_price_by_tick_size(orderbook_unit * Decimal(0.97)))
#     volume = str(Decimal(volume).quantize(
#         Decimal('1e-8'), rounding=ROUND_DOWN))

#     order_uuid = create_order(trading_pair, price_3percent_rounded, volume)
#     order_info = get_order(order_uuid)
#     print(order_info)

```

<br />

# Sibling pages

* [시장가 매도 주문 생성](https://docs.upbit.com/kr/docs/market-ask-order-creation.md)