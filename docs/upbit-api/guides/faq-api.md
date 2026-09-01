---
updatedAt: 2026-05-04T06:59:04.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# API 공통 문의

API 사용과 관련한 자주 묻는 질문 모음입니다.

## REST API 요청 시 오류가 발생합니다.

API 요청 처리 중 오류가 발생한 경우 HTTP 응답 본문(Body)에 에러 코드가 함께 반환됩니다. 주요 에러 코드는 [REST API 사용 안내](https://docs.upbit.com/kr/reference/rest-api-guide) 페이지 및 각 API Reference 문서 우측 하단 응답 예시에서 확인하실 수 있습니다.\
만약 위 문서에서 확인되지 않는 오류에 대해 발생 원인을 확인할 수 없는 경우 해당 에러 코드를 포함하여 문의주시기 바랍니다.

<br />

## 파라미터를 정확하게 입력했음에도 'Invalid parameter. Check the given value!' 에러가 발생합니다.

파라미터에 : 또는 +와 같은 특수 문자가 포함되었을 때(주로 날짜 형식), 쿼리 문자열을 URL 인코딩 없이 요청 하는 경우 위와 같은 에러가 발생할 수 있습니다. [REST API 사용 및 에러 안내](https://docs.upbit.com/kr/reference/rest-api-guide) 페이지 또는 [첫 업비트 API 호출하기](https://docs.upbit.com/kr/docs/first-quotation-api-call) 페이지를 참고하여 URL 인코딩 후 요청해주세요.

[block:html]
{
  "html": "<div class=\"callout-section\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i>  URL 인코딩이란?\n      </div>\nURL 인코딩은 통신 프로토콜에서 URL 내에 포함할 수 없는 문자를 전송 가능한 문자로 변환하는 인코딩 방식입니다. 특수 문자를 포함한 대상 문자들은 인코딩 시 % 기호와 2자리 16진수로 이루어진 문자열로 변환됩니다.<br><br>\n(예시) <code>:</code>은 인코딩 후 <code>%3A</code>로, <code>+</code>는 <code>%2B</code>로 변환\n  </div>"
}
[/block]

<br />

## API Key 허용 IP 주소 목록에 현재 사용중인 IP를 추가해도 오류가 발생합니다.

로컬 네트워크에서 확인되는 IP 주소와 실제 통신에 사용하는 IP 주소가 다른 경우 문제가 발생할 수 있습니다.

* 로컬 PC를 이용하시는 경우: 구글 등의 검색엔진에서 "what is my ip" 혹은 "내 IP 주소" 등을 검색하여 접속한 후 확인된  IP 주소
* 서버를 이용하시는 경우 : 서버의 외부망 통신에 사용하는 IP 주소

위 주소를 허용 IP 목록에 등록한 후 다시 시도해보시기 바랍니다.

<br />

## 유동 IP 환경에서 API를 사용하고 싶습니다.

API Key 기반 인증을 통한 Exchange API 호출은 반드시 고정 IP 환경에서 허용 IP 목록을 등록한 후에 사용 가능합니다. 자산 입출금 및 주문과 관련된 민감한 기능인 만큼, 이는 회원님의 자산을 안전하게 보호하기 위한 조치이므로 번거로우시더라도 클라우드 서버 혹은 고정 IP 서비스를 통해 이용해주시기를 부탁드립니다.

<br />

## 인증서 에러가 발생합니다. (SSL: CERTIFICATE\_VERIFY\_FAILED)

로컬 환경의 인증서를 업데이트를 통한 최신 인증서 반영을 권장합니다.

* Windows -> 최신 버전 업데이트
* macOS -> 최신 버전 업데이트
* Linux -> `apt-get install ca-certificates`

<br />

## 요청이 잘 처리되었었는데, 최근 에러 발생 빈도가 높아졌습니다.

DNS 캐시 문제일 가능성이 높습니다. 로컬 환경의 DNS 캐시 초기화를 통한 조치를 권장합니다.

* 윈도우 → `ipconfig /flushdns`

<br />

## CORS 에러가 발생합니다

시세 조회 관련 REST API, WebSocket 요청에 Origin 헤더가 존재하는 경우 요청 수 제한이 10초당 1회로 상향 제한됩니다.

* [관련 공지 바로가기](https://docs.upbit.com/kr/changelog/origin_rate_limit)
* 해당 정책이 적용된 경우 응답의 Remaining-Req 헤더가 `group=origin` 으로 반환됩니다. 해당 헤더가 반환되며 요청 수 제한을 초과한 요청에 대해 CORS에러가 발생 할 수 있습니다.
* 브라우저에서 요청이 필요할 경우, 제한에 맞추어 요청해주시거나 별도의 프록시 서버를 구성하여 사용하시기 바랍니다.

# Sibling pages

* [주문 관련 문의](https://docs.upbit.com/kr/docs/faq-order.md)
* [입출금 관련 문의](https://docs.upbit.com/kr/docs/faq-withdraw-deposit.md)
* [Upbit Strategy Toolkit 관련 문의](https://docs.upbit.com/kr/docs/faq-upbit-strategy-toolkit.md)