package types

// Market 마켓(거래쌍) 정보
type Market struct {
	Market      string       `json:"market"`
	KoreanName  string       `json:"korean_name"`
	EnglishName string       `json:"english_name"`
	MarketEvent *MarketEvent `json:"market_event,omitempty"`
}

// MarketEvent 마켓 경보 정보
type MarketEvent struct {
	Warning bool          `json:"warning"`
	Caution *CautionFlags `json:"caution,omitempty"`
}

// CautionFlags 주의 종목 경보 유형
type CautionFlags struct {
	PriceFluctuations        bool `json:"PRICE_FLUCTUATIONS"`
	TradingVolumeSoaring     bool `json:"TRADING_VOLUME_SOARING"`
	DepositAmountSoaring     bool `json:"DEPOSIT_AMOUNT_SOARING"`
	GlobalPriceDifferences   bool `json:"GLOBAL_PRICE_DIFFERENCES"`
	ConcentrationSmallAccounts bool `json:"CONCENTRATION_OF_SMALL_ACCOUNTS"`
}

// Ticker 현재가 정보
type Ticker struct {
	Market             string  `json:"market"`
	TradeDate          string  `json:"trade_date"`
	TradeTime          string  `json:"trade_time"`
	TradeDateKst       string  `json:"trade_date_kst"`
	TradeTimeKst       string  `json:"trade_time_kst"`
	TradeTimestamp     int64   `json:"trade_timestamp"`
	OpeningPrice       float64 `json:"opening_price"`
	HighPrice          float64 `json:"high_price"`
	LowPrice           float64 `json:"low_price"`
	TradePrice         float64 `json:"trade_price"`
	PrevClosingPrice   float64 `json:"prev_closing_price"`
	Change             string  `json:"change"` // EVEN, RISE, FALL
	ChangePrice        float64 `json:"change_price"`
	ChangeRate         float64 `json:"change_rate"`
	SignedChangePrice  float64 `json:"signed_change_price"`
	SignedChangeRate   float64 `json:"signed_change_rate"`
	TradeVolume        float64 `json:"trade_volume"`
	AccTradePrice      float64 `json:"acc_trade_price"`
	AccTradePrice24h   float64 `json:"acc_trade_price_24h"`
	AccTradeVolume     float64 `json:"acc_trade_volume"`
	AccTradeVolume24h  float64 `json:"acc_trade_volume_24h"`
	Highest52WeekPrice float64 `json:"highest_52_week_price"`
	Highest52WeekDate  string  `json:"highest_52_week_date"`
	Lowest52WeekPrice  float64 `json:"lowest_52_week_price"`
	Lowest52WeekDate   string  `json:"lowest_52_week_date"`
	Timestamp          int64   `json:"timestamp"`
}

// OrderbookUnit 호가 단위
type OrderbookUnit struct {
	AskPrice float64 `json:"ask_price"`
	BidPrice float64 `json:"bid_price"`
	AskSize  float64 `json:"ask_size"`
	BidSize  float64 `json:"bid_size"`
}

// Orderbook 호가 정보
type Orderbook struct {
	Market         string          `json:"market"`
	Timestamp      int64           `json:"timestamp"`
	TotalAskSize   float64         `json:"total_ask_size"`
	TotalBidSize   float64         `json:"total_bid_size"`
	OrderbookUnits []OrderbookUnit `json:"orderbook_units"`
	Level          float64         `json:"level"`
}

// Trade 체결 내역
type Trade struct {
	Market           string  `json:"market"`
	TradeDateUtc     string  `json:"trade_date_utc"`
	TradeTimeUtc     string  `json:"trade_time_utc"`
	Timestamp        int64   `json:"timestamp"`
	TradePrice       float64 `json:"trade_price"`
	TradeVolume      float64 `json:"trade_volume"`
	PrevClosingPrice float64 `json:"prev_closing_price"`
	ChangePrice      float64 `json:"change_price"`
	AskBid           string  `json:"ask_bid"` // ASK, BID
	SequentialID     int64   `json:"sequential_id"`
}

