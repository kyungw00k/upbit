---
name: glm-agent-builder
description: "GLM API 클라이언트 + 프롬프트 엔지니어링 + 기술 분석 에이전트. glm-agent-builder 에이전트 전용."
---

# GLM Agent Builder Skill

## 사전 확인
1. `trader/agent/` 디렉토리 확인 (없으면 생성)
2. `trader/agents/` 디렉토리 확인
3. `trader/agent.go`의 `Agent` 인터페이스 확인
4. `trader/types/types.go`의 `AgentInput`, `AnalysisReport`, `CoinAnalysis`, `SentimentData` 확인
5. 의존성: `github.com/sashabaranov/go-openai`
6. Config: `trader/config/config.go`의 `APIConfig.GLMAPIKey`, `GLMBaseURL`, `GLMModel`

## 파일 구성

| 파일 | 역할 |
|------|------|
| `trader/agent/glm.go` | GLM API 클라이언트 (OpenAI 호환) |
| `trader/agent/prompt.go` | SystemPrompt + BuildAnalysisPrompt |
| `trader/agent/glm_test.go` | httptest mock 테스트 |
| `trader/agents/technical.go` | 순수 Go 지표 기반 시그널 생성 |
| `trader/agents/technical_test.go` | 임계값 테스트 |

## 핵심 타입

```go
type GLMAgent struct {
    id     string
    client *openai.Client
    model  string
}
func NewGLMAgent(apiKey, baseURL, model string) *GLMAgent
```

BaseURL: `"https://api.z.ai/v1"`, 모델: `"glm-4.7"`, `"glm-4.7-flash"` (무료)

## 핵심 함수 시그니처
- `SystemPrompt() string` — 영어, JSON-only 출력, 스키마 포함
- `BuildAnalysisPrompt(market string, indicators map[string]any, sentiment *types.SentimentData) string` — 지표 + 센티먼트 → 프롬프트
- `func (g *GLMAgent) Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error)` — API 호출 → 파싱
- `func (g *GLMAgent) ID() string` → `"glm-analyst"`

## 순수 Go 기술 에이전트
- RSI(30%), MACD(25%), BB(20%), EMA(15%), Volume(10%) 가중치
- 총점 → STRONG_BUY/BUY/HOLD/SELL/STRONG_SELL

## 테스트 기준
- `httptest.NewServer`로 GLM API mock
- 정상 JSON → AnalysisReport 파싱, 빈 응답/잘못된 JSON → 에러
- BuildAnalysisPrompt 출력 포맷 확인
- RSI 80 → SELL, RSI 20 → BUY 확인

## 검증
```bash
go test ./trader/agent/... ./trader/agents/... -v
go build ./...
```

## 주의사항
- GLM API가 JSON 모드에서도 마크다운 펜스를 포함할 수 있음 → 응답 정제 필요
- 프롬프트 토큰: 입력 ~1500, 출력 ~500 예상
