---
updatedAt: 2026-05-04T05:49:32.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 계정주 확인 자동화

사용자가 타 거래소에 보유 중인 디지털 자산을 업비트로 입금 시 자금세탁을 방지하기 위해 가상자산사업자가 계정주 동일 여부를 확인해야 합니다. 업비트 API를 사용해 자동으로 계정주 확인을 하는 방법을 학습합니다.

[block:tutorial-tile]
{
  "backgroundColor": "#ffffff",
  "emoji": "✅",
  "id": "691be0e426e5b5b0e24ced13",
  "link": "https://docs.upbit.com/v1.6.0/recipes/python-계정주-확인-자동화",
  "slug": "python-계정주-확인-자동화",
  "title": "[Python] 계정주 확인 자동화"
}
[/block]

[block:html]
{
  "html": "<meta name=\"robots\" content=\"noindex\">\n<div class=\"callout-section\">\n  <div class=\"callout-title\">\n    <i class=\"fa-solid fa-graduation-cap\"></i> 튜토리얼 예제 코드 안내\n    </div>\n  \t본 튜토리얼은 사용자의 이해를 돕기 위해 단계별 부분 코드 및 API 호출/응답 예시를 포함하고 있습니다. \n\t\t<br>전체 코드 예제는 Recipes 메뉴에서 확인하실 수 있습니다. \n  \t<br>위 버튼을 클릭하면 본 튜토리얼의 전체 코드 Recipe 페이지로 이동합니다.\n</div>"
}
[/block]

<br />

## 시작하기

디지털 자산 입금은 다음과 같은 순서로 진행됩니다.

1. 업비트에서 디지털 자산을 입금받을 수 있는 입금 주소를 생성합니다.
2. 타 거래소에 보유 중인 디지털 자산을 업비트의 입금 주소로 전송합니다.
3. 수신자와 송신자 계정주 일치 여부를 확인합니다.
4. 확인이 완료된 후, 입금 반영 여부를 확인합니다.

이 가이드에서는 업비트 API를 사용해 위와 동일한 기능을 구현하고 타 거래소에서 업비트로 입금 시 자동으로 트래블룰 검증을 진행하는 방법을 안내합니다.

[block:html]
{
  "html": "  <div class=\"callout-section callout-section--danger\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i> 이 가이드는 실제 자산을 이용한 디지털 자산 입금 과정을 포함하고 있습니다.\n      </div>\n    따라서 가이드에 작성된 코드를 그대로 작성 및 실행할 경우, 사용자가 타 거래소에 보유한 실제 자산을 업비트로 입금합니다. 디지털 자산의 입금을 원하지 않는 경우, 업비트 API 사용 방법만 참조하시고, 실제 입금 관련 함수는 호출하지 않도록 주의해 주시기 바랍니다.\n  </div>"
}
[/block]

[block:html]
{
  "html": "<div class=\"callout-section\">\n  <div class=\"callout-title\">\n    <i class=\"fa-solid fa-graduation-cap\"></i> 트래블룰이란?\n    </div>\n  트래블룰은 자금세탁방지(AML)를 위해 시행되는 규제로, 디지털 자산을 송수신하는 가상자산사업자(VASP)가 송금인과 수취인의 정보를 검증한 후 전송해야 하는 제도입니다.\n  <br>\n  <br>\n  100만 원 이상의 디지털 자산을 입금하는 경우 트래블룰 적용 대상이며 송금인, 수취인 일치 여부를 검증하고 확인된 사용자의 입금 건만 정상적으로 반영합니다. 트래블룰에 대한 자세한 사항은 아래 링크를 통해 확인할 수 있습니다. \n  <br>\n  <br>\n<a href=\"https://support.upbit.com/hc/ko/articles/4498679629337-%ED%8A%B8%EB%9E%98%EB%B8%94%EB%A3%B0-%EC%95%8C%EC%95%84%EB%B3%B4%EA%B8%B0\"> 업비트 고객센터 > 트래블룰 알아보기<a/>\n\n  </div>"
}
[/block]