// Candle 캔들 (OHLCV) 데이터
type Candle struct {
	Market               string  `json:"market"`
	CandleDateTimeUtc    string  `json:"candle_date_time_utc"`
	CandleDateTimeKst    string  `json:"candle_date_time_kst"`
	OpeningPrice         float64 `json:"opening_price"`
	HighPrice            float64 `json:"high_price"`
	LowPrice             float64 `json:"low_price"`
	TradePrice           float64 `json:"trade_price"`
	Timestamp            int64   `json:"timestamp"`
	CandleAccTradePrice  float64 `json:"candle_acc_trade_price"`
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	// 분 캔들 전용
	Unit int `json:"unit,omitempty"`
	// 일/주/월/년 캔들 전용
	PrevClosingPrice     float64 `json:"prev_closing_price,omitempty"`
	ChangePrice          float64 `json:"change_price,omitempty"`
	ChangeRate           float64 `json:"change_rate,omitempty"`
	ConvertedTradePrice  float64 `json:"converted_trade_price,omitempty"`
	// 주/월/년 캔들 전용
	FirstDayOfPeriod string `json:"first_day_of_period,omitempty"`
}

// Account 계정 잔고
type Account struct {
	Currency            string  `json:"currency"`
	Balance             Float64 `json:"balance"`
	Locked              Float64 `json:"locked"`
	AvgBuyPrice         Float64 `json:"avg_buy_price"`
	AvgBuyPriceModified bool    `json:"avg_buy_price_modified"`
	UnitCurrency        string  `json:"unit_currency"`
}

// OrderTrade 주문 내 체결 정보
type OrderTrade struct {
	Market    string  `json:"market"`
	UUID      string  `json:"uuid"`
	Price     Float64 `json:"price"`
	Volume    Float64 `json:"volume"`
	Funds     Float64 `json:"funds"`
	Side      string  `json:"side"`
	Trend     string  `json:"trend"`
	CreatedAt string  `json:"created_at"`
}

// Order 주문 정보
type Order struct {
	Market           string       `json:"market"`
	UUID             string       `json:"uuid"`
	Side             string       `json:"side"`     // ask, bid
	OrdType          string       `json:"ord_type"` // limit, price, market, best
	Price            Float64      `json:"price"`
	State            string       `json:"state"` // wait, watch, done, cancel
	CreatedAt        string       `json:"created_at"`
	Volume           Float64      `json:"volume"`
	RemainingVolume  Float64      `json:"remaining_volume"`
	ExecutedVolume   Float64      `json:"executed_volume"`
	ExecutedFunds    Float64      `json:"executed_funds"`
	ReservedFee      Float64      `json:"reserved_fee"`
	RemainingFee     Float64      `json:"remaining_fee"`
	PaidFee          Float64      `json:"paid_fee"`
	Locked           Float64      `json:"locked"`
	PreventedLocked  Float64      `json:"prevented_locked"`
	PreventedVolume  Float64      `json:"prevented_volume,omitempty"`
	TimeInForce      string       `json:"time_in_force,omitempty"`
	SmpType          string       `json:"smp_type,omitempty"`
	Identifier       string       `json:"identifier,omitempty"`
	TradesCount      int          `json:"trades_count"`
	Trades           []OrderTrade `json:"trades"`
}

// Deposit 입금 정보
type Deposit struct {
	Type            string  `json:"type"`
	UUID            string  `json:"uuid"`
	Currency        string  `json:"currency"`
	NetType         string  `json:"net_type"`
	TxID            string  `json:"txid"`
	State           string  `json:"state"`
	Amount          Float64 `json:"amount"`
	Fee             Float64 `json:"fee"`
	TransactionType string  `json:"transaction_type"`
	CreatedAt       string  `json:"created_at"`
	DoneAt          string  `json:"done_at,omitempty"`
}

// DepositAddress 입금 주소 정보
type DepositAddress struct {
	Currency         string `json:"currency"`
	NetType          string `json:"net_type"`
	DepositAddress   string `json:"deposit_address"`
	SecondaryAddress string `json:"secondary_address,omitempty"`
}

// Withdrawal 출금 정보
type Withdrawal struct {
	Type            string  `json:"type"`
	UUID            string  `json:"uuid"`
	Currency        string  `json:"currency"`
	NetType         string  `json:"net_type"`
	TxID            string  `json:"txid"`
	State           string  `json:"state"`
	Amount          Float64 `json:"amount"`
	Fee             Float64 `json:"fee"`
	TransactionType string  `json:"transaction_type"`
	CreatedAt       string  `json:"created_at"`
	DoneAt          string  `json:"done_at,omitempty"`
	IsCancelable    bool    `json:"is_cancelable"`
}

// ServiceStatus 입출금 서비스 상태
type ServiceStatus struct {
	Currency            string `json:"currency"`
	WalletState         string `json:"wallet_state"`
	BlockState          string `json:"block_state"`
	BlockHeight         *int64 `json:"block_height"`
	BlockUpdatedAt      string `json:"block_updated_at"`
	BlockElapsedMinutes *int64 `json:"block_elapsed_minutes"`
	NetType             string `json:"net_type"`
	NetworkName         string `json:"network_name"`
}

