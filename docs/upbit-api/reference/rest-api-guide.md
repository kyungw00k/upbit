---
updatedAt: 2026-08-31T04:37:05.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

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

업비트 REST API 응답으로 반환되는 HTTP 상태 코드 목록과 각 코드의 의미는 아래와 같습니다.

<style>
th.max-width {
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  }
th.min-width {
  min-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  }
.custom-table {
  table-layout: fixed;
  width: 100%;
}
</style>

<table className="custom-table">
  <thead>
    <tr>
      <th>HTTP Status Code</th>
      <th className="max-width">관련 에러 코드</th>
      <th className="min-width">발생 이유</th>
      <th>에러 해결 방법</th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td className="code-col">200 OK</td>
      <td></td>
      <td>정상 응답</td>
      <td></td>
    </tr>

    <tr>
      <td className="code-col">201 Created</td>
      <td></td>
      <td>요청으로 인한<br />생성 완료</td>
      <td></td>
    </tr>

    <tr>
      <td className="code-col" rowspan="7">400 Bad Request</td>
      <td><code>create_ask_error</code>,<br /><code>create_bid_error</code></td>
      <td>주문 요청 정보가 올바르지 않습니다.</td>
      <td>시장가 주문임에도 가격을 입력하는 경우 발생할 수 있습니다. 주문 생성 문서를 참고해주세요.</td>
    </tr>

    <tr>
      <td><code>insufficient_funds_ask</code>,<br /><code>insufficient_funds_bid</code></td>
      <td>매수/매도 가능<br />잔고가 부족합니다.</td>
      <td>잔고를 확인해주세요.</td>
    </tr>

    <tr>
      <td><code>under_min_total_ask</code>,<br /><code>under_min_total_bid</code></td>
      <td>최소 주문 금액에<br />미달합니다.</td>
      <td>페어별 최소 주문 금액 확인 후 재요청해주세요.</td>
    </tr>

    <tr>
      <td><code>withdraw_address<br />_not_registered</code></td>
      <td>허용되지 않은<br />출금 주소입니다.</td>
      <td>등록된 출금 주소 목록에 포함되어 있는지 확인해주세요.</td>
    </tr>

    <tr>
      <td><code>validation_error</code></td>
      <td>잘못된<br />API 요청입니다.</td>
      <td>필수 파라미터 누락<br />여부를 확인해주세요.</td>
    </tr>

    <tr>
      <td><code>invaild_parameter</code></td>
      <td>잘못된<br />파라미터입니다.</td>
      <td>입력한<br />파라미터를 확인해주세요.</td>
    </tr>

    <tr>
      <td><code>duplicated_identifier</code></td>
      <td>이미 등록된<br />identifier입니다.</td>
      <td>이미 등록된 identifier입니다. 새로운 identifier를 입력해 주세요.</td>
    </tr>

    <tr>
      <td className="code-col" rowspan="6">401 Unauthorized</td>
      <td><code>invalid_query_payload</code></td>
      <td>JWT 페이로드가<br />올바르지 않습니다.</td>
      <td><a href="auth">인증</a> 가이드 문서를 참고하여 서명이 올바르게 생성 되었는지 확인해주세요.</td>
    </tr>

    <tr>
      <td><code>jwt_verification</code></td>
      <td>JWT 검증에<br />실패했습니다.</td>
      <td>토큰의 생성 및 서명 상태를 점검해주세요.</td>
    </tr>

    <tr>
      <td><code>expired_access_key</code></td>
      <td>API 키가<br />만료되었습니다.</td>
      <td>새로운 키를 발급받아 사용해주세요.</td>
    </tr>

    <tr>
      <td><code>nonce_used</code></td>
      <td>이미 사용된<br />nonce 값입니다.</td>
      <td>JWT에는 매 요청마다 새로운 nonce 값을 사용해야 합니다.</td>
    </tr>

    <tr>
      <td><code>no_authorization_ip</code></td>
      <td>등록되지 않은 IP에서<br />요청되었습니다.</td>
      <td>API 키 발급 시 등록한 IP 환경에서 호출 중인지 점검해주세요.</td>
    </tr>

	<tr>
      <td><code>no_authorization_token</code></td>
      <td>인증 토큰이<br />누락되었습니다.</td>
      <td>인증 헤더가 요청에 포함되었는지 확인해주세요.</td>
    </tr>

    <tr>
      <td className="code-col" rowspan="2">403 Forbidden</td>
      <td><code>out_of_scope</code></td>
      <td>권한이 부족합니다.</td>
      <td>현재 호출에 사용한 API 키에 권한이 없습니다. 권한이 이 부여된 API 키로 다시 요청해 주세요. 해당에 부여된 권한 조회는 <a href="https://docs.upbit.com/kr/reference/list-pocket-api-keys">포켓별 API Key 목록 조회</a>에서 확인 가능합니다. <br /> API 도메인(ex.포켓, 입출금, 주문)마다 HTTP 상태 코드(401,403)가 일관되지 않을 수 있습니다. 예외처리 시 에러코드 문자열(out_of_scope)을 기준으로 검증하는 것을 권장합니다.</td>
    </tr>
    
     <tr> <td><code>open_api_withdraw_locked</code></td>
      <td>출금 안심차단을 해지하지 않았습니다.</td>
      <td>Open API 출금이 잠금 상태입니다. 업비트 모바일 앱 > 더보기 > 인증/보안 > Open API 관리에서 출금 잠금을 해제해주세요.</td>
    </tr>



    <tr>
      <td className="code-col" rowspan="3">404 Not Found</td>
      <td></td>
      <td>존재하지 않는<br />데이터에 접근</td>
      <td>주문, 출금, 입금, 체결 등 요청 항목이 존재하지 않는 경우</td>
    </tr>

	<tr>
      <td><code>pocket_not_found</code></td>
      <td>포켓을 찾지 못했습니다.</td>
      <td>포켓 UUID를 확인해주세요. 포켓 UUID는 <a href="https://docs.upbit.com/kr/reference/list-pockets">포켓 정보 조회 API</a>에서 확인 가능합니다.</td>
    </tr>

	<tr>
      <td><code>currency_not_found</code></td>
      <td>자산을 찾지 못했습니다.</td>
      <td>입력한 currency 값이 조회 가능한 자산 코드인지 확인해주세요. 자산 코드는 대문자로 입력해야하며, 존재하지 않거나 지원되지 않는 자산 코드를 입력한 경우 해당 자산을 찾을 수 없습니다.</td>
    </tr>

    <tr>
      <td className="code-col">418 I'm a teapot</td>
      <td></td>
      <td>과도한 요청으로<br />인해 거부되었습니다.</td>
      <td>IP 차단 등으로 요청이 제한됩니다.</td>
    </tr>

    <tr>
      <td className="code-col">429 Too Many Requests</td>
      <td></td>
      <td>요청 제한을<br />초과했습니다.</td>
      <td>API 호출 한도를 초과했습니다.</td>
    </tr>

    <tr>
      <td className="code-col">500 Internal Server Error</td>
      <td></td>
      <td>서버 내부 오류</td>
      <td>서비스 점검 또는 시스템 오류로 인한 처리 불가</td>
    </tr>

  </tbody>
</table>

<br />

에러 발생 시, 응답은 다음과 같은 JSON 형식으로 반환됩니다. `name` 필드는 해당 에러의 코드를, `message` 필드는 오류와 관련된 메세지를 반환합니다. Quotation API는 정수 형식의 `name` 필드를, Exchange API는 문자열 형식의 `name` 필드를 반환합니다. 각 API 별 오류 예시는 API Reference 문서 우측 하단 응답 예시 영역을 참고하시기 바랍니다.

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
    "name": "ERRPR_CODE",
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