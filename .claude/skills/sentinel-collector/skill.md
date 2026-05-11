---
name: sentinel-collector
description: "Fear & Greed Index API 연동. 센티먼트 데이터 수집/저장/테스트. sentinel-collector 에이전트 전용."
---

# Sentinel Collector Skill

## 사전 확인
1. `trader/sentiment/` 디렉토리 존재 확인 (없으면 생성)
2. `trader/types/types.go`의 `SentimentData` 타입 확인
3. PocketBase sentiment Collection: source, value, classification, timestamp 필드

## 파일 구성

| 파일 | 역할 |
|------|------|
| `trader/sentiment/feargreed.go` | API 클라이언트 + 데이터 모델 |
| `trader/sentiment/feargreed_test.go` | httptest mock 테스트 |

## 핵심 타입

```go
type FearGreedData struct {
    Value          int    `json:"value"`
    Classification string `json:"value_classification"`
    Timestamp      string `json:"timestamp"`
}
```

API 응답: `{"data": [{"value": "45", "value_classification": "Fear", "timestamp": "1746489600"}]}`
- `value`가 string으로 오므로 `strconv.Atoi` 변환 필요

## 핵심 함수 시그니처
- `FetchFearGreedIndex(ctx context.Context, limit int) ([]FearGreedData, error)` — URL: `https://api.alternative.me/fng/?limit={limit}&format=json`
- `FetchLatestFearGreedIndex(ctx context.Context) (*FearGreedData, error)` — limit=1
- `SaveSentiment(app core.App, data FearGreedData) error` — 중복 체크 후 저장

## 테스트 기준
- `httptest.NewServer`로 mock API 구성
- 정상 응답 파싱, 빈 데이터, 잘못된 value, HTTP 500 에러 케이스

## 검증
```bash
go test ./trader/sentiment/... -v
go build ./...
```

## 주의사항
- API rate limit: 분당 ~30회 (과도한 호출 금지)
- Value 범위: 0~100, Classification: "Extreme Fear" ~ "Extreme Greed"
