# Upbit Autonomous Trading Daemon — Design Spec

**Date:** 2026-04-03
**Status:** Approved
**Architecture:** Multi-Agent Expert Pool (Anthropic Go SDK)

---

## 1. Overview

Anthropic Go SDK를 사용하여 Upbit 거래소에서 완전 자율적으로 암호화폐를 거래하는 데몬 시스템. 증권회사 패러다임의 전문가 에이전트 풀 구조로, 각 에이전트가 독립적인 전문 분야를 담당하고 CEO 에이전트가 최종 의사결정을 내린다.

### 핵심 특성

- **완전 자율:** 에이전트가 코인 탐색부터 주문 실행까지 모든 결정을 단독으로 수행
- **스윙 트레이딩:** 하루 1~5회 거래, 4시간 주기 의사결정 (이벤트 기반 임시 주기 가능)
- **빠른 턴어라운드:** 특정 종목 장기 보유 없이 기회 포착 → 진입 → 청산 사이클
- **Self-Review 루프:** 거래 후 자가 리뷰를 통해 학습, 향후 의사결정 품질 향상
- **모니터링 3채널:** TUI 대시보드 + 메신저 알림 + 웹 대시보드

### 거래 전략

- 기술적 분석 (RSI, MACD, 이동평균, 볼린저밴드)
- 뉴스/감성 분석 (긍정/부정/중립 감성 점수)
- 시장 상관관계 (BTC/ETH 주도성, 알트코인 연동성)
- 온체인/심리 지표 (시장 심리, 거래량 이상 탐지)

---

## 2. Agent Team Architecture

```
                 ┌─────────────────────┐
                 │     CEO Agent       │
                 │  (Portfolio CIO)    │
                 │  모델: Claude Sonnet│
                 └─────────┬───────────┘
                           │
        ┌──────────┬───────┼───────┬──────────┐
        ▼          ▼       ▼       ▼          ▼
   ┌─────────┐┌────────┐┌─────┐┌────────┐┌────────┐
   │ Market  ││ News   ││Tech ││ Risk   ││Exec    │
   │ Analyst ││Analyst ││Anal.││Manager ││Trader  │
   │ Sonnet  ││Haiku   ││Haiku││Sonnet  ││Sonnet  │
   └─────────┘└────────┘└─────┘└────────┘└────────┘
                                                      ▼
                                               ┌────────────┐
                                               │ Self       │
                                               │ Reviewer   │
                                               │ Sonnet     │
                                               └────────────┘
```

### 에이전트 역할 정의

| 에이전트 | 역할 | 모델 | 호출 빈도 |
|---------|------|------|----------|
| **CEO** | 최종 거래 결정, 자본 배분, 코인 탐색, 전략 조율 | Sonnet | 매 의사결정 주기 |
| **Market Analyst** | 시장 상황, 상관관계, 거래량/시가총액 변동, 돈 흐름 분석 | Sonnet | 매 주기 |
| **News Analyst** | 뉴스 수집, 감성 점수 산정, 주요 이벤트 분류 | Haiku | 1시간마다 |
| **Technical Analyst** | RSI, MACD, 이동평균 등 기술 지표 계산, 매수/매도 시그널 | Haiku | 매 주기 |
| **Risk Manager** | 포지션 사이징, 손절/익절 산정, 포트폴리오 상관관계, 거부권 행사 | Sonnet | 매 주기 |
| **Execution Trader** | 주문 실행, 시장가/지정가 자동 선택, 슬리피지 최적화 | Sonnet | 주문 시 |
| **Self-Reviewer** | 거래 결과 기록, 예상 vs 실제 비교, 개선점 도출, 장기 학습 | Sonnet | 거래 완료 후 |

---

## 3. Decision Pipeline

의사결정 주기당 6단계 파이프라인:

```
[1. 수집] ──병렬──▶ [2. 분석] ──병렬──▶ [3. 종합] ──▶ [4. 리스크] ──▶ [5. 실행] ──▶ [6. 리뷰]
```

### Phase 1: 데이터 수집 (병렬)

- **Market Analyst** → Upbit WebSocket/API: 거래량, 호가창, 시가총액 변동, 시장 전체 틱 데이터
- **News Analyst** → 뉴스 API/크롤링: 최신 코인 뉴스, 소셜 미디어 언급량
- **Technical Analyst** → OHLCV 캔들 데이터: 다중 타임프레임 (1분, 1시간, 4시간, 일봉)

