---
updatedAt: 2026-09-02T01:39:36.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 공지사항(Announcement)

업비트 공지사항을 WebSocket으로 수신하기 위한 요청 및 구독 데이터 예시를 제공합니다.

## **WebSocket Endpoint**

| 구분      | Endpoint                                   |
| ------- | ------------------------------------------ |
| Private | `wss://api.upbit.com/websocket/v1/private` |

<Callout icon="fad fa-triangle-exclamation" theme="error">
  ### **Private WebSocket 연결 관리 안내**

  **Private WebSocket은 동시에 많은 연결을 유지하는 경우 신규 연결이 거절될 수 있습니다.**

  불필요한 연결 생성을 최소화하고, 가능한 경우 하나의 연결에서 필요한 데이터 타입을 함께 구독해 주세요.
</Callout>

<Callout icon="fad fa-circle-info" theme="info">
  ### **신규 공지 또는 갱신이 없으면 데이터가 전송되지 않습니다.**

  공지사항 데이터는 신규 공지사항이 게시되거나 기존 공지사항이 갱신된 경우에만 실시간으로 전송됩니다. 따라서 WebSocket 연결 후 신규 게시 또는 갱신이 없다면 데이터가 수신되지 않는 것이 정상입니다.

  공지사항 데이터는 실시간 스트림만 지원하며, `stream_type`은 `REALTIME`로 제공됩니다. WebSocket 연결 및 인증 방법은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide">WebSocket 사용 및 에러 안내</Anchor> 문서를 참고해주세요.
</Callout>

<Callout icon="fad fa-clock" theme="info">
  ### **채널별 공지사항 수신 시점에 차이가 발생할 수 있습니다.**

  **서버에서 데이터가 정상적으로 전송되더라도** 네트워크 상태나 사용자 환경에 따라 데이터를 정상적으로 수신하지 못할 수 있습니다. 또한 WebSocket을 통한 공지사항 수신 시점은 웹·앱 등 다른 채널의 공지 노출 시점과 다를 수 있으며, 채널별 전송 경로와 시스템 처리 상태, 시스템 보호를 위한 조치 등에 따라 수신 시점에 차이가 발생할 수 있습니다.
</Callout>

<br />

## Request 메세지 형식

공지사항 데이터 수신을 요청하기 위해서는 WebSocket 연결 이후 아래 구조의 JSON Object를 생성한 뒤 요청 메세지의 Data Type Object로 포함하여 전송해야 합니다.

Ticket, Format 필드를 포함한 전체 WebSocket 데이터 요청 메세지 명세는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide">WebSocket 사용 및 에러 안내</Anchor> 문서를 참고해주세요.

| 필드명           | 타입          | 내용                                                                                                                                                        | 필수 여부    | 기본 값    |
| ------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| type          | String      | 수신할 데이터 타입.<br />공지사항 데이터는 `announcement`로 지정합니다.                                                                                                         | Required |         |
| categories    | List:String | 수신하고자 하는 공지사항 카테고리 목록.<br />생략할 경우 모든 카테고리의 공지사항을 수신합니다.<br />지원하지 않는 값을 포함하여 요청하는 경우 `INVALID_PARAM` 오류가 발생합니다.                                          | Optional | 전체      |
| include\_body | Boolean     | 공지사항 본문 포함 여부. <br />true로 요청한 경우 구독 데이터에 body 필드가 포함됩니다.<br />본문 데이터의 용량이 큰 경우 null로 제공될 수 있습니다.                                                         | Optional | `false` |
| format        | String      | 수신하고자 하는 데이터 포맷입니다. <br />`DEFAULT` : 기본 포맷.<br />`SIMPLE` : 간략한 포맷. 각 필드가 축약어 형태로 반환됩니다.<br />`JSON_LIST` : 리스트 포맷.<br />`SIMPLE_LIST` : 축약어 형태의 리스트 포맷. | Required |         |

<br />

### 공지사항 카테고리 (categories)

categories를 지정하면 원하는 카테고리의 공지사항만 구독할 수 있습니다. 별도로 지정하지 않으면 모든 카테고리의 공지사항을 수신합니다.

| 한국어    | 영문            | `categories` 값  | 설명                                          |
| ------ | ------------- | --------------- | ------------------------------------------- |
| 안내     | Notice        | `notice`        | 신규 기능, 기능 개선 등 서비스 이용 관련 공지                 |
| 거래     | Trade         | `trade`         | 신규 거래지원, 마켓 추가, 거래 유의 촉구·지정 및 거래지원 종료 관련 공지 |
| 입출금    | D/W           | `wallet`        | 디지털 자산 입출금 지원, 중단 및 재개 관련 공지                |
| 점검     | Maintenance   | `maintenance`   | 시스템 및 서비스 점검 관련 공지                          |
| 디지털 자산 | Digital Asset | `digital_asset` | 에어드랍, 유통량 등 디지털 자산 관련 공지                    |
| NFT    | NFT           | `nft`           | NFT 드롭 관련 공지                                |
| 서비스+   | Service+      | `service`       | 업비트 서비스+ 관련 공지                              |
| 이벤트    | Event         | `event`         | 업비트 이벤트 관련 공지                               |

