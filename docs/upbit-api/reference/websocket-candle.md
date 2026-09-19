---
updatedAt: 2026-09-09T06:01:05.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 캔들 (Candle)

캔들 데이터를 WebSocket으로 수신하기 위한 요청 및 구독 데이터 예시를 제공합니다.

## **WebSocket Endpoint**

| 구분     | Endpoint                           |
| ------ | ---------------------------------- |
| Public | `wss://api.upbit.com/websocket/v1` |

## 캔들 실시간 스트림 안내

### 데이터 전송 주기

캔들 데이터의 실시간 스트림 전송 주기는 **1초**입니다.

<Callout icon="fad fa-clock" theme="info">
  ### **캔들 데이터는 매초 전송되지 않을 수 있습니다.**

  **해당 시간대에 체결이 발생하여 캔들 데이터가 변경된 경우에만 전송됩니다.&#x20;**&#xB530;라서 1초가 지나더라도 체결이 발생하지 않으면 실시간 캔들 데이터가 전송되지 않습니다.
</Callout>

### 데이터 생성 방식

요청 시점에 해당 시간 단위의 캔들이 아직 생성되지 않은 경우, 가장 최근에 생성된 캔들이 스냅샷으로 반환됩니다.

**예시:&#x20;**`candle.3m`

12:00:00 기준 3분봉이 존재하고, 12:03:00 이후 체결이 발생하지 않은 상태에서 12:04:00에 `candle.3m`을 요청하는 경우:

* **12:04:00** → 12:00:00 기준 3분봉을 스냅샷으로 반환
* **12:04:05** → 첫 체결 발생 후 12:03:00 기준 3분봉 생성
* **12:04:06** → 다음 전송 주기에 12:03:00 기준 3분봉 데이터 전송

<Callout icon="fad fa-magnifying-glass-arrows-rotate" theme="info">
  ### **같은 캔들 기준 시각의 데이터가 여러 번 전송될 수 있습니다.**

  체결 시점에 따라 동일한 `candle_date_time_utc` 또는 `candle_date_time_kst`의 캔들 데이터가 여러 번 전송될 수 있습니다. 동일한 캔들 기준 시각에서는 가장 마지막으로 수신한 데이터가 최신 데이터이므로, 해당 필드를 기준으로 값을 업데이트해주세요.
</Callout>

## Request 메세지 형식

캔들 데이터 수신을 요청하기 위해서는 WebSocket 연결 이후 아래 구조의 JSON Object를 생성한 뒤 요청 메세지의 Data Type Object로 포함하여 전송해야 합니다. Ticket, Format 필드를 포함한 전체 WebSocket 데이터 요청 메세지 명세는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide">WebSocket 사용 안내</Anchor> 문서를 참고해주세요.

| 필드명                | 타입          | 내용                                                                                                                                                                                                                          | 필수 여부    | 기본 값    |
| ------------------ | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| type               | String      | 수신할 캔들 단위.<br />`candle.1s`: 초봉<br />`candle.1m`: 1분봉<br />`candle.3m`: 3분봉<br />`candle.5m`: 5분봉<br />`candle.10m`: 10분봉<br />`candle.15m`: 15분봉<br />`candle.30m`: 30분봉<br />`candle.60m`: 60분봉<br />`candle.240m`: 240분봉 | Required |         |
| codes              | List:String | 현재가 데이터를 수신할 페어 코드 목록. 페어 코드는 대문자로 입력해야 합니다.                                                                                                                                                                                | Required |         |
| is\_only\_snapshot | Boolean     | `true`로 설정하면 요청 시점의 현재가 스냅샷 데이터만 1회 수신합니다.                                                                                                                                                                                  | Optional | `false` |
| is\_only\_realtime | Boolean     | `true`로 설정하면 현재가 스냅샷 없이 실시간 스트림 데이터만 수신합니다.                                                                                                                                                                                 | Optional | `false` |
| format             | String      | 수신하고자 하는 데이터 포맷입니다. <br />`DEFAULT` : 기본 포맷.<br />`SIMPLE` : 간략한 포맷. 각 필드가 축약어 형태로 반환됩니다.<br />`JSON_LIST` : 리스트 포맷.<br />`SIMPLE_LIST` : 축약어 형태의 리스트 포맷.                                                                   | Required |         |

<br />

## 구독 데이터 명세

캔들 데이터는 스냅샷 또는 실시간 스트림 형태로 아래와 같이 반환됩니다.

