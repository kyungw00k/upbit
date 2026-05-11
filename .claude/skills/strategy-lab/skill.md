---
name: strategy-lab
description: "다중 전략 백테스트 실행 및 비교. SQLite 기반. strategy-lab 에이전트 사용."
---

# Strategy Lab — 전략 백테스트 스킬

## 데이터 로드
```go
db, _ := storage.Open("data/candles.db")
candles, _ := storage.LoadCandles(db, "KRW-BTC", "60", "2024-01-01", "2025-12-31")
```

## 구현할 전략 목록

### 1. RSI-2 Mean Reversion (기준선)
- RSI(2) < 10 매수, RSI(2) > 70 매도
- 트렌드 필터: close > SMA50
- 트레일링: 2×ATR

### 2. Bollinger Bounce
- 매수: close가 BB 하단 아래 이탈 후 되돌아옴 (close > lower)
- 매도: close가 BB 중간선 도달
- 필터: BB 밴드폭 > 20일 평균 (변동성 충분)

### 3. MACD + EMA Cross
- 매수: MACD 히스토그램 음→양 전환 + EMA20 > EMA50
- 매도: MACD 히스토그램 양→음 전환
- 필터: 거래량 > 평균 1.2배

### 4. RSI + BB 결합
- 매수: RSI(14) < 35 AND close < BB 중간선 아래 30%
- 매도: RSI(14) > 65 OR close > BB 상단
- 트레일링: 1.5×ATR

### 5. EMA Cross + ATR Trail
- 매수: EMA20이 EMA50 상향 돌파
- 매도: EMA20이 EMA50 하향 돌파 OR 2×ATR 트레일링
- 포지션: 변동성 역비례 사이징

### 6. BB Squeeze Breakout
- 매수: BB 밴드폭 20일 최저 + close가 상단 돹파 + volume > 1.5×
- 매도: 2×ATR 트레일링
- 타임프레임: 1h

## 백테스트 공통 조건
- 수수료: 0.05% (매수+매도 = 왕복 0.1%)
- 초기자본: 330,000 KRW
- 포지션: 30%
- 평가: 수익률, 승률, MDD, 거래횟수, 월별 일관성

## 출력
1. 전략별 전체 성과 테이블
2. 월별 수익률 히트맵
3. 2024년 vs 2025년 안정성 비교
4. 최종 추천 순위