<br />

## 구독 데이터 명세

신규 공지사항이 게시되거나 기존 공지사항이 갱신되는 경우 실시간 스트림 데이터가 아래와 같이 반환됩니다.

| 필드명               | 축약형  | 내용                                                                  | 타입     | 값                                              |
| ----------------- | ---- | ------------------------------------------------------------------- | ------ | ---------------------------------------------- |
| type              | ty   | 데이터 타입                                                              | String | `announcement`                                 |
| event\_type       | et   | 공지사항 이벤트 유형                                                         | String | `CREATED` : 신규 공지 게시<br />`UPDATED` : 기존 공지 갱신 |
| uuid              | uid  | 공지사항의 고유 식별자                                                        | String |                                                |
| title             | tt   | 공지사항 제목                                                             | String |                                                |
| category          | cat  | 공지사항 카테고리.                                                          | String | 예: `trade`, `wallet`, `maintenance`, `service` |
| url               | url  | 공지사항 상세 페이지 URL                                                     | String |                                                |
| first\_listed\_at | flat | 최초 게시 시각 (KST)                                                      | String | ISO 8601                                       |
| listed\_at        | lat  | 최종 게시 시각 (KST)<br />공지사항이 갱신되는 경우 함께 갱신됩니다.                         | String | ISO 8601                                       |
| timestamp         | tms  | 타임스탬프 (ms)                                                          | Long   |                                                |
| stream\_type      | st   | 스트림 타입                                                              | String | `REALTIME` : 실시간 스트림                           |
| body              | bd   | 공지사항 본문(Markdown)<br />요청 시 `include_body`를 `true`로 지정한 경우에만 포함됩니다. | String |                                                |

<br />

## 예시

#### 모든 공지사항 구독

> 카테고리를 지정하지 않고 업비트의 모든 공지사항을 구독하는 예시입니다. 공지사항 본문은 제외하고 수신합니다. ticket은 요청을 식별하기 위한 값으로 원하는 문자열을 직접 지정할 수 있습니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "announcement-monitor"
  },
  {
    "type": "announcement",
    "include_body": false
  },
  {
    "format": "DEFAULT"
  }
]
```

##### 수신 메시지 예제

```json
{
  "type": "announcement",
  "event_type": "CREATED",
  "uuid": "562593104",
  "title": "[거래] 신규 거래 지원 안내 ...",
  "category": "trade",
  "url": "https://www.upbit.com/service_center/notice?id=562593104&view=share",
  "first_listed_at": "2026-07-14T16:36:34+09:00",
  "listed_at": "2026-07-14T16:36:34+09:00",
  "timestamp": 1752566400123,
  "stream_type": "REALTIME"
}
```

<br />

#### 특정 카테고리 공지사항 구독

> categories에 trade를 지정하여 거래 관련 공지사항만 구독하는 예시입니다. 공지사항 본문은 제외하고 수신합니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "trade-announcement-monitor"
  },
  {
    "type": "announcement",
    "categories": ["trade"],
    "include_body": false
  },
  {
    "format": "DEFAULT"
  }
]
```

##### 수신 메시지 예제

```json
{
  "type": "announcement",
  "event_type": "CREATED",
  "uuid": "562593104",
  "title": "[거래] 신규 거래 지원 안내 ...",
  "category": "trade",
  "url": "https://www.upbit.com/service_center/notice?id=562593104&view=share",
  "first_listed_at": "2026-07-14T16:36:34+09:00",
  "listed_at": "2026-07-14T16:36:34+09:00",
  "timestamp": 1752566400123,
  "stream_type": "REALTIME"
}
```

<br />

#### 공지사항 본문 포함 구독

