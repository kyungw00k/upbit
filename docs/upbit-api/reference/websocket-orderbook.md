---
updatedAt: 2026-09-09T06:00:44.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 호가 (Orderbook)

호가 데이터를 WebSocket으로 수신하기 위한 요청 및 구독 데이터 예시를 제공합니다.

## **WebSocket Endpoint**

| 구분     | Endpoint                           |
| ------ | ---------------------------------- |
| Public | `wss://api.upbit.com/websocket/v1` |

## 호가 모아보기 (level)

원화마켓(KRW)에서만 지원하며, `level`을 지정하면 호가를 해당 가격 단위로 묶어 조회할 수 있습니다.

**예시:&#x20;**`KRW-BTC`**,&#x20;**`level=100000`

* `ask_price`, `bid_price`: 100,000 KRW 단위로 반환
* `ask_size`: 해당 호가 구간의 매도 주문 잔량 합계
* `bid_size`: 해당 호가 구간의 매수 주문 잔량 합계

페어별로 지원하는 모아보기 단위가 다릅니다. 지원 단위는 <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/faq-market-policy">마켓별 주문 정책</Anchor> 또는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-orderbook-instruments">호가 정책 조회</Anchor> API의 `supported_levels`에서 확인할 수 있습니다.

기본값은 `0`이며, 지원하지 않는 단위를 지정하면 데이터가 수신되지 않을 수 있습니다.

## 호가 조회 개수 지정

수신할 호가 쌍의 개수는 `codes`의 페어 코드 뒤에 마침표(`.`)와 `unit`을 붙여 지정합니다.

> `{pair_code}.{unit}`<br />예시: `KRW-BTC.15`

지원하는 `unit` 값은 `1`, `5`, `15`, `30`이며, 별도로 지정하지 않으면 30개의 호가 쌍이 반환됩니다.

<br />

## Request 메시지 형식

호가 데이터 수신을 요청하기 위해서는 WebSocket 연결 이후 아래 구조의 JSON Object를 생성한 뒤 요청 메시지의 Data Type Object로 포함하여 전송해야 합니다. Ticket, Format 필드를 포함한 전체 WebSocket 데이터 요청 메시지 명세는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide">WebSocket 사용 안내</Anchor> 문서를 참고하세요.

| 필드명                | 타입          | 내용                                                                                                                                                        | 필수 여부    | 기본 값    |
| ------------------ | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| type               | String      | 수신할 데이터 타입.<br />호가 데이터를 요청하는 경우 `orderbook`로 지정합니다.                                                                                                      | Required |         |
| codes              | List:String | 수신할 페어 코드 목록.<br />페어 코드는 대문자로 입력해야 합니다.                                                                                                                  | Required |         |
| is\_only\_snapshot | Boolean     | `true`로 설정하면 요청 시점의 현재가 스냅샷 데이터만 1회 수신합니다.                                                                                                                | Optional | `false` |
| is\_only\_realtime | Boolean     | `true`로 설정하면 현재가 스냅샷 없이 실시간 스트림 데이터만 수신합니다.                                                                                                               | Optional | `false` |
| format             | String      | 수신하고자 하는 데이터 포맷입니다. <br />`DEFAULT` : 기본 포맷.<br />`SIMPLE` : 간략한 포맷. 각 필드가 축약어 형태로 반환됩니다.<br />`JSON_LIST` : 리스트 포맷.<br />`SIMPLE_LIST` : 축약어 형태의 리스트 포맷. | Required |         |

<br />

## 구독 데이터 명세

호가 데이터 스냅샷 또는 실시간 스트림 데이터는 아래와 같이 반환됩니다.

