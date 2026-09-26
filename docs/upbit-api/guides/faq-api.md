---
updatedAt: 2026-09-23T07:15:03.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# API 공통 문의

API 이용 중 자주 발생하는 문의와 해결 방법을 안내합니다.

### REST API 요청 시 오류가 발생합니다.

API 요청 처리 중 오류가 발생한 경우 HTTP 응답 본문(Body)에 에러 코드가 함께 반환됩니다. 주요 에러 코드는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rest-api-guide#:~:text=%EB%AC%B8%EC%84%9C%EB%A5%BC%20%EC%B0%B8%EC%A1%B0%ED%95%98%EC%8B%9C%EA%B8%B0%20%EB%B0%94%EB%9E%8D%EB%8B%88%EB%8B%A4.-,%EC%9D%91%EB%8B%B5%20%EC%83%81%ED%83%9C%20%EC%BD%94%EB%93%9C%20%EB%B0%8F%20%EC%97%90%EB%9F%AC%20%EC%95%88%EB%82%B4,-%EC%97%85%EB%B9%84%ED%8A%B8%20REST%20API">REST API 사용 안내</Anchor> 페이지 및 각 API Reference 문서 우측 하단 응답 예시에서 확인하실 수 있습니다.

만약 위 문서에서 확인되지 않는 오류에 대해 발생 원인을 확인할 수 없는 경우 해당 에러 코드를 포함하여 문의([open-api@upbit.com)](open-api@upbit.com "open-api@upbit.com")주시기 바랍니다.

***

### 파라미터를 정확하게 입력했음에도 `Invalid parameter. Check the given value!` 에러가 발생합니다.

파라미터에 `:` 또는 `+`와 같은 특수 문자가 포함되었을 때(주로 날짜 형식), 쿼리 문자열을 URL 인코딩 없이 요청하는 경우 위와 같은 에러가 발생할 수 있습니다.

<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rest-api-guide">REST API 사용 및 에러 안내</Anchor> 페이지 또는 <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/first-quotation-api-call#2-%EC%B2%AB-rest-api-%ED%98%B8%EC%B6%9C%ED%95%98%EA%B8%B0-%EC%97%85%EB%B9%84%ED%8A%B8-%EC%A7%80%EC%9B%90-%ED%8E%98%EC%96%B4-%EC%A1%B0%ED%9A%8C:~:text=2.-,%EC%B2%AB%20REST%20API%20%ED%98%B8%EC%B6%9C%ED%95%98%EA%B8%B0%3A%20%EC%97%85%EB%B9%84%ED%8A%B8%20%EC%A7%80%EC%9B%90%20%ED%8E%98%EC%96%B4%20%EC%A1%B0%ED%9A%8C,-cURL%20%EC%8B%A4%ED%96%89%20%ED%99%98%EA%B2%BD%EC%9D%B4">첫 업비트 API 호출하기</Anchor> 페이지를 참고하여 URL 인코딩 후 요청해주세요.

**URL 인코딩이란?**

URL 인코딩은 통신 프로토콜에서 URL 내에 포함할 수 없는 문자를 전송 가능한 문자로 변환하는 인코딩 방식입니다. 특수 문자를 포함한 대상 문자들은 인코딩 시 `%` 기호와 2자리 16진수로 이루어진 문자열로 변환됩니다.

* `:` → `%3A`
* `+` → `%2B`

***

### API Key 허용 IP 주소 목록에 현재 사용 중인 IP를 추가해도 오류가 발생합니다.

로컬 네트워크에서 확인되는 IP 주소와 실제 통신에 사용하는 IP 주소가 다른 경우 문제가 발생할 수 있습니다.

* **로컬 PC를 이용하는 경우**: 구글 등의 검색엔진에서 `what is my ip` 또는 `내 IP 주소` 등을 검색하여 확인된 공인 IP 주소
* **서버를 이용하는 경우**: 서버의 외부망 통신에 사용하는 공인 IP 주소

<span style="color:#0066CC;">API Key에 등록한 IP 주소와 실제 API 요청에 사용되는 공인 IP 주소가 동일한지 확인한 후 다시 시도해주세요.</span>

***

### 유동 IP 환경에서 API를 사용하고 싶습니다.

API Key 기반 인증을 통한 Exchange API 호출은 고정 IP 환경에서 허용 IP 목록을 등록한 후 사용할 수 있습니다.

자산 입출금 및 주문과 관련된 민감한 기능인 만큼, 회원님의 자산을 안전하게 보호하기 위한 조치이므로 클라우드 서버 또는 고정 IP 서비스를 통해 이용해주시기 바랍니다.

***

### Public WebSocket과 Private WebSocket은 무엇이 다른가요?

업비트 WebSocket은 <span style="color:#0066CC;">**API Key 인증 필요 여부에 따라 Public과 Private으로 구분됩니다.**</span>

