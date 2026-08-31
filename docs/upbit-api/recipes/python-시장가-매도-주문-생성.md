---
updatedAt: 2026-06-02T08:40:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [Python] 시장가 매도 주문 생성

```python Python
from urllib.parse import unquote, urlencode
from collections.abc import Mapping
from typing import Any
import hashlib
import uuid
import jwt # PyJWT
import requests
from decimal import Decimal, ROUND_DOWN

# 본 코드를 실행하기 위해서는 JWT를 생성해 API를 호출해야 합니다.
# JWT 생성은 인증 토큰(JWT) 생성 가이드를 확인해 주시기 바랍니다.

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

```json Response Example
{
  "uuid": "<your order UUID>",
  "side": "ask",
  "ord_type": "market",
  "state": "done",
  "market": "KRW-BTC",
  "created_at": "2025-07-24T15:35:00+09:00",
  "volume": "0.00006251",
  "remaining_volume": "0",
  "prevented_volume": "0",
  "reserved_fee": "0",
  "remaining_fee": "0",
  "paid_fee": "5.0001749",
  "locked": "0",
  "prevented_locked": "0",
  "executed_volume": "0.00006251",
  "executed_funds": "10000.3498",
  "trades_count": 1
}
```

# 유틸 라이브러리 Import

<!-- python@1-8 -->

기능 구현을 위해 필요한 모듈을 import 합니다. 별도의 설치가 필요한 모듈의 경우  `pip install <module name>` 명령어를 실행해 설치할 수 있습니다.

# 인증 토큰 생성

<!-- python@10-11 -->

본 코드를 실행하기 위해서는 JWT를 생성해 API를 호출해야 합니다.

JWT 생성은 [인증 토큰(JWT) 생성 가이드](python-인증-토큰jwt-생성)를 확인해 주시기 바랍니다.

# 계정 잔고 조회

<!-- python@13-32 -->

사용자 포켓에 보유 중인 자산과 잔고를 조회하는 함수입니다. 원화(KRW) 자산과 디지털 자산 잔고를 모두 조회할 수 있으며 사용자가 입력한 `currency`의 보유 여부를 확인하고 보유하고 있는 경우 페어와 잔고를 반환합니다.

# 매도 가능 여부 확인

<!-- python@34-59 -->

사용자가 입력한 페어의 거래 가능 여부를 조회하는 함수입니다. 시장가 주문 생성 가능 여부와 주문 생성 시 주문당 최소 주문 금액을 반환합니다.

# 주문 생성

<!-- python@60-83 -->

거래 페어, 주문 방향(매수/매도), 주문 유형, 매도 수량을 파라미터로 입력받아 주문을 생성하고 생성된 주문의 UUID를 반환하는 함수입니다.

시장가 매도 주문을 생성할 예정이므로 거래 방향은 ask(매도), 거래 타입은 market(시장가)로 설정합니다.

# 주문 상태 조회

<!-- python@86-98 -->

UUID를 사용해 특정 주문의 정보를 조회 및 반환하는 함수입니다. 현재 주문의 상태를 확인할 수 있습니다.

# 시장가 매도 주문 생성

<!-- python@101-118 -->

앞서 정의한 함수를 사용해 시장가 매도 주문을 생성하는 함수를 구현합니다.

사용자가 입력한 페어로 해당 디지털 자산의 보유 여부와 잔고를 확인하고 잔고의 절반을 시장가로 매도하는 주문을 생성합니다. 주문 생성 후 UUID를 반환하고 해당 UUID로 주문을 조회해 주문의 상태를 조회합니다. 

주석을 해제하고 코드를 실행할 경우, 실제 주문이 생성될 수 있으므로 실행하기 전 다시 한 번 확인해 주시기 바랍니다.