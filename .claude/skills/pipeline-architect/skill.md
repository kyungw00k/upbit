---
name: pipeline-architect
description: "의사결정 파이프라인 + Risk Gate + 실행기 + 백테스팅 엔진 구현. pipeline-architect 에이전트 전용. Phase 2/2.5/3 완료 후 착수."
---

# Pipeline Architect Skill

## 사전 확인
1. Phase 2 (indicator), 2.5 (sentiment), 3 (glm-agent) 완료 확인
2. `trader/indicator/indicator.go`의 `Result` 타입 확인
3. `trader/agent/glm.go`의 `GLMAgent` 확인
4. `trader/agent.go`의 `Agent` 인터페이스 확인
5. `api/exchange/order.go`의 `CreateOrder` 시그니처 확인
6. PocketBase Collections: signals, trades, pipeline_runs

## 파일 구성

| 파일 | 역할 |
|------|------|
| `trader/upbit/client.go` | Upbit 클라이언트 래퍼 |
| `trader/risk.go` | Risk Gate (순수 Go) |
| `trader/executor.go` | Paper/Real 실행기 |
| `trader/pipeline.go` | 전체 파이프라인 오케스트레이터 |
| `trader/daemon.go` | Cron 스케줄러 등록 |
| `trader/backtest/backtest.go` | 백테스팅 엔진 |
| `trader/risk_test.go` | Risk Gate 테스트 |
| `trader/executor_test.go` | Executor 테스트 |
| `trader/pipeline_test.go` | Pipeline 통합 테스트 |
| `trader/backtest/backtest_test.go` | 백테스팅 테스트 |

## 핵심 타입

```go
type RiskGate struct { config config.RiskConfig }
type Portfolio struct { TotalValue float64; Positions map[string]Position; DailyPnL float64 }
type TradeRequest struct { Market, Side string; Price, Volume float64; Mode string }
type TradeResult struct { OrderUUID string; FilledPrice float64; Status string }
type Executor interface { Execute(ctx context.Context, req TradeRequest) (*TradeResult, error) }
```

## Risk Gate 체크 항목 (순수 Go, LLM 금지)
1. 단일 코인 ≤ MaxPositionPct (10%)
2. 일일 손실 ≤ MaxDailyLossPct (3%)
3. 동시 포지션 ≤ MaxPositions (5)
4. 손절가가 StopLossPct (5%) 이내
5. 위반 시: "rejected: {reason}" 반환

## Pipeline Run 흐름
1. pipeline_runs 생성 (status="running")
2. 지표 조회 (indicators Collection)
3. 센티먼트 조회 (sentiment Collection)
4. 에이전트 실행 → AnalysisReport 수집
5. 시그널 병합 (가중 투표 또는 만장일치)
6. Risk Gate 평가
7. 승인된 거래 실행
8. pipeline_runs 업데이트

## Cron 스케줄
- indicator_calc: `*/15 * * * *`
- sentiment_collector: `0 * * * *`
- analysis_pipeline: 설정된 decision_interval

## 백테스팅 지표
TotalReturn, SharpeRatio, MaxDrawdown, WinRate, TotalTrades

## 검증
```bash
go test ./trader/... -v
go build ./...
```

## 주의사항
- Risk Gate는 절대 LLM에 의존하지 않음
- RealExecutor는 API 키 설정 시만 사용
- 파이프라인 실패 시 부분 실행 거래 추적 필수
