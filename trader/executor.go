package trader

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kyungw00k/upbit/api/exchange"
)

// TradeRequest represents a request to execute a trade.
type TradeRequest struct {
	Market string
	Side   string  // "bid" (buy) or "ask" (sell)
	Price  float64
	Volume float64
	Mode   string  // "paper" or "real"
}

// TradeResult holds the outcome of an executed trade.
type TradeResult struct {
	OrderUUID   string
	FilledPrice float64
	Volume      float64
	Status      string
	Fee         float64
}

// Executor defines the interface for executing trades.
type Executor interface {
	Execute(ctx context.Context, req TradeRequest) (*TradeResult, error)
}

// PaperExecutor simulates trades without hitting real markets.
type PaperExecutor struct{}

// NewPaperExecutor creates a new PaperExecutor.
func NewPaperExecutor() *PaperExecutor {
	return &PaperExecutor{}
}

// Execute simulates a trade at the current price with a 0.05% fee.
func (p *PaperExecutor) Execute(ctx context.Context, req TradeRequest) (*TradeResult, error) {
	if err := validateTradeRequest(req); err != nil {
		return nil, err
	}

	orderUUID := uuid.New().String()
	filledPrice := req.Price
	volume := req.Volume
	fee := req.Price * req.Volume * DefaultFeeRate

	status := "closed"
	if req.Mode != "" {
		status = "closed"
	}

	return &TradeResult{
		OrderUUID:   orderUUID,
		FilledPrice: filledPrice,
		Volume:      volume,
		Status:      status,
		Fee:         fee,
	}, nil
}

// RealExecutor executes real trades via the Upbit exchange API.
type RealExecutor struct {
	client *exchange.ExchangeClient
}

// NewRealExecutor creates a new RealExecutor with the given exchange client.
func NewRealExecutor(client *exchange.ExchangeClient) *RealExecutor {
	return &RealExecutor{client: client}
}

// Execute places a real order via the Upbit API.
func (r *RealExecutor) Execute(ctx context.Context, req TradeRequest) (*TradeResult, error) {
	if err := validateTradeRequest(req); err != nil {
		return nil, err
	}
	if r.client == nil {
		return nil, fmt.Errorf("real executor: exchange client is not configured")
	}

	// Determine order type: for market orders use "price" for bid, "market" for ask.
	ordType := "price"
	if req.Side == "ask" {
		ordType = "market"
	}

	orderReq := &exchange.OrderRequest{
		Market:  req.Market,
		Side:    req.Side,
		OrdType: ordType,
		Price:   fmt.Sprintf("%.8f", req.Price*req.Volume),
		Volume:  fmt.Sprintf("%.8f", req.Volume),
	}

	order, err := r.client.CreateOrder(ctx, orderReq)
	if err != nil {
		return nil, fmt.Errorf("real executor: create order: %w", err)
	}

	fee := req.Price * req.Volume * DefaultFeeRate

	return &TradeResult{
		OrderUUID:   order.UUID,
		FilledPrice: req.Price,
		Volume:      req.Volume,
		Status:      order.State,
		Fee:         fee,
	}, nil
}

// validateTradeRequest checks that the trade request has valid fields.
func validateTradeRequest(req TradeRequest) error {
	if req.Market == "" {
		return fmt.Errorf("trade request: market is required")
	}
	if req.Side != "bid" && req.Side != "ask" {
		return fmt.Errorf("trade request: side must be 'bid' or 'ask', got %q", req.Side)
	}
	if req.Price <= 0 {
		return fmt.Errorf("trade request: price must be positive")
	}
	if req.Volume <= 0 {
		return fmt.Errorf("trade request: volume must be positive")
	}
	return nil
}
