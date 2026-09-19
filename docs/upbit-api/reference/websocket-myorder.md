---
updatedAt: 2026-09-09T06:07:41.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 내 주문 및 체결 (MyOrder)

내 주문 및 체결 데이터를 WebSocket으로 수신하기 위한 요청 및 구독 데이터 예시를 제공합니다.

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
  ### **주문 또는 체결이 없으면 데이터가 전송되지 않습니다.**

  내 주문 및 체결 데이터의 경우 실제 주문 또는 체결이 발생할 때만 해당 내용이 실시간 스트림으로 전송됩니다. 장시간 데이터가 수신되지 않는 경우에도 연결이 유지되도록 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide">WebSocket 사용 안내</Anchor>의 연결 유지 항목을 참고해주세요.
</Callout>

<br />

## Request 메세지 형식

내 주문 및 체결 데이터 수신을 요청하기 위해서는 WebSocket 연결 이후 아래 구조의 JSON Object를 생성한 뒤 요청 메세지의 Data Type Object로 포함하여 전송해야 합니다. Ticket, Format 필드를 포함한 전체 WebSocket 데이터 요청 메세지 명세는 [WebSocket 사용 안내](https://docs.upbit.com/kr/reference/websocket-guide) 문서를 참고해주세요.

| 필드명    | 타입          | 내용                                                                                                                                                        | 필수 여부    | 기본 값                                    |
| ------ | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | --------------------------------------- |
| type   | String      | 수신할 데이터 타입.<br />내 주문 및 체결 데이터를 요청하는 경우 `myOrder`로 지정합니다.                                                                                                 | Required |                                         |
| codes  | List:String | 수신할 페어 코드 목록.<br />페어 코드는 대문자로 입력해야 합니다.                                                                                                                  | Optional | 생략하거나 빈 배열로 요청할 경우 모든 마켓에 대한 정보를 수신합니다. |
| format | String      | 수신하고자 하는 데이터 포맷입니다. <br />`DEFAULT` : 기본 포맷.<br />`SIMPLE` : 간략한 포맷. 각 필드가 축약어 형태로 반환됩니다.<br />`JSON_LIST` : 리스트 포맷.<br />`SIMPLE_LIST` : 축약어 형태의 리스트 포맷. | Required |                                         |

<br />

## 구독 데이터 명세

<Callout icon="fad fa-circle-info" theme="info">
  ### **신규 필드 추가 안내 (2025.07.02)**

  자전거래 체결 방지 (Self-Match Prevention) 기능 추가로 인해 주문 데이터 필드가 아래와 같이 추가됩니다. 자세한 사항은 관련 공지 참고 부탁드립니다.<Anchor target="_blank" href="/kr/docs/smp"> \[SMP 상세 설명 바로가기]</Anchor>

  - smp_type : `reduce`, `cancel_maker`, `cancel_taker`
  - 주문 상태 state 필드에 `prevented` (체결 방지) 타입이 신규로 추가됩니다.
  - 신규 주문 조건 `post_only`가 추가되어 time_in_force 필드에 `post_only`타입이 신규로 추가됩니다.
</Callout>

| 필드명               | 축약형  | 내용                                                              | 타입      | 값                                                                                                                 |
| ----------------- | ---- | --------------------------------------------------------------- | ------- | ----------------------------------------------------------------------------------------------------------------- |
| type              | ty   | 타입                                                              | String  | `myOrder`                                                                                                         |
| code              | cd   | 페어 코드<br />(예시: `KRW-BTC`)                                      | String  |                                                                                                                   |
| uuid              | uid  | 주문의 유일 식별자                                                      | String  |                                                                                                                   |
| ask\_bid          | ab   | 매수/매도 구분                                                        | String  | `ASK` : 매도<br />`BID` : 매수                                                                                        |
| order\_type       | ot   | 주문 타입                                                           | String  | `limit`: 지정가 주문 `price`: 시장가 매수 주문 <br />`market`: 시장가 매도 주문 <br />`best`: 최유리 지정가 주문                             |
| state             | s    | 주문 상태                                                           | String  | `wait`: 체결 대기<br />`watch`: 예약 주문 대기 <br />`trade`: 체결 발생 `done`: 전체 체결 완료<br />`cancel`: 주문 취소 `prevented`:체결 방지 |
| trade\_uuid       | tuid | 체결의 유일 식별자                                                      | String  |                                                                                                                   |
| price             | p    | 주문 가격 또는 체결 가격(state: trade 일 때)                                | Double  |                                                                                                                   |
| avg\_price        | ap   | 평균 체결 가격                                                        | Double  |                                                                                                                   |
| volume            | v    | 주문량 또는 체결량 (state: trade 일 때)                                   | Double  |                                                                                                                   |
| remaining\_volume | rv   | 체결 후 주문 잔량                                                      | Double  |                                                                                                                   |
| executed\_volume  | ev   | 체결된 수량                                                          | Double  |                                                                                                                   |
| trades\_count     | tc   | 해당 주문에 걸린 체결 수                                                  | Integer |                                                                                                                   |
| reserved\_fee     | rsf  | 수수료로 예약된 비용                                                     | Double  |                                                                                                                   |
| remaining\_fee    | rmf  | 남은 수수료                                                          | Double  |                                                                                                                   |
| paid\_fee         | pf   | 사용된 수수료                                                         | Double  |                                                                                                                   |
| locked            | l    | 거래에 사용중인 비용                                                     | Double  |                                                                                                                   |
| executed\_funds   | ef   | 체결된 금액                                                          | Double  |                                                                                                                   |
| time\_in\_force   | tif  | IOC, FOK, POST ONLY 설정                                          | String  | `ioc`, `fok`, `post_only`                                                                                         |
| trade\_fee        | tf   | 체결 시 발생한 수수료 (state:trade가 아닐 경우 null)                          | Double  |                                                                                                                   |
| is\_maker         | im   | 체결이 발생한 주문의 메이커/테이커 여부 (state:trade가 아닐 경우 null)                | Boolean | `true` : 메이커 주문 `false` : 테이커 주문                                                                                  |
| identifier        | id   | 클라이언트 지정 주문 식별자                                                 | String  |                                                                                                                   |
| smp\_type         | smpt | 자전거래 체결 방지 타입 (동일 회원의 메이커/테이커 주문 간 체결 방지)                       | String  | `reduce`: 체결 수량만큼 주문 수량 차감 후 진행<br />`cancel_maker`: 메이커 주문 취소 `cancel_taker`: 테이커 주문 취소                          |
| prevented\_volume | pv   | 자전거래 체결 방지로 인해 취소된 주문 수량                                        | Double  |                                                                                                                   |
| prevented\_locked | pl   | (매수 시)자전거래 체결 방지 설정으로 인해 취소된 금액 (매도 시)자전거래 체결 방지 설정으로 인해 취소된 수량 | Double  |                                                                                                                   |
| trade\_timestamp  | ttms | 체결 타임스탬프 (ms)                                                   | Long    |                                                                                                                   |
| order\_timestamp  | otms | 주문 타임스탬프 (ms)                                                   | Long    |                                                                                                                   |
| timestamp         | tms  | 타임스탬프 (ms)                                                      | Long    |                                                                                                                   |
| stream\_type      | st   | 스트림 타입                                                          | String  | `REALTIME`: 실시간 스트림<br />`SNAPSHOT`: 스냅샷                                                                          |

<br />

## 예시

#### 내 전체 주문 및 체결 구독

> 모든 마켓에서 발생하는 내 주문 및 체결 데이터를 실시간으로 수신하는 예시입니다. codes를 생략하거나 빈 배열로 지정하면 모든 마켓의 데이터를 수신할 수 있습니다. ticket은 요청을 식별하기 위한 값으로 원하는 문자열을 직접 지정할 수 있습니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "my-order-monitor"
  },
  {
    "type": "myOrder"
  }
]
```

```json
[
  {
    "ticket": "my-order-monitor"
  },
  {
    "type": "myOrder",
    "codes": []
  }
]
```

##### 수신 메시지 예제

```json
{
  "type": "myOrder",
  "code": "KRW-BTC",
  "uuid": "ac2dc2a3-fce9-40a2-a4f6-5987c25c438f",
  "ask_bid": "BID",
  "order_type": "limit",
  "state": "trade",
  "trade_uuid": "68315169-fba4-4175-ade3-aff14a616657",
  "price": 0.001453,
  "avg_price": 0.00145372,
  "volume": 30925891.29839369,
  "remaining_volume": 29968038.09235948,
  "executed_volume": 30925891.29839369,
  "trades_count": 1,
  "reserved_fee": 44.23943970238218,
  "remaining_fee": 21.77177967409916,
  "paid_fee": 22.467660028283017,
  "locked": 43565.33112787242,
  "executed_funds": 44935.32005656603,
  "time_in_force": null,
  "trade_fee": 22.467660028283017,
  "is_maker": true,
  "identifier": "test-1",
  "smp_type": "cancel_maker",
  "prevented_volume": 1.174291929,
  "prevented_locked": 0.001706246173,
  "trade_timestamp": 1710751590421,
  "order_timestamp": 1710751590000,
  "timestamp": 1710751597500,
  "stream_type": "REALTIME"
}
```

<br />

#### 특정 페어의 내 주문 및 체결 구독

> KRW-BTC에서 발생하는 내 주문 및 체결 데이터를 실시간으로 수신하는 예시입니다. codes에 페어를 지정하면 해당 페어의 주문 및 체결 데이터만 구독할 수 있습니다. ticket은 요청을 식별하기 위한 값으로 원하는 문자열을 직접 지정할 수 있습니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "my-order-monitor"
  },
  {
    "type": "myOrder",
    "codes": ["KRW-BTC"]
  }
]
```

