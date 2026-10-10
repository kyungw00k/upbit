---
updatedAt: 2026-10-06T05:15:18.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# REST API 사용 및 에러 안내

업비트 REST API 사용을 위한 요청, 인증, 에러 및 gzip 지원 관련 안내입니다.

## Endpoint

> <https://api.upbit.com/v1>

<br />

## TLS

업비트 Open API는 회원님의 정보를 안전하게 보호하기 위해 TLS 1.2 이상 버전만 지원합니다.<br />TLS 1.2 미만 버전은 더 이상 지원되지 않으므로, 최소 TLS 1.2 이상으로 업그레이드해 주시기 바랍니다(TLS 1.3 권장 버전).

<br />

## Content Type

업비트 REST API는 `application/json` Content Type을 지원합니다. 특히 POST 요청의 경우, 본문(Body)을 JSON 형식으로 요청해야 하며 아래 헤더를 함께 지정하여 주시기 바랍니다.

> Content-Type: application/json; charset=utf-8

<Callout icon="fad fa-circle-info" theme="error">
  ### **POST API에 대한 Form 방식 요청은 2022년 3월 1일부로 지원이 종료되었습니다.**

  Form 방식 지원 종료에 따라 Urlencoded Form 방식으로 전송하는 POST 요청에 대한 정상적인 동작을 보장하지 않습니다. **반드시 JSON 형식으로 요청 본문(Body**)을 전송해주시기 바랍니다.
</Callout>

<br />

## 인증

