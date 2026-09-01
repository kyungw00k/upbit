Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [안내] 공지사항(Announcement) WebSocket 지원

안녕하세요. 업비트 개발자센터입니다.

업비트 공지사항을 WebSocket으로 실시간 수신할 수 있는 **공지사항(Announcement) 데이터 타입**이 추가됩니다.

신규 공지 게시 및 기존 공지 갱신 정보를 API 연동 환경에서 보다 편리하게 확인할 수 있습니다.

## 1. 적용 일시

* 적용일: 2026-08-31 (월)

## 2. 신규 기능 내용

* **기능명:** 공지사항(Announcement)
* **Endpoint:** `wss://api.upbit.com/websocket/v1/private`
* **데이터 타입:** `announcement`
* **지원 이벤트**

  * `CREATED`: 신규 공지사항 게시
  * `UPDATED`: 기존 공지사항 갱신
* `categories`를 지정하여 원하는 카테고리의 공지사항만 구독할 수 있습니다.
* `include_body=true`로 요청하면 공지사항 본문(`body`)을 함께 수신할 수 있습니다.
* 공지사항 데이터는 실시간 스트림만 지원하며, `stream_type`은 `REALTIME`로 제공됩니다.

## 3. 이용 방법

* Private WebSocket Endpoint에 연결한 후 인증을 진행합니다.
* Data Type Object의 `type`을 `announcement`로 지정하여 구독을 요청합니다.
* 필요한 경우 `categories`, `include_body`를 함께 지정할 수 있습니다.
* 자세한 요청 및 응답 명세는 아래 문서를 참고해주세요.

  * [공지사항 (Announcement)](https://docs.upbit.com/kr/reference/websocket-announcement)
  * [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide)

## 4. 유의사항

* 신규 공지사항이 게시되거나 기존 공지사항이 갱신된 경우에만 데이터가 전송됩니다. 관련 이벤트가 발생하지 않는 동안 데이터가 수신되지 않는 것은 정상 동작입니다.
* 서버에서 데이터가 정상적으로 전송되더라도, 네트워크 상태나 사용자 환경에 따라 데이터를 정상적으로 수신하지 못할 수 있습니다.
* WebSocket을 통한 공지사항 수신 시점은 웹·앱 등 다른 채널의 공지 노출 시점과 다를 수 있습니다. 채널별 전송 경로와 시스템 처리 상태, 시스템 보호를 위한 조치 등에 따라 수신 시점에 차이가 발생할 수 있습니다.
* categories에 지원하지 않는 값을 포함하여 요청하는 경우 INVALID\_PARAM 오류가 발생합니다.

***

문의 사항이 있으실 경우, <open-api@upbit.com>으로 연락해 주시기 바랍니다.

감사합니다.

업비트 개발자센터 드림