##### 수신 데이터 예제

```json
{
  "type": "myOrder",
  "code": "KRW-BTC",
  "uuid": "ac2dc2a3-fce9-40a2-a4f6-5987c25c438f",
  "ask_bid": "BID",
  "order_type": "limit",
  "state": "trade",
  "trade_uuid": "68315169-fba4-4175-ade3-aff14a616657",
  "price": 0.001453,
  "avg_price": 0.00145372,
  "volume": 30925891.29839369,
  "remaining_volume": 29968038.09235948,
  "executed_volume": 30925891.29839369,
  "trades_count": 1,
  "reserved_fee": 44.23943970238218,
  "remaining_fee": 21.77177967409916,
  "paid_fee": 22.467660028283017,
  "locked": 43565.33112787242,
  "executed_funds": 44935.32005656603,
  "time_in_force": null,
  "trade_fee": 22.467660028283017,
  "is_maker": true,
  "identifier": "test-1",
  "smp_type": "cancel_maker",
  "prevented_volume": 1.174291929,
  "prevented_locked": 0.001706246173,
  "trade_timestamp": 1710751590421,
  "order_timestamp": 1710751590000,
  "timestamp": 1710751597500,
  "stream_type": "REALTIME"
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
* [공지사항(Announcement)](https://docs.upbit.com/kr/reference/websocket-announcement.md)
* [내 자산 (MyAsset)](https://docs.upbit.com/kr/reference/websocket-myasset.md)
* [구독 중인 스트림 목록 조회(Subscriptions)](https://docs.upbit.com/kr/reference/list-subscriptions.md)