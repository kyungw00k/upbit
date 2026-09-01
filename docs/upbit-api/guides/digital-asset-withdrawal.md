---
updatedAt: 2026-07-14T09:05:34.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 디지털 자산 출금

API를 사용하여 사용자가 업비트에 보유한 디지털 자산을 타 거래소로 출금하는 방법을 알아봅니다.

<Recipe slug="python-디지털-자산-출금" title="[Python] 디지털 자산 출금" />

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

사용자가 업비트에 보유 중인 디지털 자산을 타 거래소로 출금하기 위한 순서는 다음과 같습니다.

1. 업비트 웹 사이트에서 출금하기 권한을 포함한 API Key를 준비합니다.
2. 업비트 모바일 앱에 로그인 후 \[더보기 > 인증/보안 > Open API 관리 > 활성]페이지에서 출금 안심차단을 해지해 줍니다.(최초 1회)
3. 업비트 웹 사이트에서 디지털 자산을 출금할 주소를 등록합니다.
4. 출금 받을 주소가 등록되었는지 확인합니다.
5. 해당 디지털 자산이 현재 출금 가능한 상태인지 확인합니다.
6. 디지털 자산을 출금합니다.
7. 디지털 자산 출금 건을 조회하여 출금 현황을 확인합니다.

이 가이드는 업비트 API를 사용해 위와 동일한 기능을 구현하고 업비트에서 타 거래소로 디지털 자산을 출금하는 방법을 제공합니다.

<div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> 이 가이드는 실제 자산을 이용한 디지털 자산 출금 과정을 포함하고 있습니다.
      </div>
    따라서 가이드에 작성된 코드를 그대로 작성 및 실행할 경우, 사용자가 업비트에 보유한 실제 디지털 자산을 타 거래소로 출금할 수 있습니다. 따라서 디지털 자산의 출금을 원하지 않는 경우, 업비트 API 사용 방법만 참조하시고, 실제 출금 관련 함수는 호출하지 않도록 주의해 주시기 바랍니다. 

출금 요청 시 \`open\_api\_withdraw\_locked\`에러가 반환된다면 \`출금 안심차단\`을 해지했는지 확인해주세요.

  </div>

<br />

## 인증 안내

출금 관리 API는 Exchange API로, 인증 대상입니다. [인증](https://docs.upbit.com/kr/reference/auth) 문서와 아래 레시피를 참고하여 API 호출 시 인증 헤더를 추가하시기 바랍니다.

<br />

## 디지털 자산 출금 주소 등록

업비트에서 타 거래소 혹은 지갑으로 디지털 자산을 출금하기 위해서는 사전에 출금 주소를 등록해야 합니다. 출금 허용 주소 등록 가이드를 확인해 출금 주소 등록 방법을 확인하시기 바랍니다.

* [거래소 지갑 주소 등록 가이드 바로가기](https://docs.upbit.com/kr/docs/open-api-withdraw-custodial-wallet)
* [개인지갑 주소 등록 가이드 바로가기](https://docs.upbit.com/kr/docs/open-api-withdraw-non-custodial-wallet)

<br />

## 출금 주소 등록 확인

[출금 허용 주소 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawal-addresses) API를 호출하여 사용자가 등록한 출금 허용 주소 목록 중 특정 디지털 자산에 관련된 주소 목록을 반환하는 함수를 구현합니다. 사용자가 출금 주소 등록 시 입력한 출금 주소와 네트워크 타입, 해당 주소를 발급한 거래소 이름이 반환됩니다.

```python
def get_withdrawal_address(currency: str, net_type: str, vasp_name: str) -> Sequence:
    jwt_token = _create_jwt(access_key, secret_key)
    url = "https://api.upbit.com/v1/withdraws/coin_addresses"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers).json()
    if not response:
        raise ValueError("There is no withdrawal address.")

    address_info = [{k: v for k, v in item.items() 
                    if k in ['withdraw_address', 'net_type', 'exchange_name']} 
                    for item in response if item.get('currency') == currency 
                    and item.get('net_type') == net_type 
                    and item.get('exchange_name') == vasp_name]
    
    if not address_info:
        raise ValueError("There is no withdrawal address for {currency}.".format(currency=currency))
    return address_info
