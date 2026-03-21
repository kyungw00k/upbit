# upbit — Upbit 거래소 CLI

Upbit 거래소의 시세 조회, 매수/매도 주문, 입출금 관리, 실시간 WebSocket 스트림을 터미널에서 바로 사용할 수 있는 CLI 도구입니다.

## 설치

```bash
go install github.com/kyungw00k/upbit/cmd/upbit@latest
```

또는 소스에서 빌드:

```bash
git clone https://github.com/kyungw00k/upbit.git
cd upbit
make build
```

빌드 결과물은 `./bin/upbit`에 생성됩니다.

## 인증

환경변수로 API 키를 설정합니다:

```bash
export UPBIT_ACCESS_KEY=your_access_key
export UPBIT_SECRET_KEY=your_secret_key
```

시세 조회 명령어는 인증 없이 사용할 수 있습니다. 거래, 입출금, 실시간 개인 스트림(`watch my-order`, `watch my-asset`)은 인증이 필요합니다.

## 명령어

### 시세

| 명령 | 설명 | 예시 |
|------|------|------|
| `market` | 마켓 목록 조회 | `upbit market` |
| `ticker` | 현재가 조회 | `upbit ticker KRW-BTC` |
| `candle` | 캔들 조회 | `upbit candle KRW-BTC -i 1m -c 50` |
| `trades` | 체결 내역 조회 | `upbit trades KRW-BTC` |
| `orderbook` | 호가 조회 | `upbit orderbook KRW-BTC` |
| `orderbook-levels` | 호가 모아보기 단위 조회 | `upbit orderbook-levels` |
| `tick-size` | 호가 단위 조회 | `upbit tick-size KRW-BTC` |

### 거래

| 명령 | 설명 | 예시 |
|------|------|------|
| `balance` | 계정 잔고 조회 | `upbit balance` |
| `buy` | 매수 주문 | `upbit buy KRW-BTC -p 50000000 -V 0.001` |
| `sell` | 매도 주문 | `upbit sell KRW-BTC -p 50000000 -V 0.001` |
| `order list` | 주문 목록 조회 | `upbit order list` |
| `order show` | 주문 상세 조회 | `upbit order show <uuid>` |
| `order cancel` | 주문 취소 | `upbit order cancel <uuid>` |
| `order replace` | 주문 정정 (취소 후 재주문) | `upbit order replace <uuid> -p 51000000` |
| `order chance` | 마켓별 주문 가능 정보 조회 | `upbit order chance KRW-BTC` |

### 입출금

| 명령 | 설명 | 예시 |
|------|------|------|
| `wallet` | 입출금 서비스 상태 조회 | `upbit wallet` |
| `deposit list` | 입금 목록 조회 | `upbit deposit list` |
| `deposit show` | 개별 입금 조회 | `upbit deposit show <uuid>` |
| `deposit address` | 입금 주소 조회 | `upbit deposit address BTC` |
| `withdraw list` | 출금 목록 조회 | `upbit withdraw list` |
| `withdraw show` | 개별 출금 조회 | `upbit withdraw show <uuid>` |
| `withdraw request` | 출금 요청 | `upbit withdraw request BTC -a 0.01 --address <addr>` |
| `withdraw cancel` | 출금 취소 | `upbit withdraw cancel <uuid>` |
| `travelrule` | 트래블룰 검증 관리 | `upbit travelrule` |

### 실시간

| 명령 | 설명 | 예시 |
|------|------|------|
| `watch ticker` | 현재가 실시간 스트림 | `upbit watch ticker KRW-BTC KRW-ETH` |
| `watch orderbook` | 호가 실시간 스트림 | `upbit watch orderbook KRW-BTC` |
| `watch trade` | 체결 실시간 스트림 | `upbit watch trade KRW-BTC` |
| `watch candle` | 캔들 실시간 스트림 | `upbit watch candle KRW-BTC -i 1m` |
| `watch my-order` | 내 주문 실시간 스트림 (인증 필요) | `upbit watch my-order` |
| `watch my-asset` | 내 자산 실시간 스트림 (인증 필요) | `upbit watch my-asset` |

### 유틸

| 명령 | 설명 | 예시 |
|------|------|------|
| `api-keys` | API 키 목록 조회 | `upbit api-keys` |
| `cache` | 캐시 관리 (정보 조회, 삭제) | `upbit cache --clear` |
| `tool-schema` | 명령어 JSON Schema 출력 (LLM/MCP용) | `upbit tool-schema ticker` |

## 출력 포맷

- **TTY** (터미널) &rarr; 테이블 형식 (기본)
- **Non-TTY** (파이프, 리다이렉트) &rarr; JSON (자동)
- `-o` 플래그로 명시적 지정: `table`, `json`, `jsonl`, `csv`
- `--json field1,field2`로 특정 필드만 JSON 출력

```bash
# 테이블 출력
upbit ticker KRW-BTC

# JSON 출력
upbit ticker KRW-BTC -o json

# 특정 필드만 JSON 출력
upbit ticker KRW-BTC --json market,trade_price

# CSV 출력
upbit candle KRW-BTC -o csv

# 파이프에서는 자동으로 JSON
upbit ticker KRW-BTC | jq '.trade_price'
```

## 캔들 캐시

`candle` 명령어에 `--from` 플래그를 사용하면 지정한 시작일부터 현재까지 자동으로 페이지네이션하여 데이터를 수집합니다. 수집된 캔들 데이터는 SQLite 데이터베이스에 캐시되어, 동일한 구간을 다시 조회할 때 API 호출 없이 빠르게 결과를 반환합니다.

```bash
# 2025년 1월 1일부터 현재까지 일봉 조회 (최초 실행 시 API 호출, 이후 캐시)
upbit candle KRW-BTC --from 2025-01-01

# 캐시를 무시하고 다시 조회
upbit candle KRW-BTC --from 2025-01-01 --no-cache

# 캐시 정보 확인
upbit cache

# 캐시 삭제
upbit cache --clear
```

## 라이선스

MIT
