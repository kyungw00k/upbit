# Upbit CLI v3 — 구현 계획

## 프로젝트 개요

Go로 작성하는 Upbit 거래소 CLI. 3가지 페르소나를 만족해야 함:
1. **사람 (터미널)**: 짧은 명령, 읽기 좋은 테이블 출력
2. **LLM (도구)**: 예측 가능한 JSON, 명확한 에러 코드
3. **개발자 (스크립트)**: jq 파이프, 종료 코드, --json fields

## 핵심 설계 결정

| 항목 | 결정 |
|------|------|
| 명령 패턴 | 자원 우선 + 스마트 기본동사 (C+) |
| 매수/매도 | `buy`/`sell` 최상위 명령 (→ API side=bid/ask) |
| 실시간 | `watch` (WebSocket 읽기 전용 구독) |
| 출력 | tty=테이블, non-tty=JSON 자동. `--json fields` 지원, `--jq` 미지원 |
| 위험 동작 | 기본 확인 프롬프트, `-f, --force`로 스킵 |
| 저수준 API | 미지원 (편의성 우선) |
| 인프라 코드 | v2 기반 재작성 (검증됨), v1 retry/types 참고 |
| 설정 경로 | `~/.config/upbit/` (XDG 존재 시) → `~/.upbit/` (폴백) |
| 별칭 | 개발 후 검토 |
| TUI | Phase 분리 — MVP는 테이블/JSON, Phase 7에서 Bubble Tea |
| 플래그 단축 | 모든 주요 플래그에 단축 표현 (-p, -V, -t, -c, -i, -m, -q, -o, -f) |

## 명령어 체계

```bash
# === 시세 (기본동사: show, 생략 가능) ===
upbit ticker KRW-BTC                    # 현재가
upbit orderbook KRW-BTC                 # 호가
upbit trades KRW-BTC -c 20             # 체결 내역
upbit candle KRW-BTC -i 1m -c 200     # 캔들
upbit market                            # 마켓 목록
upbit market -q KRW                    # KRW 마켓만

# === 자산 ===
upbit balance                           # 전체 잔고
upbit balance KRW                       # 특정 통화

# === 거래 (buy/sell 최상위) ===
upbit buy KRW-BTC -p 50000000 -V 0.001      # 지정가 매수
upbit sell KRW-BTC -p 60000000 -V 0.001     # 지정가 매도
upbit buy KRW-BTC -t 100000                  # 시장가 매수 (총액)
upbit sell KRW-BTC -V 0.001                  # 시장가 매도 (수량)
upbit buy KRW-BTC -p 50000000 -V 0.001 --test  # 테스트 주문

# === 주문 관리 ===
upbit order list                        # 대기 주문
upbit order list --closed               # 종료 주문
upbit order show <uuid>                 # 개별 조회
upbit order cancel <uuid>               # 취소
upbit order cancel --all -m KRW-BTC    # 일괄 취소
upbit order replace <uuid> -p 51000000 # 취소 후 재주문

# === 입금 ===
upbit deposit list
upbit deposit show <uuid>
upbit deposit address BTC
upbit deposit address create BTC

# === 출금 ===
upbit withdraw list
upbit withdraw request BTC --amount 0.5 --to <addr>
upbit withdraw cancel <uuid>

# === 실시간 (WebSocket 읽기 전용) ===
upbit watch ticker KRW-BTC KRW-ETH
upbit watch orderbook KRW-BTC
upbit watch trade KRW-BTC
upbit watch candle KRW-BTC -i 1m
upbit watch my-order                    # Private (인증 자동)
upbit watch my-asset                    # Private (인증 자동)

# === 유틸 ===
upbit status                            # 서비스 상태
upbit api-keys                          # API 키 목록
upbit config set access-key xxx
upbit config show
```

## 프로젝트 구조

