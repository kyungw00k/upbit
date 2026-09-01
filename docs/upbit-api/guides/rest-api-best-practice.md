---
updatedAt: 2026-05-04T05:48:12.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# REST API 연동 Best Practice

업비트 REST API 연동 구현을 위한 가이드라인 문서로서 인증, 요청 수 제한, 에러 처리 등 실제 연동 시 참고해야 하는 구현 요구사항을 안내합니다.

## 시작하기

개발자 센터의 튜토리얼 및 API Reference에 포함된 예제 코드는 API의 기능과 호출 방식을 빠르게 이해할 수 있도록, API 호출 위주의 단순한 코드로 작성되어 있습니다. 이 가이드는 단순한 호출 예제를 안정적이고 실전적인 전략 구현 단계로 발전시키기 위해 필요한 모범 사례를 제시합니다.

<br />

## 업비트 API 연동 시 고려해야 할 사항

<br />

### 인증

인증이 필요한 Exchange API 호출 시, 요청의 파라미터 또는 본문(Body)을 기반으로 업비트 API  [인증](https://docs.upbit.com/kr/reference/auth) 문서의 안내에 따라 유효한 인증 토큰을 생성하여 요청 헤더에 포함해야 합니다. 실제 구현 시 API의 Path를 기반으로 적절한 인증 핸들러를 구현하여 적용할 수 있습니다.

<br />

### 정상 응답 구분 및 에러 처리

REST API 응답 수신 시 HTTP 상태 코드를 기반으로 정상 응답(2xx)과 에러 응답(4xx, 5xx)을 구분한 뒤, 아래와 같이 적절한 응답 객체 Parsing을 수행해야 합니다.  [REST API 사용 및 에러 안내](https://docs.upbit.com/kr/reference/rest-api-guide)  문서에서 업비트 REST API에서 반환되는 HTTP 상태코드별 에러 코드 목록을 확인할 수 있습니다.

* 정상 응답이 반환된 경우 API Reference를 참고하여 응답 객체로부터 필요한 정보를 참조합니다.
* 에러 응답이 반환된 경우 error.name과 error.message를 기반으로 원인을 식별하고, 재시도 여부나 대응 방안을 결정할 수 있어야 합니다. 일반적으로 코드의 수정 또는 잔고 갱신과 같이 별도의 사용자 행동 없이 해결되기 어려운 4xx 오류는 즉시 호출을 중단하고 원인을 식별해야 하며, 5xx와 같은 일시적 네트워크 장애 상황에서는 API 호출 목적에 따라 적절히 재시도해볼 수 있습니다.

<br />

### 요청 수 제한(Rate Limit) 정책 준수

업비트 REST API는 초당 최대 [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits) 정책을 적용하고 있습니다. 반복적으로 제한을 초과하여 요청하는 경우, 임시 또는 영구적으로 REST API 호출이 제한될 수 있으므로 적절한 Throttling을 통해 호출 속도를 조정해야 합니다.

* 허용 정책을 초과한 API 호출로 429 Too Many Requests 오류가 반환되는 경우, 즉시 해당 그룹에 속한 API 호출을 중단하고 적절한 시간 이후 재시도해야 합니다.
* 모든 API 응답의 헤더에는 해당 API가 속한 Rate Limit 그룹과 잔여 요청 수가 반환되므로, 해당 정보를 참고하여 호출 속도를 조절하도록 구현할 수 있습니다.

구현 가이드 라인은 본 문서 하단 예제를 참고해주세요.

<br />

### 보안 및 운영을 고려한 Logging

API 기반 매매 시스템 구현 시 호출 이력과 결과를 추적할 수 있는 구조화재시도해야된 로그를 남기는 것을 권장합니다. 로그에는 요청 시각, API Endpoint, 쿼리 파라미터(민감 정보 제외), 응답 코드, 주요 응답 헤더(Remaining-Req 값 포함), 응답 지연 시간 등을 기록할 수 있습니다. 로그를 통해 주문 및 입출금 이력 확인과 문제 발생 시 신속한 원인 분석이 가능합니다. **단, 인증 키나 사용자 개인정보 등 민감한 데이터는 로그에 기록하지 않도록 주의**해야 하며, 로그는 적절한 보관 기간과 접근 제어 정책 하에 관리해야 합니다.

<br />

## Best Practice - Python 예제로 알아보는 구현 가이드라인

튜토리얼 및 예제 코드에서 볼 수 있는 간단한 호출 예제로 시작하여 실 연동 시 고려해야 하는 사항들을 구현에 적용하는 과정을 Python 예제를 통해 순차적으로 안내합니다.

[block:html]
{
  "html": "<div class=\"callout-section\">\n  <div class=\"callout-title\">\n    <i class=\"fa-solid fa-circle-exclamation\"></i> 가이드라인 예제 코드 \n    </div>\n    <li>본 가이드라인에서 사용된 Python 코드는 보다 다양한 버전의 Python 환경을 지원하기 위해 타입 힌트와 최신 문법을 적극적으로 사용하기 보다는 <b>호환성을 고려하여 구현</b>되었습니다. </li>\n    <li>requests 라이브러리에서 동기 방식의 HTTP 요청 방식에 한정한 예제로, <b>절대적인 구현 제약사항이 아닌 하나의 구현 예시</b>로 작성되었음을 사용 시 참고 부탁드립니다. </li>\n    <li>더 간결한 예제 코드 작성을 위해 docstring 및 주석은 최소화하였으며, <b>주요 코드 설명은 코드 상단 문서 영역을 참조</b> 부탁드립니다.</li>\n</div>"
}
[/block]

<br />

### 기본 요청 예제

[개발 환경 설정 가이드](https://docs.upbit.com/kr/docs/dev-environment)에서 안내한 `requests` 라이브러리를 기반으로 기본적인 GET API를 호출하기 위한 몇몇 유틸리티 코드를 다음과 같이 작성할 수 있습니다.

* `UpbitClient` : 업비트 API 요청을 처리하기 위한 기능 클래스 입니다. 간단한 예시로 GET API 만을 호출할 수 있는 `request_get` 메서드를 제공합니다.
* `_build_query_string`:  요청 파라미터 리스트로부터 인코딩 된 쿼리 문자열을 생성하는 메서드입니다.
* `_build_url`: API Path와 쿼리 문자열을 조합하여 전체 URL을 생성하는 메서드입니다.
* main 함수에서는 UpbitClient를 초기화하고, 업비트에서 지원하는 페어 목록을 조회하기 위한 GET API를 호출합니다. 간단한 예시로서, 응답을 따로 Parsing하거나 처리하지 않고 출력하는 예제로 작성되었습니다.

```python
from collections.abc import Mapping
from urllib.parse import unquote, urlencode
import requests

class UpbitClient(object):
    def __init__(self, base_url):
        # type: (str) -> None
        self.base_url = base_url.rstrip("/")

    def _build_url(self, path, query_string=""):
        # type: (str, str) -> str
        url = "{0}/{1}".format(self.base_url, path.lstrip("/"))
        if query_string:
            url += "?{0}".format(query_string)
        return url

    def _build_query_string(self, params):
        # type: (object) -> str
        data = params if isinstance(params, Mapping) else params
        return unquote(urlencode(data, doseq=True))

    def request_get(self, path, params=None):
        # type: (str, object) -> object
        query_str = self._build_query_string(params) if params is not None else ""
        url = self._build_url(path, query_str)
        resp = requests.get(url)
        try:
            return resp.json()
        except ValueError:
            return resp.text

if __name__ == "__main__":
    client = UpbitClient("https://api.upbit.com")

    # Example:List Pairs
    data = client.request_get("/v1/market/all")
    print(data)

```

<br />

### 인증 구현

request\_get 메서드에서 입력받은 API Path를 기준으로 인증 필요 여부를 판단하고, 필요에 따라 인증 헤더를 추가하도록 아래와 같이 코드를 개선할 수 있습니다.

* UpbitClient 초기화 시 API Key의 Access Key와 Secret Key를 초기화합니다. (12-13라인)
* `_requires_auth`: 인증이 필요하지 않은 Public API의 API Path Prefix를 기반으로 간단하게 인증 여부를 판단하도록 구현하였습니다. 실제 구현 시 API별 정의를 포함하는 별도의 클래스로 분리하여 인증 여부를 해당 클래스에서 제공하도록 구현할 수 있습니다. (14라인, 38-40라인)
* `_create_jwt_token`: Access Key와 Secret Key, 그리고 쿼리 문자열을 기반으로 JWT 인증 토큰을 생성합니다. (28-36라인)
* 이제 `request_get` 메서드 내부에서 API 요청을 전송하기 전 필요에 따라 인증 토큰을 생성하고 인증 헤더를 추가합니다. (47-53라인)

```python
from collections.abc import Mapping
from urllib.parse import unquote, urlencode
import hashlib
import uuid
import jwt  # PyJWT
import requests

class UpbitClient(object):
    def __init__(self, base_url, access_key, secret_key):
        # type: (str, str, str) -> None
        self.base_url = base_url.rstrip("/")
        self.access_key = access_key
        self.secret_key = secret_key
        self.public_prefixes = ("/v1/market", "/v1/ticker", "/v1/trades", "/v1/candles", "/v1/orderbook")

    def _build_url(self, path, query_string=""):
        # type: (str, str) -> str
        url = "{0}/{1}".format(self.base_url, path.lstrip("/"))
        if query_string:
            url += "?{0}".format(query_string)
        return url

    def _build_query_string(self, params):
        # type: (object) -> str
        data = params if isinstance(params, Mapping) else params
        return unquote(urlencode(data, doseq=True))
    
    def _create_jwt_token(self, query_string=None):
        # type: (str) -> str
        payload = {"access_key": self.access_key, "nonce": str(uuid.uuid4())}
        if query_string:
            query_hash = hashlib.sha512(query_string.encode("utf-8")).hexdigest()
            payload["query_hash"] = query_hash
            payload["query_hash_alg"] = "SHA512"
        token = jwt.encode(payload, self.secret_key, algorithm="HS512")
        return token if isinstance(token, str) else token.decode("utf-8")

    def _requires_auth(self, path):
        # type: (str) -> bool
        return not any(path.startswith(pub) for pub in self.public_prefixes)

    def request_get(self, path, params=None):
        # type: (str, object) -> object
        query_str = self._build_query_string(params) if params is not None else ""
        url = self._build_url(path, query_str)
        
        headers = {}
        if self._requires_auth(path):
            if not self.access_key or not self.secret_key:
                raise ValueError("인증이 필요한 API입니다. access_key와 secret_key를 설정하세요.")
            headers["Authorization"] = "Bearer {0}".format(self._create_jwt_token(query_str))
            
        resp = requests.get(url, headers=headers)
        try:
            return resp.json()
        except ValueError:
            return resp.text

if __name__ == "__main__":
    client = UpbitClient(
        "https://api.upbit.com",
        access_key="YOUR_ACCESS_KEY",
        secret_key="YOUR_SECRET_KEY"
    )

    # Example:List Pairs
    data = client.request_get("/v1/market/all")
    print(data)

    # Example:List Open Orders
    params = [("market", "KRW-BTC"), ("states[]", "wait"), ("states[]", "watch")]
    data = client.request_get("/v1/orders/open", params=params)
    print(data)
```

<br />

### 정상 응답 구분 및 에러 처리

응답의 HTTP 상태코드를 기반으로 정상 응답과 에러 응답을 구분하고, 에러 응답인 경우 상태코드를 포함한 에러 객체를 반환하도록 개선할 수 있습니다. 실제 서비스에서는 에러 응답이 반환된 경우 적절한 처리 로직을 구현할 수 있습니다. (54-79라인)

```python
from collections.abc import Mapping
from urllib.parse import unquote, urlencode
import hashlib
import uuid
import jwt  # PyJWT
import requests

class UpbitClient(object):
    def __init__(self, base_url, access_key, secret_key):
        # type: (str, str, str) -> None
        self.base_url = base_url.rstrip("/")
        self.access_key = access_key
        self.secret_key = secret_key
        self.public_prefixes = ("/v1/market", "/v1/ticker", "/v1/trades", "/v1/candles", "/v1/orderbook")

    def _build_url(self, path, query_string=""):
        # type: (str, str) -> str
        url = "{0}/{1}".format(self.base_url, path.lstrip("/"))
        if query_string:
            url += "?{0}".format(query_string)
        return url

    def _build_query_string(self, params):
        # type: (object) -> str
        data = params if isinstance(params, Mapping) else params
        return unquote(urlencode(data, doseq=True))
    
    def _create_jwt_token(self, query_string=None):
        # type: (str) -> str
        payload = {"access_key": self.access_key, "nonce": str(uuid.uuid4())}
        if query_string:
            query_hash = hashlib.sha512(query_string.encode("utf-8")).hexdigest()
            payload["query_hash"] = query_hash
            payload["query_hash_alg"] = "SHA512"
        token = jwt.encode(payload, self.secret_key, algorithm="HS512")
        return token if isinstance(token, str) else token.decode("utf-8")

    def _requires_auth(self, path):
        # type: (str) -> bool
        return not any(path.startswith(pub) for pub in self.public_prefixes)

    def request_get(self, path, params=None):
        # type: (str, object) -> object
        query_str = self._build_query_string(params) if params is not None else ""
        url = self._build_url(path, query_str)
        
        headers = {}
        if self._requires_auth(path):
            if not self.access_key or not self.secret_key:
                raise ValueError("인증이 필요한 API입니다. access_key와 secret_key를 설정하세요.")
            headers["Authorization"] = "Bearer {0}".format(self._create_jwt_token(query_str))
            
        resp = requests.get(url, headers=headers)
        if 200 <= resp.status_code < 300:
            try:
                return resp.json()
            except ValueError:
                return resp.text

        try:
            ej = resp.json()
            if isinstance(ej, dict) and "error" in ej:
                e = ej["error"]
                return {
                    "status_code": resp.status_code,
                    "name": e.get("name"),
                    "message": e.get("message")
                }
            return {
                "status_code": resp.status_code,
                "name": None,
                "message": ej
            }
        except ValueError:
            return {
                "status_code": resp.status_code,
                "name": None,
                "message": resp.text
            }

if __name__ == "__main__":
    client = UpbitClient(
        "https://api.upbit.com",
        access_key="YOUR_ACCESS_KEY",
        secret_key="YOUR_SECRET_KEY"
    )

    # Example:List Pairs
    data = client.request_get("/v1/market/all")
    print(data)

    # Example:List Open Orders
    params = [("market", "KRW-BTC"), ("states[]", "wait"), ("states[]", "watch")]
    data = client.request_get("/v1/orders/open", params=params)
    print(data)
```

<br />

### 요청 수 제한(Rate Limit) 관련 처리

요청 수 제한 처리를 위해 RateLimiter 클래스를 추가한 예제입니다. 요청 전 잔여 요청 가능 횟수를 확인 후 전송하거나 일정 시간 대기하고, 응답 수신 후 헤더의 잔여 요청 수 정보를 참고하여 갱신합니다. 추가된 세부 구현은 다음과 같습니다.

* 업비트 요청 수 제한 정책을 반영하여 고정된 초 단위의 window 내에서 최대 RPS(Requests Per Seconds)를 확인하는 Rate Limiter 클래스를 구현하였습니다. 특정 초 내에서 요청 가능 횟수(토큰)을 모두 소진하면, 다음 초 Window가 될 때까지 대기하는 간단한 모델입니다. Exchange API의 경우 요청 수 제한 정책이 토큰-버킷 모델로 동작하나, Quotation API를 아우르는 안정적인 호출 속도 제어를 위해 이와 같이 구현합니다. (7-68라인)
* [요청 수 제한](https://docs.upbit.com/kr/docs/rate-limits)rate-limits 문서에 명시된 바와 같이, 각 Rate Limit 그룹 별 최대 호출 수를 설정값(cfg)으로 사전에 정의하였습니다. (10-19라인)
* `acquire` 메서드는 실제 API 요청을 전송하기 전 잔여 호출 토큰이 있는지 확인하고 차감하는 동작을 수행합니다. 마지막 확인 시각이 이미 지난 window인 경우, 최대 요청 가능 수를 최대 호출 수로 충전한 뒤 차감합니다. 사용 가능한 토큰이 없는 경우 다음 window 시각이 될 때 까지 sleep을 통해 대기합니다. (26-45라인, 128라인)
* `update_from_header` 메서드는 응답 헤더로부터 현재 window 내에서의 잔여 요청 가능 횟수를 조회하여 잔여 토큰 수를 업데이트합니다. (47-62라인, 137라인)
* 응답 수신 시 상태 코드가 429인 경우, 해당 window 내에서 추가 요청이 발생하지 않도록 `mark_exhausted` 메서드를 호출하여 잔여 요청 수를 0으로 초기화합니다. (134-135라인)

아래와 같은 구현을 통해, 최대 요청 수를 초과한 request\_get 메서드 호출시에도 적절한 대기를 통해 429 Too Many Requests 응답 수신을 최소화 할 수 있습니다. 아래 예시는 간단한 Rate Limiter 모델로, 실제 구현 시에는 서비스 특성에 맞춰 백오프(backoff) 정책이나 지터(jitter) 를 적용하고, 요청이 임계치에 도달하기 전에 적절한 시간 버퍼를 두어 호출 간격을 조정해 보다 안정적인 운영 환경을 구현해보시기 바랍니다.

```python
from collections.abc import Mapping
from urllib.parse import unquote, urlencode
import time, hashlib, uuid
import jwt
import requests

class RateLimiter(object):
    def __init__(self):
        # group -> (capacity, window_sec)
        self.cfg = {
            "market": (10, 1),
            "ticker": (10, 1),
            "trades": (10, 1),
            "candles": (10, 1),
            "orderbook": (10, 1),
            "default": (30, 1),
            "order": (8, 1),
            "order-cancel-all": (1, 2),
        }
        # group -> (remaining, window_start_epoch)
        self.state = {}

    def _win_start(self, now_sec, win):
        return now_sec - (now_sec % win)

    def acquire(self, group):
        cap, win = self.cfg.get(group, (10, 1))
        now = time.time()
        now_sec = int(now)
        win_start = self._win_start(now_sec, win)
        remaining, cur_win_start = self.state.get(group, (cap, win_start))

        if cur_win_start != win_start:
            remaining, cur_win_start = cap, win_start

        if remaining <= 0:
            sleep_for = (cur_win_start + win) - now + 0.01
            if sleep_for > 0:
                time.sleep(sleep_for)
            now = time.time()
            now_sec = int(now)
            cur_win_start = self._win_start(now_sec, win)
            remaining = cap

        self.state[group] = (remaining - 1, cur_win_start)

    def update_from_header(self, header_value):
        # Remaining-Req: "group=default; min=1800; sec=29"
        if not header_value:
            return
        g, sec = "default", None
        try:
            for p in [s.strip() for s in header_value.split(";")]:
                if p.startswith("group="): g = p.split("=", 1)[1].strip()
                elif p.startswith("sec="): sec = int(p.split("=", 1)[1].strip())
        except Exception:
            return
        if g in self.cfg and sec is not None:
            cap, win = self.cfg[g]
            now_sec = int(time.time())
            win_start = self._win_start(now_sec, win)
            self.state[g] = (min(cap, sec), win_start)

    def mark_exhausted(self, group):
        cap, win = self.cfg.get(group, (10, 1))
        now_sec = int(time.time())
        win_start = self._win_start(now_sec, win)
        self.state[group] = (0, win_start)

class UpbitClient(object):
    def __init__(self, base_url, access_key=None, secret_key=None, limiter=None):
        # type: (str, str, str, RateLimiter) -> None
        self.base_url = base_url.rstrip("/")
        self.access_key = access_key
        self.secret_key = secret_key
        self.limiter = limiter or RateLimiter()
        self.public_prefixes = ("/v1/market", "/v1/ticker", "/v1/trades", "/v1/candles", "/v1/orderbook")

    def _build_url(self, path, query_string=""):
        # type: (str, str) -> str
        url = "{0}/{1}".format(self.base_url, path.lstrip("/"))
        if query_string:
            url += "?{0}".format(query_string)
        return url

    def _build_query_string(self, params):
        # type: (object) -> str
        data = params if isinstance(params, Mapping) else params
        return unquote(urlencode(data, doseq=True))
    
    def _create_jwt_token(self, query_string=None):
        # type: (str) -> str
        payload = {"access_key": self.access_key, "nonce": str(uuid.uuid4())}
        if query_string:
            query_hash = hashlib.sha512(query_string.encode("utf-8")).hexdigest()
            payload["query_hash"] = query_hash
            payload["query_hash_alg"] = "SHA512"
        token = jwt.encode(payload, self.secret_key, algorithm="HS512")
        return token if isinstance(token, str) else token.decode("utf-8")

    def _requires_auth(self, path):
        # type: (str) -> bool
        return not any(path.startswith(pub) for pub in self.public_prefixes)
    
    def _group_for(self, method, path):
        # type: (str, str) -> str
        if path.startswith("/v1/market"): return "market"
        if path.startswith("/v1/ticker"): return "ticker"
        if path.startswith("/v1/trades"): return "trades"
        if path.startswith("/v1/candles"): return "candles"
        if path.startswith("/v1/orderbook"): return "orderbook"
        if path.startswith("/v1/orders/open") and method.upper() == "DELETE": return "order-cancel-all"
        if path.startswith("/v1/orders") and method.upper() == "POST": return "order"
        return "default"

    def request_get(self, path, params=None):
        # type: (str, object) -> object
        query_str = self._build_query_string(params) if params is not None else ""
        url = self._build_url(path, query_str)
        
        headers = {}
        if self._requires_auth(path):
            if not self.access_key or not self.secret_key:
                raise ValueError("인증이 필요한 API입니다. access_key와 secret_key를 설정하세요.")
            headers["Authorization"] = "Bearer {0}".format(self._create_jwt_token(query_str))
            
        group = self._group_for("GET", path)
        self.limiter.acquire(group)

        resp = requests.get(url, headers=headers)
        
        if resp.status_code == 429:
            self.limiter.mark_exhausted(group)

        self.limiter.update_from_header(resp.headers.get("Remaining-Req"))


        if 200 <= resp.status_code < 300:
            try:
                return resp.json()
            except ValueError:
                return resp.text

        try:
            ej = resp.json()
            if isinstance(ej, dict) and "error" in ej:
                e = ej["error"]
                return {
                    "status_code": resp.status_code,
                    "name": e.get("name"),
                    "message": e.get("message")
                }
            return {
                "status_code": resp.status_code,
                "name": None,
                "message": ej
            }
        except ValueError:
            return {
                "status_code": resp.status_code,
                "name": None,
                "message": resp.text
            }

if __name__ == "__main__":
    client = UpbitClient(
        "https://api.upbit.com",
        access_key="YOUR_ACCESS_KEY",
        secret_key="YOUR_SECRET_KEY"
    )

    # Example:List Pairs
    data = client.request_get("/v1/market/all")
    print(data)

    # Example:List Open Orders
    params = [("market", "KRW-BTC"), ("states[]", "wait"), ("states[]", "watch")]
    data = client.request_get("/v1/orders/open", params=params)
    print(data)
```

<br />

### Logging 추가

마지막으로, Logger를 이용하여 사용중인 환경에 맞게 적절한 로깅을 추가할 수 있습니다. 실제 운영 환경에서는 Logger 설정 및 로깅 레벨은 프로젝트 단위로 설정하여 참고할 수 있으나, 예제 코드 작성을 위해 코드 상단부에 명시하는 형태로 작성했습니다. 주요 적용 사항은 다음과 같습니다.

* Logger를 정의하고, DEBUG로 로깅 레벨을 설정하였습니다. 이에 따라 DEBUG, INFO를 포함한 상위 로깅 레벨로 남겨진 모든 로그가 출력 또는 파일에 저장됩니다. (8-9라인)
* 로깅 시간과 레벨, 메세지를 출력하기 위한 간단한 로그 포맷을 정의하였습니다. (10-13라인)
* Rate Limiter 내부 및 요청 전후 필요한 곳에서 로깅하는 부분을 추가하였습니다. 응답의 상태 코드는 INFO 레벨로 기록합니다. Body 본문은 로그 사이즈를 최적화하기 위해 오류 상황에서만 INFO 레벨로 기록하도록 작성했습니다.  (46라인, 71라인, 78라인, 127-141라인, 157라인, 161라인)

```python
import logging
from collections.abc import Mapping
from urllib.parse import unquote, urlencode
import time, hashlib, uuid, json
import jwt
import requests

logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)  # INFO 
handler = logging.StreamHandler()
formatter = logging.Formatter('[%(asctime)s] [%(levelname)s] %(message)s')
handler.setFormatter(formatter)
logger.addHandler(handler)

class RateLimiter(object):
    def __init__(self):
        # group -> (capacity, window_sec)
        self.cfg = {
            "market": (10, 1),
            "ticker": (10, 1),
            "trades": (10, 1),
            "candles": (10, 1),
            "orderbook": (10, 1),
            "default": (30, 1),
            "order": (8, 1),
            "order-cancel-all": (1, 2),
        }
        # group -> (remaining, window_start_epoch)
        self.state = {}

    def _win_start(self, now_sec, win):
        return now_sec - (now_sec % win)

    def acquire(self, group):
        cap, win = self.cfg.get(group, (10, 1))
        now = time.time()
        now_sec = int(now)
        win_start = self._win_start(now_sec, win)
        remaining, cur_win_start = self.state.get(group, (cap, win_start))

        if cur_win_start != win_start:
            remaining, cur_win_start = cap, win_start

        if remaining <= 0:
            sleep_for = (cur_win_start + win) - now + 0.01
            logger.debug("RateLimiter: group=%s exhausted, sleeping %.3fs", group, sleep_for)
            if sleep_for > 0:
                time.sleep(sleep_for)
            now = time.time()
            now_sec = int(now)
            cur_win_start = self._win_start(now_sec, win)
            remaining = cap

        self.state[group] = (remaining - 1, cur_win_start)

    def update_from_header(self, header_value):
        # Remaining-Req: "group=default; min=1800; sec=29"
        if not header_value:
            return
        g, sec = "default", None
        try:
            for p in [s.strip() for s in header_value.split(";")]:
                if p.startswith("group="): g = p.split("=", 1)[1].strip()
                elif p.startswith("sec="): sec = int(p.split("=", 1)[1].strip())
        except Exception:
            return
        if g in self.cfg and sec is not None:
            cap, win = self.cfg[g]
            now_sec = int(time.time())
            win_start = self._win_start(now_sec, win)
            logger.debug("RateLimiter: update group=%s remaining=%s", g, sec)
            self.state[g] = (min(cap, sec), win_start)

    def mark_exhausted(self, group):
        cap, win = self.cfg.get(group, (10, 1))
        now_sec = int(time.time())
        win_start = self._win_start(now_sec, win)
        logger.warning("RateLimiter: mark exhausted for group=%s", group)
        self.state[group] = (0, win_start)

class UpbitClient(object):
    def __init__(self, base_url, access_key, secret_key, limiter=None):
        # type: (str, str, str, RateLimiter) -> None
        self.base_url = base_url.rstrip("/")
        self.access_key = access_key
        self.secret_key = secret_key
        self.limiter = limiter or RateLimiter()
        self.public_prefixes = ("/v1/market", "/v1/ticker", "/v1/trades", "/v1/candles", "/v1/orderbook")

    def _build_url(self, path, query_string=""):
        # type: (str, str) -> str
        url = "{0}/{1}".format(self.base_url, path.lstrip("/"))
        if query_string:
            url += "?{0}".format(query_string)
        return url

    def _build_query_string(self, params):
        # type: (object) -> str
        data = params if isinstance(params, Mapping) else params
        return unquote(urlencode(data, doseq=True))
    
    def _create_jwt_token(self, query_string=None):
        # type: (str) -> str
        payload = {"access_key": self.access_key, "nonce": str(uuid.uuid4())}
        if query_string:
            query_hash = hashlib.sha512(query_string.encode("utf-8")).hexdigest()
            payload["query_hash"] = query_hash
            payload["query_hash_alg"] = "SHA512"
        token = jwt.encode(payload, self.secret_key, algorithm="HS512")
        return token if isinstance(token, str) else token.decode("utf-8")

    def _requires_auth(self, path):
        # type: (str) -> bool
        return not any(path.startswith(pub) for pub in self.public_prefixes)
    
    def _group_for(self, method, path):
        # type: (str, str) -> str
        if path.startswith("/v1/market"): return "market"
        if path.startswith("/v1/ticker"): return "ticker"
        if path.startswith("/v1/trades"): return "trades"
        if path.startswith("/v1/candles"): return "candles"
        if path.startswith("/v1/orderbook"): return "orderbook"
        if path.startswith("/v1/orders/open") and method.upper() == "DELETE": return "order-cancel-all"
        if path.startswith("/v1/orders") and method.upper() == "POST": return "order"
        return "default"
    
    def _log_response(self, resp):
        status = resp.status_code
        logger.info("HTTP Response | status=%s", status)
        logger.debug("HTTP Headers: %s", dict(resp.headers))

        try:
            body_obj = resp.json()
            body_str = json.dumps(body_obj, ensure_ascii=False)
        except ValueError:
            body_str = resp.text

        if 200 <= status < 300:
            logger.debug("HTTP Body: %s", body_str)
        else:
            logger.info("HTTP Error Body: %s", body_str)

    def request_get(self, path, params=None):
        # type: (str, object) -> object
        query_str = self._build_query_string(params) if params is not None else ""
        url = self._build_url(path, query_str)
        
        headers = {}
        if self._requires_auth(path):
            if not self.access_key or not self.secret_key:
                raise ValueError("인증이 필요한 API입니다. access_key와 secret_key를 설정하세요.")
            headers["Authorization"] = "Bearer {0}".format(self._create_jwt_token(query_str))
            
        group = self._group_for("GET", path)
        self.limiter.acquire(group)

        logger.info("Sending GET request: url=%s", url)
        resp = requests.get(url, headers=headers)
        
        if resp.status_code == 429:
            logger.warning("Rate limit exceeded for group=%s", group)
            self.limiter.mark_exhausted(group)

        self.limiter.update_from_header(resp.headers.get("Remaining-Req"))
        self._log_response(resp)


        if 200 <= resp.status_code < 300:
            try:
                return resp.json()
            except ValueError:
                return resp.text

        try:
            ej = resp.json()
            if isinstance(ej, dict) and "error" in ej:
                e = ej["error"]
                return {
                    "status_code": resp.status_code,
                    "name": e.get("name"),
                    "message": e.get("message")
                }
            return {
                "status_code": resp.status_code,
                "name": None,
                "message": ej
            }
        except ValueError:
            return {
                "status_code": resp.status_code,
                "name": None,
                "message": resp.text
            }

if __name__ == "__main__":
    client = UpbitClient(
        "https://api.upbit.com",
        access_key="YOUR_ACCESS_KEY",
        secret_key="YOUR_SECRET_KEY"
    )

    # Example:List Pairs
    data = client.request_get("/v1/market/all")
    print(data)

    # Example:List Open Orders
    params = [("market", "KRW-BTC"), ("states[]", "wait"), ("states[]", "watch")]
    data = client.request_get("/v1/orders/open", params=params)
    print(data)
```

<br />

## 마치며

본 가이드에서는 업비트 REST API를 실제 시스템 수준에서 연동하기 위해 구현해야 하는 권장 요구사항들을 살펴보고, 간단한 코드 예제를 통해 구현 가이드라인을 안내했습니다. 본 문서를 통해 이해한 업비트 API 사용 방법을 사용 중인 언어와 지원 버전, 프레임워크, 시스템 특성에 맞게 적용하여 안정적인 연동 환경을 구축해보시기 바랍니다.

* [24시간 누적 거래대금 확인](https://docs.upbit.com/kr/docs/24-hour-accumulated-trade-volume) 튜토리얼 페이지로 이동하여 업비트 API를 활용한 더 많은 사용 예제를 살펴보세요.
* 또는 [REST API 사용 및 에러 안내](https://docs.upbit.com/kr/reference/rest-api-guide) 페이지와 [인증](https://docs.upbit.com/kr/reference/auth), [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits) 페이지로 이동하여 업비트 API 정책을 확인하고 나만의 Best Practice 구현을 즉시 시작할 수 있습니다.

# Sibling pages

* [WebSocket 연동 Best Practice](https://docs.upbit.com/kr/docs/websocket-best-practice.md)