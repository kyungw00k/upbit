---
updatedAt: 2026-09-01T06:40:51.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Upbit Skills

Upbit Skills에 대하여 안내합니다.

Upbit Skills은 AI 에이전트가 `upbit` CLI를 통해 업비트 API를 더 쉽고 정확하게 사용할 수 있도록 돕는 Skills 패키지입니다.

시세 조회, 잔고 확인, 주문, 입출금, 포켓, 트래블룰 검증(Travel Rule) 관련 작업을 자연어 요청에 맞게 처리할 수 있습니다.

예시:

> “KRW-BTC 현재가 알려줘”<br />“내 잔고 확인해줘”<br />“서브포켓 잔고 확인해줘”<br />“BTC 1만 원 시장가 매수 명령 만들어줘”<br />“USDT 출금 가능한 네트워크 확인해줘”<br />“Travel Rule 확인이 필요한 입금 내역 조회해줘”

<br />

## 설치

Node.js 18 이상이 필요합니다.

```bash
npx skills add upbit-official/upbit-agent-skills
npm install -g @upbit-official/upbit-cli
upbit --version
```

전역 설치가 필요한 경우:

```bash
npx skills add upbit-official/upbit-agent-skills -g
```

<br />

## 주요 기능

| 구분      | 할 수 있는 일                                          | 인증    |
| ------- | ------------------------------------------------- | ----- |
| 시세      | 현재가, 호가, 체결, 캔들, 마켓 목록 조회                         | 불필요   |
| 계정      | 잔고, 보유 자산, API 키 정보 조회                            | 필요    |
| 주문      | 주문 가능 정보 조회, 주문 생성, 주문 조회, 주문 취소, 주문 테스트          | 필요    |
| 입금      | 입금 주소 조회·생성, 입금 내역 조회                             | 필요    |
| 출금      | 출금 가능 정보 조회, 출금 요청, 출금 취소                         | 필요    |
| 포켓      | 포켓 목록, API Key, 잔고 조회, 메인포켓·서브포켓 자산 이전 및 이전 내역 조회 | 필요    |
| 트래블룰 검증 | VASP 목록 조회, 계정주 확인이 필요한 입금건에 대한 검증                | 필요    |
| 지갑 상태   | 자산별 입출금 지원 상태 조회                                  | 불필요   |
| 예제      | 주문, 입출금, DCA, TP/SL 시나리오 실행                       | 일부 필요 |

권장 설정 방식:

```bash
upbit config set
```

또는 환경 변수 사용:

```bash
export UPBIT_ACCESS_KEY=<your-access-key>
export UPBIT_SECRET_KEY=<your-secret-key>
```

API 키와 Secret Key는 응답, 로그, 코드 예제에 노출하지 않아야 합니다.

<br />

## 자주 쓰는 명령

> 포켓 자산 조회는 메인포켓 API Key로만 이용할 수 있습니다. 서브포켓 API Key로 요청하는 경우 out\_of\_scope 오류가 발생합니다.

```bash
# 현재가 조회
upbit tickers list-by-trading-pairs --markets KRW-BTC

# 마켓 목록 조회
upbit trading-pairs list

# 호가 조회
upbit orderbooks list --markets KRW-BTC

# 잔고 조회
upbit accounts list

# 포켓 목록 조회
upbit pockets list

# 서브포켓 잔고 조회
upbit pockets retrieve-balance 
  --uuid "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"

# 주문 가능 정보 조회
upbit orders retrieve-chance --market KRW-BTC

# 실제 주문 없이 주문 형식 테스트
upbit orders test-create --market KRW-BTC --side bid --ord-type price --price 10000
```

<br />

## 업비트 마켓 표기 방식

업비트는 `{QUOTE}-{BASE}` 형식의 마켓 코드를 사용합니다.

예: `BTC-KRW`, `BTC/KRW` → `KRW-BTC`

| 마켓         | 의미            |
| ---------- | ------------- |
| `KRW-BTC`  | BTC를 KRW로 거래  |
| `KRW-ETH`  | ETH를 KRW로 거래  |
| `BTC-ETH`  | ETH를 BTC로 거래  |
| `USDT-XRP` | XRP를 USDT로 거래 |