// AvailableDeposit 입금 가능 정보
type AvailableDeposit struct {
	Currency                     string `json:"currency"`
	NetType                      string `json:"net_type"`
	IsDepositPossible            bool   `json:"is_deposit_possible"`
	DepositImpossibleReason      string `json:"deposit_impossible_reason"`
	MinimumDepositAmount         string `json:"minimum_deposit_amount"`
	MinimumDepositConfirmations  int    `json:"minimum_deposit_confirmations"`
	DecimalPrecision             int    `json:"decimal_precision"`
}

// CreateDepositAddressResult 입금 주소 생성 결과 (비동기)
type CreateDepositAddressResult struct {
	Success          *bool  `json:"success,omitempty"`
	Message          string `json:"message,omitempty"`
	Currency         string `json:"currency,omitempty"`
	NetType          string `json:"net_type,omitempty"`
	DepositAddress   string `json:"deposit_address,omitempty"`
	SecondaryAddress string `json:"secondary_address,omitempty"`
}

// WithdrawalAddress 출금 허용 주소
type WithdrawalAddress struct {
	Currency              string `json:"currency"`
	NetType               string `json:"net_type"`
	NetworkName           string `json:"network_name"`
	WithdrawAddress       string `json:"withdraw_address"`
	SecondaryAddress      string `json:"secondary_address,omitempty"`
	BeneficiaryName       string `json:"beneficiary_name,omitempty"`
	BeneficiaryCompanyName string `json:"beneficiary_company_name,omitempty"`
	BeneficiaryType       string `json:"beneficiary_type,omitempty"`
	ExchangeName          string `json:"exchange_name,omitempty"`
	WalletType            string `json:"wallet_type,omitempty"`
}

// WithdrawalChance 출금 가능 정보
type WithdrawalChance struct {
	MemberLevel   WithdrawalChanceMemberLevel   `json:"member_level"`
	Currency      WithdrawalChanceCurrency      `json:"currency"`
	Account       WithdrawalChanceAccount       `json:"account"`
	WithdrawLimit WithdrawalChanceWithdrawLimit `json:"withdraw_limit"`
}

// WithdrawalChanceMemberLevel 사용자 보안 등급 정보
type WithdrawalChanceMemberLevel struct {
	SecurityLevel           int  `json:"security_level"`
	FeeLevel                int  `json:"fee_level"`
	EmailVerified           bool `json:"email_verified"`
	IdentityAuthVerified    bool `json:"identity_auth_verified"`
	BankAccountVerified     bool `json:"bank_account_verified"`
	TwoFactorAuthVerified   bool `json:"two_factor_auth_verified"`
	Locked                  bool `json:"locked"`
	WalletLocked            bool `json:"wallet_locked"`
}

// WithdrawalChanceCurrency 통화 정보
type WithdrawalChanceCurrency struct {
	Code          string   `json:"code"`
	WithdrawFee   string   `json:"withdraw_fee"`
	IsCoin        bool     `json:"is_coin"`
	WalletState   string   `json:"wallet_state"`
	WalletSupport []string `json:"wallet_support"`
}

// WithdrawalChanceAccount 자산 잔고 정보
type WithdrawalChanceAccount struct {
	Currency            string  `json:"currency"`
	Balance             Float64 `json:"balance"`
	Locked              Float64 `json:"locked"`
	AvgBuyPrice         Float64 `json:"avg_buy_price"`
	AvgBuyPriceModified bool    `json:"avg_buy_price_modified"`
	UnitCurrency        string  `json:"unit_currency"`
}

// WithdrawalChanceWithdrawLimit 출금 제약 조건
type WithdrawalChanceWithdrawLimit struct {
	Currency              string  `json:"currency"`
	Onetime               string  `json:"onetime"`
	Daily                 *string `json:"daily"`
	RemainingDaily        string  `json:"remaining_daily"`
	RemainingDailyFiat    string  `json:"remaining_daily_fiat"`
	FiatCurrency          string  `json:"fiat_currency"`
	Minimum               string  `json:"minimum"`
	Fixed                 *int    `json:"fixed"`
	WithdrawDelayedFiat   string  `json:"withdraw_delayed_fiat"`
	CanWithdraw           bool    `json:"can_withdraw"`
	RemainingDailyKRW     string  `json:"remaining_daily_krw"`
}

