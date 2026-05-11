package trader

import (
	"context"

	"github.com/kyungw00k/upbit/trader/types"
)

// Agent is the interface that all trading agents must implement.
type Agent interface {
	ID() string
	Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error)
}

// MockAgent is a test double for Agent.
type MockAgent struct {
	id      string
	Report  types.AnalysisReport
	Err     error
	RunFunc func(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error)
}

// NewMockAgent creates a MockAgent with the given ID.
func NewMockAgent(id string) *MockAgent {
	return &MockAgent{id: id}
}

// ID returns the agent's identifier.
func (m *MockAgent) ID() string {
	return m.id
}

// Run executes the agent's analysis. It uses RunFunc if set,
// otherwise returns Report or Err.
func (m *MockAgent) Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	if m.RunFunc != nil {
		return m.RunFunc(ctx, input)
	}
	if m.Err != nil {
		return types.AnalysisReport{}, m.Err
	}
	return m.Report, nil
}