<br />

## 주문 기본 개념

| 값     | 의미 |
| ----- | -- |
| `bid` | 매수 |
| `ask` | 매도 |

| 주문 유형    | 의미     | 필요한 값             |
| -------- | ------ | ----------------- |
| `limit`  | 지정가 주문 | `price`, `volume` |
| `price`  | 시장가 매수 | `price`           |
| `market` | 시장가 매도 | `volume`          |
| `best`   | 최유리 주문 | 조건에 따라 다름         |

시장가 매수는 수량이 아니라 사용할 총 금액을 `price`에 입력합니다.

```bash
upbit orders test-create --market KRW-BTC --side bid --ord-type price --price 10000
```

시장가 매도는 매도할 수량을 `volume`에 입력합니다.

```bash
upbit orders test-create --market KRW-BTC --side ask --ord-type market --volume 0.001
```

처음 주문하는 마켓에서는 주문 전 가능 정보를 먼저 확인하는 것이 좋습니다.

```bash
upbit orders retrieve-chance --market KRW-BTC
```

<br />

## 안전한 실행 규칙

다음 작업은 실제 자산에 영향을 줄 수 있으므로 바로 실행하면 안 됩니다.

* 주문 생성
* 주문 취소
* 출금 요청
* 출금 취소
* 원화 입금 요청
* 입금 주소 생성
* 메인포켓 자산 이전
* 서브포켓 자산 이전
* Travel Rule 입금 검증
* 자동 매매 예제의 실제 실행

AI 에이전트는 먼저 실행 내용을 요약하고, 사용자가 아래 문구를 단독으로 입력한 경우에만 실행해야 합니다.

```text
CONFIRM
```

`orders test-create`와 Dry run은 실제 주문을 생성하지 않으므로 `CONFIRM` 없이 사용할 수 있습니다.

<br />

## 출금 전 확인 사항

출금은 되돌릴 수 없으므로 실행 전 아래 항목을 확인해야 합니다.

* 자산 코드
* 출금 수량
* 출금 주소
* 네트워크 `net_type`
* 보조 주소 또는 메모 필요 여부
* 출금 안심차단 해지 여부
* 출금 주소 등록 여부

멀티체인 자산은 출금 네트워크를 반드시 확인해야 합니다.

```bash
upbit withdraws list-coin-addresses
```

<br />

## 입출금 상태

| 상태                      | 의미                |
| ----------------------- | ----------------- |
| `PROCESSING`            | 처리 중              |
| `ACCEPTED`              | 완료                |
| `CANCELLED`             | 취소                |
| `REJECTED`              | 거절                |
| `TRAVEL_RULE_SUSPECTED` | Travel Rule 확인 필요 |
| `REFUNDING`             | 반환 진행 중           |
| `REFUNDED`              | 반환 완료             |

`TRAVEL_RULE_SUSPECTED` 상태의 입금은 `travel-rule` 명령으로 검증할 수 있습니다.

<br />

## 예제 코드

```bash
# 인증 불필요 예제
bash examples/quotation_kr.sh

# 거래대금 상위 페어 조회
bash examples/indicators_kr.sh

# 주문 흐름 Dry run
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders_kr.sh

# 실제 주문 실행
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders_kr.sh
```

실제 실행은 반드시 사용자가 `CONFIRM`한 뒤에만 진행해야 합니다.

<br />

## 예제 목록

| 예제                  | 설명                        | 인증  | 기본 동작    |
| ------------------- | ------------------------- | --- | -------- |
| `quotation_kr.sh`   | 시세, 캔들, 체결, 호가 조회         | 불필요 | 조회 전용    |
| `indicators_kr.sh`  | KRW 마켓 24시간 거래대금 상위 페어 조회 | 불필요 | 조회 전용    |
| `orders_kr.sh`      | 주문 생성, 조회, 취소 흐름          | 필요  | Dry run  |
| `orders_test_kr.sh` | 주문 생성 테스트 API로 주문 유형 검증   | 필요  | 실제 주문 없음 |
| `deposits_kr.sh`    | 입금 주소 및 입금 내역 관리          | 필요  | Dry run  |
| `withdrawals_kr.sh` | 출금 정보 및 출금 흐름 확인          | 필요  | Dry run  |
| `dca_kr.sh`         | 정기 시장가 매수 자동화             | 필요  | Dry run  |
| `tp_sl_kr.sh`       | 익절·손절 자동 매도               | 필요  | Dry run  |