```

[block:html]
{
  "html": "  <div class=\"callout-section\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i>  네트워크 타입(net_type)과 네트워크 이름(network_name)\n      </div>\n      네트워크 타입(net_type)은 디지털 자산 입출금시 실제 자산이 이동되는 블록체인 네트워크(대상 체인)를 지정하기 위한 식별자 필드(예: BTC)입니다. 디지털 자산 출금 시 필수 파라미터로, 정상적인 입출금 진행을 위해 정확한 식별자 값을 사용해야 합니다.\n      디지털 자산 출금 API 호출 시 사전에 출금 허용 주소 목록 조회 API를 호출한 뒤 응답으로부터 정확한 네트워크 타입 값을 참조하여 사용하시기 바랍니다. \n      <br><br>\n      네트워크 이름(network_name)은 블록체인 네트워크의 전체 이름(예: Bitcoin)을 나타내는 필드로서, 사람이 인식할 수 있는 정보이며 식별자로 사용할 수 없습니다. 서비스 UI 등에서 블록체인 네트워크를 표현하는 용도로 사용할 수 있습니다.\n  </div>"
}
[/block]

<br />

함수 실행 시 사용자가 등록한 출금 주소의 정보를 반환합니다.

```json
[
  {
    "net_type": "TRX",
    "withdraw_address": "<wtihdrawal_address>",
    "exchange_name": "<거래소 이름 A>"
  }
]

ValueError: Unsupported currency USDTz.

ValueError: There is no withdrawal address. Please register your withdrawal address.
```

<br />

## 입출금 서비스 상태 확인

[입출금 서비스 상태 조회](https://docs.upbit.com/kr/reference/get-service-status) API로 현재 특정 자산의 출금 가능 여부를 확인하는 함수를 구현합니다. 사용자가 지정한 디지털 자산의 네트워크별 입출금 지갑 지원 상태를 확인할 수 있습니다.

```python
def check_withdrawal_status(currency: str, net_type: str) -> str: 
    jwt_token = _create_jwt(access_key, secret_key)
    url = "https://api.upbit.com/v1/status/wallet"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers).json()
    wallets = [item for item in response if item.get('currency') == currency]
    print("{wallets}\n".format(wallets=wallets))
    wallet = next((item for item in wallets if item.get('net_type') == net_type), None)
    if wallet is None:
        raise ValueError("There is no withdrawal address for {currency}.".format(currency=currency))
    print("The {currency}-{net_type} wallet status is {wallet_state}.".format(currency=currency, net_type=net_type, wallet_state=wallet.get('wallet_state')))   
    return wallet.get('wallet_state')
```

<br />

함수 실행 시 디지털 자산이 배포된 블록체인 네트워크의 현황을 확인할 수 있으며 `wallet_state` 필드를 통해 해당 디지털 자산의 입출금 가능 여부를 판단합니다. USDT와 같이 멀티체인 환경의 디지털 자산의 경우, 업비트가 지원하는 모든 블록체인 네트워크의 입출금 지갑의 상태를 확인할 수 있습니다.

* 단일 네트워크의 입출금 지갑 상태

```json
{
  "currency": "BTC",
  "wallet_state": "working",
  "block_state": "normal",
  "block_height": 908692,
  "block_updated_at": "2025-08-05T07:38:57.889+00:00",
  "block_elapsed_minutes": 10,
  "net_type": "BTC",
  "network_name": "Bitcoin"
}

