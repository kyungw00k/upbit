package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/kyungw00k/upbit/internal/i18n"
)

// ToolSchema tool-schema 출력용 구조체
type ToolSchema struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Parameters  ToolSchemaParams    `json:"parameters"`
	Response    *ToolSchemaResponse `json:"response,omitempty"`
}

// ToolSchemaResponse 응답 스키마 정의
type ToolSchemaResponse struct {
	Type        string                       `json:"type"`
	Description string                       `json:"description,omitempty"`
	Items       *ToolSchemaResponseObject    `json:"items,omitempty"`
	Properties  map[string]ToolSchemaProperty `json:"properties,omitempty"`
}

// ToolSchemaResponseObject 배열 내 객체 정의
type ToolSchemaResponseObject struct {
	Type       string                       `json:"type"`
	Properties map[string]ToolSchemaProperty `json:"properties"`
}

// ToolSchemaParams JSON Schema 파라미터 정의
type ToolSchemaParams struct {
	Type       string                        `json:"type"`
	Properties map[string]ToolSchemaProperty  `json:"properties,omitempty"`
	Required   []string                      `json:"required,omitempty"`
}

// ToolSchemaProperty 개별 파라미터 속성
type ToolSchemaProperty struct {
	Type        string              `json:"type"`
	Description string              `json:"description,omitempty"`
	Items       *ToolSchemaProperty `json:"items,omitempty"`
	Default     interface{}         `json:"default,omitempty"`
	Enum        []string            `json:"enum,omitempty"`
}

var toolSchemaCmd = &cobra.Command{
	Use:     "tool-schema [command]",
	Short:   i18n.T(i18n.MsgToolSchemaShort),
	GroupID: "util",
	Args:    cobra.MaximumNArgs(1),
	Example: `  upbit tool-schema              # 전체 명령 스키마
  upbit tool-schema ticker       # ticker 명령 스키마만
  upbit tool-schema buy          # buy 명령 스키마만`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var schemas []ToolSchema

		if len(args) == 1 {
			// 특정 명령만
			target := args[0]
			found := findCommandByName(rootCmd, target)
			if found == nil {
				return fmt.Errorf("%s", i18n.Tf(i18n.ErrToolSchemaNotFound, target))
			}
			schema := buildSchema(found, "")
			schemas = append(schemas, schema)
		} else {
			// 전체 명령 순회
			schemas = collectSchemas(rootCmd, "")
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(schemas)
	},
}

// findCommandByName 이름으로 명령 검색 (재귀, 첫 번째 매칭)
func findCommandByName(parent *cobra.Command, name string) *cobra.Command {
	for _, cmd := range parent.Commands() {
		if cmd.Name() == name {
			return cmd
		}
		if found := findCommandByName(cmd, name); found != nil {
			return found
		}
	}
	return nil
}

// collectSchemas 명령 트리를 재귀 순회하여 스키마 수집
func collectSchemas(parent *cobra.Command, prefix string) []ToolSchema {
	var schemas []ToolSchema

	for _, cmd := range parent.Commands() {
		if cmd.Hidden || cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "tool-schema" {
			continue
		}

		// 서브커맨드가 있는 경우 (order, deposit, withdraw, watch, config 등) 하위만 수집
		if cmd.HasSubCommands() {
			sub := collectSchemas(cmd, buildPrefix(prefix, cmd.Name()))
			schemas = append(schemas, sub...)
			continue
		}

		schema := buildSchema(cmd, prefix)
		schemas = append(schemas, schema)
	}

	return schemas
}

// buildPrefix 명령 접두사 생성
func buildPrefix(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "_" + name
}

// buildSchema 단일 명령의 스키마 생성
func buildSchema(cmd *cobra.Command, prefix string) ToolSchema {
	name := cmd.Name()
	if prefix != "" {
		name = prefix + "_" + name
	}
	// upbit_ 접두사 추가
	fullName := "upbit_" + name

	schema := ToolSchema{
		Name:        fullName,
		Description: cmd.Short,
		Parameters: ToolSchemaParams{
			Type:       "object",
			Properties: make(map[string]ToolSchemaProperty),
		},
	}

	// 위치 인자 (Args)를 파라미터로 변환
	addArgsToSchema(cmd, &schema)

	// 로컬 플래그를 파라미터로 변환
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		prop := flagToProperty(f)
		schema.Parameters.Properties[f.Name] = prop
	})

	// 글로벌 플래그 (output, json, force) 포함
	globalFlags := []string{"output", "json", "force"}
	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		for _, gf := range globalFlags {
			if f.Name == gf {
				prop := flagToProperty(f)
				schema.Parameters.Properties[f.Name] = prop
				return
			}
		}
	})

	// 응답 스키마 매핑
	if resp, ok := responseSchemas[fullName]; ok {
		schema.Response = resp
	}

	return schema
}

