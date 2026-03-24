package websocket

// TickerStream is a WebSocket ticker (current price) stream message.
type TickerStream struct {
	Type              string  `json:"type"`
	Code              string  `json:"code"`
	OpeningPrice      float64 `json:"opening_price"`
	HighPrice         float64 `json:"high_price"`
	LowPrice          float64 `json:"low_price"`
	TradePrice        float64 `json:"trade_price"`
	PrevClosingPrice  float64 `json:"prev_closing_price"`
	Change            string  `json:"change"`
	ChangePrice       float64 `json:"change_price"`
	SignedChangePrice float64 `json:"signed_change_price"`
	ChangeRate        float64 `json:"change_rate"`
	SignedChangeRate  float64 `json:"signed_change_rate"`
	TradeVolume       float64 `json:"trade_volume"`
	AccTradeVolume    float64 `json:"acc_trade_volume"`
	AccTradeVolume24h float64 `json:"acc_trade_volume_24h"`
	AccTradePrice     float64 `json:"acc_trade_price"`
	AccTradePrice24h  float64 `json:"acc_trade_price_24h"`
	TradeDate         string  `json:"trade_date"`
	TradeTime         string  `json:"trade_time"`
	TradeTimestamp    int64   `json:"trade_timestamp"`
	AskBid            string  `json:"ask_bid"`
	AccAskVolume      float64 `json:"acc_ask_volume"`
	AccBidVolume      float64 `json:"acc_bid_volume"`
	Highest52WeekPrice float64 `json:"highest_52_week_price"`
	Highest52WeekDate  string  `json:"highest_52_week_date"`
	Lowest52WeekPrice  float64 `json:"lowest_52_week_price"`
	Lowest52WeekDate   string  `json:"lowest_52_week_date"`
	MarketState       string  `json:"market_state"`
	Timestamp         int64   `json:"timestamp"`
	StreamType        string  `json:"stream_type"`
}

// TradeStream is a WebSocket trade (executed order) stream message.
type TradeStream struct {
	Type             string  `json:"type"`
	Code             string  `json:"code"`
	TradePrice       float64 `json:"trade_price"`
	TradeVolume      float64 `json:"trade_volume"`
	AskBid           string  `json:"ask_bid"`
	PrevClosingPrice float64 `json:"prev_closing_price"`
	Change           string  `json:"change"`
	ChangePrice      float64 `json:"change_price"`
	TradeDate        string  `json:"trade_date"`
	TradeTime        string  `json:"trade_time"`
	TradeTimestamp   int64   `json:"trade_timestamp"`
	Timestamp        int64   `json:"timestamp"`
	SequentialID     int64   `json:"sequential_id"`
	BestAskPrice     float64 `json:"best_ask_price"`
	BestAskSize      float64 `json:"best_ask_size"`
	BestBidPrice     float64 `json:"best_bid_price"`
	BestBidSize      float64 `json:"best_bid_size"`
	StreamType       string  `json:"stream_type"`
}

// OrderbookStream is a WebSocket orderbook stream message.
type OrderbookStream struct {
	Type           string              `json:"type"`
	Code           string              `json:"code"`
	TotalAskSize   float64             `json:"total_ask_size"`
	TotalBidSize   float64             `json:"total_bid_size"`
	OrderbookUnits []OrderbookUnitItem `json:"orderbook_units"`
	Timestamp      int64               `json:"timestamp"`
	Level          float64             `json:"level"`
	StreamType     string              `json:"stream_type"`
}

// OrderbookUnitItem represents a single orderbook price level.
type OrderbookUnitItem struct {
	AskPrice float64 `json:"ask_price"`
	BidPrice float64 `json:"bid_price"`
	AskSize  float64 `json:"ask_size"`
	BidSize  float64 `json:"bid_size"`
}

// CandleStream is a WebSocket candle stream message.
type CandleStream struct {
	Type                string  `json:"type"`
	Code                string  `json:"code"`
	CandleDateTimeUTC   string  `json:"candle_date_time_utc"`
	CandleDateTimeKST   string  `json:"candle_date_time_kst"`
	OpeningPrice        float64 `json:"opening_price"`
	HighPrice           float64 `json:"high_price"`
	LowPrice            float64 `json:"low_price"`
	TradePrice          float64 `json:"trade_price"`
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	CandleAccTradePrice  float64 `json:"candle_acc_trade_price"`
	Timestamp           int64   `json:"timestamp"`
	StreamType          string  `json:"stream_type"`
}

// MyOrderStream is a WebSocket my-order stream message.
type MyOrderStream struct {
	Type            string   `json:"type"`
	Code            string   `json:"code"`
	UUID            string   `json:"uuid"`
	AskBid          string   `json:"ask_bid"`
	OrderType       string   `json:"order_type"`
	State           string   `json:"state"`
	TradeUUID       string   `json:"trade_uuid"`
	Price           float64  `json:"price"`
	AvgPrice        float64  `json:"avg_price"`
	Volume          float64  `json:"volume"`
	RemainingVolume float64  `json:"remaining_volume"`
	ExecutedVolume  float64  `json:"executed_volume"`
	TradesCount     int      `json:"trades_count"`
	ReservedFee     float64  `json:"reserved_fee"`
	RemainingFee    float64  `json:"remaining_fee"`
	PaidFee         float64  `json:"paid_fee"`
	Locked          float64  `json:"locked"`
	ExecutedFunds   float64  `json:"executed_funds"`
	TimeInForce     *string  `json:"time_in_force"`
	TradeFee        *float64 `json:"trade_fee"`
	IsMaker         *bool    `json:"is_maker"`
	Identifier      string   `json:"identifier"`
	SMPType         string   `json:"smp_type"`
	PreventedVolume float64  `json:"prevented_volume"`
	PreventedLocked float64  `json:"prevented_locked"`
	TradeTimestamp  int64    `json:"trade_timestamp"`
	OrderTimestamp  int64    `json:"order_timestamp"`
	Timestamp       int64    `json:"timestamp"`
	StreamType      string   `json:"stream_type"`
}

// MyAssetStream is a WebSocket my-asset stream message.
type MyAssetStream struct {
	Type           string           `json:"type"`
	AssetUUID      string           `json:"asset_uuid"`
	Assets         []MyAssetItem    `json:"assets"`
	AssetTimestamp int64            `json:"asset_timestamp"`
	Timestamp      int64            `json:"timestamp"`
	StreamType     string           `json:"stream_type"`
}

// MyAssetItem represents a single asset entry in a MyAssetStream.
type MyAssetItem struct {
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Locked   float64 `json:"locked"`
}

// StreamMessage is a minimal message used to identify the stream type.
type StreamMessage struct {
	Type string `json:"type"`
}
