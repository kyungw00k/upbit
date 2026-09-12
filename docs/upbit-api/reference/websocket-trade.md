---
updatedAt: 2026-09-02T08:00:35.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 체결 (Trade)

체결(Trade) 데이터를 WebSocket으로 수신하기 위한 요청 및 구독 데이터 예시를 제공합니다.

## **WebSocket Endpoint**

| 구분     | Endpoint                           |
| ------ | ---------------------------------- |
| Public | `wss://api.upbit.com/websocket/v1` |

<br />

## Request 메세지 형식

체결 데이터 수신을 요청하기 위해서는 WebSocket 연결 이후 아래 구조의 JSON Object를 생성한 뒤 요청 메세지의 Data Type Object로 포함하여 전송해야 합니다. Ticket, Format 필드를 포함한 전체 WebSocket 데이터 요청 메세지 명세는 <Anchor target="_blank" href="ref:websocket-guide">WebSocket 사용 안내</Anchor> 문서를 참고해주세요.

| **필드명**            | **타입**      | **내용**                                                                                                                                                    | **필수 여부** | **기본 값** |
| ------------------ | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- | -------- |
| type               | String      | 수신할 데이터 타입.<br />체결 데이터를 요청하는 경우 `trade`로 지정합니다.                                                                                                          | Required  |          |
| codes              | List:String | 수신할 페어 코드 목록. <br />페어 코드는 대문자로 입력해야 합니다.                                                                                                                 | Required  |          |
| is\_only\_snapshot | Boolean     | `true`로 설정하면 요청 시점의 현재가 스냅샷 데이터만 1회 수신합니다.                                                                                                                | Optional  | `false`  |
| is\_only\_realtime | Boolean     | `true`로 설정하면 현재가 스냅샷 없이 실시간 스트림 데이터만 수신합니다.                                                                                                               | Optional  | `false`  |
| format             | String      | 수신하고자 하는 데이터 포맷입니다. <br />`DEFAULT` : 기본 포맷.<br />`SIMPLE` : 간략한 포맷. 각 필드가 축약어 형태로 반환됩니다.<br />`JSON_LIST` : 리스트 포맷.<br />`SIMPLE_LIST` : 축약어 형태의 리스트 포맷. | Required  |          |

<br />

## 구독 데이터 명세

체결 데이터 스냅샷 또는 실시간 스트림 데이터는 아래와 같이 반환됩니다.

| **필드명**              | **축약형** | **내용**                     | **타입** | **값**                                                        |
| -------------------- | ------- | -------------------------- | ------ | ------------------------------------------------------------ |
| type                 | ty      | 데이터 타입                     | String | `trade`                                                      |
| code                 | cd      | 페어 코드<br />(예시: `KRW-BTC`) | String |                                                              |
| trade\_price         | tp      | 체결 가격                      | Double |                                                              |
| trade\_volume        | tv      | 체결량                        | Double |                                                              |
| ask\_bid             | ab      | 매수/매도 구분                   | String | `ASK`<br />: 매도<br />`BID`<br />: 매수                         |
| prev\_closing\_price | pcp     | 전일 종가                      | Double |                                                              |
| change               | c       | 전일 종가 대비 가격 변동 방향          | String | `RISE`<br />: 상승<br />`EVEN`<br />: 보합<br />`FALL`<br />: 하락 |
| change\_price        | cp      | 전일 대비 가격 변동의 절대값           | Double |                                                              |
| trade\_date          | td      | 체결 일자(UTC 기준)              | String | `yyyy-MM-dd`                                                 |
| trade\_time          | ttm     | 체결 시각(UTC 기준)              | String | `HH:mm:ss`                                                   |
| trade\_timestamp     | ttms    | 체결 타임스탬프(ms)               | Long   |                                                              |
| timestamp            | tms     | 타임스탬프(ms)                  | Long   |                                                              |
| sequential\_id       | sid     | 체결 번호(Unique)              | Long   |                                                              |
| best\_ask\_price     | bap     | 최우선 매도 호가                  | Double |                                                              |
| best\_ask\_size      | bas     | 최우선 매도 잔량                  | Double |                                                              |
| best\_bid\_price     | bbp     | 최우선 매수 호가                  | Double |                                                              |
| best\_bid\_size      | bbs     | 최우선 매수 잔량                  | Double |                                                              |
| stream\_type         | st      | 스트림 타입                     | String | `SNAPSHOT`<br />: 스냅샷<br />`REALTIME`<br />: 실시간             |

