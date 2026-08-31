---
updatedAt: 2026-08-31T12:33:13.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 내 자산 (MyAsset)

내 자산 데이터를 WebSocket으로 수신하기 위한 요청 및 구독 데이터 예시를 제공합니다.

## **WebSocket Endpoint**

| 구분      | Endpoint                                   |
| ------- | ------------------------------------------ |
| Private | `wss://api.upbit.com/websocket/v1/private` |

<Callout icon="fad fa-circle-info" theme="info">
  ### **자산 변동이 없으면 데이터가 전송되지 않습니다.**

  내 자산 데이터는 **보유 자산에 변동이 발생한 경우에만** 실시간으로 전송됩니다. 따라서 WebSocket 연결 후 자산 변동이 없다면 데이터가 수신되지 않는 것이 정상입니다. 데이터가 없는 동안에도 연결이 유지될 수 있도록 [WebSocket 사용 안내의 연결 유지 항목](https://docs.upbit.com/kr/v1.6.3_add_announcement_websocket/reference/websocket-guide)을 참고해주세요.
</Callout>

<Callout icon="fad fa-clock" theme="warn">
  ### **최초 구독 시 데이터 수신이 지연될 수 있습니다.**

  내 자산 WebSocket을 처음 구독하는 경우에는 자산에 변동이 발생하더라도 **초기 수분 동안 데이터 수신이 지연될 수 있습니다.**

  **예시**

  - **00:00** → 최초 연결 및 자산 변동 발생
  - **00:05** → 실시간 스트림 데이터 전송 시작
  - **00:10** → 재연결 후에는 자산 변동 발생 시 실시간으로 데이터 수신

  최초 연결 후에는 재연결 등을 통해 **데이터가 정상적으로 수신되는지 확인한 뒤 사용해주세요.**
</Callout>

<br />

## Request 메세지 형식

내 자산 데이터 수신을 요청하기 위해서는 WebSocket 연결 이후 아래 구조의 JSON Object를 생성한 뒤 요청 메세지의 Data Type Object로 포함하여 전송해야 합니다. Ticket, Format 필드를 포함한 전체 WebSocket 데이터 요청 메세지 명세는 [WebSocket 사용 안내](https://docs.upbit.com/kr/reference/websocket-guide) 문서를 참고해주세요.

<Callout icon="fad fa-circle-info" theme="error">
  ### **내 자산 타입 데이터 구독 요청은 페어 코드("codes") 파라미터를 지원하지 않습니다.**

  다른 데이터 항목 구독 요청과 달리 내 자산 타입 데이터 구독 요청에 페어 코드 파라미터를 포함하는 경우 "WRONG_FORMAT" 에러가 발생하오니 사용에 주의바랍니다.
</Callout>

| 필드명    | 타입     | 내용                                                                                                                                                        | 필수 여부    | 기본 값 |
| ------ | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---- |
| type   | String | 수신할 데이터 타입.<br />내 자산 데이터를 요청하는 경우 `myAsset`으로 지정합니다.                                                                                                     | Required |      |
| format | String | 수신하고자 하는 데이터 포맷입니다. <br />`DEFAULT` : 기본 포맷.<br />`SIMPLE` : 간략한 포맷. 각 필드가 축약어 형태로 반환됩니다.<br />`JSON_LIST` : 리스트 포맷.<br />`SIMPLE_LIST` : 축약어 형태의 리스트 포맷. | Required |      |

<br />

## 구독 데이터 명세

| 필드명              | 축약형    | 내용            | 타입              | 값                |
| ---------------- | ------ | ------------- | --------------- | ---------------- |
| type             | ty     | 타입            | String          | `myAsset`        |
| asset\_uuid      | astuid | 자산 고유 식별자     | String          |                  |
| assets           | ast    | 자산 목록         | List of Objects |                  |
| assets.currency  | ast.cu | 화폐 코드         | String          |                  |
| assets.balance   | ast.b  | 주문가능 수량       | Double          |                  |
| assets.locked    | ast.l  | 주문 중 묶여있는 수량  | Double          |                  |
| asset\_timestamp | asttms | 자산 타임스탬프 (ms) | Long            |                  |
| timestamp        | tms    | 타임스탬프 (ms)    | Long            |                  |
| stream\_type     | st     | 스트림 타입        | String          | `REALTIME` : 실시간 |

<br />

## 예시

#### 내 자산 변동 구독

> 내 자산의 잔액 또는 주문 가능 금액에 변동이 발생할 때 자산 정보를 실시간으로 수신하는 예시입니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "my-asset-monitor"
  },
  {
    "type": "myAsset"
  }
]
```

##### 수신 메시지 예제

```json
{
    "type": "myAsset",
    "asset_uuid": "e635f223-1609-4969-8fb6-4376937baad6",
    "assets": [
      {
        "currency": "KRW",
        "balance": 1386929.37231066771348207123,
        "locked": 10329.670127489597585685
      }
    ],
    "asset_timestamp": 1710146517259,
    "timestamp": 1710146517267,
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
* [내 주문 및 체결 (MyOrder)](https://docs.upbit.com/kr/reference/websocket-myorder.md)
* [구독 중인 스트림 목록 조회(Subscriptions)](https://docs.upbit.com/kr/reference/list-subscriptions.md)