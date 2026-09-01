---
updatedAt: 2026-07-07T15:36:48.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [Python] 디지털 자산 출금

```python Python
from urllib.parse import unquote, urlencode
from typing import Any, Optional
from collections.abc import Mapping, Sequence
import hashlib
import uuid
import jwt # PyJWT
import requests
from decimal import Decimal, ROUND_DOWN

# 본 코드를 실행하기 위해서는 JWT를 생성해 API를 호출해야 합니다.
# JWT 생성은 인증 토큰(JWT) 생성 가이드를 확인해 주시기 바랍니다.

def get_withdrawal_address(currency: str, net_type: str, vasp_name: str) -> Sequence:
    """
    업비트에 등록된 출금 주소 조회  
    """
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
    """
    특정 currency의 입출금 가능 여부 확인
    """
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
    """
    디지털 자산 출금
    """
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
    """
    출금 상태 조회
    """
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

```json Response Example
{
  "type": "deposit",
  "uuid": "<your_withdrawal_uuid>",
  "currency": "USDT",
  "net_type": "TRX",
  "txid": "<your_withdrawal_txid>",
  "state": "ACCEPTED",
  "created_at": "2025-07-04T15:00:02+09:00",
  "done_at": "2025-07-04T15:00:48+09:00",
  "amount": "2500.363189",
  "fee": "0.0",
  "transaction_type": "default"
}

```

# 디지털 자산 출금 주소 등록



본 가이드를 원활하게 진행하기 위해서는 사전에 출금 안심차단 해지 및 출금 주소가 등록되어 있어야 합니다. 아래 링크를 통해 출금 안심차단 해지 및 디지털 출금 주소 관리 페이지로 이동하실 수 있으며, 해당 페이지에서 출금 대상 거래소 또는 개인 지갑 주소를 등록하실 수 있습니다.

- [출금 안심해지 안내 공지](https://docs.upbit.com/kr/changelog/withdrawal-safety-lock)
- [디지털 출금 주소 관리 페이지 바로가기](https://upbit.com/mypage/open_api_management?tab=fund_source)

# 유틸 라이브러리 Import

<!-- python@1-8 -->

기능 구현을 위해 필요한 모듈을 import 합니다. 별도의 설치가 필요한 모듈의 경우  `pip install <module name>` 명령어를 실행해 설치할 수 있습니다.

# 인증 토큰 생성

<!-- python@10-11 -->

본 코드를 실행하기 위해서는 JWT를 생성해 API를 호출해야 합니다.

JWT 생성은 [인증 토큰(JWT) 생성 가이드](python-인증-토큰jwt-생성)를 확인해 주시기 바랍니다.

# 출금 허용 주소 조회

<!-- python@13-35 -->

사용자가 등록한 출금 허용 주소 목록 중 특정 디지털 자산에 관련된 주소 목록을 반환하는 함수입니다. 

사용자가 출금 주소 등록 시 입력한 출금 주소와 네트워크 타입, 해당 주소를 발급한 거래소 이름이 반환됩니다.

# 출금 가능 여부 확인

<!-- python@37-55 -->

사용자가 입력한 디지털 자산의 출금 가능 여부를 확인하는 함수입니다. 입출금 지갑의 정보를 반환합니다. 다양한 블록체인에서 동일한 `currency`를 지원하는 경우, 관련된 모든 입출금 지갑의 상태를 확인할 수 있으며 사용자가 입력한 `net_type`으로 정확한 블록체인 네트워크를 선택해 입출금 지갑의 정보를 반환합니다.

# 디지털 자산 출금 요청

<!-- python@56-88 -->

디지털 자산 출금을 요청하는 함수입니다. 출금할 디지털 자산의 통화, 네트워크 타입, 출금 수량, 출금 주소, 2차 주소(해당 네트워크가 2차 주소를 지원하는 경우)를 입력하여 디지털 자산 출금을 신청합니다. 출금 요청이 성공적으로 완료된 경우 해당 출금 건에 대한 `UUID`를 반환합니다.

# 출금 상태 확인

<!-- python@90-105 -->

출금 요청 시 반환된 `UUID`로 해당 출금 처리 상태를 조회하기 위한 함수입니다. 응답의 `state` 필드를 확인합니다.

# 디지털 자산 출금 진행

<!-- python@107-127 -->

앞서 정의한 함수들을 사용해 디지털 자산 출금을 실행하는 함수입니다.

사용자가 입력한 디지털 자산 통화와 거래소 이름을 사용해 출금 주소를 특정하고 입출금 지갑의 상태를 확인합니다. 

출금이 가능한 경우 사용자가 입력한 `amount` 만큼을 해당 주소로 출금한 뒤 출금 상태를 확인합니다.

주석 처리된 코드의 주석을 삭제하고 호출 시 실제 출금이 진행될 수 있으므로 실행 전 반드시 확인해 주시기 바랍니다.