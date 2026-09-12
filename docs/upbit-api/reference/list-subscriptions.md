---
updatedAt: 2026-09-09T06:08:43.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 구독 중인 스트림 목록 조회(Subscriptions)

WebSocket 연결을 통해 구독중인 데이터 스트림 항목을 확인할 수 있습니다.

## **WebSocket Endpoint**

| 구분      | Endpoint                                   |
| ------- | ------------------------------------------ |
| Public  | `wss://api.upbit.com/websocket/v1`         |
| Private | `wss://api.upbit.com/websocket/v1/private` |

<br />

## 구독 목록 조회 안내

일반적인 데이터 구독 요청은 `type` 필드를 사용하지만, 구독 중인 스트림 목록 조회는 `method` 필드에 `LIST_SUBSCRIPTIONS`를 지정하여 요청합니다.

조회 결과는 **현재 요청을 전송한 WebSocket 연결에서 구독 중인 스트림**을 기준으로 반환됩니다.

<Callout icon="fad fa-code" theme="warn">
  **현재 구독 중인 스트림과 동일한 Format을 사용해주세요.**

  구독 목록 조회 요청에서 Format을 변경하면 기존에 구독 중인 실시간 스트림의 응답 Format도 함께 변경됩니다.

  예를 들어 SIMPLE 포맷으로 실시간 데이터를 수신하는 중 LIST_SUBSCRIPTIONS 요청의 Format을 DEFAULT로 지정하면, 이후 구독 중인 실시간 스트림도 DEFAULT 포맷으로 반환됩니다.
</Callout>

<Callout icon="fad fa-gauge-high" theme="info">
  **요청 수 제한이 적용됩니다.**

  구독 목록 조회 요청도 WebSocket 요청 수 제한에 포함됩니다.
</Callout>

<br />

현재 WebSocket 연결에서 구독 중인 스트림 목록을 조회하려면 아래 형식의 Operation Object를 요청 메시지에 포함하여 전송합니다. Ticket, Format 필드를 포함한 전체 WebSocket 요청 메시지 명세는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide">WebSocket 사용 안내</Anchor> 문서를 참고해주세요.

| 필드명    | 타입     | 내용                                                         | 필수 여부    | 기본 값 |
| ------ | ------ | ---------------------------------------------------------- | -------- | ---- |
| method | String | 요청 메서드.<br />구독 중인 스트림 목록 조회는 `LIST_SUBSCRIPTIONS`로 지정합니다. | Required |      |

### 예시

```json
[
  {
    "ticket": "0e66c0ac-7e13-43ef-91fb-2a87c2956c49"
  },
  {
    "method": "LIST_SUBSCRIPTIONS"
  }
]
```

***

## 응답 명세

현재 WebSocket 연결에서 구독 중인 스트림 목록이 `result`에 반환됩니다.

| 필드명          | 축약형      | 내용             | 타입              | 값                    |
| ------------ | -------- | -------------- | --------------- | -------------------- |
| method       | mthd     | 요청 메서드         | String          | `LIST_SUBSCRIPTIONS` |
| result       | rslt     | 구독 중인 스트림 목록   | List of Objects |                      |
| result.type  | rslt.ty  | 데이터 타입         | String          |                      |
| result.codes | rslt.cds | 구독 중인 페어 코드 목록 | List of String  |                      |
| result.level | rslt.lv  | 호가 모아보기 단위     | Double          |                      |
| `ticket`     | tck`t`   | 요청 식별 값        | String          |                      |

<br />

## 예시

#### **Public WebSocket**

##### 구독 요청 예제

```json
{
  "method": "LIST_SUBSCRIPTIONS",
  "result": [
    {
      "type": "ticker",
      "codes": [
        "KRW-BTC",
        "KRW-ETH"
      ]
    },
    {
      "type": "orderbook",
      "codes": [
        "KRW-BTC",
        "KRW-ETH"
      ],
      "level": 0
    }
  ],
  "ticket": "unique uuid"
}
```

#### **Public WebSocket**

##### 구독 요청 예제

```json
{
  "method": "LIST_SUBSCRIPTIONS",
  "result": [
    {
      "type": "myAsset"
    },
    {
      "type": "myOrder",
      "codes": [
        "KRW-BTC",
        "KRW-ETH"
      ]
    }
  ],
  "ticket": "unique uuid"
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

| Rate Limit 그룹       | 정책                | 적용 단위                     |
| ------------------- | ----------------- | ------------------------- |
| `websocket-connect` | 초당 최대 5회          | Public - IP, Private - 포켓 |
| `websocket-message` | 초당 최대 5회, 분당 100회 | 커넥션                       |

<Callout icon="fad fa-gauge" theme="warn">
  ### **WebSocket 요청 수 제한 관리**

  WebSocket은 REST API와 달리 잔여 요청 수를 별도로 제공하지 않습니다. 클라이언트에서 WebSocket 연결 및 데이터 요청 메시지의 전송 횟수를 관리하여 요청 수 제한을 준수해 주세요.<br />요청 수 제한에 도달한 경우 일정 시간 대기한 후 다시 요청해 주세요.
</Callout>

# Sibling pages

* [현재가 (Ticker)](https://docs.upbit.com/kr/reference/websocket-ticker.md)
* [체결 (Trade)](https://docs.upbit.com/kr/reference/websocket-trade.md)
* [호가 (Orderbook)](https://docs.upbit.com/kr/reference/websocket-orderbook.md)
* [캔들 (Candle)](https://docs.upbit.com/kr/reference/websocket-candle.md)
* [공지사항(Announcement)](https://docs.upbit.com/kr/reference/websocket-announcement.md)
* [내 주문 및 체결 (MyOrder)](https://docs.upbit.com/kr/reference/websocket-myorder.md)
* [내 자산 (MyAsset)](https://docs.upbit.com/kr/reference/websocket-myasset.md)