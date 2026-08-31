---
updatedAt: 2026-01-07T10:53:31.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [Python] 지정가 매수 주문 생성

```python Python
from urllib.parse import unquote, urlencode
from typing import Any, Union
from collections.abc import Mapping
import hashlib
import uuid
import jwt # PyJWT
import requests
from decimal import Decimal, getcontext, ROUND_FLOOR

# 본 코드를 실행하기 위해서는 JWT를 생성해 API를 호출해야 합니다.
# JWT 생성은 인증 토큰(JWT) 생성 가이드를 확인해 주시기 바랍니다.

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

```json Response Example
{
  "uuid": "<your_order_uuid>",
  "side": "bid",
  "ord_type": "limit",
  "price": "162211000",
  "state": "wait",
  "market": "KRW-BTC",
  "created_at": "2025-07-28T14:04:11+09:00",
  "volume": "0.00006165",
  "remaining_volume": "0.00006165",
  "prevented_volume": "0",
  "reserved_fee": "5.000154075",
  "remaining_fee": "5.000154075",
  "paid_fee": "0",
  "locked": "10005.308304075",
  "prevented_locked": "0",
  "executed_volume": "0",
  "trades_count": 0,
  "trades": []
}

```

# 유틸 라이브러리 Import

<!-- python@1-8 -->

기능 구현을 위해 필요한 모듈을 import 합니다. 별도의 설치가 필요한 모듈의 경우  `pip install <module name>` 명령어를 실행해 설치할 수 있습니다.

# 인증 토큰 생성

<!-- python@10-11 -->

본 코드를 실행하기 위해서는 JWT를 생성해 API를 호출해야 합니다.

[JWT 생성은 인증 토큰(JWT) 생성 가이드](python-인증-토큰jwt-생성)를 확인해 주시기 바랍니다.

# 거래 가능 여부 확인

<!-- python@13-23 -->

업비트가 사용자가 입력한 페어(예: KRW-BTC)의 거래를 지원하는지 확인하는 함수입니다. 사용자가 입력한 페어가 거래 지원 목록에 존재할 경우 해당 페어의 이름을 반환하고 존재하지 않을 경우 에러를 반환합니다.

# 최고 호가 조회

<!-- python@25-42 -->

사용자가 입력한 페어의 호가를 조회하여 현재 존재하는 매수 호가 중 가장 높은 호가를 반환하는 함수입니다. 이 가이드에서는 가장 높은 호가에서 3% 낮은 주문 단가로 지정가 매수 주문을 생성합니다.

# 원화 마켓의 호가 가격 단위 반환

<!-- python@45-63 -->

사용자가 입력한 가격을 기반으로 해당 가격에 적합한 호가 가격 단위를 반환합니다. 호가 가격 단위는 적합한 주문 단가를 계산할 때 사용합니다.

# 주문 단가 조정

<!-- python@66-68 -->

사용자가 입력한 주문 단가를 해당 호가 단위에 맞춰 계산합니다. 단가는 호가 단위에 따라 내림 처리되며, 이를 통해 안정적인 가격으로 주문을 생성할 수 있습니다.

# 주문 생성

<!-- python@71-96 -->

거래 페어, 주문 방향(매수/매도), 주문 유형, 주문 단가, 매수 수량을 파라미터로 입력받아 주문을 생성하고 해당 주문의 식별자인 UUID를 반환하는 함수입니다. 이 가이드에서는 지정가 매수 주문을 생성할 예정이므로 거래 방향은 bid(매수), 거래 유형은 limit(지정가)로 설정합니다.

# 주문 상태 조회

<!-- python@99-111 -->

UUID를 사용해 특정 주문 정보를 조회 및 반환하는 함수입니다. 이 함수를 호출하여 현재 주문의 상태를 확인할 수 있습니다.

# 실제 지정가 주문 생성

<!-- python@113-128 -->

앞서 정의한 함수를 사용해 지정가 매수 주문을 생성하는 함수를 구현합니다.

사용자가 입력한 페어의 거래 지원 여부와 최고 매수 호가를 확인하고 최고 매수 호가 대비 3% 낮은 주문 단가로 지정가 매수 주문을 생성합니다. 주문 생성 후 반환된 UUID로 주문의 현재 상태를 조회할 수 있습니다.

주석을 해제하고 코드를 실행할 경우, 실제 주문이 생성될 수 있으므로 실행하기 전 다시 한 번 확인해 주시기 바랍니다.