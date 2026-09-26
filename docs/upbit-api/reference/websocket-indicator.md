---
updatedAt: 2026-09-23T07:13:07.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 지표 값(Indicator)

원하는 기술적 지표를 구독하면 계산된 지표 값을 받아볼 수 있습니다.

## WebSocket Endpoint

| 구분      | Endpoint                                |
| ------- | --------------------------------------- |
| Private | `wss://api.upbit.com/websocket/v1/info` |

<Callout icon="fad fa-circle-info" theme="info">
  ### **기술적 지표는 Private WebSocket입니다.**

  연결 시 발급받은 API Key로 생성한 JWT를 Authorization 헤더에 포함해야 합니다. 다만 본 기능을 이용하기 위해 API Key에 별도의 권한(scope)을 설정할 필요는 없습니다. Private WebSocket에 대한 내용은 <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/faq-api">FAQ 문서</Anchor>를 확인해 주세요.
</Callout>

<Callout icon="fad fa-triangle-exclamation" theme="error">
  ### **Private WebSocket 연결 관리 안내**

  **Private WebSocket은 동시에 많은 연결을 유지하는 경우 신규 연결이 거절될 수 있습니다.&#x20;**<br />불필요한 연결 생성을 최소화하고, 가능한 경우 하나의 연결에서 필요한 데이터 타입을 함께 구독해 주세요.
</Callout>

