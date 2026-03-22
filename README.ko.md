# upbit — AI 친화적 Upbit 거래소 CLI

[English](README.md)

> 사람과 AI 에이전트 모두를 위해 설계된 유일한 Upbit CLI.

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![API Coverage](https://img.shields.io/badge/Upbit_API-35_commands-brightgreen)]()

## 빠른 데모

![demo](docs/demo.ko.gif)

## 왜 upbit인가?

### AI 퍼스트 설계
- `upbit tool-schema`로 LLM 도구 호출용 JSON Schema 내보내기 (OpenAI/MCP) — 35개 명령 지원
- Non-TTY 환경 자동 감지 후 JSON 출력 — AI 에이전트 파이프라인에 적합
- 프로그래밍 방식의 에러 처리를 위한 구조화된 에러 코드
- 사람의 터미널 사용과 AI 에이전트 도구로의 사용을 동시에 고려한 설계

### 스마트 트레이딩
- 호가 단위 자동 보정으로 KRW/BTC/USDT 마켓 주문 실패 방지
- SQLite 캔들 캐시와 자동 페이지네이션 (`--from 2025-01-01`로 전체 이력 조회)
- WebSocket 구독 전 마켓 유효성 검증
- 매수/매도 확인 프롬프트 제공, `--force`로 자동화 지원

### 다국어 지원
- `LANG=ko_KR` → 한국어 인터페이스
- 기본값 → 영어 인터페이스
- POSIX 로케일 표준 (LC_ALL > LC_MESSAGES > LANG)

### 독립 실행형
- 단일 정적 바이너리, 런타임 의존성 없음
- 자체 업데이트: `upbit update`
- 크로스 플랫폼: Linux, macOS, Windows (amd64, arm64)

## 설치

### 릴리스에서 설치 (권장)

```bash
# macOS / Linux
curl -sSL https://github.com/kyungw00k/upbit/releases/latest/download/upbit_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m).tar.gz | tar xz
mv upbit ~/.local/bin/
```

### 소스에서 빌드

```bash
git clone https://github.com/kyungw00k/upbit.git
cd upbit
make install   # ~/.local/bin/upbit에 설치
```

### Go로 설치

```bash
go install github.com/kyungw00k/upbit/cmd/upbit@latest
```

## 빠른 시작

```bash
# 인증 불필요 — 시세 데이터
upbit ticker KRW-BTC
upbit candle KRW-BTC -i 1d -c 5

# API 키 설정 (환경변수만 사용 — 디스크에 저장하지 않음)
export UPBIT_ACCESS_KEY=your_key
export UPBIT_SECRET_KEY=your_secret

# 거래
upbit buy KRW-BTC -p 100000000 -V 0.001
upbit sell KRW-BTC -p 110000000 -V 0.001
upbit balance
```

## AI 연동

```bash
# 전체 명령어 스키마를 LLM 함수 호출용으로 내보내기
upbit tool-schema

# 특정 명령어 스키마만 내보내기
upbit tool-schema buy

# AI 에이전트는 자동으로 JSON 출력 (Non-TTY)
upbit ticker KRW-BTC | jq '.trade_price'

# 어떤 환경에서든 JSON 강제 출력
upbit ticker KRW-BTC -o json

# 특정 필드만 선택
upbit ticker KRW-BTC --json market,trade_price,signed_change_rate
```

### 예시: tool-schema 출력

```json
[
  {
    "name": "upbit_buy",
    "description": "매수 주문",
    "parameters": {
      "type": "object",
      "properties": {
        "market": { "type": "string", "description": "매수할 마켓 코드 (예: KRW-BTC)" },
        "price": { "type": "string", "description": "주문 단가" },
        "volume": { "type": "string", "description": "주문 수량" },
        "total": { "type": "string", "description": "주문 총액 (시장가 매수)" },
        "force": { "type": "boolean", "description": "확인 프롬프트 스킵" },
        "test": { "type": "boolean", "description": "테스트 주문 (실제 체결 안됨)" }
      },
      "required": ["market"]
    }
  }
]
```

## 명령어

### 시세 (인증 불필요)

| 명령 | 설명 |
|------|------|
| `ticker <market...>` | 현재가 조회 |
| `candle <market>` | OHLCV 캔들 조회 |
| `orderbook <market...>` | 호가 조회 |
| `trades <market>` | 체결 내역 조회 |
| `market` | 마켓 목록 조회 |
| `tick-size <market...>` | 호가 단위 조회 |
| `orderbook-levels <market...>` | 호가 모아보기 단위 조회 |

### 거래 (인증 필요)

| 명령 | 설명 |
|------|------|
| `buy <market>` | 매수 주문 (`-p` 단가, `-V` 수량, `-t` 총액) |
| `sell <market>` | 매도 주문 |
| `balance [currency]` | 계정 잔고 조회 (KRW 평가액 포함) |
| `order list` | 주문 목록 조회 |
| `order show <uuid>` | 주문 상세 조회 |
| `order cancel <uuid>` | 주문 취소 |
| `order replace <uuid>` | 주문 정정 (취소 후 재주문) |
| `order chance <market>` | 주문 가능 정보 조회 |

### 입출금 (인증 필요)

| 명령 | 설명 |
|------|------|
| `wallet` | 입출금 서비스 상태 조회 |
| `deposit list` | 입금 목록 조회 |
| `deposit show <uuid>` | 개별 입금 조회 |
| `deposit address [currency]` | 입금 주소 조회 |
| `deposit address create` | 입금 주소 생성 요청 |
| `withdraw list` | 출금 목록 조회 |
| `withdraw show <uuid>` | 개별 출금 조회 |
| `withdraw request` | 출금 요청 |
| `withdraw cancel <uuid>` | 출금 취소 |
| `travelrule vasps` | 트래블룰 지원 거래소 목록 |
| `travelrule verify-txid` | TxID 기반 트래블룰 검증 |
| `travelrule verify-uuid` | UUID 기반 트래블룰 검증 |

### 실시간 (WebSocket)

| 명령 | 설명 |
|------|------|
| `watch ticker <market...>` | 현재가 실시간 스트림 |
| `watch orderbook <market...>` | 호가 실시간 스트림 |
| `watch trade <market...>` | 체결 실시간 스트림 |
| `watch candle <market...>` | 캔들 실시간 스트림 |
| `watch my-order` | 내 주문 실시간 스트림 (인증 필요) |
| `watch my-asset` | 내 자산 실시간 스트림 (인증 필요) |

### 유틸리티

| 명령 | 설명 |
|------|------|
| `tool-schema [cmd]` | LLM/MCP용 JSON Schema 출력 |
| `api-keys` | API 키 목록 조회 |
| `cache` | 캐시 정보 조회 / `--clear`로 삭제 |
| `update` | 자체 업데이트 / `--check`로 확인만 |

## 실시간 TUI

Watch 명령은 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 기반 전체 화면 터미널 UI를 제공합니다:

- **watch ticker** — 실시간 가격 테이블 (상승=초록, 하락=빨강)
- **watch orderbook** — 스프레드 중심 매수/매도 호가 차트
- **watch trade** — 체결 스트림 스크롤
- **watch candle** — ASCII 캔들스틱 차트 + 거래량 바

복수 마켓: Tab 또는 ←/→로 전환

```bash
upbit watch ticker KRW-BTC KRW-ETH KRW-XRP
upbit watch orderbook KRW-BTC
upbit watch candle KRW-BTC -i 1m
```

## 출력 포맷

| 컨텍스트 | 기본값 | 변경 |
|---------|--------|------|
| 터미널 (TTY) | 정렬된 테이블 | `-o json`, `-o csv` |
| 파이프 (Non-TTY) | 압축 JSON | `-o table`, `-o jsonl` |

```bash
upbit ticker KRW-BTC              # 테이블 (터미널)
upbit ticker KRW-BTC | jq .       # 자동 JSON (파이프)
upbit ticker KRW-BTC -o csv       # CSV
upbit ticker KRW-BTC --json price # 특정 필드만
```

## 자체 업데이트

```bash
upbit update          # 최신 버전 다운로드 및 설치
upbit update --check  # 확인만 하고 다운로드하지 않음
```

## 릴리스 워크플로우

릴리스는 [go-semantic-release](https://github.com/go-semantic-release/semantic-release) + [goreleaser](https://goreleaser.com)를 통해 완전 자동화됩니다:

1. [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:` 등)로 `main`에 푸시
2. 버전이 자동으로 결정됨 (semver)
3. Linux, macOS, Windows (amd64/arm64) 크로스 컴파일 바이너리 생성
4. 변경 이력과 체크섬이 포함된 GitHub Release 생성
5. `upbit update`로 업데이트 가능

## 캔들 캐시

```bash
# 지정 날짜부터 자동 페이지네이션 (SQLite 캐시)
upbit candle KRW-BTC --from 2025-01-01

# 캐시 무시
upbit candle KRW-BTC --from 2025-01-01 --no-cache

# 캐시 정보 확인 및 관리
upbit cache
upbit cache --clear
```

## 인증

API 키는 **환경변수에서만** 읽으며, 디스크에 저장하지 않습니다.

```bash
export UPBIT_ACCESS_KEY=your_access_key
export UPBIT_SECRET_KEY=your_secret_key
```

시세 조회 명령어는 인증 없이 사용할 수 있습니다. 거래, 입출금, 실시간 개인 스트림(`watch my-order`, `watch my-asset`)은 인증이 필요합니다.

## 감사의 글

- 캔들스틱 차트 렌더링은 [cli-candlestick-chart](https://github.com/Julien-R44/cli-candlestick-chart)에서 영감을 받았습니다

## 라이선스

MIT