<br />

## 예시

#### KRW-BTC, KRW-ETH 실시간 체결 데이터 구독

> KRW-BTC와 KRW-ETH의 체결 가격, 체결 수량, 매수·매도 구분 등 실시간 체결 데이터를 수신하는 예시입니다. ticket은 요청을 식별하기 위한 값으로 원하는 문자열을 직접 지정할 수 있습니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "trade-monitor"
  },
  {
    "type": "trade",
    "codes": ["KRW-BTC", "KRW-ETH"]
  },
  {
    "format": "DEFAULT"
  }
]
```

##### 수신 메시지 예제

```json
{
  "type": "trade",
  "code": "KRW-BTC",
  "timestamp": 1787728554853,
  "trade_date": "2026-08-26",
  "trade_time": "07:15:54",
  "trade_timestamp": 1787728554797,
  "trade_price": 109847000.0,
  "trade_volume": 0.00509799,
  "ask_bid": "ASK",
  "prev_closing_price": 109165000.0,
  "change": "RISE",
  "change_price": 682000.0,
  "sequential_id": 17877285547970000,
  "best_ask_price": 109867000,
  "best_ask_size": 0.0099545,
  "best_bid_price": 109847000,
  "best_bid_size": 0.1583673,
  "stream_type": "SNAPSHOT"
}
{
  "type": "trade",
  "code": "KRW-ETH",
  "timestamp": 1787728551752,
  "trade_date": "2026-08-26",
  "trade_time": "07:15:51",
  "trade_timestamp": 1787728551696,
  "trade_price": 3427000.0,
  "trade_volume": 0.00883869,
  "ask_bid": "BID",
  "prev_closing_price": 3394000.0,
  "change": "RISE",
  "change_price": 33000.0,
  "sequential_id": 17877285516960000,
  "best_ask_price": 3427000,
  "best_ask_size": 1.12778516,
  "best_bid_price": 3426000,
  "best_bid_size": 10.05419423,
  "stream_type": "SNAPSHOT"
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
| `websocket-connect` | 초당 최대 5회          | IP    |
| `websocket-message` | 초당 최대 5회, 분당 100회 | 커넥션   |

<Callout icon="fad fa-gauge" theme="warn">
  ### **WebSocket 요청 수 제한 관리**

  WebSocket은 REST API와 달리 잔여 요청 수를 별도로 제공하지 않습니다. 클라이언트에서 WebSocket 연결 및 데이터 요청 메시지의 전송 횟수를 관리하여 요청 수 제한을 준수해 주세요. 요청 수 제한에 도달한 경우 일정 시간 대기한 후 다시 요청해 주세요.
</Callout>

# Sibling pages

* [현재가 (Ticker)](https://docs.upbit.com/kr/reference/websocket-ticker.md)
* [호가 (Orderbook)](https://docs.upbit.com/kr/reference/websocket-orderbook.md)
* [캔들 (Candle)](https://docs.upbit.com/kr/reference/websocket-candle.md)
* [공지사항(Announcement)](https://docs.upbit.com/kr/reference/websocket-announcement.md)
* [내 주문 및 체결 (MyOrder)](https://docs.upbit.com/kr/reference/websocket-myorder.md)
* [내 자산 (MyAsset)](https://docs.upbit.com/kr/reference/websocket-myasset.md)
* [구독 중인 스트림 목록 조회(Subscriptions)](https://docs.upbit.com/kr/reference/list-subscriptions.md)