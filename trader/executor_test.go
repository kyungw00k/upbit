package trader

import (
	"context"
	"strings"
	"testing"
)

func TestPaperExecutorBuy(t *testing.T) {
	executor := NewPaperExecutor()
	req := TradeRequest{
		Market: "KRW-BTC",
		Side:   "bid",
		Price:  75_000_000,
		Volume: 0.001,
		Mode:   "paper",
	}
	result, err := executor.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderUUID == "" {
		t.Error("expected non-empty order UUID")
	}
	if result.FilledPrice != req.Price {
		t.Errorf("expected filled price %f, got %f", req.Price, result.FilledPrice)
	}
	if result.Volume != req.Volume {
		t.Errorf("expected volume %f, got %f", req.Volume, result.Volume)
	}
	if result.Status != "closed" {
		t.Errorf("expected status 'closed', got %q", result.Status)
	}
	expectedFee := req.Price * req.Volume * DefaultFeeRate
	if result.Fee != expectedFee {
		t.Errorf("expected fee %f, got %f", expectedFee, result.Fee)
	}
}

func TestPaperExecutorSell(t *testing.T) {
	executor := NewPaperExecutor()
	req := TradeRequest{
		Market: "KRW-ETH",
		Side:   "ask",
		Price:  3_000_000,
		Volume: 0.5,
		Mode:   "paper",
	}
	result, err := executor.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.FilledPrice != req.Price {
		t.Errorf("expected filled price %f, got %f", req.Price, result.FilledPrice)
	}
	if result.Fee <= 0 {
		t.Errorf("expected positive fee, got %f", result.Fee)
	}
}

func TestPaperExecutorFeeCalculation(t *testing.T) {
	executor := NewPaperExecutor()
	// 75,000,000 * 0.01 = 750,000 KRW trade value.
	// Fee at 0.05% = 750,000 * 0.0005 = 375 KRW.
	req := TradeRequest{
		Market: "KRW-BTC",
		Side:   "bid",
		Price:  75_000_000,
		Volume: 0.01,
	}
	result, err := executor.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedFee := 375.0
	if result.Fee != expectedFee {
		t.Errorf("expected fee %.2f, got %.2f", expectedFee, result.Fee)
	}
}

func TestPaperExecutorValidationEmptyMarket(t *testing.T) {
	executor := NewPaperExecutor()
	req := TradeRequest{
		Market: "",
		Side:   "bid",
		Price:  75_000_000,
		Volume: 0.001,
	}
	_, err := executor.Execute(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty market")
	}
	if !strings.Contains(err.Error(), "market") {
		t.Errorf("error should mention market: %v", err)
	}
}

func TestPaperExecutorValidationInvalidSide(t *testing.T) {
	executor := NewPaperExecutor()
	req := TradeRequest{
		Market: "KRW-BTC",
		Side:   "invalid",
		Price:  75_000_000,
		Volume: 0.001,
	}
	_, err := executor.Execute(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid side")
	}
}

func TestPaperExecutorValidationZeroPrice(t *testing.T) {
	executor := NewPaperExecutor()
	req := TradeRequest{
		Market: "KRW-BTC",
		Side:   "bid",
		Price:  0,
		Volume: 0.001,
	}
	_, err := executor.Execute(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for zero price")
	}
}

func TestPaperExecutorValidationZeroVolume(t *testing.T) {
	executor := NewPaperExecutor()
	req := TradeRequest{
		Market: "KRW-BTC",
		Side:   "bid",
		Price:  75_000_000,
		Volume: 0,
	}
	_, err := executor.Execute(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for zero volume")
	}
}

func TestRealExecutorNilClient(t *testing.T) {
	executor := NewRealExecutor(nil)
	req := TradeRequest{
		Market: "KRW-BTC",
		Side:   "bid",
		Price:  75_000_000,
		Volume: 0.001,
	}
	_, err := executor.Execute(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for nil exchange client")
	}
	if !strings.Contains(err.Error(), "exchange client") {
		t.Errorf("error should mention exchange client: %v", err)
	}
}