```

* 멀티체인 네트워크의 입출금 지갑 상태

```json
[
  {
    "currency": "USDT",
    "wallet_state": "working",
    "block_state": "normal",
    "block_height": 23089140,
    "block_updated_at": "2025-08-07T12:20:00.928+00:00",
    "block_elapsed_minutes": 1,
    "net_type": "ETH",
    "network_name": "Ethereum"
  },
  {
    "currency": "USDT",
    "wallet_state": "working",
    "block_state": "normal",
    "block_height": 74636780,
    "block_updated_at": "2025-08-07T12:21:51.663+00:00",
    "block_elapsed_minutes": 0,
    "net_type": "TRX",
    "network_name": "Tron"
  },
  {
    "currency": "USDT",
    "wallet_state": "working",
    "block_state": "normal",
    "block_height": 3191769000,
    "block_updated_at": "2025-08-07T12:21:50.948+00:00",
    "block_elapsed_minutes": 0,
    "net_type": "APT",
    "network_name": "Aptos"
  }
]

```

<br />

## 디지털 자산 출금

디지털 자산 출금을 요청하는 함수를 구현합니다. 출금할 디지털 자산의 통화, 네트워크 타입, 출금 수량, 출금 주소, 2차 주소(해당 네트워크가 2차 주소를 지원하는 경우)를 입력한 뒤 [디지털 자산 출금 요청](https://docs.upbit.com/kr/reference/withdraw) API를 호출하여 출금을 요청합니다.

```python
def withdraw_digital_asset(
        currency: str, 
        net_type: str, 
        amount: str, 
        address: str, 
        secondary_address: Optional[str] = None, 
        ) -> str:
    params = {
        "currency": currency,
        "net_type": net_type,
        "amount": amount,
        "address": address,
        "transaction_type": "default"
    }

    if secondary_address:
        params["secondary_address"] = secondary_address  
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/withdraws/coin"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=params).json()
    uuid = response.get('uuid')
    if uuid is None:
        raise ValueError(f"Please check the withdrawal issue. {response}")
    else:
        return uuid
```

<br />

성공적으로 출금이 요청된 경우 출금의 uuid가 반환됩니다.

```json
{
  "type": "withdraw",
  "uuid": "<withdrawal uuid>",
  "currency": "USDT",
  "net_type": "TRX",
  "txid": None,
  "state": "WAITING",
  "created_at": "2025-07-17T13:53:31+09:00",
  "done_at": None,
  "amount": "13.0",
  "fee": "0.0",
  "transaction_type": "default",
  "is_cancelable": False
}
```

<br />

## 출금 상태 조회

출금 UUID로 [단일 출금 조회](https://docs.upbit.com/kr/reference/get-withdrawal) API를 호출하여 출금 내역 및 결과를 조회할 수 있습니다.

```python
def get_withdrawal_state(uuid: str) -> Mapping:
    params = {
        "uuid": uuid
    }   
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/withdraw"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    return response
```

<br />

출금 정보의 state 필드의 값이 DONE 인 경우 출금이 완료됨을 의미합니다. 출금 상태에 따라 다음과 같은 값을 확인할 수 있습니다.

* `WAITING`: 출금 대기 중인 상태를 의미합니다.
* `PROCESSING`: 출금 진행 중인 상태를 의미합니다.
* `DONE`: 출금이 완료된 것을 의미합니다.
* `FAILED`: 출금에 실패한 것을 의미합니다.
* `CANCELLED`: 출금이 취소된 것을 의미합니다.
* `REJECTED`: 출금이 거절된 것을 의미합니다.

블록체인 트랜잭션 처리로 인해 실제 출금까지는 일정 시간이 소요될 수 있습니다. 출금 요청 이후 단일 출금 조회시 state가 변경되지 않은 경우 일정 시간 이후 재시도해주세요.

```json
{
  "type": "withdraw",
  "uuid": "<your_withdraw_uuid>",
  "currency": "USDT",
  "net_type": "TRX",
  "txid": "<your_withdraw_txid>",
  "state": "DONE",
  "created_at": "2025-07-17T13:53:31+09:00",
  "done_at": "2025-07-17T13:56:00+09:00",
  "amount": "13.0",
  "fee": "0.0",
  "transaction_type": "default",
  "is_cancelable": False
}

```

<br />

## 전체 코드

디지털 자산을 출금하는 전체 코드는 다음과 같습니다.

```python
from urllib.parse import unquote, urlencode
from typing import Any, Optional
from collections.abc import Mapping, Sequence
import hashlib
import uuid
import jwt # PyJWT
import requests
from decimal import Decimal, ROUND_DOWN