### Phase 2: 분석 (병렬)

- **Market Analyst** → 시장 상관관계 맵, 이상 거래량 탐지, 시장 센티먼트
- **News Analyst** → 감성 점수 (각 코인별 -1.0 ~ +1.0), 주요 이벤트 분류, FUD/FOMO 감지
- **Technical Analyst** → 기술 지표 계산 결과, 매수/매도 시그널, 지지/저항선

### Phase 3: 종합 (CEO)

- 세 분석가 리포트 취합
- 신규 코인 기회 탐색 (시장 전체 스캔)
- 매수 후보군 + 매도 후보군 우선순위 결정
- 자본 배분 제안 (각 코인별 비중)

### Phase 4: 리스크 검토 (Risk Manager)

- 각 거래의 포지션 사이징 (고정 비율 또는 Kelly Criterion 변형)
- 손절/익절 가격 산정
- 포트폴리오 상관관계 검토 (상관성 높은 코인 과다 집중 방지)
- 최대 일일 손실 한도 체크 → 초과 시 거래 거부
- **독립 거부권:** CEO 결정이라도 리스크 기준 미충족 시 거부 가능

### Phase 5: 실행 (Execution Trader)

- 리스크 승인 후 주문 실행
- 시장가/지정가 자동 선택 (거래량과 슬리피지 기반)
- 분할 주문 옵션 (대규모 포지션 시)
- 실행 결과 보고 (체결가, 수수료, 슬리피지)

### Phase 6: 기록 & Self-Review

- 매매 로그 SQLite 저장 (진입/청산 가격, 사유, 에이전트 분석 요약)
- CEO Self-Review: 예상 vs 실제 결과 비교, 의사결정 오류 분석
- 개선점 도출 → 프롬프트 컨텍스트 업데이트
- 장기 학습: 주간/월간 성과 분석, 효과적 전략/비효율 코인 식별

---

## 4. Project Structure

별도 Go 프로젝트로 생성. 기존 `upbit/` 프로젝트를 Go 모듈 의존성으로 사용.

```
Sandbox/
├── upbit/                     # 기존 Upbit CLI 프로젝트 (github.com/kyungw00k/upbit)
├── upbit-trader/              # [신규] 자율 거래 데몬 (github.com/kyungw00k/upbit-trader)
│   ├── go.mod                 # module github.com/kyungw00k/upbit-trader
│   ├── cmd/
│   │   └── trader/
│   │       └── main.go       # 엔트리포인트
│   │
│   ├── trader/
│   │   ├── daemon.go         # 데몬 메인 루프, 스케줄러
│   │   ├── pipeline.go       # 의사결정 파이프라인 오케스트레이션
│   │   ├── agent.go          # 에이전트 인터페이스 정의
│   │   │
│   │   ├── agents/           # 에이전트 구현
│   │   │   ├── ceo.go        # CEO Agent
│   │   │   ├── market.go     # Market Analyst
│   │   │   ├── news.go       # News Analyst
│   │   │   ├── technical.go  # Technical Analyst
│   │   │   ├── risk.go       # Risk Manager
│   │   │   ├── execution.go  # Execution Trader
│   │   │   └── reviewer.go   # Self-Reviewer
│   │   │
│   │   ├── memory/           # 학습 메모리 시스템
│   │   │   ├── store.go      # SQLite 기반 거래 기록
│   │   │   ├── reviewer.go   # Self-review 로직
│   │   │   └── context.go    # 프롬프트 컨텍스트 생성
│   │   │
│   │   ├── monitor/          # 모니터링
│   │   │   ├── tui.go        # Bubble Tea TUI 대시보드
│   │   │   ├── webhook.go    # Telegram/Slack 알림
│   │   │   └── api.go        # 웹 대시보드 HTTP 서버
│   │   │
│   │   └── config/
│   │       └── config.go     # 데몬 설정 (리스크 파라미터, 에이전트 설정)
│   │
│   └── configs/
│       └── config.yaml       # 사용자 설정 파일
```

### 핵심 의존성

- `github.com/anthropics/anthropic-sdk-go` — Claude API
- `modernc.org/sqlite` — SQLite (메모리/거래 기록)
- `github.com/charmbracelet/bubbletea` — TUI 대시보드
- `go.uber.org/zap` — 구조화 로깅

### 기존 Upbit API 활용