<br />

## AI 에이전트 원칙

* 사용자 언어에 맞춰 응답합니다.
* 공개 조회는 인증 없이 수행합니다.
* 개인 계정 및 포켓 관련 요청은 인증 설정 여부를 먼저 확인합니다.
* 실제 자산에 영향을 주는 작업은 실행 전 반드시 요약하고 확인을 받습니다.
* 메인포켓 및 서브포켓 자산 이전은 실행 전에 출발 포켓, 대상 포켓, 이전 자산 및 수량을 확인합니다.
* 실제 거래 전에는 `orders test-create` 또는 Dry run을 우선 사용합니다.
* API 키와 Secret Key는 절대 출력하지 않습니다.
* 사용자의 마켓 표기가 업비트 형식과 다르면 올바른 형식으로 변환합니다.
* 명령 옵션이 불확실하면 `upbit <resource> <command> --help`를 먼저 확인합니다.
* 한국어 응답에서는 주요 API 필드를 자연스러운 용어로 설명합니다.

| API 필드                | 한국어 설명          |
| --------------------- | --------------- |
| `bid`                 | 매수              |
| `ask`                 | 매도              |
| `balance`             | 보유 잔고           |
| `locked`              | 주문 또는 출금에 묶인 수량 |
| `trade_price`         | 현재 체결가          |
| `acc_trade_price_24h` | 24시간 누적 거래대금    |

<br />

## 상세 명령 확인

```bash
upbit <resource> <command> --help
```

예시:

```bash
upbit orders create --help
```

포켓 명령을 확인하려면:

```bash
upbit pockets --help
```

<br />

## 참조 문서

| 문서                            | 내용          |
| ----------------------------- | ----------- |
| `references/setup.md`         | 설치 및 인증 설정  |
| `references/orders.md`        | 주문 명령       |
| `references/tickers.md`       | 현재가 조회      |
| `references/candles.md`       | 캔들 조회       |
| `references/orderbooks.md`    | 호가 조회       |
| `references/trades.md`        | 체결 조회       |
| `references/trading-pairs.md` | 마켓 목록       |
| `references/withdraws.md`     | 출금          |
| `references/deposits.md`      | 입금          |
| `references/pockets.md`       | 포켓          |
| `references/travel-rule.md`   | Travel Rule |
| `references/account.md`       | 계정 및 지갑 상태  |
| `references/output.md`        | 출력 형식 및 필터링 |
| `references/glossary.md`      | 한국어·영어 용어집  |

<br />

## 문의 및 주의

버그 제보 및 피드백:

```text
open-api@upbit.com
```

Upbit CLI Skill은 업비트 Open API와 `upbit` CLI 사용을 보조하는 도구입니다.
이 Skill을 통해 실행되는 모든 주문, 출금, 포켓 간 자산 이전, 자동 매매 및 기타 계정 관련 작업의 최종 책임은 사용자에게 있습니다.

<br />

# Sub pages

* [Cursor](https://docs.upbit.com/kr/docs/cursor.md)
* [Claude Code](https://docs.upbit.com/kr/docs/claude-code.md)

# Sibling pages

* [개요](https://docs.upbit.com/kr/docs/upbit-ai-agent-overview.md)
* [Upbit SDK](https://docs.upbit.com/kr/docs/upbit-sdk.md)
* [Upbit CLI](https://docs.upbit.com/kr/docs/upbit-cli.md)
* [Upbit Strategy Toolkit](https://docs.upbit.com/kr/docs/upbit-strategy-toolkit.md)
* [CCXT 라이브러리 연동 안내](https://docs.upbit.com/kr/docs/ccxt-library-guide.md)
* [업비트 어시스턴트 이용 안내](https://docs.upbit.com/kr/docs/assistant-guide.md)