<br />

## 인증 안내

입금 관리 API는 Exchange API로, 인증 대상입니다. [인증](https://docs.upbit.com/kr/reference/auth) 문서와 아래 레시피를 참고하여 API 호출 시 인증 헤더를 추가하시기 바랍니다.

<br />

## 입금 주소 생성

타 거래소에서 출금한 자산을 업비트로 입금받기 위해서는 입금 주소가 필요합니다. [입금 주소 생성 요청](https://docs.upbit.com/kr/reference/create-deposit-address) API를 호출하여 지정한 디지털 자산의 입금 주소를 생성하는 함수를 구현합니다.

```python
def create_deposit_address(currency: str, net_type: str) -> Mapping:
    body = {
        "currency": currency,
        "net_type": net_type
    }
    query_string = _build_query_string(body)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/deposits/generate_coin_address"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=body).json()
    return response
```

<br />

해당 API는 비동기로 처리되기 때문에 최초 실행 시 "생성 중" 응답을 확인할 수 있습니다.

```json
{
  "success": true,
  "message": "USDT 입금주소를 생성중입니다."
}
```

<br />

일정 시간 이후 해당 함수를 다시 한 번 호출하거나 이미 입금 주소가 존재하는 경우, 다음과 같이 생성된 입금 주소가 반환됩니다.

```json
{
  "currency": "USDT",
  "net_type": "TRX",
  "deposit_address": "<your_deposit_address>",
  "secondary_address": None
}

```

<br />

## 타 거래소에서 업비트로 입금

타 거래소에서 업비트 입금 주소로 디지털 자산을 입금합니다.

[block:html]
{
  "html": "  <div class=\"callout-section callout-section--danger\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i> 입금 전 반드시 확인하세요!\n      </div>\n    입금 시, 디지털 자산의 블록체인 네트워크가 일치하지 않거나, 생성한 입금 주소와 실제 디지털 자산을 전송받을 주소를 다르게 입력한 경우, 입금이 반영되지 않을 수 있습니다. 입금 전 반드시 전송받는 입금 주소와 해당 디지털 자산의 네트워크를 확인하시기 바랍니다.\n  </div>"
}
[/block]

<br />

[단일 입금 조회](https://docs.upbit.com/kr/reference/get-deposit) API를 호출하여 특정 입금 건의 정보를 조회하는 함수를 구현합니다. 타 거래소에서 업비트로 디지털 자산을 출금할 때 받은 UUID 또는 TxID 로 해당 입금 건의 상태를 조회할 수 있습니다.

* UUID로 특정 입금 건 정보 조회

```python
def get_deposit_by_uuid(uuid: str) -> Mapping: 
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

* TxID로 특정 입금 건 정보 조회

```python
def get_deposit_by_txid(txid: str) -> Mapping: 
    params = {
        "txid": txid
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

입금 정보의 state 필드의 값이 TRAVEL\_RULE\_SUSPECTED인 경우 트래블룰 계정주 확인 진행이 필요한 상태입니다.

```json
{
  "type": "deposit",
  "uuid": "<your_deposit_uuid>",
  "currency": "USDT",
  "net_type": "TRX",
  "txid": "<your_deposit_txid>",
  "state": "TRAVEL_RULE_SUSPECTED",
  "created_at": "2025-07-04T15:00:02+09:00",
  "done_at": "2025-07-04T15:00:48+09:00",
  "amount": "2500.363189",
  "fee": "0.0",
  "transaction_type": "default"
}
```

<br />

## 계정주 확인

### 트래블룰 계정주 확인 서비스 지원 거래소 목록 조회

트래블룰 계정주 확인 서비스를 요청하기 위해서는 상대 거래소 UUID가 필요합니다. 아래 함수를 구현하여 [계정주 확인 서비스 지원 거래소 목록 조회](https://docs.upbit.com/kr/reference/list-travelrule-vasps) API를 호출한 뒤 거래소의 이름에 상응하는 UUID를 조회할 수 있습니다.

```python
def get_vasp_uuid(vasp_name: str) -> str:
    params = {
        "vasp_name": vasp_name
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/travel_rule/vasps"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    vasp_uuid = next((item.get('vasp_uuid') for item in response if item.get('vasp_name') == vasp_name), None)
    if vasp_uuid is None:
        raise ValueError("{vasp_name} is NOT_FOUND".format(vasp_name=vasp_name))
    return vasp_uuid
```

<br />

### 계정주 확인 요청

트래블룰 검증을 실행하는 함수를 구현합니다. 입금한 거래소의 UUID와 입금 건에 대한 UUID 혹은 TxID로 검증을 실행합니다.

* UUID로 트래블룰 검증 - [입금 UUID로 계정주 검증 요청](https://docs.upbit.com/kr/reference/verify-travelrule-by-uuid)

```python
def verify_travel_rule_by_uuid(deposit_uuid: str, vasp_uuid: str) -> str:
    params = {
        "deposit_uuid": deposit_uuid,
        "vasp_uuid": vasp_uuid
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/travel_rule/deposit/uuid"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=params).json()
    verification_result = response.get('verification_result')
    if verification_result is None:
        raise ValueError("Please check the response. {response}".format(response=response))
    else:
        return verification_result
```

* TxID로 트래블룰 검증 - [입금 TxID로 계정주 검증 요청](https://docs.upbit.com/kr/reference/verify-travelrule-by-txid)

```python
def verify_travel_rule_by_txid(deposit_txid: str, vasp_uuid: str, currency: str, net_type: str) -> str:
    params = {
        "txid": deposit_txid,
        "vasp_uuid": vasp_uuid,
        "currency": currency,
        "net_type": net_type
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/travel_rule/deposit/txid"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=params).json()
    verification_result = response.get('verification_result')
    if verification_result is None:
        raise ValueError("Please check the response. {response}".format(response=response))
    else:
        return verification_result
```

[block:html]
{
  "html": "  <div class=\"callout-section callout-section--danger\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i> 동일한 입금 건에 대한 트래블룰 검증 요청은 10분당 1회로 제한됩니다.\n      </div>\n따라서 호출하기 전 입력 정보를 다시 한 번 확인한 후 호출하시기 바랍니다.  </div>"
}
[/block]

<br />

함수 실행 시 다음과 같이 검증에 대한 결과를 확인할 수 있습니다.

```json
{
  "deposit_uuid": "00000000-0000-0000-0000-000000000000",
  "verification_result": "verified",
  "deposit_state": "PROCESSING"
}
```

<br />

## 입금 반영 확인

트래블룰 검증을 완료한 후, 앞서 구현한 입금 정보 조회 함수로 입금 건의 상태를 확인합니다.

```json
{
  "type": "deposit",
  "uuid": "<your deposit uuid>",
  "currency": "USDT",
  "net_type": "TRX",
  "txid": "<your deposit transaction txid>",
  "state": "ACCEPTED",
  "created_at": "2025-07-04T15:00:02+09:00",
  "done_at": "2025-07-04T15:00:48+09:00",
  "amount": "2500.363189",
  "fee": "0.0",
  "transaction_type": "default"
}
```

<br />

## 전체 코드

계정주를 자동으로 확인하는 전체 코드는 다음과 같습니다.

```python
from urllib.parse import unquote, urlencode
from typing import Any
from collections.abc import Mapping
import hashlib
import uuid
import jwt # PyJWT
import requests

def create_deposit_address(currency: str, net_type: str) -> Mapping:
    body = {
        "currency": currency,
        "net_type": net_type
    }
    query_string = _build_query_string(body)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/deposits/generate_coin_address"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=body).json()
    return response

def get_deposit_by_uuid(uuid: str) -> Mapping: 
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

def get_deposit_by_txid(txid: str) -> Mapping: 
    params = {
        "txid": txid
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

def get_vasp_uuid(vasp_name: str) -> str:
    params = {
        "vasp_name": vasp_name
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/travel_rule/vasps"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    vasp_uuid = next((item.get('vasp_uuid') for item in response if item.get('vasp_name') == vasp_name), None)
    if vasp_uuid is None:
        raise ValueError("{vasp_name} is NOT_FOUND".format(vasp_name=vasp_name))
    return vasp_uuid

def verify_travel_rule_by_uuid(deposit_uuid: str, vasp_uuid: str) -> str:
    params = {
        "deposit_uuid": deposit_uuid,
        "vasp_uuid": vasp_uuid
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/travel_rule/deposit/uuid"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=params).json()
    verification_result = response.get('verification_result')
    if verification_result is None:
        raise ValueError("Please check the response. {response}".format(response=response))
    else:
        return verification_result

def verify_travel_rule_by_txid(deposit_txid: str, vasp_uuid: str, currency: str, net_type: str) -> str:
    params = {
        "txid": deposit_txid,
        "vasp_uuid": vasp_uuid,
        "currency": currency,
        "net_type": net_type
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/travel_rule/deposit/txid"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=params).json()
    verification_result = response.get('verification_result')
    if verification_result is None:
        raise ValueError("Please check the response. {response}".format(response=response))
    else:
        return verification_result

# 주석을 해제하고 실행할 경우 트래블룰 검증 프로세스가 실행될 수 있습니다. 실행 전 다시 한 번 확인해 주시기 바랍니다.
# if __name__ == "__main__":
# ####### uuid 기반 트래블룰 검증 ########

#     currency = "<Enter_your_currency>"
#     net_type = "<Enter_your_net_type>"
#     deposit_address_dict = create_deposit_address(currency=currency, net_type=net_type)
#     deposit_address = deposit_address_dict.get('deposit_address')
    
#     if deposit_address is not None:
#         deposit_uuid = "<Enter_your_deposit_uuid>"
#         deposit_state = get_deposit_by_uuid(deposit_uuid)['state']
#         print(deposit_state)
#         if deposit_state == "ACCEPTED":
#             vasp_uuid = get_vasp_uuid("바이낸스")

#             result = verify_travel_rule_by_uuid(deposit_uuid=deposit_uuid, vasp_uuid=vasp_uuid)

#             if result == "verified":
#                 print(get_deposit_by_uuid(deposit_uuid))
#             else: 
#                 raise ValueError("Check the travel rule verification result.")

#         else: 
#             raise ValueError("This deposit does not require verification.")

#     else:
#         raise ValueError("Check the deposit address.")


# ######## txid 기반 트래블룰 검증 ########

#     currency = "<Enter_your_currency>"
#     net_type = "<Enter_your_net_type>"
#     deposit_address_dict = create_deposit_address(currency=currency, net_type=net_type)
#     deposit_address = deposit_address_dict.get('deposit_address')

#     if deposit_address is not None:
#         deposit_txid = "<Enter_your_deposit_txid>"
#         deposit_state = get_deposit_by_txid(deposit_txid)['state']

#         if deposit_state == "TRAVEL_RULE_SUSPECTED":
#             vasp_uuid = get_vasp_uuid("바이낸스")

#             result = verify_travel_rule_by_txid(deposit_txid=deposit_txid, vasp_uuid=vasp_uuid, currency=currency, net_type=net_type)

#             if result == "verified":
#                 print(get_deposit_by_txid(deposit_txid))
#             else: 
#                 raise ValueError("Check the travel rule verification result.")

#         else: 
#             raise ValueError("This deposit does not require verification.")

#     else:
#         raise ValueError("Check the deposit address.")
```