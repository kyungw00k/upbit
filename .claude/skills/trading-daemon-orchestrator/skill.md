---
name: trading-daemon-orchestrator
description: "트레이딩 데몬 전체 개발 오케스트레이터. Phase별 에이전트 구성, 병렬/순차 워크플로우, 데이터 흐름 관리."
---

# Trading Daemon Orchestrator

## 현재 상태
- Phase 0 완료: types, agent interface, config, 테스트 통과
- Phase 1 완료: PocketBase 임베딩, 6개 Collection, migrations
- Phase 2~6: 에이전트 팀이 subagent로 분기 실행

## 에이전트 구성표

| 에이전트 | Phase | 병렬 그룹 | 담당 스킬 | 산출물 |
|---------|-------|----------|----------|--------|
| indicator-engineer | 2 | A (병렬) | indicator-engineer | trader/indicator/*.go, trader/collector.go |
| sentinel-collector | 2.5 | A (병렬) | sentinel-collector | trader/sentiment/*.go |
| glm-agent-builder | 3 | B (Phase 2 후) | glm-agent-builder | trader/agent/glm.go, trader/agents/technical.go |
| pipeline-architect | 4-5 | C (Phase 3 후) | pipeline-architect | trader/pipeline.go, trader/risk.go, trader/executor.go, trader/backtest/ |

## 워크플로우

```
Phase 2 (병렬 그룹 A — subagent 동시 분기):
  ├── indicator-engineer ──→ trader/indicator/*.go, trader/collector.go
  └── sentinel-collector ──→ trader/sentiment/*.go

Phase 3 (순차, 그룹 A 완료 후):
  └── glm-agent-builder ──→ trader/agent/glm.go, trader/agent/prompt.go, trader/agents/technical.go

Phase 4-5 (순차, Phase 3 완료 후):
  └── pipeline-architect ──→ trader/pipeline.go, trader/risk.go, trader/executor.go, trader/backtest/
```

## Subagent 분기 방법

### Phase 2 (병렬)
Agent 툴로 indicator-engineer와 sentinel-collector를 동시에 실행:

```
Agent(subagent_type="general-purpose", name="indicator-engineer", prompt="...")
Agent(subagent_type="general-purpose", name="sentinel-collector", prompt="...")
```

각 에이전트 프롬프트에 해당 스킬의 핵심 내용만 포함 (context 절약).

### Phase 3 (순차)
Phase 2 완료 후 glm-agent-builder를 subagent로 실행.

### Phase 4-5 (순차)
Phase 3 완료 후 pipeline-architect를 subagent로 실행.

## 데이터 흐름

```
Upbit API (candles)
    ↓
indicator-engineer: CalculateAll() → Result{RSI, MACD, BB, EMA, ATR}
    ↓ PocketBase indicators Collection

Alternative.me API (fear & greed)
    ↓
sentinel-collector: FetchFearGreedIndex() → FearGreedData
    ↓ PocketBase sentiment Collection

지표 + 센티먼트
    ↓
glm-agent-builder: BuildAnalysisPrompt() → AnalysisReport
    ↓ PocketBase signals Collection

시그널 + 리스크 평가
    ↓
pipeline-architect: RiskGate → Executor → trades Collection
```

## 검증 체크포인트

| Phase | 명령어 |
|-------|--------|
| 2 | `go test ./trader/indicator/... -v` |
| 2.5 | `go test ./trader/sentiment/... -v` |
| 3 | `go test ./trader/agent/... ./trader/agents/... -v` |
| 4 | `go test ./trader/... -v` |
| 전체 | `go build ./...` |
