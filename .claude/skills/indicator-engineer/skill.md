---
name: indicator-engineer
description: "순수 Go 기술 지표 계산 엔진. RSI, MACD, BB, EMA, ATR 계산 + 테스트 + collector. indicator-engineer 에이전트 전용."
---

# Indicator Engineer Skill

## 사전 확인
1. `trader/indicator/` 디렉토리 존재 확인 (없으면 생성)
2. `trader/types/types.go`의 `AgentInput` 타입 확인
3. `api/quotation/candle.go`의 `GetCandlesAll` 시그니처 확인
4. PocketBase: `app.DB()` 사용

## 파일 구성

| 파일 | 역할 |
|------|------|
| `trader/indicator/math.go` | SMA, StdDev, TrueRange 헬퍼 |
| `trader/indicator/indicator.go` | RSI, EMA, MACD, BB, ATR + CalculateAll |
| `trader/indicator/indicator_test.go` | 각 지표별 단위 테스트 |
| `trader/collector.go` | 캔들 수집 → 지표 계산 → DB 저장 |

## 핵심 타입

```go
type CandleData struct {
    Open, High, Low, Close, Volume float64
    Time string
}
type Result struct {
    RSI, MACDLine, MACDSignal, MACDHist float64
    BBUpper, BBMiddle, BBLower          float64
    EMA20, EMA50, ATR, VolRatio         float64
    Trend                                string // "up", "down", "neutral"
}
```

## 핵심 함수 시그니처
- `SMA(values []float64, period int) []float64`
- `EMA(values []float64, period int) []float64`
- `RSI(closes []float64, period int) []float64` — Wilder's smoothing
- `MACD(closes []float64, fast, slow, signal int) (macdLine, signalLine, histogram []float64)` — 기본: 12, 26, 9
- `BollingerBands(closes []float64, period int, stdDev float64) (upper, middle, lower []float64)` — 기본: 20, 2.0
- `ATR(candles []CandleData, period int) []float64` — TrueRange + Wilder's
- `CalculateAll(candles []CandleData) Result` — Trend: close > EMA20 > EMA50 = "up"
- `CollectAndCalculate(app core.App, qClient, markets, interval)` — 수집→계산→저장

## 테스트 기준
- 하드코딩된 50개 데이터 포인트로 검증
- RSI: 마지막 값이 0~100
- BB: upper > middle > lower
- ATR: 양수
- CalculateAll: 모든 필드 합리적 범위

## 검증
```bash
go test ./trader/indicator/... -v
go build ./...
```

## 주의사항
- `len(closes) < period` 케이스 처리 필수
- 동시성 안전: CalculateAll은 순수 함수