`go get github.com/kyungw00k/upbit`으로 모듈 의존성 추가 후 임포트:
- `github.com/kyungw00k/upbit/api/quotation` — 시장 데이터 (티커, 캔들, 호가창, 체결)
- `github.com/kyungw00k/upbit/api/exchange` — 주문 (매수, 매도, 취소)
- `github.com/kyungw00k/upbit/api/wallet` — 계좌 잔고
- `github.com/kyungw00k/upbit/api/websocket` — 실시간 스트리밍
- `github.com/kyungw00k/upbit/ratelimit` — API 호출 제한
- `github.com/kyungw00k/upbit/retry` — 재시도 로직
- `github.com/kyungw00k/upbit/types` — 데이터 타입

---

## 5. Agent Communication

### 인터페이스 정의

```go
// Agent 는 모든 에이전트의 기본 인터페이스
type Agent interface {
    ID() string
    Run(ctx context.Context, input AgentInput) (AgentOutput, error)
}

// AgentInput 은 에이전트에 전달되는 입력
type AgentInput struct {
    // 동적 데이터 (시장, 뉴스, 기술 지표)
    MarketData   *MarketSnapshot
    NewsData     *NewsDigest
    TechData     *TechnicalDigest
    // 이전 단계 에이전트들의 분석 결과
    Reports      []AnalysisReport
    // 현재 포트폴리오 상태
    Portfolio    *PortfolioState
    // 과거 학습 컨텍스트
    Memory       *MemoryContext
}

// AnalysisReport 는 에이전트 분석 결과
type AnalysisReport struct {
    AgentID   string        `json:"agent_id"`
    Timestamp time.Time     `json:"timestamp"`
    Signal    string        `json:"signal"` // BUY | SELL | HOLD | NO_ACTION
    Coins     []CoinAnalysis `json:"coins"`
    Summary   string        `json:"summary"`
}

// CoinAnalysis 는 코인별 분석 결과
type CoinAnalysis struct {
    Symbol     string  `json:"symbol"`     // "KRW-BTC"
    Signal     string  `json:"signal"`     // STRONG_BUY | BUY | HOLD | SELL | STRONG_SELL
    Confidence float64 `json:"confidence"` // 0.0 ~ 1.0
    EntryPrice float64 `json:"entry_price"`
    StopLoss   float64 `json:"stop_loss"`
    TakeProfit float64 `json:"take_profit"`
    Position   float64 `json:"position"`   // 제안 포지션 크기 (KRW)
    Reasoning  string  `json:"reasoning"`
}
```

### Claude API 프롬프트 구조

각 에이전트는 고정 시스템 프롬프트 + 동적 사용자 메시지 구조:

- **System Prompt:** 역할, 분석 기준, 출력 형식 (JSON 스키마) — 에이전트별 고정
- **User Message:** 실시간 데이터 (JSON 형태로 직렬화된 MarketSnapshot 등) — 매 주기마다 갱신
- **Output:** 구조화 JSON (`AnalysisReport`) — Go에서 파싱하여 다음 단계로 전달

---

## 6. Risk Management

### 하드 리밋 (코드 레벨 강제 — 에이전트가 초과 불가)

| 파라미터 | 기본값 | 설명 |
|---------|--------|------|
| max_daily_loss_pct | 5% | 초과 시 당일 자동 거래 중지 |
| max_position_pct | 10% | 단일 코인 최대 포지션 비율 (총 자산 대비) |
| max_simultaneous_positions | 5 | 동시 보유 코인 수 제한 |
| max_leverage | 1x | 현물만 (선물/레버리지 금지) |
| min_trade_interval | 1시간 | 동일 코인 최소 거래 간격 |
| emergency_keywords | ["flash crash", "hack", "exploit"] | 뉴스 감지 시 전 포지션 청산 |

### 소프트 가이드 (에이전트가 준수하되 예외 가능)

- 추천 손절가: -3% (Risk Manager가 시장 변동성에 따라 조정)
- 추천 익절가: +5~15% (기술적 분석 기반)
- 신규 코인 첫 거래: 최대 자산 1% 소액 테스트

### Risk Manager 독립 거부권

Risk Manager는 CEO의 결정을 독립적으로 검토하며, 다음 기준을 충족하지 못하면 거부:
- 일일 손실 한도 초과
- 포지션 과다 집중
- 포트폴리오 상관관계 위험
- 시장 변동성 과도 (ATR 기준)