```
upbit-cli-v3/
├── cmd/upbit/
│   └── main.go
├── internal/
│   ├── cli/                         # Cobra 명령어
│   │   ├── root.go                  # 루트, 글로벌 플래그, 카테고리 그룹
│   │   ├── buy.go                   # upbit buy
│   │   ├── sell.go                  # upbit sell
│   │   ├── ticker.go
│   │   ├── orderbook.go
│   │   ├── trades.go
│   │   ├── candle.go
│   │   ├── market.go
│   │   ├── balance.go
│   │   ├── status.go
│   │   ├── apikeys.go
│   │   ├── order/                   # list, show, cancel, replace
│   │   ├── deposit/                 # list, show, address
│   │   ├── withdraw/                # request, cancel, list
│   │   ├── watch/                   # ticker, orderbook, trade, candle, my-order, my-asset
│   │   └── config/                  # set, show
│   ├── api/                         # API 클라이언트
│   │   ├── client.go                # HTTP, Rate Limit 헤더, 에러 파싱
│   │   ├── auth.go                  # JWT HS512 (v2 기반)
│   │   ├── errors.go                # 구조화된 에러
│   │   ├── quotation/               # market, ticker, orderbook, trade, candle
│   │   ├── exchange/                # account, order, order_batch, chance
│   │   ├── wallet/                  # deposit, withdraw, status
│   │   └── websocket/               # client, subscribe, types
│   ├── config/                      # 설정 관리 (XDG)
│   ├── output/                      # formatter, table, json, csv, confirm
│   ├── ratelimit/                   # Token bucket (v2 기반)
│   ├── retry/                       # Exponential backoff (v1 참고)
│   └── types/                       # 공용 모델, Custom JSON 타입
├── docs/upbit-api/                  # API 문서 (103개 md)
├── go.mod
├── go.sum
└── Makefile
```

## Phase별 구현 계획

### Phase 1: 기반 인프라 — cc:완료
> 목표: API 호출 가능한 최소 골격

- [x] 1. Go 모듈 초기화 + 의존성 + Makefile (모듈: github.com/kyungw00k/upbit)
- [x] 2. `internal/config/` — 설정 파일 로드/저장, XDG, 환경변수, 0600 권한
- [x] 3. `internal/types/` — 공용 모델, Custom JSON 타입 (Float64/Timestamp)
- [x] 4. `internal/api/auth.go` — JWT HS512 토큰 (v2 기반)
- [x] 5. `internal/ratelimit/` — Token bucket, Remaining-Req 헤더 (v2 기반)
- [x] 6. `internal/retry/` — Exponential backoff (v1 참고)
- [x] 7. `internal/api/client.go` — HTTP 클라이언트, Rate Limit, 에러 파싱
- [x] 8. `internal/api/errors.go` — 구조화된 에러, 종료 코드 (0/1/2/3/4)
- [x] 9. `internal/output/` — tty 감지, table/json/csv/confirm 포맷터
- [x] 10. `internal/cli/root.go` — 루트 명령, 글로벌 플래그
- [x] 11. `cmd/upbit/main.go` — 엔트리포인트

### Phase 2: 시세 명령어 — cc:완료
> 목표: 인증 불필요 명령 완성, 출력 포맷 검증
> 참고: 실제 API 경로가 Plans.md 초안과 다름 — /market/all, /ticker, /trades/ticks 사용

- [x] 12. market (GET /market/all) — `-q` 필터
- [x] 13. ticker (GET /ticker) — 복수 마켓
- [x] 14. orderbook (GET /orderbooks)
- [x] 15. trades (GET /trades/ticks) — `-c` count
- [x] 16. candle (GET /candles/{unit}) — `-i` interval, `-c` count

### Phase 3: 거래 명령어 — cc:완료
> 목표: 인증 명령, buy/sell, 확인 프롬프트

- [x] 17. balance (GET /accounts) — 통화 필터
- [x] 18. buy + sell (POST /orders) — 주문 유형 자동 판별, 확인 프롬프트
- [x] 19. order chance (GET /orders/chance) — 주문 가능 정보
- [x] 20. order list + show
- [x] 21. order cancel (단건/일괄)
- [x] 22. order replace (cancel_and_new)

### Phase 4: 입출금 명령어 — cc:완료
> 목표: 입금/출금 전체 커버

