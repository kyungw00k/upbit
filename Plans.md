# Upbit CLI v3 — 구현 계획

## 프로젝트 개요

Go로 작성하는 Upbit 거래소 CLI. AI 에이전트와 사람이 함께 사용하는 도구.
- **사람 (터미널)**: 짧은 명령, 테이블 출력, i18n (ko/en)
- **LLM (도구)**: non-tty 자동 JSON, tool-schema, 구조화된 에러
- **개발자 (스크립트)**: jq 파이프, CSV, 종료 코드

## 핵심 설계 결정

| 항목 | 결정 |
|------|------|
| 명령 패턴 | 자원 우선 + 스마트 기본동사 |
| 매수/매도 | `buy`/`sell` 최상위 명령 |
| 실시간 | `watch` (WebSocket 읽기 전용) |
| 출력 | tty=테이블, non-tty=JSON 자동. runewidth CJK 정렬 |
| 인증 | 환경변수 전용 (파일 저장 없음) |
| 호가 보정 | tick-size API 기반 자동 보정 |
| 캔들 캐시 | SQLite, --from 자동 페이지네이션 |
| i18n | POSIX locale (LC_ALL > LC_MESSAGES > LANG) |
| 릴리스 | go-semantic-release + goreleaser |

## 명령어 체계

```bash
# 시세 (인증 불필요)
upbit ticker KRW-BTC                    # 현재가
upbit ticker -q KRW                     # KRW 전체 시세
upbit orderbook KRW-BTC                 # 호가
upbit trades KRW-BTC -c 20             # 체결
upbit candle KRW-BTC -i 1d --from 2025-01-01  # 캔들 (캐시)
upbit market -q KRW                    # 마켓
upbit tick-size KRW-BTC                # 호가 단위
upbit orderbook-levels KRW-BTC         # 호가 모아보기

# 거래 (인증 필요)
upbit buy KRW-BTC -p 50000000 -V 0.001   # 매수
upbit sell KRW-BTC -p 60000000 -V 0.001  # 매도
upbit balance                             # 잔고 + 평가금액
upbit order list/show/cancel/replace/chance

# 입출금 (인증 필요)
upbit wallet                           # 서비스 상태
upbit deposit list/show/address
upbit withdraw list/show/request/cancel
upbit travelrule vasps/verify-txid/verify-uuid

# 실시간 (WebSocket)
upbit watch ticker/orderbook/trade/candle/my-order/my-asset

# 유틸
upbit api-keys                         # API 키 (만료 필터)
upbit cache / cache --clear            # 캐시 관리
upbit tool-schema [cmd]                # LLM용 JSON Schema
upbit update / update --check          # 자가 업데이트
```

## Phase 완료 현황

### Phase 1~6: 기반 → 시세 → 거래 → 입출금 → WebSocket → 유틸 — 완료
### Phase 7a: 캔들 캐시 + 호가 보정 + 폴리싱 — 완료
### Phase 8: i18n + 테이블 정렬 + 테스트 — 완료
### Phase 9: 배포 + 문서 — 완료

- go-semantic-release + goreleaser (자동 버저닝)
- GitHub Actions (release, sync-upbit-docs)
- README (en/ko) + 데모 GIF
- MIT LICENSE
- 자가 업데이트 (upbit update)

---

### Phase 7b: Bubble Tea TUI — TODO

> 목표: watch 명령의 실시간 TUI, 확인 프롬프트 UI

#### 7b-1. TUI 기반 구조
- [ ] bubbletea + lipgloss 의존성 추가
- [ ] `internal/tui/` 패키지 생성
- [ ] tty 감지 → TUI 모드 자동 전환 (non-tty면 기존 동작)

#### 7b-2. watch ticker TUI
- [ ] 실시간 갱신 테이블 (화면 덮어쓰기, 깜빡임 없음)
- [ ] 복수 마켓 행 표시 (가격, 변동률, 거래량)
- [ ] 변동 시 색상 하이라이트 (상승=초록, 하락=빨강)
- [ ] q/Ctrl+C로 종료

#### 7b-3. watch orderbook TUI
- [ ] 매도/매수 호가 양방향 바 차트
- [ ] 실시간 갱신 (깊이 변화 시각화)
- [ ] 스프레드 중앙 표시

#### 7b-4. watch trade TUI
- [ ] 체결 스트림 스크롤 (최근 N건 유지)
- [ ] 매수=초록, 매도=빨강 색상

#### 7b-5. 확인 프롬프트 TUI
- [ ] buy/sell/cancel/withdraw 확인 시 styled 프롬프트
- [ ] 주문 정보 하이라이트 (마켓, 가격, 수량)
- [ ] Y/N 키보드 입력

#### 7b-6. 인터랙티브 모드 (선택적)
- [ ] `upbit` 인자 없이 → TUI 메뉴
- [ ] 카테고리 탐색 → 명령 선택 → 인자 입력
