---
updatedAt: 2026-07-09T13:23:10.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 기술 레퍼런스

이 문서는 전략 JSON 스키마, CLI 사용법, 백테스트 엔진의 상세 동작을 다룹니다. 처음 사용하는 경우 Upbit Strategy Toolkit — 시작 가이드를 먼저 참고하세요.

> 본 도구는 **백테스트(과거 데이터 시뮬레이션) 전용**이며, 실시간 주문이나 실제 자금 이동은 발생하지 않습니다. 실시간 주문·계좌 변경·자동 매매가 필요한 경우 별도의 방법([Upbit Skills](https://docs.upbit.com/kr/docs/upbit-agent-skills) 등)을 이용하길 권장드립니다.

> 본 도구는 **베타(Beta)**&#xB85C; 제공됩니다. 정식 출시 전까지 기능·사양 및 전략 JSON 스키마가 사전 예고 없이 변경될 수 있으며, 이용 중 일부 동작이 원활하지 않을 수 있는 점 양해 부탁드립니다.

## 스킬과 도구 구성

### 제공 스킬

자연어로 요청해도 에이전트가 적절히 선택하지만, 아래 명령으로 직접 호출할 수도 있습니다.

* `/create-strategy` — 아이디어를 전략으로 정리·저장
* `/backtest` — 백테스트 실행
* `/setup` — 첫 환경 설정 또는 환경 진단 (CLI 오류 시)

### Upbit Strategy Toolkit과 Upbit Skills

두 도구 모두 AI 에이전트 기반이지만 역할이 다릅니다. 본 도구는 전략을 과거 데이터로 검증(백테스트)하는 데 쓰이며, 실시간 조회·주문 등 실거래는 Upbit Skills에서 지원합니다.

> 전략을 충분히 검증한 뒤 실거래 여부를 판단하시기 바랍니다.

| 구분          | Upbit Strategy Toolkit | Upbit Skills             |
| ----------- | ---------------------- | ------------------------ |
| 목적          | 매매 전략의 과거 성과 검증(백테스트)  | 실시간 시세 조회 및 매매 실행        |
| 사용 데이터      | 업비트 과거 캔들 데이터          | 실시간 시세 + 실제 계좌           |
| 자금 이동       | 없음 (모의 시뮬레이션)          | 있음 (실제 주문·체결)            |
| 실거래 지원      | 미지원                    | 지원                       |
| 주요 산출물      | 전략 JSON, 백테스트 리포트(CSV) | 주문 접수·체결 결과, 계좌·잔고 조회 결과 |
| 인증(API Key) | 불필요 (공개 캔들 데이터만 사용)    | 필요 (주문·조회 권한)            |
| 리스크         | 자금 손실 없음               | 실제 자금 손실 가능              |

## 전략 스키마

전략은 JSON으로 표현되며, Pydantic 모델로 엄격히 검증됩니다. 에이전트가 자동으로 생성하므로 대화로만 사용하는 경우 이 섹션은 참고하지 않아도 됩니다. 전략 파일을 직접 확인하거나 수정하려는 경우에만 아래 레퍼런스를 참고합니다.

### 전략 JSON 예시

```json
{
  "name": "SMA5/20 golden cross",
  "market": "KRW-BTC",
  "timeframe": "1d",
  "stop_loss": 8,
  "take_profit": 20,
  "indicators": [
    { "type": "moving_average", "ref": "sma5", "params": { "type": "SMA", "period": 5 } },
    { "type": "moving_average", "ref": "sma20", "params": { "type": "SMA", "period": 20 } }
  ],
  "buy": {
    "operator": "AND",
    "conditions": [
      {
        "left": { "type": "indicator", "ref": "sma5.value" },
        "op": "cross_above",
        "right": { "type": "indicator", "ref": "sma20.value" }
      }
    ]
  },
  "sell": {
    "operator": "AND",
    "conditions": [
      {
        "left": { "type": "indicator", "ref": "sma5.value" },
        "op": "cross_below",
        "right": { "type": "indicator", "ref": "sma20.value" }
      }
    ]
  }
}
```

### 필드·지표·연산자

지원되는 필드·지표·연산자는 다음과 같습니다.

#### 최상위 필드

| 필드            | 타입             |  필수 | 기본값       | 설명                      |
| ------------- | -------------- | :-: | --------- | ----------------------- |
| `name`        | string         |  O  | —         | 전략 이름 (한글 가능)           |
| `market`      | string         |  X  | `KRW-BTC` | 거래 페어. 예: `KRW-BTC`     |
| `timeframe`   | enum           |  X  | `1d`      | 캔들 주기 (아래 목록)           |
| `exchange`    | enum           |  X  | `kr`      | 거래소 (`kr`)              |
| `indicators`  | array          |  O  | —         | 사용할 지표 목록 (`ref` 중복 불가) |
| `buy`         | object         |  O  | —         | 매수 조건 그룹                |
| `sell`        | object         |  O  | —         | 매도 조건 그룹                |
| `stop_loss`   | number \| null |  X  | `5` (%)   | 손절 퍼센트. 양수 또는 null      |
| `take_profit` | number \| null |  X  | `15` (%)  | 익절 퍼센트. 양수 또는 null      |

> 백테스트 실행 시 초기 자본 기본값은 KRW: 1,000,000 / BTC: 0.01 / USDT: 1,000 입니다.

<br />

#### 지원 타임프레임 (12종)

```text
1s  1m  3m  5m  10m  15m  30m  1h  4h  1d  1w  1M
```

#### 지원 지표 (15종)

각 지표는 조건식에서 참조할 출력 키(`ref`)를 제공합니다.

| type              | 출력 키(ref)                                               | 주요 파라미터                                              |
| ----------------- | ------------------------------------------------------- | ---------------------------------------------------- |
| `moving_average`  | `value`                                                 | `type`(SMA/EMA), `period`(1–200)                     |
| `bollinger_bands` | `bb_upper`, `bb_middle`, `bb_lower`                     | `period`, `multiplier`                               |
| `envelopes`       | `upper`, `center`, `lower`                              | `period`, `percent`                                  |
| `ichimoku_cloud`  | `Conversion`, `Base`, `Lagging`, `Leading1`, `Leading2` | `conversion`, `base`, `leading_span2`                |
| `rsi`             | `rsi`, `rsi_signal`                                     | `period`, `signal_period`, `signal_type`             |
| `macd`            | `macd`, `macd_signal`, `histogram`                      | `fast`, `slow`, `signal_period`                      |
| `atr`             | `atr`                                                   | `period`                                             |
| `stochastic_slow` | `slow_k`, `slow_d`                                      | `period`, `k_period`, `d_period`                     |
| `williams_r`      | `williams_r`                                            | `period`                                             |
| `adx`             | `adx`, `adx_pdi`, `adx_mdi`                             | `period`                                             |
| `obv`             | `obv`, `obv_signal`                                     | `signal_period`, `signal_type`                       |
| `cci`             | `cci`, `cci_signal`                                     | `period`, `signal_period`, `signal_type`             |
| `stochastic_rsi`  | `stoch_rsi_k`, `stoch_rsi_d`                            | `rsi_period`, `stoch_period`, `k_period`, `d_period` |
| `mfi`             | `mfi`                                                   | `period`                                             |
| `disparity`       | `disp_5`, `disp_10`, `disp_20`, `disp_60`               | `periods`                                            |

> 각 지표의 정의와 활용은 **지표 정의 문서**를 참고합니다.

#### 지원 연산자 (7종)

| 연산자           | 의미                            |
| ------------- | ----------------------------- |
| `gt` / `lt`   | 초과 / 미만                       |
| `gte` / `lte` | 이상 / 이하                       |
| `eq`          | 같음                            |
| `cross_above` | 상향 돌파 (양쪽 모두 시계열이어야 함, 상수 불가) |
| `cross_below` | 하향 돌파 (양쪽 모두 시계열이어야 함, 상수 불가) |

### 검증 정책

전략 JSON은 저장·실행 전에 다음 규칙으로 엄격히 검증됩니다. 하나라도 위반하면 거부됩니다.

* 카탈로그(15종)에 없는 지표 `type` 거부
* 지표 `ref` 중복 거부
* 정의되지 않은 지표·출력 `ref` 참조 거부
* `cross_above` / `cross_below`는 양변이 모두 시계열이어야 함 (상수 불가)
* 빈 `conditions` 거부
* `stop_loss` · `take_profit`는 양수 또는 null만 허용
* 스키마에 없는 필드 거부
* 실제 업비트 마켓 리스트에 존재하는 `market`인지 검증

### 자가수정 루프

에이전트가 생성한 전략이 위 검증을 통과하지 못하면, 9개 오류 패턴을 기반으로 LLM이 스스로 최대 3회까지 자동 패치합니다. 3회를 초과해도 통과하지 못하면, 임의로 고치지 않고 사용자에게 전략 가설을 다시 협의하도록 요청합니다.

### 백테스트 기본 기간

`backtest` 스킬은 사용자가 기간을 명시하지 않으면 타임프레임별 기본 기간을 자동 적용합니다 (예: `1d` → 1년, `1h` → 3개월, `1m` → 1주). 시작일·종료일을 직접 지정하면 그 값으로 덮어씁니다.

## CLI 직접 사용

에이전트 워크플로 없이 CLI를 직접 호출할 수도 있습니다.

### 전략 유효성 검사

```bash
uv run upbit-strategy-toolkit strategy validate strategies/examples/sma_cross.json
```

성공하면 `OK`를, 실패하면 검증 오류를 출력합니다.

### 백테스트 실행

```bash
uv run upbit-strategy-toolkit backtest run strategies/examples/sma_cross.json \
  --start 2024-01-01 --end 2024-12-31
```

추가 옵션: `--capital <float>`, `--fee-rate <float>`, `--force-refresh`

## 백테스트 해석 시 주의사항

### 포지션 진입

이 도구는 전략의 매수 조건이 충족된 **다음 봉 시가**에 보유 자본 전액으로 매수합니다.

매수 수량은 진입 수수료를 포함한 총비용이 보유 현금을 넘지 않도록 계산됩니다(`현금 ÷ (1 + 수수료율) ÷ 체결가`). 또한 체결가는 마켓별 호가 단위(tick)로 반올림되며, 주문 금액이 최소 주문 금액에 미달하면 진입하지 않습니다.

### 청산 우선순위

같은 캔들 안에서 여러 조건이 동시에 충족되면 다음 순서로 적용됩니다.

1. **손절** — 캔들 저가가 손절 가격 이하
2. **익절** — 캔들 고가가 익절 가격 이상
3. **매도** — 전략의 매도 조건 충족
4. **종료** — 백테스트 마지막 캔들에서 포지션 강제 청산

### 체결 가격

| 청산 유형 | 체결 가격                               |
| ----- | ----------------------------------- |
| 손절    | 손절 가격 (단, 갭 하락으로 시가 < 손절 가격인 경우 시가) |
| 익절    | 익절 가격 (단, 갭 상승으로 시가 > 익절 가격인 경우 시가) |
| 매도    | 시그널 발생 다음 캔들 시가                     |
| 종료    | 마지막 캔들 종가                           |

손절·익절 주문은 진입 시 설정됩니다. 일반적으로 해당 가격에서 체결되지만, 갭이 발생하면 시가에서 체결됩니다(슬리피지). 손절·익절은 진입 캔들에서는 평가되지 않고 다음 캔들부터 적용됩니다.

### 수수료

수수료율은 거래소 + 마켓 접두어(페어의 `-` 앞부분: KRW/BTC/USDT 등) 기준으로 내장 마켓룰에서 자동 적용됩니다. 업비트 한국(`kr`) 기준 예시는 다음과 같습니다.

| 마켓 접두어 | 수수료율  |
| ------ | ----- |
| KRW    | 0.05% |
| BTC    | 0.25% |
| USDT   | 0.25% |

같은 자산이라도 어느 마켓(KRW/BTC/USDT)에서 거래하느냐에 따라 수수료가 다릅니다. 승률과 Profit Factor는 수수료 적용 전 수치입니다.

> 슬리피지, 스프레드, 시장 충격은 반영되지 않습니다. 진입과 매도 시그널 청산은 캔들 시가에서 전량 체결을 가정하므로, 단기 타임프레임이나 유동성이 낮은 마켓에서는 실제 성과와 차이가 있을 수 있습니다.

### 수치 정밀도

모든 계산은 IEEE 754 배정밀도(float64)를 사용합니다. 결과값은 백테스트용 근사치이며, 실거래나 회계 목적으로는 사용하지 않습니다.

## 여러 전략 결과 분석

지표가 많아 결과를 직접 판단하기 어렵다면, 여러 전략을 백테스트한 뒤 결과를 모아 AI 에이전트에게 비교·분석을 맡길 수 있습니다.

1. 조건을 바꿔가며 여러 전략을 백테스트합니다. 결과 CSV는 `reports/` 폴더에 순차적으로 저장됩니다.
2. 저장된 리포트를 하나의 전략집(예: `strategy-report.md`)으로 정리합니다. 파일마다 전략 설정과 성과 요약(Total Return·Benchmark·MDD·Win Rate 등)을 함께 담습니다.
3. 정리한 전략집을 에이전트에 전달하고 분석을 요청합니다.

예시 프롬프트:

```text
"reports/ 폴더의 백테스트 결과들을 비교해서 전략집 strategy-report.md로 정리해줘.
각 전략의 설정, Total Return, Benchmark, MDD, Win Rate를 표로 만들고,
벤치마크 대비 초과 성과 기준으로 정렬해줘."

"strategy-report.md를 읽고, 어떤 전략이 하락장에서 손실을 잘 방어했는지,
어떤 전략이 과최적화로 의심되는지 분석해줘. 지표만 보고 단정하지 말고
거래 횟수와 MDD도 함께 근거로 들어줘."
```

> 에이전트의 분석은 참고 의견이며 투자 판단의 근거가 되지 않습니다. AI는 결과를 잘못 해석할 수 있으므로 원본 리포트 수치를 함께 확인하세요.

***

### 유의사항

* **백테스트 결과는 투자 조언이 아닙니다.** 수익률·승률·Profit Factor·Sharpe는 과거 데이터 기반 참고 수치이며, 미래 수익을 보장하지 않습니다.
* **전략 JSON만 지원합니다.** 임의 Python 코드 실행은 지원하지 않습니다. 새로운 전략 표현이 필요한 경우 Pydantic 모델과 명시적 스키마로 정의합니다.
* **실전 매매는 지원하지 않습니다.** 실시간 주문·계좌 변경·자동 매매가 필요한 경우 `upbit-agent-skills`를 사용합니다.
* 모든 계산은 IEEE 754 배정밀도(float64)를 사용하며, 실거래 주문 실행이나 정산 목적이 아닌 백테스트 참고용 근삿값입니다.

<Accordion title="책임 한정 및 이용 주의사항" icon="fa-info-circle">
  - Upbit Strategy Toolkit(이하 '본 툴킷')은 이용자가 사용하는 제3자 제공 AI 서비스(이하 '이용자 AI')를 활용하여 기술지표를 기반으로 매매 전략을 설계하고 백테스트를 할 수 있도록 지원하는 도구입니다.

- 본 툴킷은 이용자 AI를 활용한 투자 전략 수립 및 백테스트 수행 시 편의를 제공하기 위한 보조 도구에 불과합니다. 본 툴킷은 법적·재정적·투자적 자문 또는 조언을 목적으로 하지 않으며, 회사는 어떠한 경우에도 특정 종목, 거래 시점, 투자 방식, 투자 전략 등을 추천·권유하지 않습니다.

- 이용자는 본 툴킷을 참고 목적으로만 활용하여야 하며, 본 툴킷을 활용하여 이용자 AI가 생성·제공하는 결과물, 답변, 전략, 백테스트 결과, 전략 개선 제안 또는 그 밖의 산출물(이하 총칭하여 'AI 산출물')을 자신의 책임과 판단하에 독립적으로 검토·검증하여야 합니다. 이용자는 AI 산출물이 실제 시장 상황, 거래 환경, 수수료, 슬리피지, 유동성, 체결 가능성, 가격 변동성 등 다양한 요인을 충분히 반영하지 못할 수 있음을 충분히 인지하고 본 툴킷을 이용하여야 하며, 회사는 과거 데이터를 기초로 한 백테스트 결과가 장래 또는 실제 거래에서 동일하거나 유사한 투자 성과로 이어질 것임을 보장하지 않습니다. 이용자는 AI 산출물을 투자 판단의 유일한 또는 주요한 근거로 삼지 말아야 하며, 이용자 AI를 활용한 거래·투자 행위와 관련된 모든 책임은 이용자에게 귀속됩니다.

- 회사는 본 툴킷을 활용한 AI 산출물의 정확성, 진실성, 완전성, 무결성, 최신성 등을 보증하지 않으며, AI 산출물에는 회사의 어떠한 의사나 입장도 반영되지 않습니다.

- 본 툴킷은 실전 매매(Live Trading)와 직접 연동되지 않으며, 이용자 AI를 통한 실전 매매를 희망하는 경우 별도의 방법(Upbit Agent Skills 등)을 이용하여야 합니다. 이러한 회사의 입장에도 불구하고, 이용자가 본 툴킷을 임의로 수정·변경·우회하거나 제3자가 제공하는 프로그램, 플러그인, 서비스 등과 결합하여 실전 매매, 자동 매매 또는 이에 준하는 행위를 하는 경우, 그로 인하여 발생하는 모든 결과와 책임은 이용자에게 귀속됩니다.

- 회사는 이용자 AI 등 제3자가 제공하는 서비스의 가용성, 보안성, 적법성, 정확성 및 성능 등에 대하여 책임을 부담하지 않습니다. 이용자는 제3자 서비스 이용에 따른 약관, 정책 및 위험을 스스로 확인하여야 합니다. 회사는 본 툴킷을 검증되지 않은 제3자 서비스와 연동 또는 조합하여 사용하는 것을 권장하지 않으며, 이용자의 제3자 서비스 사용으로 발생한 보안 사고(API Key 노출 등을 포함)나 금전적 손해 등에 대하여 책임을 부담하지 않습니다.

- 이용자는 본 툴킷을 활용함에 있어 관련 이용약관(업비트 이용약관, Open API 이용약관 등)을 준수하여야 합니다.

- 회사는 이용자 간 또는 이용자와 제3자 간에 본 툴킷과 관련하여 발생한 분쟁에 관여할 법적 의무가 없으며, 이와 관련하여 어떠한 책임도 지지 않습니다.

- 본 툴킷은 베타 버전으로 제공됩니다. 회사의 사정에 의하여 별도의 사전 고지 없이 제공 중단되거나 수정·변경될 수 있습니다.

- 본 툴킷 및 이에 포함된 문서, 코드, 예시, 가이드, 명령어 등 관련 리소스에 관한 권리 일체는 회사 또는 회사가 지정한 정당한 권리자에게 귀속됩니다. 본 툴킷 및 관련 리소스의 복제·배포·변경 및 영리 목적 이용은 관련 저작권 및 약관의 적용을 받으며, 본 툴킷을 활용한 프로그램 및 관련 리소스 등의 상업적 배포·판매·이용은 허용되지 않습니다.
</Accordion>

***

### 함께 보면 좋은 문서

* 처음 시작한다면 → <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/upbit-strategy-toolkit">**Upbit Strategy Toolkit — 시작 가이드**</Anchor>
* 지표 정의가 궁금하다면 → <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/upbit-strategy-toolkit-indicators">**지표 정의 문서**</Anchor>

<br />

# Sibling pages

* [지표 정의](https://docs.upbit.com/kr/docs/upbit-strategy-toolkit-indicators.md)