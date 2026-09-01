---
updatedAt: 2026-07-14T09:05:34.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 시장가 매도 주문 생성

업비트 API를 사용해 사용자가 보유한 디지털 자산을 시장가로 매도하는 주문을 생성하는 방법을 학습합니다.

<Recipe slug="python-시장가-매도-주문-생성" title="[Python] 시장가 매도 주문 생성" />

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

업비트 사용자가 디지털 자산을 매도할 때 일반적으로 수행하는 절차를 순서대로 정리하면 다음과 같습니다.

1. 사용자가 보유 중인 디지털 자산 중 매도할 디지털 자산을 선택합니다.
2. 선택한 디지털 자산의 주문 가능 정보를 확인해 주문 조건을 확인합니다.
3. 매도 수량을 입력하여 시장가 매도 주문을 생성합니다.
4. 주문이 체결됩니다.
5. 체결된 주문을 확인합니다.

이 가이드에서는 사용자가 보유한 디지털 자산과 수량을 조회한 뒤, 선택한 자산의 보유 수량의 50%를 시장가로 매도하는 시나리오를 구현합니다. Python 예제 코드를 함께 제공합니다.

  <div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> 본 가이드는 실제 자산을 이용한 디지털 자산 매도 과정을 포함하고 있습니다.
      </div>
    따라서 가이드에 작성된 코드를 그대로 작성 및 실행할 경우, 사용자의 실제 자산을 사용해 디지털 자산을 매도할 수 있습니다. 디지털 자산 매도를 원하지 않는 사용자는 업비트 API 사용 방법만 참조하시고, 실제 매도 관련 함수는 호출하지 않도록 주의해 주시기 바랍니다.
  </div>

  <div className="callout-section">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i>  매도 주문을 생성하기 위해서는 디지털 자산 매수가 선행되어야 합니다.
      </div>
    이 가이드를 따라하기 전, 디지털 자산 매수가 필요한 경우 <a href="limit-bid-order-creation">지정가 매수 주문 생성 가이드</a>를 따라 디지털 자산을 매수한 후 이 가이드를 진행해 주시기 바랍니다. 
  </div>

<br />

## 인증 안내