| 구분             | Public WebSocket                   | Private WebSocket                                                                       |
| -------------- | ---------------------------------- | --------------------------------------------------------------------------------------- |
| **API Key 인증** | 불필요                                | 필요                                                                                      |
| **Endpoint**   | `wss://api.upbit.com/websocket/v1` | `wss://api.upbit.com/websocket/v1/private`<br />`wss://api.upbit.com/websocket/v1/info` |
| **주요 데이터**     | 현재가, 체결, 호가, 캔들 등                  | 내 주문 및 체결, 내 자산, 공지사항, 기술적 지표 등                                                         |
| **API Key 권한** | 불필요                                | 데이터 타입에 따라 다름                                                                           |

**Public WebSocket**은 API Key 인증 없이 이용할 수 있으며, 현재가·체결·호가 등 공개된 시세 데이터를 구독할 때 사용합니다.

**Private WebSocket**은 API Key를 이용한 인증이 필요하며, 인증이 필요한 데이터 타입을 구독할 때 사용합니다.

<span style="color:#0066CC;">Private WebSocket이라고 해서 모든 데이터 타입에 별도의 API Key 권한이 필요한 것은 아닙니다.</span> API Key 권한 필요 여부는 구독하는 데이터 타입에 따라 다르게 적용됩니다.

자세한 연결 및 인증 방법은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide#%EC%9A%94%EC%B2%AD-%EC%98%88%EC%A0%9C:~:text=%EC%A1%B0%ED%9A%8C%ED%95%98%EC%97%AC%20%EC%82%AC%EC%9A%A9%ED%95%98%EC%8B%9C%EA%B8%B0%20%EB%B0%94%EB%9E%8D%EB%8B%88%EB%8B%A4.-,%EC%9A%94%EC%B2%AD%20%EB%A9%94%EC%84%B8%EC%A7%80%EC%9D%98%20%EA%B5%AC%EC%A1%B0%EC%99%80%20%ED%98%95%EC%8B%9D,-%EB%8D%B0%EC%9D%B4%ED%84%B0%20%EC%A0%84%EC%86%A1%20%EC%9A%94%EC%B2%AD">WebSocket 사용 및 에러</Anchor> 안내를 참고해주세요.

***

### 기존 API Key로 REST API와 Private WebSocket을 함께 이용해도 되나요?

<span style="color:#0066CC;">**네, 함께 이용할 수 있습니다.**</span>

예를 들어 주문 권한이 설정된 기존 API Key로 REST API를 통해 주문을 요청하면서, 동일한 API Key를 Private WebSocket 인증에도 사용할 수 있습니다.

이때 동일한 API Key를 사용하더라도 <span style="color:#0066CC;">REST API와 WebSocket의 요청 수 제한이 하나로 합산되는 것은 아닙니다.</span>

요청 수 제한(Rate Limit)은 각 API에 정의된 **Rate Limit 그룹과 측정 단위**에 따라 적용됩니다. REST API 요청과 WebSocket 연결 및 데이터 요청은 각각 해당하는 요청 수 제한 정책에 따라 관리됩니다.

따라서 Private WebSocket을 구독하는 것만으로 동일한 API Key를 이용한 주문 API의 요청 가능 횟수가 차감되지는 않습니다.

자세한 내용은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rate-limits">요청 수 제한(Rate Limits)</Anchor>을 참고해주세요.

***

### 인증서 에러가 발생합니다. (`SSL: CERTIFICATE_VERIFY_FAILED`)

로컬 환경의 인증서 업데이트를 통한 최신 인증서 반영을 권장합니다.

* **Windows** → 최신 버전 업데이트
* **macOS** → 최신 버전 업데이트
* **Linux** → `apt-get install ca-certificates`

***

### 요청이 잘 처리되었었는데, 최근 에러 발생 빈도가 높아졌습니다.

DNS 캐시 문제일 가능성이 높습니다. 로컬 환경의 DNS 캐시 초기화를 통한 조치를 권장합니다.

* **Windows** → `ipconfig /flushdns`

***

### CORS 에러가 발생합니다.

시세 조회 관련 REST API, WebSocket 요청에 `Origin` 헤더가 존재하는 경우 요청 수 제한이 10초당 1회로 상향 제한됩니다.

* <Anchor target="_blank" href="https://docs.upbit.com/kr/changelog/origin_rate_limit">관련 공지 바로가기</Anchor>
* 해당 정책이 적용된 경우 응답의 `Remaining-Req` 헤더가 `group=origin`으로 반환됩니다.
* 요청 수 제한을 초과한 요청에 대해 CORS 에러가 발생할 수 있습니다.
* 브라우저에서 요청이 필요한 경우 제한에 맞추어 요청하거나 별도의 프록시 서버를 구성하여 사용하시기 바랍니다.

# Sibling pages

* [주문 관련 문의](https://docs.upbit.com/kr/docs/faq-order.md)
* [입출금 관련 문의](https://docs.upbit.com/kr/docs/faq-withdraw-deposit.md)
* [Upbit Strategy Toolkit 관련 문의](https://docs.upbit.com/kr/docs/faq-upbit-strategy-toolkit.md)