인증이 필요한 Exchange API 요청 시 [인증](https://docs.upbit.com/kr/reference/auth) 가이드를 참고하여 생성한 JWT 토큰을 `Authorization` 헤더에 반드시 포함하여 요청해야 합니다. 업비트 REST API는 Bearer 인증을 지원하며, 아래와 같은 형식으로 입력합니다.

> Authorization: Bearer eyJhb...d8sTw

<br />

## 요청 수 제한

REST API 요청 수 제한 정책은 [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits)문서를 참조하시기 바랍니다.

<br />

## 응답 상태 코드 및 에러 안내

업비트 REST API에서 반환하는 주요 HTTP 상태 코드와 의미는 다음과 같습니다.

| HTTP Status Code            | 설명                                      |
| --------------------------- | --------------------------------------- |
| `200 OK`                    | 요청이 정상적으로 처리되었습니다.                      |
| `201 Created`               | 요청한 리소스가 정상적으로 생성되었습니다.                 |
| `400 Bad Request`           | 요청 형식이나 파라미터가 올바르지 않습니다.                |
| `401 Unauthorized`          | API 인증에 실패했습니다.                         |
| `403 Forbidden`             | API Key 권한 또는 보안 설정으로 인해 요청이 허용되지 않습니다. |
| `404 Not Found`             | 요청한 리소스를 찾을 수 없습니다.                     |
| `418 I'm a teapot`          | 과도한 요청이 반복되어 요청이 제한되었습니다.               |
| `429 Too Many Requests`     | 요청 수 제한을 초과했습니다.                        |
| `500 Internal Server Error` | 서버 내부 오류로 요청을 처리할 수 없습니다.               |

<br />

### 400 Bad Request

요청 형식이나 파라미터가 올바르지 않은 경우 발생합니다.

| 에러 코드                             | 원인 및 해결 방법                                                                                                                                                                              |
| --------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `create_ask_error`                | **발생 이유**<br />매도 주문 요청 정보가 올바르지 않습니다.<br /><br />**해결 방법**<br />주문 유형에 필요한 파라미터와 입력값을 확인해 주세요. 시장가 주문에 가격을 입력하는 등 주문 조건에 맞지 않는 파라미터 조합으로 요청한 경우 발생할 수 있습니다.                            |
| `create_bid_error`                | **발생 이유**<br />매수 주문 요청 정보가 올바르지 않습니다.<br /><br />**해결 방법**<br />주문 유형에 필요한 파라미터와 입력값을 확인해 주세요. 시장가 주문에 가격을 입력하는 등 주문 조건에 맞지 않는 파라미터 조합으로 요청한 경우 발생할 수 있습니다.                            |
| `insufficient_funds_ask`          | **발생 이유**<br />매도 가능한 잔고가 부족합니다.<br /><br />**해결 방법**<br />주문 가능한 자산 잔고를 확인해 주세요.                                                                                                       |
| `insufficient_funds_bid`          | **발생 이유**<br />매수 가능한 잔고가 부족합니다.<br /><br />**해결 방법**<br />주문 가능한 자산 잔고를 확인해 주세요.                                                                                                       |
| `under_min_total_ask`             | **발생 이유**<br />최소 매도 주문 금액에 미달합니다.<br /><br />**해결 방법**<br />마켓별 최소 주문 금액을 확인한 후 다시 요청해 주세요.                                                                                            |
| `under_min_total_bid`             | **발생 이유**<br />최소 매수 주문 금액에 미달합니다.<br /><br />**해결 방법**<br />마켓별 최소 주문 금액을 확인한 후 다시 요청해 주세요.                                                                                            |
| `withdraw_address_not_registered` | **발생 이유**<br />허용되지 않은 출금 주소입니다.<br /><br />**해결 방법**<br />Open API 출금 허용 주소로 등록된 주소인지 확인해 주세요.                                                                                         |
| `validation_error`                | **발생 이유**<br />API 요청 형식이 올바르지 않습니다.<br /><br />**해결 방법**<br />필수 파라미터 누락 여부와 요청 데이터 형식을 확인해 주세요.                                                                                       |
| `invalid_parameter`               | **발생 이유**<br />요청한 파라미터가 올바르지 않습니다.<br /><br />**해결 방법**<br />파라미터의 형식, 허용값 및 입력 범위를 확인해 주세요.                                                                                           |
| `invalid_post_only`               | **발생 이유**<br />`post_only`와 함께 사용할 수 없는 `smp_type` 값을 지정한 경우 발생합니다.<br /><br />**해결 방법**<br />`time_in_force=post_only` 주문에서는 `smp_type`을 지정할 수 없습니다. `smp_type` 파라미터를 제외하고 다시 요청해 주세요. |
| `duplicated_identifier`           | **발생 이유**<br />이미 사용 중인 `identifier`입니다.<br /><br />**해결 방법**<br />이전 주문에 사용하지 않은 새로운 `identifier`로 요청해 주세요.                                                                            |

<br />

### 401 Unauthorized

JWT 토큰, API Key 또는 허용 IP 인증에 실패한 경우 발생합니다.

| 에러 코드                    | 원인 및 해결 방법                                                                                                                                               |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `invalid_query_payload`  | **발생 이유**<br />JWT 페이로드가 올바르지 않습니다.<br /><br />**해결 방법**<br />[인증](https://docs.upbit.com/kr/reference/auth) 문서를 참고하여 JWT 페이로드와 서명이 올바르게 생성되었는지 확인해 주세요. |
| `jwt_verification`       | **발생 이유**<br />JWT 검증에 실패했습니다.<br /><br />**해결 방법**<br />JWT 생성 과정과 서명에 사용한 Secret Key를 확인해 주세요.                                                         |
| `expired_access_key`     | **발생 이유**<br />API Key가 만료되었습니다.<br /><br />**해결 방법**<br />새로운 API Key를 발급받아 사용해 주세요.                                                                    |
| `nonce_used`             | **발생 이유**<br />이미 사용된 `nonce` 값입니다.<br /><br />**해결 방법**<br />JWT를 생성할 때마다 새로운 `nonce` 값을 사용해 주세요.                                                       |
| `no_authorization_ip`    | **발생 이유**<br />API Key에 등록되지 않은 IP에서 요청했습니다.<br /><br />**해결 방법**<br />요청을 전송한 IP가 해당 API Key의 허용 IP로 등록되어 있는지 확인해 주세요.                                  |
| `no_authorization_token` | **발생 이유**<br />인증 토큰이 요청에 포함되지 않았습니다.<br /><br />**해결 방법**<br />요청의 `Authorization` 헤더에 `Bearer {JWT}` 형식의 인증 토큰이 포함되어 있는지 확인해 주세요.                      |

<br />

### 403 Forbidden

API Key 권한 또는 Open API 보안 설정으로 인해 요청이 허용되지 않은 경우 발생합니다.

| 에러 코드                      | 원인 및 해결 방법                                                                                                                                                                                                                        |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `out_of_scope`             | **발생 이유**<br />요청에 필요한 API Key 권한이 없습니다.<br /><br />**해결 방법**<br />사용 중인 API Key에 해당 API를 호출할 수 있는 권한이 부여되어 있는지 확인해 주세요. API Key별 권한은 [포켓별 API Key 목록 조회](https://docs.upbit.com/kr/reference/list-pocket-api-keys)에서 확인할 수 있습니다. |
| `open_api_withdraw_locked` | **발생 이유**<br />Open API 출금이 잠금 상태입니다.<br /><br />**해결 방법**<br />업비트 모바일 앱의 **더보기 > 인증/보안 > Open API 관리**에서 출금 잠금을 해제해 주세요.                                                                                                        |

<Callout icon="fad fa-circle-info" theme="info">
  API 영역에 따라 `out_of_scope` 에러와 함께 반환되는 HTTP 상태 코드가 `401 Unauthorized` 또는 `403 Forbidden`으로 다를 수 있습니다.

  에러 처리 시 HTTP 상태 코드만으로 분기하지 말고, 응답의 `error.name` 값도 함께 확인해 주세요.
</Callout>

<br />

### 404 Not Found

요청한 리소스가 존재하지 않거나 조회할 수 없는 경우 발생합니다.

| 에러 코드                | 원인 및 해결 방법                                                                                                                                                               |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| —                    | **발생 이유**<br />요청한 주문, 출금, 입금 또는 체결 등의 데이터를 찾을 수 없습니다.<br /><br />**해결 방법**<br />요청에 사용한 UUID, `identifier` 또는 조회 조건을 확인해 주세요.                                           |
| `pocket_not_found`   | **발생 이유**<br />요청한 포켓을 찾을 수 없습니다.<br /><br />**해결 방법**<br />포켓 UUID를 확인해 주세요. 포켓 UUID는 [포켓 정보 조회](https://docs.upbit.com/kr/reference/list-pockets)에서 확인할 수 있습니다.        |
| `currency_not_found` | **발생 이유**<br />요청한 자산을 찾을 수 없습니다.<br /><br />**해결 방법**<br />`currency` 값이 조회 가능한 자산 코드인지 확인해 주세요. 자산 코드는 대문자로 입력해야 하며, 존재하지 않거나 지원되지 않는 자산 코드를 입력한 경우 해당 에러가 발생할 수 있습니다. |

<br />

### 418 I'm a teapot

과도한 요청이 반복되어 요청이 제한된 경우 발생합니다.

| 에러 코드 | 원인 및 해결 방법                                                                                                                                                                          |
| ----- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| —     | **발생 이유**<br />과도한 요청으로 인해 요청 제한이 적용되었습니다.<br /><br />**해결 방법**<br />추가 요청을 중단하고, 이후에는 [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits) 정책에 맞게 호출량을 조정해 주세요. |

<br />

### 429 Too Many Requests

허용된 요청 수 제한을 초과한 경우 발생합니다.

| 에러 코드 | 원인 및 해결 방법                                                                                                                                                        |
| ----- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| —     | **발생 이유**<br />API 호출 한도를 초과했습니다.<br /><br />**해결 방법**<br />[요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits) 문서를 참고하여 요청 간격과 호출량을 조정해 주세요. |

<br />

### 500 Internal Server Error

서버 내부 오류로 요청을 처리할 수 없는 경우 발생합니다.

| 에러 코드 | 원인 및 해결 방법                                                                                                                             |
| ----- | -------------------------------------------------------------------------------------------------------------------------------------- |
| —     | **발생 이유**<br />서비스 점검 또는 일시적인 시스템 오류로 요청을 처리할 수 없습니다.<br /><br />**해결 방법**<br />잠시 후 다시 요청해 주세요. 오류가 계속 발생하는 경우 요청 시각과 응답 내용을 확인해 주세요. |

<br />

### 에러 응답 형식

에러가 발생하면 다음과 같은 JSON 형식으로 응답이 반환됩니다.

* `error.name`: 에러 코드
* `error.message`: 에러와 관련된 메시지

시세(Quotation) API는 정수 형식의 `name` 값을 반환하며, 거래 및 자산 관리(Exchange) API는 문자열 형식의 `name` 값을 반환합니다.

각 API의 상세 오류 예시는 API Reference 문서 우측 하단의 응답 예시 영역에서 확인할 수 있습니다.

```json Quotation API Error Response
{
  "error": {
    "name": 400,
    "message": "ERROR_MESSAGE"
  }
}
```

```json Exchange API Error Response
{
  "error": {
    "name": "ERROR_CODE",
    "message": "ERROR_MESSAGE"
  }
}
```

에러 발생 시, 응답은 다음과 같은 JSON 형식으로 반환됩니다. `name` 필드는 해당 에러의 코드를, `message` 필드는 오류와 관련된 메세지를 반환합니다. Quotation API는 정수 형식의 `name` 필드를, Exchange API는 문자열 형식의 `name` 필드를 반환합니다. 각 API 별 오류 예시는 API Reference 문서 우측 하단 응답 예시 영역을 참고하시기 바랍니다.

```json Quotation API Error Response
{
  "error": {
    "name": 400,
    "message": "ERROR_MESSAGE"
  }
}
```

<br />

## 인코딩

GET 또는 DELETE API에 대해 쿼리 파라미터를 포함한 요청을 전송하는 경우 모든 쿼리 파라미터를 URL 인코딩 한 후 요청해야 합니다. 인코딩이 정상적으로 이루어지지 않은 요청에 대해 응답으로 400 `Invalid parameter` 에러가 발생할 수 있습니다. 단, Exchange API의 파라미터 중 배열 형식의 파라미터가 이름에 \[]를 포함하고 있는 경우, '\[',']' 문자는 인코딩 대상에서 제외합니다.

<Callout icon="fad fa-circle-info" theme="info">
  ### URL 인코딩이란?

  URL 인코딩은 통신 프로토콜에서 URL 내에 포함할 수 없는 문자를 전송 가능한 문자로 변환하는 인코딩 방식입니다. 특수 문자를 포함한 대상 문자들은 인코딩 시 % 기호와 2자리 16진수로 이루어진 문자열로 변환됩니다.

  (예시) :은 인코딩 후 %3A로, +는 %2B로 변환
</Callout>

<br />

## gzip 응답 지원

인코딩 옵션을 `gzip`으로 요청하여 REST API 응답을 gzip 압축된 형태로 수신할 수 있습니다. gzip 형식으로 요청 시 API를 통해 주고 받는 데이터 크기를 줄여 트래픽 비용 및 응답 시간을 절감할 수 있습니다. gzip 인코딩은 **시세(Quotation) API만 지원**합니다. gzip 옵션 사용 시 아래와 같이 헤더를 지정합니다.

> Accept-Encoding: gzip

<br />

## API Reference 예제 코드 안내

보다 쉬운 API 사용을 위해 각 API Reference 우측 상단에서 해당 API 호출을 위한 예제 코드를 제공합니다. Shell(cURL), Python, Java, Node.js의 네가지 도구/언어 예시를 제공하며 Java와 Node.js는 사용하는 HTTP 클라이언트 라이브러리에 따라 아래와 같이 세분화하여 제공합니다. (코드 박스 상단 아래 화살표를 클릭하여 라이브러리별 예제를 변경할 수 있습니다.)

* **Java** - AsyncHttp, java.net.http, OkHttp, Unirest
* **Node.js** - Axios, fetch, https

문서의 요청 쿼리 파라미터 및 본문(Body) 영역에서 각 파라미터의 값을 예시값 또는 임의의 값으로 입력할 수 있습니다. 입력된 값에 따라 우측 예제 코드 또한 실시간으로 변경됩니다.

거래 및 자산 관리(Exchange) API의 경우 인증 토큰을 생성하는 부분은 예제 코드에서 제외되어 있습니다. [인증](https://docs.upbit.com/kr/reference/auth)문서의 인증 토큰 생성 예제 코드를 참고하여 실제 연동 시에는 반드시 인증 토큰을 요청에 포함하도록 구현해주시기 바랍니다.

# Sibling pages

* [개요](https://docs.upbit.com/kr/reference/api-overview.md)
* [인증](https://docs.upbit.com/kr/reference/auth.md)
* [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits.md)
* [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide.md)