| **필드명**                     | **축약형** | **내용**                                                                                            | **타입**          | **값**                                            |
| --------------------------- | ------- | ------------------------------------------------------------------------------------------------- | --------------- | ------------------------------------------------ |
| type                        | ty      | 타입                                                                                                | String          | `orderbook`                                      |
| code                        | cd      | 페어 코드<br />(예시: `KRW-BTC`)                                                                        | String          |                                                  |
| total\_ask\_size            | tas     | 호가 매도 총 잔량                                                                                        | Double          |                                                  |
| total\_bid\_size            | tbs     | 호가 매수 총 잔량                                                                                        | Double          |                                                  |
| orderbook\_units            | obu     | 호가                                                                                                | List of Objects |                                                  |
| orderbook\_units.ask\_price | obu.ap  | 매도 호가                                                                                             | Double          |                                                  |
| orderbook\_units.bid\_price | obu.bp  | 매수 호가                                                                                             | Double          |                                                  |
| orderbook\_units.ask\_size  | obu.as  | 매도 잔량                                                                                             | Double          |                                                  |
| orderbook\_units.bid\_size  | obu.bs  | 매수 잔량                                                                                             | Double          |                                                  |
| timestamp                   | tms     | 타임스탬프 (ms)                                                                                        | Long            |                                                  |
| level                       | lv      | 호가 모아보기 단위 (default: 0, 기본 호가단위)<br />\*호가 모아보기 기능은 원화마켓(KRW)에서만 지원하므로 BTC, USDT 마켓의 경우 0만 존재합니다. | Double          | 모아보기 단위                                          |
| stream\_type                | st      | 스트림 타입                                                                                            | String          | `SNAPSHOT`<br />: 스냅샷<br />`REALTIME`<br />: 실시간 |

<br />

## 예시

#### 여러 페어의 호가 깊이 모니터링