[인증](https://docs.upbit.com/kr/reference/auth) 문서와 아래 레시피를 참고하여 모든 Exchange API 호출 시 인증 헤더를 추가하시기 바랍니다.

<Recipe slug="python-인증-토큰jwt-생성" title="[Python] 인증 토큰(JWT) 생성" />

<Recipe slug="java-인증-토큰jwt-생성" title="[Java] 인증 토큰(JWT) 생성" />

<Recipe slug="nodejs-인증-토큰jwt-생성" title="[Node.js] 인증 토큰(JWT) 생성" />

<br />

## 보유 자산 확인 및 매도 자산 선택

<Anchor target="_blank" href="ref:get-balance">포켓 잔고 조회</Anchor> API를 호출하여 사용자 포켓에 보유중인 자산 잔고를 조회합니다.

```python
def get_pair_and_balance_from_account(currency: str) -> Mapping:
    jwt_token = _create_jwt(access_key, secret_key)
    url = "https://api.upbit.com/v1/accounts"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers).json()
    trading_pair_list = [item for item in response if item.get(
        "currency") == currency]
    if len(trading_pair_list) == 0:
        raise ValueError(
            "Currency {currency} is not found".format(currency=currency))
    else:
        pair = trading_pair_list[0]
        return {
            "pair": "{unit_currency}-{currency}".format(unit_currency=pair.get('unit_currency'), currency=pair.get('currency')),
            "balance": pair.get('balance')
        }
```

함수 실행 시 포켓이 보유하고 있는 모든 자산 정보를 조회하고 사용자가 입력한 `currency`와 일치하는 자산의 페어와 수량의 정보가 반환됩니다.

```json
[
  {
    "currency": "BTC",
    "balance": "0.00006617",
    "locked": "0",
    "avg_buy_price": "151109000",
    "avg_buy_price_modified": False,
    "unit_currency": "KRW"
  },
  {
    "currency": "KRW",
    "balance": "35264.15938484",
    "locked": "10005",
    "avg_buy_price": "0",
    "avg_buy_price_modified": True,
    "unit_currency": "KRW"
  },
  {
    "currency": "USDT",
    "balance": "28.363189",
    "locked": "0",
    "avg_buy_price": "1361.5",
    "avg_buy_price_modified": False,
    "unit_currency": "KRW"
  }
]

```

<br />

## 주문 조건 확인

선택한 디지털 자산의 거래 가능 여부와 시장가 매도 주문을 생성하기 위해 [페어별 주문 가능 정보 조회](https://docs.upbit.com/kr/reference/available-order-information) API를 호출하여 주문 가능 정보를 확인하는 함수를 구현합니다.

```python
def get_order_chance(trading_pair: str) -> Mapping:
    params = {
        "market": trading_pair
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/orders/chance"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    min_total = response.get("market").get("ask").get("min_total")
    ask_types = response.get("market").get("ask_types")
    if "market" not in ask_types:
        raise ValueError("This pair does not support market order. {ask_types}".format(
            ask_types=ask_types))
    else:
        order_type_market = True

    return {
        "order_type_market": order_type_market,
        "min_total": min_total
    }
```

<br />

함수 실행 시 사용자가 입력한 페어의 주문 가능 정보를 조회하여 시장가 매도 주문 가능 여부와 최소 주문 금액을 반환합니다.

```json
{
  "bid_fee": "0.0005",
  "ask_fee": "0.0005",
  "maker_bid_fee": "0.0005",
  "maker_ask_fee": "0.0005",
  "market": {
    "id": "KRW-BTC",
    "name": "BTC/KRW",
    "order_types": [
      "limit"
    ],
    "order_sides": [
      "ask",
      "bid"
    ],
    "bid_types": [
      "best_fok",
      "best_ioc",
      "limit",
      "limit_fok",
      "limit_ioc",
      "price"
    ],
    "ask_types": [
      "best_fok",
      "best_ioc",
      "limit",
      "limit_fok",
      "limit_ioc",
      // 시장가 매도 가능 확인
      "market"
    ],
    "bid": {
      "currency": "KRW",
      "min_total": "5000"
    },
    "ask": {
      "currency": "BTC",
      // 최소 매도 금액 5,000 원 확인
      "min_total": "5000"
    },
    "max_total": "1000000000",
    "state": "active"
  },
  "bid_account": {
    "currency": "KRW",
    "balance": "35264.15938484",
    "locked": "10005",
    "avg_buy_price": "0",
    "avg_buy_price_modified": True,
    "unit_currency": "KRW"
  },
  "ask_account": {
    "currency": "BTC",
    "balance": "0.00006617",
    "locked": "0",
    "avg_buy_price": "151109000",
    "avg_buy_price_modified": False,
    "unit_currency": "KRW"
  }
}

```

<br />

## 주문 생성

현재가에 디지털 자산을 매도하는 시장가 매도 주문을 생성하기 위해서는 다음과 같은 값이 필요합니다.

1. 매도하고자 하는 디지털 자산의 페어(예: KRW-BTC)
2. 주문 방향(매수, 매도)
3. 주문 유형(지정가, 시장가 등)
4. 매도할 디지털 자산의 수량

이 튜토리얼은 **시장가 매도** 주문을 생성하므로 주문 방향, 주문 유형은 각각 ask, market으로 설정합니다. 매도할 디지털 자산의 페어와 수량을 입력해 시장가 매도 주문을 생성할 수 있습니다.

시장가 매도 주문을 생성하는 함수를 구현합니다. 이 함수는 시장가 매도 주문을 생성한 후, [주문 생성](https://docs.upbit.com/kr/reference/new-order) API를 호출한 뒤 생성된 해당 주문의 UUID를 반환합니다.

```python
def create_order(
    trading_pair: str,
    volume: str
) -> str:
    body = {
        "market": trading_pair,
        "side": "ask",
        "ord_type": "market",
        "volume": volume
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

## 주문 상태 조회

[단일 주문 조회](https://docs.upbit.com/kr/reference/get-order) API를 호출하여 주문을 생성 시 반환된 UUID로 주문의 현재 상태를 조회하는 함수를 구현합니다.

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

함수 실행 시 사용자가 입력한 UUID 값과 일치하는 주문의 정보를 반환합니다. `state` 필드를 통해 주문 체결 대기, 체결 완료 등 현재 주문 상태를 확인할 수 있습니다.

```json
{
  "uuid": "<your_order_uuid>",
  "side": "ask",
  "ord_type": "market",
  "price": "163494000",
  "state": "done",
  "market": "KRW-BTC",
  "created_at": "2025-07-30T20:06:54+09:00",
  "volume": "0.00006116",
  "remaining_volume": "0.00006116",
  "prevented_volume": "0",
  "reserved_fee": "4.99964652",
  "remaining_fee": "4.99964652",
  "paid_fee": "0",
  "locked": "10004.29268652",
  "prevented_locked": "0",
  "executed_volume": "0",
  "trades_count": 0,
  "trades": []
}

```

<br />

## 전체 코드

시장가 매도 주문을 생성하는 전체 코드는 다음과 같습니다. 사용자가 입력한 currency 값으로 해당 디지털 자산의 보유 여부와 수량을 확인합니다. 이후 시장가 매도 주문 가능 여부를 검증한 뒤, 가능할 경우 보유 수량의 50%를 시장가로 매도하는 주문을 생성합니다.

```python
from urllib.parse import unquote, urlencode
from collections.abc import Mapping
from typing import Any
import hashlib
import uuid
import jwt # PyJWT
import requests
from decimal import Decimal, ROUND_DOWN

access_key = "<YOUR_ACCESS_KEY>"
secret_key = "<YOUR_SECRET_KEY>"

# Auth Token을 생성하는 로직은 여기에 추가해 주시기 바랍니다.

def get_pair_and_balance_from_account(currency: str) -> Mapping:
    jwt_token = _create_jwt(access_key, secret_key)
    url = "https://api.upbit.com/v1/accounts"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers).json()
    trading_pair_list = [item for item in response if item.get(
        "currency") == currency]
    if len(trading_pair_list) == 0:
        raise ValueError(
            "Currency {currency} is not found".format(currency=currency))
    else:
        pair = trading_pair_list[0]
        return {
            "pair": "{unit_currency}-{currency}".format(unit_currency=pair.get('unit_currency'), currency=pair.get('currency')),
            "balance": pair.get('balance')
        }


def get_order_chance(trading_pair: str) -> Mapping:
    params = {
        "market": trading_pair
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/orders/chance"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    min_total = response.get("market").get("ask").get("min_total")
    ask_types = response.get("market").get("ask_types")
    if "market" not in ask_types:
        raise ValueError("This pair does not support market order. {ask_types}".format(
            ask_types=ask_types))
    else:
        order_type_market = True

    return {
        "order_type_market": order_type_market,
        "min_total": min_total
    }


def create_order(
    trading_pair: str,
    volume: str
) -> str:
    body = {
        "market": trading_pair,
        "side": "ask",
        "ord_type": "market",
        "volume": volume
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


# 주석을 해제하고 실행할 경우 실제 시장가 매도 주문이 생성될 수 있습니다. 실행 전 다시 한 번 확인해 주시기 바랍니다.
# if __name__ == "__main__":
#     currency = "<Enter your currency>"
#     account_info = get_pair_and_balance_from_account(currency)
#     trading_pair = account_info.get("pair")

#     balance = account_info.get("balance")
#     order_chance = get_order_chance(trading_pair)
#     fifty_percent_volume = str(
#         (Decimal(balance) * Decimal("0.5")).quantize(Decimal("1e-8"), rounding=ROUND_DOWN))

#     if order_chance["order_type_market"]:
#         order_uuid = create_order(trading_pair, fifty_percent_volume)
#         order_info = get_order(order_uuid)
#         print("order_info: {order_info}".format(order_info=order_info))
#     else:
#         raise ValueError("This pair does not support market order. {order_chance}".format(
#             order_chance=order_chance))

```

<br />

# Sibling pages

* [지정가 매수 주문 생성](https://docs.upbit.com/kr/docs/limit-bid-order-creation.md)