<Callout icon="fad fa-gauge" theme="warn">
  ### **데이터 수신 시점 및 연결 유지 안내**

  * 구독 시점의 최근 지표 값을 제공하는 스냅샷은 지원하지 않습니다. `indicator`는 구독 이후 선택한 시간 단위(`timeframe`)의 캔들이 처음 마감된 후부터 지표 값을 전송합니다.
  * 기술적 지표 값은 선택한 캔들이 마감될 때 전송되므로 데이터 수신 간격이 길어질 수 있습니다. 일정 시간 동안 데이터 송수신이 없는 경우 WebSocket 연결이 종료될 수 있으므로, 클라이언트에서 주기적으로 PING 프레임을 전송하는 등 연결 유지 처리를 적용해 주세요. 자세한 내용은 [WebSocket 공통 가이드의 연결 관리 항목](https://docs.upbit.com/kr/reference/websocket-guide)을 참고해 주세요.
</Callout>

## Request 메시지 형식

지표 데이터를 수신하려면 WebSocket 연결 및 인증 후, 아래 필드로 구성한 Data Type Object를 요청 메시지에 포함하여 전송해야 합니다.

<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide#:~:text=%EB%93%A4%EC%9D%84%20%ED%8F%AC%ED%95%A8%ED%95%B4%EC%95%BC%20%ED%95%A9%EB%8B%88%EB%8B%A4.-,Ticket%20Object,-%EB%B0%B0%EC%97%B4%EC%9D%98%20%EC%B2%AB%EB%B2%88%EC%A7%B8%20%EC%9A%94%EC%86%8C%EB%A1%9C">Ticket</Anchor>, <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide#:~:text=%EC%82%AC%EC%9A%A9%ED%95%A0%20%EC%88%98%20%EC%9E%88%EC%8A%B5%EB%8B%88%EB%8B%A4.-,Format%20Object,-%EB%B0%B0%EC%97%B4%EC%9D%98%20%EB%A7%88%EC%A7%80%EB%A7%89%20%EC%9A%94%EC%86%8C%EB%A1%9C">Format 필드</Anchor>를 포함한 전체 WebSocket 데이터 요청 메시지 명세는 <Anchor target="_blank" href="ref:websocket-guide">WebSocket 사용 안내</Anchor> 문서를 참고해 주세요.

| 필드명                  | 타입           | 내용                                                                                       | 필수 여부    | 기본 값 |
| -------------------- | ------------ | ---------------------------------------------------------------------------------------- | -------- | ---- |
| type                 | String       | 수신할 데이터 타입.<br />`indicator` : 지표 값. 구독 이후 캔들이 마감될 때마다 전송                                | Required | -    |
| codes                | List: String | 수신할 마켓 코드 목록.<br />KRW 마켓을 지원합니다.<br />예: `["KRW-BTC", "KRW-ETH"]`                       | Required | -    |
| timeframe            | String       | 지표 계산에 사용할 캔들 단위.<br />`5m`, `15m`, `60m`, `240m`, `1d`                                  | Required | -    |
| indicators           | List         | 구독할 지표 목록.                                                                               | Required | -    |
| indicators\[].type   | String       | 지표 타입.<br />`volume`, `rsi`, `moving_average`, `ma_crossover`, `bollinger_bands`, `macd` | Required | -    |
| indicators\[].params | Object       | 지표 계산에 사용할 파라미터.<br />지표에 따라 지원하는 필드와 값이 다릅니다.                                           | Required | -    |

<br />

### 지표별 지원 범위(`indicators`)

아래에서 **원하는 지표를 선택해 계산 파라미터와 반환 값을 확인해 주세요.**

<Accordion title="거래량 (volume)" icon="fad fa-chart-column">
  | 항목                          | 지원 값                                                                                   |
  | --------------------------- | -------------------------------------------------------------------------------------- |
  | 지표명                         | 거래량                                                                                    |
  | 지표 타입 (`indicators[].type`) | `volume`                                                                               |
  | 계산 파라미터 (`params`)          | 비교 기간(`period`): 5, 10, 15, 20, 25, 30, 35, 40, 45, 50                                 |
  | 지표 값 (`value`)              | `volume` : 기준 캔들의 거래량<br />`avg_volume` :비교 기간의 평균 거래량<br />`ratio` : 평균 거래량 대비 거래량 비율 |

  > 반환되는 지표 값의 설명은 문서 아래 지표 값(value) 항목을 참고해 주세요.

  **구독 예시**

  ```text 지표 값 구독
  [
    {
      "ticket": "indicator-value-subscription" // 요청을 식별하기 위한 값
    },
    {
      "type": "indicator", // 지표 값 구독
      "codes": ["KRW-BTC"], // 구독할 페어 코드
      "timeframe": "5m", // 지표 계산에 사용할 캔들 단위
      "indicators": [
        {
          "type": "volume", // 거래량 지표
          "params": {
            "period": 20 // 평균 거래량 계산 기간
          }
        }
      ]
    }
  ]
  ```
</Accordion>

<Accordion title="RSI (rsi)" icon="fad fa-chart-line">
  | 항목                          | 지원 값                                                             |
  | --------------------------- | ---------------------------------------------------------------- |
  | 지표명                         | RSI                                                              |
  | 지표 타입 (`indicators[].type`) | `rsi`                                                            |
  | 계산 파라미터 (`params`)          | 기간(`period`): 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 21, 25 |
  | 지표 값 (`value`)              | `rsi` : RSI 값                                                    |

  > 반환되는 지표 값은 문서 아래 지표 값(value) 항목을 참고해 주세요.

  **구독 예시**

  ```text 지표 값 구독
  [
    {
      "ticket": "rsi-value-subscription" // 요청을 식별하기 위한 값
    },
    {
      "type": "indicator", // 지표 값 구독
      "codes": ["KRW-BTC"], // 구독할 페어 코드
      "timeframe": "5m", // 지표 계산에 사용할 캔들 단위
      "indicators": [
        {
          "type": "rsi", // RSI 지표
          "params": {
            "period": 14 // RSI 계산 기간
          }
        }
      ]
    }
  ]
  ```
</Accordion>

<Accordion title="이동평균선(단일) (moving_average)" icon="fad fa-chart-line-up">
  <Table>
    <thead>
      <tr>
        <th>
          항목
        </th>

        <th>
          지원 값
        </th>
      </tr>
    </thead>

    <tbody>
      <tr>
        <td>
          지표명
        </td>

        <td>
          이동평균선(단일)
        </td>
      </tr>

      <tr>
        <td>
          지표 타입 (`indicators[].type`)
        </td>

        <td>
          `moving_average`
        </td>
      </tr>

      <tr>
        <td>
          계산 파라미터 (`params`)
        </td>

        <td>
          - 종류(`type`) : `SMA`, `EMA`
          - 기간(`period`): 3, 5, 7, 10, 15, 20, 25, 30, 50, 60, 90, 99, 100, 120, 200
        </td>
      </tr>

      <tr>
        <td>
          지표 값 (`value`)
        </td>

        <td>
          `value` : 이동평균선 값
        </td>
      </tr>
    </tbody>
  </Table>

  > 반환되는 지표 값은 문서 아래 지표 값(value) 항목을 참고해 주세요.

  **구독 예시**

  ```text 지표 값 구독
  [
    {
      "ticket": "moving-average-value-subscription" // 요청을 식별하기 위한 값
    },
    {
      "type": "indicator", // 지표 값 구독
      "codes": ["KRW-BTC"], // 구독할 페어 코드
      "timeframe": "5m", // 지표 계산에 사용할 캔들 단위
      "indicators": [
        {
          "type": "moving_average", // 이동평균선 지표
          "params": {
            "type": "SMA", // 이동평균선 종류
            "period": 20 // 이동평균선 계산 기간
          }
        }
      ]
    }
  ]
  ```
</Accordion>

<Accordion title="이동평균선(교차) (ma_crossover)" icon="fad fa-arrows-cross">
  <Table>
    <thead>
      <tr>
        <th>
          항목
        </th>

        <th>
          지원 값
        </th>
      </tr>
    </thead>

    <tbody>
      <tr>
        <td>
          지표명
        </td>

        <td>
          이동평균선(교차)
        </td>
      </tr>

      <tr>
        <td>
          지표 타입 (`indicators[].type`)
        </td>

        <td>
          `ma_crossover`
        </td>
      </tr>

      <tr>
        <td>
          계산 파라미터 (`params`)
        </td>

        <td>
          종류(`type`) : `SMA`, `EMA`<br />교차 조합(`short_period`/`long_period`) : <br />5/20, 20/60, 30/90, 50/200

          > 위 4개 교차 조합만 지원합니다.
          > 그 외 조합은 moving_average를 두 개 구독하여 직접 비교해 주세요.
        </td>
      </tr>

      <tr>
        <td>
          지표 값 (`value`)
        </td>

        <td>
          `ma_short` : 단기 이동평균선 값<br />`ma_long` : 장기 이동평균선 값
        </td>
      </tr>
    </tbody>
  </Table>

  > 반환되는 지표 값은 문서 아래 지표 값(value) 항목을 참고해 주세요.

  **구독 예시**

  ```text 지표 값 구독
  [
    {
      "ticket": "ma-crossover-value-subscription" // 요청을 식별하기 위한 값
    },
    {
      "type": "indicator", // 지표 값 구독
      "codes": ["KRW-BTC"], // 구독할 페어 코드
      "timeframe": "5m", // 지표 계산에 사용할 캔들 단위
      "indicators": [
        {
          "type": "ma_crossover", // 이동평균선 교차 지표
          "params": {
            "type": "SMA", // 이동평균선 종류
            "short_period": 5, // 단기 이동평균선 계산 기간
            "long_period": 20 // 장기 이동평균선 계산 기간
          }
        }
      ]
    }
  ]
  ```
</Accordion>

<Accordion title="볼린저밴드 (bollinger_bands)" icon="fad fa-chart-line">
  <Table>
    <thead>
      <tr>
        <th>
          항목
        </th>

        <th>
          지원 값
        </th>
      </tr>
    </thead>

    <tbody>
      <tr>
        <td>
          지표명
        </td>

        <td>
          볼린저밴드
        </td>
      </tr>

      <tr>
        <td>
          지표 타입 (`indicators[].type`)
        </td>

        <td>
          `bollinger_bands`
        </td>
      </tr>

      <tr>
        <td>
          계산 파라미터 (`params`)
        </td>

        <td>
          기간(`period`): 10, 20, 50

          배수(`multiplier`): 1.0, 1.5, 2.0, 2.5, 3.0
        </td>
      </tr>

      <tr>
        <td>
          지표 값 (`value`)
        </td>

        <td>
          `bb_upper` : 상단 밴드 값<br />`bb_middle` : 중심선 값<br />`bb_lower` : 하단 밴드 값
        </td>
      </tr>
    </tbody>
  </Table>

  > 반환되는 지표 값은 문서 아래 지표 값(value) 항목을 참고해 주세요.

  **구독 예시**

  ```text 지표 값 구독
  [
    {
      "ticket": "bollinger-bands-value-subscription" // 요청을 식별하기 위한 값
    },
    {
      "type": "indicator", // 지표 값 구독
      "codes": ["KRW-BTC"], // 구독할 페어 코드
      "timeframe": "5m", // 지표 계산에 사용할 캔들 단위
      "indicators": [
        {
          "type": "bollinger_bands", // 볼린저밴드 지표
          "params": {
            "period": 20, // 볼린저밴드 계산 기간
            "multiplier": 2.0 // 표준편차 배수
          }
        }
      ]
    }
  ]
  ```
</Accordion>

<Accordion title="MACD (macd)" icon="fad fa-chart-mixed">
  | 항목                          | 지원 값                                                                            |
  | --------------------------- | ------------------------------------------------------------------------------- |
  | 지표명                         | MACD                                                                            |
  | 지표 타입 (`indicators[].type`) | `macd`                                                                          |
  | 계산 파라미터 (`params`)          | `fast` / `slow` / `signal_period`: 12/26/9, 8/17/9, 5/35/5                      |
  | 지표 값 (`value`)              | `macd` : MACD 값<br />`macd_signal` : MACD 시그널 값<br />`histogram` : MACD 히스토그램 값 |

  > 반환되는 지표 값은 문서 아래 지표 값(value) 항목을 참고해 주세요.

  **구독 예시**

  ```text 지표 값 구독
  [
    {
      "ticket": "macd-value-subscription" // 요청을 식별하기 위한 값
    },
    {
      "type": "indicator", // 지표 값 구독
      "codes": ["KRW-BTC"], // 구독할 페어 코드
      "timeframe": "5m", // 지표 계산에 사용할 캔들 단위
      "indicators": [
        {
          "type": "macd", // MACD 지표
          "params": {
            "fast": 12, // 단기 기간
            "slow": 26, // 장기 기간
            "signal_period": 9 // 시그널선 기간
          }
        }
      ]
    }
  ]
  ```
</Accordion>

#### 파라미터 입력 규칙

* 기간·개수·기준선에 해당하는 파라미터는 정수로 입력합니다. 볼린저밴드의 표준편차 배수는 지원 값 중에서 입력하며, 소수점 첫째 자리까지 지원합니다.
* 동일한 값(`2`, `2.0`, `2.00`)을 나타내는 입력은 허용하며 서버에서 `2.0`으로 정규화합니다. 지원하지 않는 값(예: `2.04`)은 반올림하지 않고 오류를 반환하며, 숫자 파라미터는 문자열이 아닌 JSON number 타입으로 입력해야 합니다.

<Callout icon="fad fa-circle-info" theme="info">
  ### **구독 관리**

  * 한 요청 메시지에서 `indicator` 구독 조합은 최대 100개까지 요청할 수 있습니다. 같은 요청 메시지에 `indicator_signal`을 함께 포함하는 경우 두 데이터 타입의 구독 조합 수를 합산합니다. 전체 구독 조합 수는 **각 Data Type Object의 (**`codes`**&#x20;개수 ×&#x20;**`indicators`**&#x20;개수)의 합**으로 계산하며, 100개를 초과하면 일부만 등록하지 않고 요청 전체를 `TOO_MANY_REQUEST` 오류로 거부합니다.
  * 같은 WebSocket 연결에서 새로운 구독 요청을 전송하면 기존 구독 정보가 초기화되고, 새로운 요청 내용으로 전체 대체됩니다. 기존 구독을 유지하면서 항목을 추가하려면 유지할 항목과 추가할 항목을 모두 포함해 요청해 주세요. 동일한 조합은 같은 연결에서 하나만 구독할 수 있습니다.
  * 구독 요청 전 지원 대상 마켓과 입력한 마켓 코드를 확인해 주세요. 지원하지 않거나 존재하지 않는 마켓 코드를 요청하더라도 별도로 오류가 반환되지 않으며, 해당 마켓의 데이터는 전송되지 않습니다. 만약, 해당 마켓이 차후 지원 대상에 포함되고 WebSocket 연결과 구독이 유지되고 있는 경우에는, 지원 시점 이후 생성되는 데이터부터 수신할 수 있습니다.
</Callout>

<br />

## 구독 데이터 명세

선택한 캔들이 마감되면 아래 형식의 실시간 스트림 데이터가 반환됩니다.

### 공통 필드

| 필드명                    | 축약형      | 내용                             | 타입     | 값                                                                            |
| ---------------------- | -------- | ------------------------------ | ------ | ---------------------------------------------------------------------------- |
| `type`                 | `ty`     | 구독 데이터 타입                      | String | `indicator`                                                                  |
| `code`                 | `cd`     | 마켓 코드                          | String | 예: `KRW-BTC`                                                                 |
| `timeframe`            | `tf`     | 지표 계산에 사용된 캔들 단위               | String | `5m`, `15m`, `60m`, `240m`, `1d`                                             |
| `indicator`            | `ind`    | 지표 타입                          | String | `volume`, `rsi`, `moving_average`, `ma_crossover`, `bollinger_bands`, `macd` |
| `params`               | `p`      | 지표 계산에 사용된 파라미터                | Object | 지표별 구성 상이                                                                    |
| `value`                | `v`      | 계산된 지표 값                       | Object | 아래 지표별 필드 참고                                                                 |
| `candle_date_time_utc` | `cdttmu` | 지표 계산의 기준이 된 캔들 구간의 시작 시각(UTC) | String | 예: `2026-09-03T04:50:00`                                                     |
| `timestamp`            | `tms`    | 응답 생성 시각                       | Long   | Unix timestamp(ms)                                                           |
| `stream_type`          | `st`     | 스트림 타입                         | String | `REALTIME`                                                                   |

<br />

### 계산 파라미터(`params`)

아래는 `params` 객체 내부 필드입니다. 해당 지표에서 사용하는 필드만 포함됩니다.

| 필드명             | 축약형   | 내용             | 타입      | 적용 지표                      |
| --------------- | ----- | -------------- | ------- | -------------------------- |
| `period`        | `p`   | 지표 계산 또는 비교 기간 | Integer | 거래량, RSI, 이동평균선(단일), 볼린저밴드 |
| `type`          | `t`   | 이동평균선 종류       | String  | 이동평균선(단일·교차): `SMA`, `EMA` |
| `short_period`  | `shp` | 단기 이동평균선 기간    | Integer | 이동평균선(교차)                  |
| `long_period`   | `lop` | 장기 이동평균선 기간    | Integer | 이동평균선(교차)                  |
| `multiplier`    | `m`   | 표준편차 배수        | Double  | 볼린저밴드                      |
| `fast`          | `f`   | 단기 계산 기간       | Integer | MACD                       |
| `slow`          | `s`   | 장기 계산 기간       | Integer | MACD                       |
| `signal_period` | `sgp` | 시그널 계산 기간      | Integer | MACD                       |

<br />

### 지표 값(`value`)

아래 필드는 `value` 객체 내부에 포함됩니다.

| 지표        | 필드명           | 축약형     | 내용               | 타입     |
| --------- | ------------- | ------- | ---------------- | ------ |
| 거래량       | `volume`      | `vol`   | 기준 캔들의 거래량       | Double |
| 거래량       | `avg_volume`  | `av`    | 비교 기간의 평균 거래량    | Double |
| 거래량       | `ratio`       | `rto`   | 평균 거래량 대비 거래량 비율 | Double |
| RSI       | `rsi`         | `rsi`   | RSI 값            | Double |
| 이동평균선(단일) | `value`       | `value` | 이동평균선 값          | Double |
| 이동평균선(교차) | `ma_short`    | `mas`   | 단기 이동평균선 값       | Double |
| 이동평균선(교차) | `ma_long`     | `mal`   | 장기 이동평균선 값       | Double |
| 볼린저밴드     | `bb_upper`    | `bbu`   | 상단 밴드 값          | Double |
| 볼린저밴드     | `bb_middle`   | `bbm`   | 중심선 값            | Double |
| 볼린저밴드     | `bb_lower`    | `bbl`   | 하단 밴드 값          | Double |
| MACD      | `macd`        | `macd`  | MACD 값           | Double |
| MACD      | `macd_signal` | `macds` | MACD 시그널 값       | Double |
| MACD      | `histogram`   | `hist`  | MACD 히스토그램 값     | Double |

`SIMPLE` 포맷은 최상위 필드뿐 아니라 `params`, `value` 내부 필드에도 위 축약형을 적용합니다. 데이터 타입은 바뀌지 않으며, 숫자 값은 JSON number로 제공됩니다.

예를 들어 최상위 `value`는 `v`로 축약되지만, 이동평균선 값에 해당하는 내부 필드 `value`는 그대로 유지됩니다. 볼린저밴드의 `params.multiplier`는 `m`으로 축약됩니다.

<br />

### 지표 계산 기준

<Accordion title="마감된 캔들 기준" icon="fad fa-circle-question">
  * 사용자가 선택한 시간 단위(`timeframe`)의 캔들이 마감된 후 지표를 계산합니다. 가격 비교에는 해당 캔들의 종가를 사용합니다.

  - 진행 중인 캔들의 가격 변화에 따라 지표 값을 갱신하지 않습니다.
</Accordion>

<br />

## 예시

아래의 수신 메시지는 각 지표 값이 전송되는 경우의 예시이며, 각각 독립된 JSON 메시지입니다.

#### 여러 지표 값 함께 구독

> KRW-BTC와 KRW-ETH의 5분봉 RSI, 이동평균선, 이동평균선 교차, MACD 값을 함께 구독하는 예시입니다. 전체 구독 조합 수는 8개입니다. `ticket`은 요청 식별값으로 직접 지정할 수 있습니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "indicator-monitor"
  },
  {
    "type": "indicator",
    "codes": ["KRW-BTC", "KRW-ETH"],
    "timeframe": "5m",
    "indicators": [
      {
        "type": "rsi",
        "params": {
          "period": 14
        }
      },
      {
        "type": "moving_average",
        "params": {
          "type": "SMA",
          "period": 20
        }
      },
      {
        "type": "ma_crossover",
        "params": {
          "type": "EMA",
          "short_period": 5,
          "long_period": 20
        }
      },
      {
        "type": "macd",
        "params": {
          "fast": 12,
          "slow": 26,
          "signal_period": 9
        }
      }
    ]
  },
  {
    "format": "DEFAULT"
  }
]
```

##### 수신 메시지 예제 - RSI

위 요청 중 KRW-BTC의 RSI 지표 값이 수신된 예시입니다.

```json
{
  "type": "indicator",
  "code": "KRW-BTC",
  "timeframe": "5m",
  "indicator": "rsi",
  "params": {
    "period": 14
  },
  "candle_date_time_utc": "2026-09-03T04:50:00",
  "value": {
    "rsi": 58.75323648
  },
  "timestamp": 1788411301095,
  "stream_type": "REALTIME"
}
```

##### 수신 메시지 예제 - 이동평균선

```json
{
  "type": "indicator",
  "code": "KRW-BTC",
  "timeframe": "5m",
  "indicator": "moving_average",
  "params": {
    "type": "SMA",
    "period": 20
  },
  "candle_date_time_utc": "2026-09-03T04:50:00",
  "value": {
    "value": 106791950.0
  },
  "timestamp": 1788411301095,
  "stream_type": "REALTIME"
}
```

##### 수신 메시지 예제 - MACD

```json
{
  "type": "indicator",
  "code": "KRW-BTC",
  "timeframe": "5m",
  "indicator": "macd",
  "params": {
    "fast": 12,
    "slow": 26,
    "signal_period": 9
  },
  "candle_date_time_utc": "2026-09-03T04:50:00",
  "value": {
    "macd": 33554.85769535,
    "macd_signal": 55795.38214693,
    "histogram": -22240.52445158
  },
  "timestamp": 1788411301096,
  "stream_type": "REALTIME"
}
```

#### SIMPLE 포맷으로 지표 값 구독

> KRW-BTC의 MACD 값을 축약된 필드명으로 수신하는 예시입니다. 요청 필드명은 기본 형식으로 작성하고, Format Object의 `format`을 `SIMPLE`로 지정합니다.

##### 구독 요청 예제

```json
[
  {
    "ticket": "macd-simple-monitor"
  },
  {
    "type": "indicator",
    "codes": ["KRW-BTC"],
    "timeframe": "5m",
    "indicators": [
      {
        "type": "macd",
        "params": {
          "fast": 12,
          "slow": 26,
          "signal_period": 9
        }
      }
    ]
  },
  {
    "format": "SIMPLE"
  }
]
```

##### 수신 메시지 예제

```json
{
  "ty": "indicator",
  "cd": "KRW-BTC",
  "tf": "5m",
  "ind": "macd",
  "p": {
    "f": 12,
    "s": 26,
    "sgp": 9
  },
  "cdttmu": "2026-09-03T04:50:00",
  "v": {
    "macd": 33554.85769535,
    "macds": 55795.38214693,
    "hist": -22240.52445158
  },
  "tms": 1788411301096,
  "st": "REALTIME"
}
```

#### 구독 목록 조회

신규 구독 요청에 대한 별도의 확인 응답은 제공되지 않습니다. 현재 연결에 등록된 구독 정보는 `LIST_SUBSCRIPTIONS`로 확인할 수 있습니다. 구독 목록 조회의 공통 명세는 [구독 중인 스트림 목록 조회(Subscriptions)](https://docs.upbit.com/kr/reference/list-subscriptions) 문서를 참고해 주세요.

<br />

## 에러 안내

요청 검증에 실패하면 아래 형식으로 오류 응답을 반환합니다.

```json
{
  "error": {
    "name": "ERROR_CODE",
    "message": "ERROR_MESSAGE"
  }
}
```

반환될 수 있는 주요 에러 코드는 아래와 같습니다.

| `error.name`       | 발생 이유                                    | 권장 조치                                                           |
| ------------------ | ---------------------------------------- | --------------------------------------------------------------- |
| INVALID\_AUTH      | Private WebSocket 인증 정보가 없거나 인증에 실패한 경우  | Private Endpoint와 `Authorization` 헤더의 인증 토큰을 확인해 주세요.           |
| WRONG\_FORMAT      | JSON 구조, 필드 타입 또는 파라미터 입력 형식이 올바르지 않은 경우 | 요청 배열과 Object 구조를 확인해 주세요. 숫자 필드에 문자열이나 배열을 입력하지 않았는지도 확인해 주세요. |
| NO\_TICKET         | `ticket`이 누락된 경우                         | Ticket Object에 `ticket`을 지정해 주세요.                               |
| NO\_TYPE           | `type`이 누락된 경우                           | Data Type Object에 `indicator`를 지정해 주세요.                         |
| NO\_CODES          | `codes`가 누락된 경우                          | 수신할 마켓 코드 목록을 지정해 주세요.                                          |
| INVALID\_PARAM     | 지원하지 않는 캔들 단위, 지표, 계산 파라미터 또는 조합을 요청한 경우 | `indicator`의 지원 범위를 확인해 주세요.                                    |
| TOO\_MANY\_REQUEST | 구독 조합 상한 또는 WebSocket 공통 요청 제한을 초과한 경우   | 전체 구독 조합을 100개 이하로 줄이거나, 요청 제한에 도달한 경우 일정 시간 후 다시 요청해 주세요.      |

<Callout icon="fad fa-triangle-exclamation" theme="warn">
  ### **오류 응답 후에는 다시 연결해야 합니다.**

  * 요청 검증에 실패하면 서버는 오류 응답을 전송한 뒤 WebSocket 연결을 종료합니다. 하나의 요청에 여러 지표를 포함했더라도 일부만 등록하지 않고 요청 전체를 거부합니다.
  * 오류 원인을 수정한 뒤 다시 연결하고, 필요한 구독 항목 전체를 요청해 주세요.
  * 미지원·미존재 마켓 코드라는 이유만으로 `INVALID_PARAM` 오류가 반환되지는 않습니다. 해당 마켓의 데이터가 없는 동안에는 데이터가 전송되지 않으므로, 구독할 마켓 코드를 확인해 주세요.
</Callout>

<br />

## 요청 수 제한

WebSocket 연결 요청과 데이터 요청 메시지에는 아래 제한이 적용됩니다. 같은 Rate Limit 그룹의 요청은 해당 적용 단위에서 한도를 함께 사용합니다.

| Rate Limit 그룹       | 정책                   | 적용 단위 |
| ------------------- | -------------------- | ----- |
| `websocket-connect` | 초당 최대 5회             | 포켓    |
| `websocket-message` | 초당 최대 5회, 분당 최대 100회 | 커넥션   |

요청 수 제한과 구독 조합 수 제한은 별개입니다. **요청 횟수가 제한 이내여도 하나의 요청에 포함한 지표 구독 조합은 최대 100개까지 허용됩니다.**

허용 요청 수와 구독 조합 상한은 운영 정책에 따라 변경될 수 있습니다. 요청 수 제한 정책의 자세한 내용은 [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits) 문서를 참고해 주세요.

<Callout icon="fad fa-gauge" theme="warn">
  ### **WebSocket 요청 수 제한 관리**

  WebSocket은 REST API와 달리 잔여 요청 수를 별도로 제공하지 않습니다. 클라이언트에서 WebSocket 연결 및 데이터 요청 메시지의 전송 횟수를 관리하여 요청 수 제한을 준수해 주세요. 요청 수 제한에 도달한 경우 일정 시간 대기한 후 다시 요청해 주세요.
</Callout>

<Accordion title="이용 시 유의사항" icon="fa-info-circle">
  - 지표 값은 시장 상황에 따라 일시적으로 실제 종가를 반영하지 못할 수 있으며, 네트워크 환경 등에 따라 클라이언트의 데이터 수신 시각에 차이가 발생할 수 있습니다.

  - 거래지원 종료 등으로 구독 중인 마켓의 데이터 생성이 중단되면 별도의 WebSocket 알림 없이 지표 데이터 전송이 중단될 수 있습니다. 종목명 변경 등으로 마켓 코드가 유지되는 경우에는 기존 구독이 유지됩니다.

  - 제공되는 지표 값은 참고 정보이며, 특정 디지털 자산의 매수 또는 매도를 권유하는 정보가 아닙니다.
</Accordion>

# Sibling pages

* [지표 조건 이벤트 (Indicator Signal)](https://docs.upbit.com/kr/reference/websocket-indicator-signal.md)