- [x] 23. deposit list + show + address (+ create)
- [x] 24. withdraw request + cancel + list
- [x] 25. status (GET /status/wallet) — 최상위 `upbit status`로 구현

### Phase 5: Watch (WebSocket) — cc:완료
> 목표: 실시간 데이터 구독. 리뷰 수정: pingDone 재사용, 재연결, ReadDeadline

- [x] 26. WebSocket 클라이언트 (연결, ping/pong, 재연결)
- [x] 27. subscribe 메시지 빌더 + 스트림 타입
- [x] 28. watch 6개 서브커맨드 (ticker, orderbook, trade, candle, my-order, my-asset)

### Phase 6: 유틸 + 마무리 — cc:완료
> 목표: 완성도

- [x] 29. status (Phase 4에서 완료) + api-keys
- [x] 30. config set/show
- [x] 31. tool-schema (LLM용 JSON Schema 덤프)
- [x] 32. 도움말 체계 (Examples, SuggestFor, 카테고리 그룹)
- [x] 33. 전체 빌드/테스트 검증 — go build + go vet 통과, 63개 Go 파일

### Phase 7a: 캔들 캐시 + 자동 페이지네이션 + 호가 보정 — cc:완료
> SQLite 캐시, --from, 자동 페이지네이션, tick-size 조회, buy/sell 호가 자동 보정

- [x] 34. modernc.org/sqlite + jmoiron/sqlx 의존성
- [x] 35. `internal/cache/candle.go` — SQLite 캐시 (PK: market+timeframe+datetime)
- [x] 36. `internal/api/quotation/candle.go` — GetCandlesWithTo, GetCandlesAll
- [x] 37. `internal/cli/candle.go` — --from, --asc/--desc, --no-cache
- [x] 38. 대량 조회 시 예상 횟수 안내 + 확인 프롬프트
- [x] 39. `upbit config cache clear`
- [x] 40. `tick-size` 조회 + buy/sell 호가 자동 보정
- [x] 41. POST JWT 인증 버그 수정 (JSON body 키 순서 보존)
- [x] 42. --force 로컬 플래그 (필요한 명령만)
- [x] 43. SilenceUsage (에러 시 도움말 미출력)

### Phase 7b: Bubble Tea TUI — cc:TODO
> 목표: 사람용 UX 강화

- [ ] 44. watch 명령 TUI 전환
- [ ] 45. 확인 프롬프트 TUI
- [ ] 46. 인터랙티브 모드 (선택적)

### Phase 8: 언어 통일 + 폴리싱 — cc:TODO
> 목표: 컬럼 헤더/메시지 한국어 통일, i18n 기반

- [ ] 47. 테이블 헤더 한/영 통일 (deposit address 등 영어 대문자 → 한글)
- [ ] 48. 에러 메시지 일관성 (Cobra 기본 영어 제거)
- [ ] 49. language 설정 (ko/en) 기반 다국어 지원
- [ ] 50. API 에러 코드 → 영문 메시지 매핑

## 의존성

```
github.com/spf13/cobra          # CLI 프레임워크
github.com/golang-jwt/jwt/v5    # JWT
github.com/google/uuid          # UUID
github.com/gorilla/websocket    # WebSocket
golang.org/x/term               # tty 감지
# Phase 7:
github.com/charmbracelet/bubbletea
github.com/charmbracelet/lipgloss
```

## 참고 자원

- API 문서: `docs/upbit-api/` (103개 마크다운)
- v2 코드 (검증됨): `../upbit-cli-v2/pkg/` — auth, client, ratelimit, websocket
- v1 코드 (참고): `../upbit-cli/internal/` — retry, types, errors

## 하네스 팀

| 에이전트 | 담당 Phase | 역할 |
|---------|-----------|------|
| api-architect | 1 | API 클라이언트, 인증, Rate Limit, 모델, 에러 |
| cli-builder | 2, 3, 4, 6 | Cobra 명령, 출력 포맷터, 도움말, config |
| ws-engineer | 5 | WebSocket 클라이언트, watch 명령 |
| qa-tester | 전체 | 테스트, 통합 검증 |