// orderResponseProps buy/sell/order show/order list 공용 주문 응답 필드
var orderResponseProps = map[string]ToolSchemaProperty{
	"uuid":             {Type: "string", Description: "Order UUID"},
	"market":           {Type: "string", Description: "Market code"},
	"side":             {Type: "string", Description: "Order side: bid (buy) or ask (sell)"},
	"ord_type":         {Type: "string", Description: "Order type: limit, price (market buy), market (market sell), best"},
	"price":            {Type: "string", Description: "Order price (for limit orders)"},
	"volume":           {Type: "string", Description: "Order volume"},
	"remaining_volume": {Type: "string", Description: "Remaining unfilled volume"},
	"executed_volume":  {Type: "string", Description: "Filled volume"},
	"state":            {Type: "string", Description: "Order state: wait, watch, done, cancel"},
	"created_at":       {Type: "string", Description: "Order creation time (ISO 8601)"},
}

// responseSchemas 명령별 응답 스키마 (키: fullName, 예: "upbit_ticker")
var responseSchemas = map[string]*ToolSchemaResponse{
	"upbit_ticker": {Type: "array", Items: &ToolSchemaResponseObject{Type: "object", Properties: map[string]ToolSchemaProperty{
		"market":               {Type: "string", Description: "Market code (e.g. KRW-BTC)"},
		"trade_price":          {Type: "number", Description: "Current trade price"},
		"opening_price":        {Type: "number", Description: "Opening price of the day"},
		"high_price":           {Type: "number", Description: "Highest price of the day"},
		"low_price":            {Type: "number", Description: "Lowest price of the day"},
		"prev_closing_price":   {Type: "number", Description: "Previous closing price"},
		"change":               {Type: "string", Description: "Price change direction: RISE, EVEN, FALL"},
		"change_price":         {Type: "number", Description: "Absolute price change from previous close"},
		"change_rate":          {Type: "number", Description: "Price change rate (unsigned, e.g. 0.023 = 2.3%)"},
		"signed_change_price":  {Type: "number", Description: "Signed price change (positive=rise, negative=fall)"},
		"signed_change_rate":   {Type: "number", Description: "Signed change rate (e.g. 0.023 = +2.3%, -0.01 = -1%)"},
		"trade_volume":         {Type: "number", Description: "Last trade volume"},
		"acc_trade_volume_24h": {Type: "number", Description: "24h accumulated trade volume"},
		"acc_trade_price_24h":  {Type: "number", Description: "24h accumulated trade price (KRW)"},
		"timestamp":            {Type: "number", Description: "Timestamp in milliseconds"},
	}}},
	"upbit_candle": {Type: "array", Items: &ToolSchemaResponseObject{Type: "object", Properties: map[string]ToolSchemaProperty{
		"market":                   {Type: "string", Description: "Market code"},
		"candle_date_time_utc":     {Type: "string", Description: "Candle start time (UTC, ISO 8601)"},
		"candle_date_time_kst":     {Type: "string", Description: "Candle start time (KST, ISO 8601)"},
		"opening_price":            {Type: "number", Description: "Open price"},
		"high_price":               {Type: "number", Description: "High price"},
		"low_price":                {Type: "number", Description: "Low price"},
		"trade_price":              {Type: "number", Description: "Close price"},
		"candle_acc_trade_volume":  {Type: "number", Description: "Volume during candle period"},
		"candle_acc_trade_price":   {Type: "number", Description: "Trade amount during candle period"},
	}}},
	"upbit_balance": {Type: "array", Items: &ToolSchemaResponseObject{Type: "object", Properties: map[string]ToolSchemaProperty{
		"currency":      {Type: "string", Description: "Currency code (e.g. KRW, BTC)"},
		"balance":       {Type: "string", Description: "Available balance"},
		"locked":        {Type: "string", Description: "Locked balance (in pending orders)"},
		"avg_buy_price": {Type: "string", Description: "Average buy price"},
		"eval_krw":      {Type: "string", Description: "Estimated KRW value (balance * current price)"},
	}}},
	"upbit_buy":  {Type: "object", Properties: orderResponseProps},
	"upbit_sell": {Type: "object", Properties: orderResponseProps},
	"upbit_market": {Type: "array", Items: &ToolSchemaResponseObject{Type: "object", Properties: map[string]ToolSchemaProperty{
		"market":       {Type: "string", Description: "Market code (e.g. KRW-BTC)"},
		"korean_name":  {Type: "string", Description: "Korean name of the asset"},
		"english_name": {Type: "string", Description: "English name of the asset"},
	}}},
	"upbit_orderbook": {Type: "array", Items: &ToolSchemaResponseObject{Type: "object", Properties: map[string]ToolSchemaProperty{
		"market":          {Type: "string", Description: "Market code"},
		"total_ask_size":  {Type: "number", Description: "Total ask (sell) volume"},
		"total_bid_size":  {Type: "number", Description: "Total bid (buy) volume"},
		"orderbook_units": {Type: "array", Description: "List of price levels with ask_price, bid_price, ask_size, bid_size"},
	}}},
	"upbit_trades": {Type: "array", Items: &ToolSchemaResponseObject{Type: "object", Properties: map[string]ToolSchemaProperty{
		"market":          {Type: "string", Description: "Market code"},
		"trade_price":     {Type: "number", Description: "Trade price"},
		"trade_volume":    {Type: "number", Description: "Trade volume"},
		"ask_bid":         {Type: "string", Description: "Trade side: ASK (sell) or BID (buy)"},
		"trade_timestamp": {Type: "number", Description: "Trade timestamp in milliseconds"},
	}}},
	"upbit_order_list": {Type: "array", Items: &ToolSchemaResponseObject{Type: "object", Properties: orderResponseProps}},
	"upbit_order_show": {Type: "object", Properties: orderResponseProps},
	"upbit_tick-size": {Type: "array", Items: &ToolSchemaResponseObject{Type: "object", Properties: map[string]ToolSchemaProperty{
		"market":         {Type: "string", Description: "Market code"},
		"quote_currency": {Type: "string", Description: "Quote currency (KRW, BTC, USDT)"},
		"tick_size":      {Type: "string", Description: "Minimum price unit for orders"},
	}}},
}

