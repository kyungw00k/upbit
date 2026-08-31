---
updatedAt: 2026-01-07T10:53:29.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [Python] 인증 토큰(JWT) 생성

```python Python
from urllib.parse import unquote, urlencode
from typing import Any
from collections.abc import Mapping
import hashlib
import uuid
import jwt # PyJWT
import requests

def _create_jwt(access_key: str, secret_key: str, query_string: str = None) -> str:
    """
    JWT 토큰 생성
    """
    payload = {"access_key": access_key, "nonce": str(uuid.uuid4())}

    if query_string is not None:
        query_hash = hashlib.sha512(query_string.encode("utf-8")).hexdigest()
        payload["query_hash"] = query_hash
        payload["query_hash_alg"] = "SHA512"

    token = jwt.encode(payload, secret_key, algorithm="HS512")
    return token if isinstance(token, str) else token.decode('utf-8')

def _build_query_string(params: Mapping) -> str:
    """
    Dictionary 파라미터를 쿼리 문자열 형식으로 변환
    """
    return unquote(urlencode(params, doseq=True))

if __name__ == "__main__":
    base_url = "https://api.upbit.com"
    access_key = "<YOUR_ACCESS_KEY>"
    secret_key = "<YOUR_SECRET_KEY>" # 실제로는 안전한 방식으로 로드하거나 주입하세요.
    
    # 파라미터가 없는 요청 예시
    jwt_token = _create_jwt(access_key, secret_key)
    headers = {"Authorization": f"Bearer {jwt_token}"}
                
    response = requests.get(f"{base_url}/v1/accounts", headers=headers).json()
        
    print(response)
    
    # 파라미터가 있는 GET 요청 예시
    params = {
        "market": "KRW-BTC",
        "states[]": ["wait", "watch"],
        "limit": 10
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    headers = {"Authorization": f"Bearer {jwt_token}"}
        
    response = requests.get(f"{base_url}/v1/orders/open?{query_string}", headers=headers).json()
    print(response)    
    
    # POST 요청 예시
    order_data = {
        "market": "KRW-BTC",
        "side": "bid",
        "volume": "0.001",
        "price": "50000000",
        "ord_type": "limit"
    }
        
    query_string = _build_query_string(order_data)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    headers = {
        "Authorization": f"Bearer {jwt_token}",
        "Content-Type": "application/json"
    }
    
    # 아래 주석처리된 부분 실행시 실제 주문이 발생하므로 실행 전 반드시 확인하세요.
    # response = requests.post(f"{base_url}/v1/orders", json=order_data, headers=headers).json()
    # print(response)   
```

```json Response Example
# 파라미터가 없는 GET 요청 응답 예시
[
  {
    "currency": "BTC",
    "balance": "0.00000104",
    "locked": "0",
    "avg_buy_price": "160210789.42247014",
    "avg_buy_price_modified": False,
    "unit_currency": "KRW"
  },
...
]
  
  # 파라미터가 있는 GET 요청 응답 예시

{
  "uuid": "<your order UUID>",
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


  # 바디 파라미터가 있는 POST 요청

[
  {
    "uuid": "cad93adb-6e05-4426-ba28-b1418b7f7809",
    "side": "bid",
    "ord_type": "limit",
    "price": "50000000",
    "state": "wait",
    "market": "KRW-BTC",
    "created_at": "2025-07-28T13:42:07+09:00",
    "volume": "0.001",
    "remaining_volume": "0.001",
    "prevented_volume": "0",
    "reserved_fee": "25.00032921",
    "remaining_fee": "25.00032921",
    "paid_fee": "0",
    "locked": "500025.65874921",
    "prevented_locked": "0",
    "executed_volume": "0",
    "executed_funds": "0",
    "trades_count": 0
  }
]

```

# 유틸 라이브러리 Import

<!-- python@1-7 -->

인증 토큰을 생성하기 위해 필요한 모듈을 import 합니다. 별도의 설치가 필요한 모듈의 경우  `pip install <module name>` 명령어를 실행해 설치할 수 있습니다.

# JWT 생성

<!-- python@9-21 -->

JWT는 payload를 사용해 생성합니다. 기본 payload는 Access Key와 nonce를 가진 객체로 사용자의 파라미터 입력 여부에 따라 payload의 값이 달라집니다. 

1. 사용자가 파라미터를 입력하지 않은 경우, 기본 payload를 사용해 JWT를 생성합니다.
2. 사용자가 파라미터를 입력한 경우, 파라미터를 쿼리 스트링으로 인코딩 한 후 이를 해시합니다. 해시의 결과로 반환받은 값과 해시 알고리즘 타입을 기본 payload에 추가하고 이 payload를 사용해 JWT를 생성합니다.

# 파라미터 인코딩 설정

<!-- python@23-27 -->

사용자가 입력한 파마리터를 URL 인코딩 처리하기 위한 함수입니다. 파라미터의 값 중 배열이 있는 경우, [] 문자열은 인코딩에서 제외해 업비트 API 요청에 적합한 형식으로 인코딩합니다.

# API 호출로 JWT 동작 확인

<!-- python@29-73 -->

생성한 JWT가 정상적으로 동작하는지 확인할 수 있는 예시 코드 입니다. 다음 3가지 요청을 통해 JWT의 동작을 확인할 수 있습니다. 

1.  파라미터 없는 GET 요청
2. Query 파라미터를 입력하는 GET 요청
3. Body 파라미터를 입력하는 POST 요청

위 3가지 API 호출을 실행해 볼 수 있습니다. 

단, POST 요청은 주석을 해제하고 실행해야 합니다. 또한 POST 요청 시 실제 주문이 생성될 수 있으므로 실행하기 전 반드시 확인 후 실행하시기 바랍니다.