| **필드명**                    | **축약형** | **내용**                     | **타입** | **값**                                                                                                                                                                                                       |
| -------------------------- | ------- | -------------------------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| type                       | ty      | 데이터 타입                     | String | `candle.1s`: 초봉<br />`candle.1m`: 1분봉<br />`candle.3m`: 3분봉<br />`candle.5m`: 5분봉<br />`candle.10m`: 10분봉<br />`candle.15m`: 15분봉<br />`candle.30m`: 30분봉<br />`candle.60m`: 60분봉<br />`candle.240m`: 240분봉 |
| code                       | cd      | 마켓 코드<br />(예시: `KRW-BTC`) | String |                                                                                                                                                                                                             |
| candle\_date\_time\_utc    | cdttmu  | 캔들 기준 시각(UTC 기준)           | String | `yyyy-MM-dd'T'HH:mm:ss`                                                                                                                                                                                     |
| candle\_date\_time\_kst    | cdttmk  | 캔들 기준 시각(KST 기준)           | String | `yyyy-MM-dd'T'HH:mm:ss`                                                                                                                                                                                     |
| opening\_price             | op      | 시가                         | Double |                                                                                                                                                                                                             |
| high\_price                | hp      | 고가                         | Double |                                                                                                                                                                                                             |
| low\_price                 | lp      | 저가                         | Double |                                                                                                                                                                                                             |
| trade\_price               | tp      | 종가                         | Double |                                                                                                                                                                                                             |
| candle\_acc\_trade\_volume | catv    | 누적 거래량                     | Double |                                                                                                                                                                                                             |
| candle\_acc\_trade\_price  | catp    | 누적 거래 금액                   | Double |                                                                                                                                                                                                             |
| timestamp                  | tms     | 타임스탬프 (ms)                 | Long   |                                                                                                                                                                                                             |
| stream\_type               | st      | 스트림 타입                     | String | `SNAPSHOT`: 스냅샷<br />`REALTIME`: 실시간                                                                                                                                                                        |

<br />

## 예시

#### KRW-BTC 1분봉·15분봉 구독

> KRW-BTC의 1분봉과 15분봉을 함께 구독하여 단기 가격 움직임과 보다 긴 시간 단위의 가격 흐름을 동시에 확인하는 예시입니다. ticket은 요청을 식별하기 위한 값으로 원하는 문자열을 직접 지정할 수 있습니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "candle-monitor"
  },
  {
    "type": "candle.1m",
    "codes": ["KRW-BTC"]
  },
  {
    "type": "candle.15m",
    "codes": ["KRW-BTC"]
  },
  {
    "format": "DEFAULT"
  }
]
```

##### 수신 메시지 예제

```json
{
  "type": "candle.1m",
  "code": "KRW-BTC",
  "candle_date_time_utc": "2026-08-26T08:08:00",
  "candle_date_time_kst": "2026-08-26T17:08:00",
  "opening_price": 109715000.0,
  "high_price": 109756000.0,
  "low_price": 109714000.0,
  "trade_price": 109756000.0,
  "candle_acc_trade_volume": 0.71326612,
  "candle_acc_trade_price": 78269528.64757,
  "timestamp": 1787731717049,
  "stream_type": "SNAPSHOT"
}
{
  "type": "candle.15m",
  "code": "KRW-BTC",
  "candle_date_time_utc": "2026-08-26T08:00:00",
  "candle_date_time_kst": "2026-08-26T17:00:00",
  "opening_price": 109753000.0,
  "high_price": 109756000.0,
  "low_price": 109667000.0,
  "trade_price": 109756000.0,
  "candle_acc_trade_volume": 3.63363563,
  "candle_acc_trade_price": 398651814.32878,
  "timestamp": 1787731717049,
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
* [체결 (Trade)](https://docs.upbit.com/kr/reference/websocket-trade.md)
* [호가 (Orderbook)](https://docs.upbit.com/kr/reference/websocket-orderbook.md)
* [공지사항(Announcement)](https://docs.upbit.com/kr/reference/websocket-announcement.md)
* [내 주문 및 체결 (MyOrder)](https://docs.upbit.com/kr/reference/websocket-myorder.md)
* [내 자산 (MyAsset)](https://docs.upbit.com/kr/reference/websocket-myasset.md)
* [구독 중인 스트림 목록 조회(Subscriptions)](https://docs.upbit.com/kr/reference/list-subscriptions.md)