> KRW-BTC는 10,000 KRW 단위로 15개 호가 쌍을, KRW-ETH는 1,000 KRW 단위로 5개 호가 쌍을 구독하여, 각 페어의 가격대별 매수·매도 물량과 시장 깊이를 실시간으로 모니터링하는 예시입니다. ticket은 요청을 식별하기 위한 값으로 원하는 문자열을 직접 지정할 수 있습니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "orderbook-monitor"
  },
  {
    "type": "orderbook",
    "codes": ["KRW-BTC.15"],
    "level": 10000
  },
  {
    "type": "orderbook",
    "codes": ["KRW-ETH.5"], 
    "level": 1000
  },
  {
    "format": "SIMPLE_LIST"
  }
]
```

##### 수신 메시지 예제

```json
{
  "type": "orderbook",
  "code": "KRW-BTC",
  "timestamp": 1787727947526,
  "total_ask_size": 18.93431797,
  "total_bid_size": 11.48721726,
  "orderbook_units": [
    {
      "ask_price": 109950000.0,
      "bid_price": 109880000.0,
      "ask_size": 0.00342134,
      "bid_size": 0.18819864
    },
    {
      "ask_price": 109960000.0,
      "bid_price": 109840000.0,
      "ask_size": 0.00077053,
      "bid_size": 0.08051758
    },
    {
      "ask_price": 109970000.0,
      "bid_price": 109830000.0,
      "ask_size": 0.09644533,
      "bid_size": 0.03782836
    },
    {
      "ask_price": 109980000.0,
      "bid_price": 109820000.0,
      "ask_size": 0.17772082,
      "bid_size": 0.26368128
    },
    {
      "ask_price": 109990000.0,
      "bid_price": 109810000.0,
      "ask_size": 0.38395685,
      "bid_size": 0.41145608
    },
    {
      "ask_price": 110000000.0,
      "bid_price": 109800000.0,
      "ask_size": 1.60197798,
      "bid_size": 0.36972501
    },
    {
      "ask_price": 110010000.0,
      "bid_price": 109790000.0,
      "ask_size": 7.83138819,
      "bid_size": 0.52858073
    },
    {
      "ask_price": 110020000.0,
      "bid_price": 109780000.0,
      "ask_size": 0.06749248,
      "bid_size": 0.37087058
    },
    {
      "ask_price": 110030000.0,
      "bid_price": 109770000.0,
      "ask_size": 0.68572776,
      "bid_size": 3.90963
    },
    {
      "ask_price": 110040000.0,
      "bid_price": 109760000.0,
      "ask_size": 0.42258073,
      "bid_size": 0.1357252
    },
    {
      "ask_price": 110050000.0,
      "bid_price": 109750000.0,
      "ask_size": 2.07411804,
      "bid_size": 0.18884944
    },
    {
      "ask_price": 110060000.0,
      "bid_price": 109740000.0,
      "ask_size": 0.66264869,
      "bid_size": 0.07768186
    },
    {
      "ask_price": 110070000.0,
      "bid_price": 109730000.0,
      "ask_size": 0.29757603,
      "bid_size": 0.07027257
    },
    {
      "ask_price": 110080000.0,
      "bid_price": 109720000.0,
      "ask_size": 1.65671197,
      "bid_size": 0.05635197
    },
    {
      "ask_price": 110090000.0,
      "bid_price": 109710000.0,
      "ask_size": 0.28271915,
      "bid_size": 0.01891968
    },
    {
      "ask_price": 110100000.0,
      "bid_price": 109700000.0,
      "ask_size": 0.43170875,
      "bid_size": 0.0816967
    },
    {
      "ask_price": 110110000.0,
      "bid_price": 109690000.0,
      "ask_size": 0.11668995,
      "bid_size": 0.01862403
    },
    {
      "ask_price": 110120000.0,
      "bid_price": 109680000.0,
      "ask_size": 0.69267557,
      "bid_size": 0.33017599
    },
    {
      "ask_price": 110130000.0,
      "bid_price": 109670000.0,
      "ask_size": 0.13633769,
      "bid_size": 2.05484164
    },
    {
      "ask_price": 110140000.0,
      "bid_price": 109660000.0,
      "ask_size": 0.01092664,
      "bid_size": 0.67617863
    },
    {
      "ask_price": 110150000.0,
      "bid_price": 109650000.0,
      "ask_size": 0.35735246,
      "bid_size": 0.14343052
    },
    {
      "ask_price": 110160000.0,
      "bid_price": 109640000.0,
      "ask_size": 0.02359097,
      "bid_size": 0.0906867
    },
    {
      "ask_price": 110170000.0,
      "bid_price": 109630000.0,
      "ask_size": 0.11649106,
      "bid_size": 0.41100224
    },
    {
      "ask_price": 110180000.0,
      "bid_price": 109620000.0,
      "ask_size": 0.0172558,
      "bid_size": 0.08229236
    },
    {
      "ask_price": 110190000.0,
      "bid_price": 109610000.0,
      "ask_size": 0.03970812,
      "bid_size": 0.06079286
    },
    {
      "ask_price": 110200000.0,
      "bid_price": 109600000.0,
      "ask_size": 0.36126766,
      "bid_size": 0.23611511
    },
    {
      "ask_price": 110210000.0,
      "bid_price": 109590000.0,
      "ask_size": 0.00529791,
      "bid_size": 0.03976073
    },
    {
      "ask_price": 110220000.0,
      "bid_price": 109580000.0,
      "ask_size": 0.37377172,
      "bid_size": 0.08267387
    },
    {
      "ask_price": 110230000.0,
      "bid_price": 109570000.0,
      "ask_size": 0.00540378,
      "bid_size": 0.14031581
    },
    {
      "ask_price": 110240000.0,
      "bid_price": 109560000.0,
      "ask_size": 0.000584,
      "bid_size": 0.33034109
    }
  ],
  "stream_type": "SNAPSHOT",
  "level": 10000
}
{
  "type": "orderbook",
  "code": "KRW-ETH",
  "timestamp": 1787727947404,
  "total_ask_size": 1075.24704782,
  "total_bid_size": 612.45477499,
  "orderbook_units": [
    {
      "ask_price": 3426000.0,
      "bid_price": 3425000.0,
      "ask_size": 80.27690046,
      "bid_size": 36.87952898
    },
    {
      "ask_price": 3427000.0,
      "bid_price": 3424000.0,
      "ask_size": 1.82355307,
      "bid_size": 16.55651371
    },
    {
      "ask_price": 3428000.0,
      "bid_price": 3423000.0,
      "ask_size": 1.35540193,
      "bid_size": 166.23455141
    },
    {
      "ask_price": 3429000.0,
      "bid_price": 3422000.0,
      "ask_size": 4.60910209,
      "bid_size": 14.41822979
    },
    {
      "ask_price": 3430000.0,
      "bid_price": 3421000.0,
      "ask_size": 46.76962749,
      "bid_size": 14.01087489
    },
    {
      "ask_price": 3431000.0,
      "bid_price": 3420000.0,
      "ask_size": 19.01673999,
      "bid_size": 33.9704142
    },
    {
      "ask_price": 3432000.0,
      "bid_price": 3419000.0,
      "ask_size": 223.13566619,
      "bid_size": 10.60064954
    },
    {
      "ask_price": 3433000.0,
      "bid_price": 3418000.0,
      "ask_size": 42.7925073,
      "bid_size": 42.13494345
    },
    {
      "ask_price": 3434000.0,
      "bid_price": 3417000.0,
      "ask_size": 43.9387063,
      "bid_size": 12.37324427
    },
    {
      "ask_price": 3435000.0,
      "bid_price": 3416000.0,
      "ask_size": 66.05576667,
      "bid_size": 2.02232306
    },
    {
      "ask_price": 3436000.0,
      "bid_price": 3415000.0,
      "ask_size": 40.2432618,
      "bid_size": 7.83508018
    },
    {
      "ask_price": 3437000.0,
      "bid_price": 3414000.0,
      "ask_size": 21.02402734,
      "bid_size": 4.28965157
    },
    {
      "ask_price": 3438000.0,
      "bid_price": 3413000.0,
      "ask_size": 15.93951649,
      "bid_size": 8.52924235
    },
    {
      "ask_price": 3439000.0,
      "bid_price": 3412000.0,
      "ask_size": 64.99778051,
      "bid_size": 3.14905254
    },
    {
      "ask_price": 3440000.0,
      "bid_price": 3411000.0,
      "ask_size": 98.14247184,
      "bid_size": 11.14896713
    },
    {
      "ask_price": 3441000.0,
      "bid_price": 3410000.0,
      "ask_size": 6.64691027,
      "bid_size": 20.665478
    },
    {
      "ask_price": 3442000.0,
      "bid_price": 3409000.0,
      "ask_size": 33.02223357,
      "bid_size": 4.43895484
    },
    {
      "ask_price": 3443000.0,
      "bid_price": 3408000.0,
      "ask_size": 12.21645763,
      "bid_size": 32.46564449
    },
    {
      "ask_price": 3444000.0,
      "bid_price": 3407000.0,
      "ask_size": 12.51008476,
      "bid_size": 3.26811127
    },
    {
      "ask_price": 3445000.0,
      "bid_price": 3406000.0,
      "ask_size": 12.69428031,
      "bid_size": 4.29047757
    },
    {
      "ask_price": 3446000.0,
      "bid_price": 3405000.0,
      "ask_size": 4.17409886,
      "bid_size": 4.93794905
    },
    {
      "ask_price": 3447000.0,
      "bid_price": 3404000.0,
      "ask_size": 18.23709183,
      "bid_size": 0.67174836
    },
    {
      "ask_price": 3448000.0,
      "bid_price": 3403000.0,
      "ask_size": 32.2535389,
      "bid_size": 4.67272598
    },
    {
      "ask_price": 3449000.0,
      "bid_price": 3402000.0,
      "ask_size": 26.27420171,
      "bid_size": 6.63830662
    },
    {
      "ask_price": 3450000.0,
      "bid_price": 3401000.0,
      "ask_size": 70.3652008,
      "bid_size": 9.84892265
    },
    {
      "ask_price": 3451000.0,
      "bid_price": 3400000.0,
      "ask_size": 9.50997994,
      "bid_size": 107.77116639
    },
    {
      "ask_price": 3452000.0,
      "bid_price": 3399000.0,
      "ask_size": 42.77626242,
      "bid_size": 13.92396746
    },
    {
      "ask_price": 3453000.0,
      "bid_price": 3398000.0,
      "ask_size": 9.56671911,
      "bid_size": 9.63773543
    },
    {
      "ask_price": 3454000.0,
      "bid_price": 3397000.0,
      "ask_size": 4.97285159,
      "bid_size": 2.534455
    },
    {
      "ask_price": 3455000.0,
      "bid_price": 3396000.0,
      "ask_size": 9.90610665,
      "bid_size": 2.53586481
    }
  ],
  "stream_type": "SNAPSHOT",
  "level": 0
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
* [캔들 (Candle)](https://docs.upbit.com/kr/reference/websocket-candle.md)
* [공지사항(Announcement)](https://docs.upbit.com/kr/reference/websocket-announcement.md)
* [내 주문 및 체결 (MyOrder)](https://docs.upbit.com/kr/reference/websocket-myorder.md)
* [내 자산 (MyAsset)](https://docs.upbit.com/kr/reference/websocket-myasset.md)
* [구독 중인 스트림 목록 조회(Subscriptions)](https://docs.upbit.com/kr/reference/list-subscriptions.md)