// APIKey API 키 정보
type APIKey struct {
	AccessKey  string `json:"access_key"`
	ExpireAt   string `json:"expire_at"`
}

// OrderChance 주문 가능 정보
type OrderChance struct {
	BidFee      Float64            `json:"bid_fee"`
	AskFee      Float64            `json:"ask_fee"`
	MakerBidFee Float64            `json:"maker_bid_fee"`
	MakerAskFee Float64            `json:"maker_ask_fee"`
	Market      OrderChanceMarket  `json:"market"`
	BidAccount  OrderChanceAccount `json:"bid_account"`
	AskAccount  OrderChanceAccount `json:"ask_account"`
}

// OrderChanceMarket 주문 가능 마켓 정보
type OrderChanceMarket struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	OrderTypes []string          `json:"order_types"` // deprecated
	OrderSides []string          `json:"order_sides"`
	BidTypes   []string          `json:"bid_types"`
	AskTypes   []string          `json:"ask_types"`
	Bid        OrderChanceBidAsk `json:"bid"`
	Ask        OrderChanceBidAsk `json:"ask"`
	MaxTotal   Float64           `json:"max_total"`
	State      string            `json:"state"`
}

// OrderChanceBidAsk 매수/매도 제약 정보
type OrderChanceBidAsk struct {
	Currency  string  `json:"currency"`
	PriceUnit Float64 `json:"price_unit"`
	MinTotal  Float64 `json:"min_total"`
}

// OrderChanceAccount 주문 가능 계좌 정보
type OrderChanceAccount struct {
	Currency            string  `json:"currency"`
	Balance             Float64 `json:"balance"`
	Locked              Float64 `json:"locked"`
	AvgBuyPrice         Float64 `json:"avg_buy_price"`
	AvgBuyPriceModified bool    `json:"avg_buy_price_modified"`
	UnitCurrency        string  `json:"unit_currency"`
}

// VASP 트래블룰 지원 거래소 정보
type VASP struct {
	VaspName     string `json:"vasp_name"`
	VaspEnName   string `json:"en_name,omitempty"`
	VaspUUID     string `json:"vasp_uuid"`
	Depositable  bool   `json:"depositable"`
	Withdrawable bool   `json:"withdrawable"`
}

// TravelRuleVerification 트래블룰 검증 결과
type TravelRuleVerification struct {
	DepositUUID        string `json:"deposit_uuid"`
	VerificationResult string `json:"verification_result"`
	DepositState       string `json:"deposit_state"`
}

// BatchCancelOrderInfo 일괄 취소 결과 내 개별 주문 정보
type BatchCancelOrderInfo struct {
	UUID       string `json:"uuid"`
	Market     string `json:"market"`
	Identifier string `json:"identifier,omitempty"`
}

// BatchCancelGroup 일괄 취소 성공/실패 그룹
type BatchCancelGroup struct {
	Count  int                    `json:"count"`
	Orders []BatchCancelOrderInfo `json:"orders"`
}

// BatchCancelResult 일괄 취소 결과 (CancelOrdersByIDs, BatchCancelOrders 공용)
type BatchCancelResult struct {
	Success BatchCancelGroup `json:"success"`
	Failed  BatchCancelGroup `json:"failed"`
}

// CancelAndNewOrderResult 취소 후 재주문 결과
// 취소된 주문 정보 + 신규 주문 UUID
type CancelAndNewOrderResult struct {
	Market          string  `json:"market"`
	UUID            string  `json:"uuid"`     // 취소된 주문 UUID
	Identifier      string  `json:"identifier,omitempty"`
	Side            string  `json:"side"`     // ask, bid
	OrdType         string  `json:"ord_type"` // limit, price, market, best
	Price           Float64 `json:"price"`
	State           string  `json:"state"`
	CreatedAt       string  `json:"created_at"`
	Volume          Float64 `json:"volume"`
	RemainingVolume Float64 `json:"remaining_volume"`
	ExecutedVolume  Float64 `json:"executed_volume"`
	ReservedFee     Float64 `json:"reserved_fee"`
	RemainingFee    Float64 `json:"remaining_fee"`
	PaidFee         Float64 `json:"paid_fee"`
	Locked          Float64 `json:"locked"`
	PreventedVolume Float64 `json:"prevented_volume"`
	PreventedLocked Float64 `json:"prevented_locked"`
	SmpType         string  `json:"smp_type,omitempty"`
	TradesCount     int     `json:"trades_count"`
	NewOrderUUID    string  `json:"new_order_uuid"` // 신규 생성된 주문 UUID
}