// argDescriptions 위치 인자별 설명 (명령 Use 이름 → 인자 이름 → 설명)
var argDescriptions = map[string]map[string]string{
	"ticker":    {"market": i18n.T(i18n.ArgMarketCode)},
	"orderbook": {"market": i18n.T(i18n.ArgOrderbookMarket)},
	"trades":   {"market": i18n.T(i18n.ArgTradesMarket)},
	"candle":   {"market": i18n.T(i18n.ArgCandleMarket)},
	"buy":      {"market": i18n.T(i18n.ArgBuyMarket)},
	"sell":     {"market": i18n.T(i18n.ArgSellMarket)},
	"balance":  {"currency": i18n.T(i18n.ArgBalanceCurrency)},
	"show":     {"uuid": i18n.T(i18n.ArgShowUUID)},
	"cancel":   {"uuid": i18n.T(i18n.ArgCancelUUID)},
	"replace":  {"uuid": i18n.T(i18n.ArgReplaceUUID)},
	"request":  {"currency": i18n.T(i18n.ArgRequestCurrency)},
	"address":  {"currency": i18n.T(i18n.ArgAddressCurrency)},
	"chance":   {"market": i18n.T(i18n.ArgChanceMarket)},
}

// addArgsToSchema Use 필드에서 위치 인자를 추출하여 스키마에 추가
func addArgsToSchema(cmd *cobra.Command, schema *ToolSchema) {
	use := cmd.Use
	// Use에서 명령어 이름 제거 후 인자 파싱
	parts := strings.Fields(use)
	if len(parts) <= 1 {
		return
	}

	cmdName := parts[0]

	for _, part := range parts[1:] {
		// [arg...] → array, optional
		// <arg> → string, required
		// [arg] → string, optional
		clean := strings.Trim(part, "<>[]")
		isRequired := strings.HasPrefix(part, "<")
		isArray := strings.HasSuffix(clean, "...")

		if isArray {
			clean = strings.TrimSuffix(clean, "...")
		}

		// 설명 조회
		desc := ""
		if descs, ok := argDescriptions[cmdName]; ok {
			desc = descs[clean]
		}

		if isArray {
			schema.Parameters.Properties[clean] = ToolSchemaProperty{
				Type:        "array",
				Items:       &ToolSchemaProperty{Type: "string"},
				Description: desc,
			}
		} else {
			schema.Parameters.Properties[clean] = ToolSchemaProperty{
				Type:        "string",
				Description: desc,
			}
		}

		if isRequired {
			schema.Parameters.Required = append(schema.Parameters.Required, clean)
		}
	}
}

// flagToProperty 플래그를 JSON Schema 속성으로 변환
func flagToProperty(f *pflag.Flag) ToolSchemaProperty {
	prop := ToolSchemaProperty{
		Description: f.Usage,
	}

	switch f.Value.Type() {
	case "bool":
		prop.Type = "boolean"
		if f.DefValue != "" && f.DefValue != "false" {
			prop.Default = true
		}
	case "int", "int32", "int64":
		prop.Type = "integer"
	case "float32", "float64":
		prop.Type = "number"
	case "stringSlice", "stringArray":
		prop.Type = "array"
		prop.Items = &ToolSchemaProperty{Type: "string"}
	default:
		prop.Type = "string"
		if f.DefValue != "" {
			prop.Default = f.DefValue
		}
	}

	return prop
}

func init() {
	rootCmd.AddCommand(toolSchemaCmd)
}