> 카테고리 구분 없이 모든 공지사항을 구독하며, include\_body를 true로 설정하면 공지 본문까지 함께 수신합니다. ticket은 요청을 식별하기 위한 값으로 원하는 문자열을 직접 지정할 수 있습니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "announcement-subscription"
  },
  {
    "type": "announcement",
    "include_body": true
  },
  {
    "format": "DEFAULT"
  }
]
```

#####

##### 수신 메시지 예제 - 신규 공지 게시

```json
{
  "type": "announcement",
  "event_type": "CREATED",
  "uuid": "562593104",
  "title": "[거래] 신규 거래 지원 안내 ...",
  "category": "trade",
  "url": "https://www.upbit.com/service_center/notice?id=562593104&view=share",
  "first_listed_at": "2026-07-14T16:36:34+09:00",
  "listed_at": "2026-07-14T16:36:34+09:00",
  "timestamp": 1752566400123,
  "stream_type": "REALTIME",
  "body": "공지사항 본문 ..."
}
```

<br />

## 에러 안내&#x20;

WebSocket 연결 후 요청에 대한 에러 발생 시, 응답은 다음과 같은 JSON 형식으로 반환됩니다.

```json Error Response
{
  "error": {
    "name": "ERRPR_CODE",
    "message": "ERROR_MESSAGE"
  }
}
```

반환될 수 있는 주요 에러 코드 목록은 아래와 같습니다.

| error.name        | 발생 이유                                    | 권장 조치                                                                                               |
| ----------------- | ---------------------------------------- | --------------------------------------------------------------------------------------------------- |
| INVALID\_AUTH     | 인증 정보 누락 또는 인증 토큰 검증 실패                  | Private WebSocket을 사용하는 경우 올바른 Endpoint에 연결했는지 확인하고, Authorization 헤더에 유효한 인증 토큰이 포함되어 있는지 확인해 주세요. |
| WRONG\_FORMAT     | 요청 메시지 형식 오류                             | 요청 메시지가 WebSocket 요청 형식에 맞게 작성되었는지 확인해 주세요. Object 구성과 각 필드의 타입 및 값을 함께 확인해 주세요.                    |
| NO\_TICKET        | ticket 필드 누락                             | 요청 메시지에 Ticket Object와 ticket 필드가 포함되어 있는지 확인해 주세요.                                                 |
| NO\_TYPE          | type 필드 누락                               | Data Type Object에 type 필드가 포함되어 있는지 확인하고, 구독할 데이터 타입을 지정해 주세요.                                      |
| NO\_CODES         | codes 필드 누락                              | 구독하려는 데이터 타입에서 codes 필드가 필요한지 확인하고, 수신할 페어 코드 목록을 지정해 주세요.                                          |
| INVALID\_PARAM    | 필수 요청 필드 누락 또는 지원하지 않는 값 요청              | 요청 메시지에 필요한 필드가 포함되어 있는지, 각 필드에 지원하는 값이 지정되었는지 확인해 주세요.                                             |
| Too Many Requests | 요청 한도 초과                                 | 다음 요청이 가능한 시점까지 대기한 후 다시 요청해 주세요. 요청 한도와 잔여 요청 수 확인 방법은 아래 잔여 요청 수 확인 방법을 참고해 주세요.                  |
| I'm a teapot      | Too Many Requests가 반복되어 일정 시간 요청이 제한된 상태 | 응답에 포함된 제한 시간을 확인하고, 안내된 시간이 지난 후 다시 요청해 주세요.                                                       |

<br />

## 요청 수 제한

API는 Rate Limit 그룹으로 묶입니다. 같은 그룹의 API는 초당 한도를 함께 차감합니다. Rate Limit 그룹별 초당 최대 허용 요청 수는 서비스 정책에 따라 공지 후 변경되거나, 서비스 상황에 따라 추가 제한이 발생할 수 있습니다. 자세한 설명은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rate-limits">요청 수 제한(Rate Limits)</Anchor>를 참고해주세요

| Rate Limit 그룹       | 정책                | 적용 단위 |
| ------------------- | ----------------- | ----- |
| `websocket-connect` | 초당 최대 5회          | 포켓    |
| `websocket-message` | 초당 최대 5회, 분당 100회 | 커넥션   |

<Callout icon="fad fa-gauge" theme="warn">
  ### **WebSocket 요청 수 제한 관리**

  WebSocket은 REST API와 달리 잔여 요청 수를 별도로 제공하지 않습니다. 클라이언트에서 WebSocket 연결 및 데이터 요청 메시지의 전송 횟수를 관리하여 요청 수 제한을 준수해 주세요. 요청 수 제한에 도달한 경우 일정 시간 대기한 후 다시 요청해 주세요.
</Callout>

# Sibling pages

* [현재가 (Ticker)](https://docs.upbit.com/kr/reference/websocket-ticker.md)
* [체결 (Trade)](https://docs.upbit.com/kr/reference/websocket-trade.md)
* [호가 (Orderbook)](https://docs.upbit.com/kr/reference/websocket-orderbook.md)
* [캔들 (Candle)](https://docs.upbit.com/kr/reference/websocket-candle.md)
* [내 주문 및 체결 (MyOrder)](https://docs.upbit.com/kr/reference/websocket-myorder.md)
* [내 자산 (MyAsset)](https://docs.upbit.com/kr/reference/websocket-myasset.md)
* [구독 중인 스트림 목록 조회(Subscriptions)](https://docs.upbit.com/kr/reference/list-subscriptions.md)