access_key = "<YOUR_ACCESS_KEY>"
secret_key = "<YOUR_SECRET_KEY>"

# Auth Token을 생성하는 로직은 여기에 추가해 주시기 바랍니다.

def get_withdrawal_address(currency: str, net_type: str, vasp_name: str) -> Sequence:
    jwt_token = _create_jwt(access_key, secret_key)
    url = "https://api.upbit.com/v1/withdraws/coin_addresses"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers).json()
    if not response:
        raise ValueError("There is no withdrawal address.")

    address_info = [{k: v for k, v in item.items() 
                    if k in ['withdraw_address', 'net_type', 'exchange_name']} 
                    for item in response if item.get('currency') == currency 
                    and item.get('net_type') == net_type 
                    and item.get('exchange_name') == vasp_name]
    
    if not address_info:
        raise ValueError("There is no withdrawal address for {currency}.".format(currency=currency))
    return address_info

def check_withdrawal_status(currency: str, net_type: str) -> str: 
    jwt_token = _create_jwt(access_key, secret_key)
    url = "https://api.upbit.com/v1/status/wallet"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers).json()
    wallets = [item for item in response if item.get('currency') == currency]
    print("{wallets}\n".format(wallets=wallets))
    wallet = next((item for item in wallets if item.get('net_type') == net_type), None)
    if wallet is None:
        raise ValueError("There is no withdrawal address for {currency}.".format(currency=currency))
    print("The {currency}-{net_type} wallet status is {wallet_state}.".format(currency=currency, net_type=net_type, wallet_state=wallet.get('wallet_state')))   
    return wallet.get('wallet_state')


def withdraw_digital_asset(
        currency: str, 
        net_type: str, 
        amount: str, 
        address: str, 
        secondary_address: Optional[str] = None, 
        ) -> str:
    params = {
        "currency": currency,
        "net_type": net_type,
        "amount": amount,
        "address": address,
        "transaction_type": "default"
    }

    if secondary_address:
        params["secondary_address"] = secondary_address  
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/withdraws/coin"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.post(url, headers=headers, json=params).json()
    uuid = response.get('uuid')
    if uuid is None:
        raise ValueError(f"Please check the withdrawal issue. {response}")
    else:
        return uuid

def get_withdrawal_state(uuid: str) -> Mapping:
    params = {
        "uuid": uuid
    }   
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    url = "https://api.upbit.com/v1/withdraw"
    headers = {
        "Authorization": "Bearer {jwt_token}".format(jwt_token=jwt_token),
        "Content-Type": "application/json"
    }
    response = requests.get(url, headers=headers, params=params).json()
    return response

# 주석을 해제하고 실행할 경우 실제 디지털 자산 출금 프로세스가 실행될 수 있습니다. 실행 전 다시 한 번 확인해 주시기 바랍니다.
# if __name__ == "__main__":  
    # currency = "<Enter_your_currency>"
    # net_type = "<Enter_your_net_type>"
    # vasp_name = "<Enter_your_vasp_name>"
    # amount = "13.241"
    # amount = str(Decimal(amount).quantize(Decimal("1e-8"), rounding=ROUND_DOWN))

    # withdraw_addresses = get_withdrawal_address(currency, net_type, vasp_name)
    # withdrawal_status = check_withdrawal_status(currency, net_type)
    # if len(withdraw_addresses) == 0:
    #     raise ValueError("There is no withdrawal address for {currency}.".format(currency=currency))
    # if withdrawal_status != "working":
    #     raise ValueError("The withdrawal is not working for {withdrawal_status}.".format(withdrawal_status=withdrawal_status))

#     for item in withdraw_addresses:
#         withdraw_address = item['withdraw_address']
#         response = withdraw_digital_asset(currency, net_type, amount, withdraw_address)
#         deposit_uuid = response
#         withdrawal_state = get_withdrawal_state(deposit_uuid)
#         print(withdrawal_state)
```

<br />