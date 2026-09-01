---
updatedAt: 2026-06-01T04:26:20.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Claude Code

Claude Code 환경에서 Upbit CLI를 활용하는 방법을 안내합니다.

[Claude Code](https://code.claude.com/docs/ko/overview)는 터미널 환경에서 동작하는 AI 코딩 에이전트입니다.<br />Upbit CLI Skill를 설치하면 Claude Code가 업비트 Open API와 `upbit` CLI 사용 방법을 이해하고, 사용자의 자연어 요청에 맞는 명령을 더 정확하게 작성하거나 실행할 수 있습니다.

이 가이드는 Claude Code에서 Upbit Agent Skills를 설정하고 안전하게 사용하는 방법을 설명합니다.

<br />

## 1. Upbit Agent Skill 설치

Claude Code에서 업비트 관련 작업을 자주 수행한다면 Skill을 전역 설치하는 방식을 권장합니다.

```bash bash
npx skills add upbit-official/upbit-agent-skills -g
```

<br />

## 2. Upbit CLI 설치

Upbit CLI Skill은 `upbit` CLI와 함께 동작합니다.<br />아래 명령으로 CLI를 설치합니다.

```bash bash
npm install -g @upbit-official/upbit-cli
```

설치 확인:

```bash bash
upbit --version
```

<br />

## 3. 인증 정보 설정

시세, 호가, 체결, 캔들, 마켓 목록 등 공개 API는 인증 없이 사용할 수 있습니다.

잔고 조회, 주문, 입출금, 트래블룰 등 계정 관련 API를 사용하려면 업비트 Open API 키가 필요합니다.

권장 방식:

```bash bash
upbit config set
```

설정한 인증 정보는 `~/.upbit/config`에 저장되며, 이후 CLI 명령에서 자동으로 사용됩니다.

환경 변수로도 설정할 수 있습니다.

```bash bash
export UPBIT_ACCESS_KEY=<your-access-key>
export UPBIT_SECRET_KEY=<your-secret-key>
```

API Key와 Secret Key는 Claude Code 채팅, 코드 파일, 로그에 그대로 노출하지 마세요.

## 4. Claude Code 프로젝트 규칙 추가(선택)

Claude Code에서 Upbit CLI Skill이 정상적으로 동작하는 경우 이 설정은 필수는 아니며,

프로젝트 단위로 Claude Code를 사용할 경우, 루트 디렉터리에 `CLAUDE.md` 파일을 만들고 아래 규칙을 추가하는 것을 권장합니다.

````markdown
# Upbit CLI 규칙

사용자가 업비트 시세, 잔고, 주문, 입금, 출금, 지갑 상태, Travel Rule 관련 작업을 요청하면 Upbit CLI Skill을 사용합니다.

업비트 REST API 작업은 `upbit` CLI를 사용합니다.

## 공개 API

다음 공개 API는 인증 없이 사용할 수 있습니다.

- `tickers`
- `orderbooks`
- `trades`
- `candles`
- `trading-pairs`
- `wallet-status`

## 인증 필요 API

다음 API는 인증이 필요합니다.

- `accounts`
- `api-keys`
- `orders`
- `withdraws`
- `deposits`
- `travel-rule`

## 보안 규칙

API Key와 Secret Key는 응답, 코드, 로그에 그대로 출력하거나 저장하지 않습니다.

## 실행 확인 규칙

실제 자산에 영향을 줄 수 있는 쓰기 작업을 실행하기 전에는 전체 명령을 보여주고, 사용자가 `CONFIRM`을 입력하도록 요청합니다.

쓰기 작업 예시:

- 주문 생성
- 주문 취소
- 출금 요청
- 출금 취소
- 원화 입금 요청
- 입금 주소 생성
- Travel Rule 입금 검증

`orders test-create`는 실제 주문을 생성하지 않으므로 `CONFIRM`이 필요하지 않습니다.

## 마켓 형식

업비트 마켓 형식은 `{QUOTE}-{BASE}`입니다.

예시:

- `BTC/KRW` → `KRW-BTC`
- `BTC-KRW` → `KRW-BTC`
- `ETH/KRW` → `KRW-ETH`

## 주문 규칙

시장가 매수는 `ord_type=price`를 사용하고, `price`에는 사용할 총 금액을 입력합니다.

시장가 매도는 `ord_type=market`을 사용하고, `volume`에는 매도할 자산 수량을 입력합니다.

명령 옵션이 불확실하면 아래 도움말을 먼저 확인합니다.

```bash
upbit <resource> <command> --help
````

<br />

## 5. Claude Code에서 사용할 수 있는 요청 예시

현재가 조회:

```Text bash
KRW-BTC 현재가 조회해줘
```

Claude Code가 사용할 수 있는 명령:

```bash bash
upbit tickers list-by-trading-pairs --markets KRW-BTC
```

잔고 조회:

```Text bash
내 잔고 확인해줘
```

```bash bash
upbit accounts list
```

주문 가능 정보 조회:

```Text bash
KRW-BTC 주문 가능 정보 확인해줘
```

```bash bash
upbit orders retrieve-chance --market KRW-BTC
```

실제 주문 없이 주문 형식 테스트:

```Text bash
BTC 1만 원 시장가 매수 테스트 명령 만들어줘
```

```bash bash
upbit orders test-create --market KRW-BTC --side bid --ord-type price --price 10000
```

<br />

## 6. 실제 주문 또는 출금 전 확인

Claude Code가 실제 자산에 영향을 줄 수 있는 명령을 실행하려는 경우, 바로 실행하지 않고 먼저 실행 내용을 요약해야 합니다.

예시:

```Text bash
[실행 전 확인]
작업: KRW-BTC 시장가 매수
금액: 10,000 KRW
실행 명령:
upbit orders create --market KRW-BTC --side bid --ord-type price --price 10000

주의 사항: 실제 주문이 생성됩니다.

실행하려면 CONFIRM을 입력하세요.
```

사용자가 `CONFIRM`을 단독으로 입력한 경우에만 실행합니다.

<br />

## 7. 안전한 사용을 위한 권장 사항

* 처음 사용하는 마켓은 주문 전 `orders retrieve-chance`로 주문 가능 정보와 최소 주문 금액을 확인합니다.
* 실제 주문 전에는 `orders test-create`로 주문 형식을 먼저 검증합니다.
* 출금 전에는 자산, 네트워크, 주소, 보조 주소 또는 메모 필요 여부를 반드시 확인합니다.
* API Key와 Secret Key는 절대 코드에 하드코딩하지 않습니다.
* 자동 매매 예제는 Dry run으로 먼저 검토한 뒤 실제 실행 여부를 결정합니다.

<br />

# Sibling pages

* [Cursor](https://docs.upbit.com/kr/docs/cursor.md)