### 비상 정지 (Kill Switch)

```
CLI:  trader stop --emergency
API:  POST /api/stop?type=emergency
TUI:  Ctrl+E 단축키
```

동작: 모든 미체결 주문 취소 → 당일 거래 중지 → 모든 채널에 알림 발송

---

## 7. Monitoring

### TUI 대시보드 (Bubble Tea)

- 실시간 포트폴리오 현황 (보유 코인, 평가손익)
- 현재 의사결정 주기 상태 (Phase 진행률)
- 최근 거래 내역
- 에이전트 활동 로그
- 단축키: Ctrl+E (긴급 정지), Ctrl+L (거래 로그), Ctrl+P (포트폴리오)

### 메신저 알림 (Telegram/Slack)

- 거래 체결 알림 (매수/매도)
- 일일 성과 리포트 (매일 자정)
- 긴급 상황 알림 (비상 정지, 손실 한도 근접)
- Self-Review 요약 (주간)

### 웹 대시보드 (HTTP API + SPA)

에이전트가 왜 그 결정을 내렸는지 투명하게 추적할 수 있는 대시보드. 각 의사결정 주기마다 전체 파이프라인의 사고 과정을 시각화한다.

**API 엔드포인트:**

| 엔드포인트 | 설명 |
|---------|------|
| `/api/status` | 데몬 상태, 현재 포지션, 최근 결정 주기 요약 |
| `/api/trades` | 거래 내역 (페이지네이션, 필터링) |
| `/api/performance` | 성과 지표 (수익률, 샤프비율, 최대 낙폭, 일별 PnL) |
| `/api/agents` | 에이전트 목록, 각 에이전트별 최근 활동 요약 |
| `/api/agents/{id}/reports` | 특정 에이전트의 전체 분석 리포트 히스토리 |
| `/api/decision-cycles` | 의사결정 주기 목록 (페이지네이션) |
| `/api/decision-cycles/{id}` | 특정 주기의 전체 파이프라인 상세 |
| `/api/decision-cycles/{id}/flow` | 주기별 데이터 흐름 시각화용 JSON |

**의사결정 주기 상세 (`/api/decision-cycles/{id}`):**
```json
{
  "id": "dc-20260404-080000",
  "started_at": "2026-04-04T08:00:00Z",
  "duration_ms": 12500,
  "phase": "completed",
  "agents": [
    {
      "id": "technical-analyst",
      "model": "claude-3-5-haiku-20241022",
      "latency_ms": 3200,
      "input_summary": "148 KRW markets, top 20 candle data",
      "output": {
        "signal": "BUY",
        "summary": "BTC RSI 과매도 반등, MACD 골든크로스 확",
        "coins": [
          {
            "symbol": "KRW-BTC",
            "signal": "STRONG_BUY",
            "confidence": 0.85,
            "reasoning": "RSI 35 과매도 구간, MACD 라인이 시그널 라인을 상향 돌파, 거래량 전일 대비 180% 증가"
          }
        ]
      }
    },
    {
      "id": "market-analyst",
      "model": "claude-sonnet-4-20250514",
      "latency_ms": 4500,
      "input_summary": "BTC 주도성 강화, 알트코인 연동 상승",
      "output": { ... }
    },
    {
      "id": "ceo",
      "model": "claude-sonnet-4-20250514",
      "latency_ms": 2800,
      "input_summary": "3개 분석가 리포트 종합",
      "output": {
        "decision": "BUY KRW-BTC",
        "reasoning": "Technical: RSI 과매도 반등 시그널 강함. Market: BTC 주도성 확인. News: 특별 이벤트 없음. Risk: 500KRW 소액 진입으로 리스크 한도 내.",
        "considered_alternatives": "ETH는 신호 약함(홀드), SOL은 변동성 과다"
      }
    },
    {
      "id": "risk-manager",
      "model": "claude-sonnet-4-20250514",
      "latency_ms": 800,
      "output": {
        "approved": true,
        "position_size": 500000,
        "stop_loss": 72000000,
        "take_profit": 82000000,
        "risk_ratio": "0.8%",
        "reasoning": "포지션 10% 미만, 일일 손실 한도 여유 충분"
      }
    },
    {
      "id": "execution-trader",
      "model": "claude-sonnet-4-20250514",
      "latency_ms": 500,
      "output": {
        "action": "BUY",
        "order_type": "limit",
        "price": 75000000,
        "filled": true,
        "execution_report": { ... }
      }
    }
  ],
  "final_outcome": "KRW-BTC 매수 주문 체결 완료"
}
```

