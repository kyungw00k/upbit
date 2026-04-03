package trader

import (
	"context"
	"errors"
	"testing"
	"time"
	"github.com/kyungw00k/upbit-trader/trader/types"
)

func TestMockAgentID(t *testing.T) {
	agent := NewMockAgent("test-agent")
	if agent.ID() != "test-agent" {
		t.Errorf("expected ID 'test-agent', got %s", agent.ID())
	}
}

func TestMockAgentReturnsReport(t *testing.T) {
	expected := types.AnalysisReport{
		AgentID: "test", Timestamp: time.Now().UTC(),
		Signal: types.SignalBuy, Summary: "Test",
		Coins: []types.CoinAnalysis{{Symbol: "KRW-BTC", Signal: types.SignalStrongBuy, Confidence: 0.9}},
	}
	agent := NewMockAgent("test")
	agent.Report = expected
	report, err := agent.Run(context.Background(), types.AgentInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.AgentID != expected.AgentID || report.Signal != expected.Signal {
		t.Errorf("mismatch")
	}
}

func TestMockAgentReturnsError(t *testing.T) {
	agent := NewMockAgent("failing")
	agent.Err = errors.New("agent failure")
	_, err := agent.Run(context.Background(), types.AgentInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMockAgentCustomRunFunc(t *testing.T) {
	agent := NewMockAgent("custom")
	agent.RunFunc = func(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
		return types.AnalysisReport{AgentID: "custom", Signal: types.SignalHold}, nil
	}
	report, err := agent.Run(context.Background(), types.AgentInput{})
	if err != nil || report.Signal != types.SignalHold {
		t.Error("RunFunc not working")
	}
}

func TestAgentInterfaceSatisfaction(t *testing.T) {
	var _ Agent = (*MockAgent)(nil)
}