**대시보드 UI 주요 화면:**
- **Overview**: 현재 포지션, 일별 PnL, 최근 결정 요약
- **Decision Timeline**: 주기별 시각 타임라인, 각 주기의 최종 결정과 결과
- **Agent Detail**: 선택한 에이전트의 입력 데이터 요약 → 추론(reasoning) → 출력 시각화
- **Trade History**: 체결 내역, 예상 vs 실제 성과 비교 (self-review 포함)

**기술 스택:** Go `net/http` API 서버 + 정적 HTML/CSS/JS (Vanilla 또는 가벤은 프론트엔드 프레임워크). 프론트엠드 빌드 없이 서빙 가능하도록 임베디드.

- 포트: 기본 8080 (설정 변경 가능)

---

## 8. Decision Audit Trail

각 의사결정 주기의 전체 데이터를 SQLite에 영구 저장하여 추후 분석 가능.

### 저장 데이터

| 테이블 | 내용 |
|-------|------|
| `decision_cycles` | 주기 ID, 시작/종료 시간, 최종 결과, 에러 |
| `agent_reports` | 주기 ID, 에이전트 ID, 모델, 입력 요약, 전체 출력 JSON, 지연시간 |
| `trades` | 주기 ID, 마켓, 사이드, 가격, 수량, 수수료, 상태 |
| `self_reviews` | 거래 ID, 예상 vs 실제, 교훈, 프롬프트 업데이트 내용 |

### 활용

- 에이전트 성과 비교 (동일 조건에서 어떤 에이전트가 더 나은 예측을 했는지)
- 잘못된 결정의 원인 분석 (어떤 에이전트의 분석이 오류였는지)
- Self-Review가 프롬프트 개선에 어떤 영향을 주었는지 측정
- 백테스팅을 위한 과거 의사결정 데이터 재생

---

## 8. Configuration

```yaml
# configs/config.yaml
daemon:
  decision_interval: 4h        # 의사결정 주기
  log_level: info

api:
  anthropic_api_key: ${ANTHROPIC_API_KEY}
  model_ceo: claude-sonnet-4-20250514
  model_analyst: claude-haiku-4-20251001
  model_risk: claude-sonnet-4-20250514
  model_execution: claude-sonnet-4-20250514

risk:
  max_daily_loss_pct: 5
  max_position_pct: 10
  max_simultaneous_positions: 5
  min_trade_interval: 1h
  emergency_keywords: ["flash crash", "hack", "exploit"]

monitoring:
  tui_enabled: true
  telegram_bot_token: ${TELEGRAM_BOT_TOKEN}
  telegram_chat_id: ${TELEGRAM_CHAT_ID}
  web_port: 8080

news:
  sources: ["coindesk", "cointelegraph", "decrypt"]
  refresh_interval: 1h

memory:
  db_path: ./data/memory.db
```

---

## 9. API Cost Estimate

| 에이전트 | 모델 | 호출/일 | 예상 토큰/호출 | 일일 비용 |
|---------|------|---------|---------------|----------|
| CEO | Sonnet | 6 | ~4,000 | ~$0.24 |
| Market Analyst | Sonnet | 6 | ~3,000 | ~$0.15 |
| News Analyst | Haiku | 24 | ~2,000 | ~$0.03 |
| Technical Analyst | Haiku | 6 | ~3,000 | ~$0.02 |
| Risk Manager | Sonnet | 6 | ~2,000 | ~$0.12 |
| Execution Trader | Sonnet | ~5 | ~1,000 | ~$0.05 |
| Self-Reviewer | Sonnet | ~5 | ~2,500 | ~$0.12 |
| **합계** | | | | **~$0.73/일** |

*4시간 주기 기준, 뉴스는 1시간 주기, Claude API 2025년 요금 기준*

---

## 10. CLI Commands

```
trader start              # 데몬 시작 (기본: 포그라운드)
trader start -d           # 데몬 모드 (백그라운드)
trader stop               # 정상 중지
trader stop --emergency   # 비상 정지
trader status             # 상태 조회
trader trades             # 거래 내역
trader performance        # 성과 지표
trader review             # 최근 Self-Review 요약
trader config             # 설정 확인/수정
```
