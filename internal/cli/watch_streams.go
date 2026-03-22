package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gorillaWS "github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/quotation"
	ws "github.com/kyungw00k/upbit/internal/api/websocket"
	"github.com/kyungw00k/upbit/internal/cache"
	"github.com/kyungw00k/upbit/internal/config"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
	"github.com/kyungw00k/upbit/internal/tui"
)

// --- watch ticker ---

var watchTickerCmd = &cobra.Command{
	Use:   "ticker <market...>",
	Short: i18n.T(i18n.MsgWatchTickerShort),
	Args:  RequireMinArgs(1, i18n.T(i18n.ErrCandleMarketRequired)),
	Example: `  upbit watch ticker KRW-BTC
  upbit watch ticker KRW-BTC KRW-ETH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: "ticker", Codes: codes}}
		return runPublicStream(cmd, sub, formatTicker, func() tea.Model { return tui.NewTickerModel() })
	},
}

// --- watch orderbook ---

var watchOrderbookCmd = &cobra.Command{
	Use:   "orderbook <market...>",
	Short: i18n.T(i18n.MsgWatchOrderbookShort),
	Args:  RequireMinArgs(1, i18n.T(i18n.ErrCandleMarketRequired)),
	Example: `  upbit watch orderbook KRW-BTC
  upbit watch orderbook KRW-BTC KRW-ETH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: "orderbook", Codes: codes}}
		return runPublicStream(cmd, sub, formatOrderbook, func() tea.Model { return tui.NewOrderbookModel() })
	},
}

// --- watch trade ---

var watchTradeCmd = &cobra.Command{
	Use:   "trade <market...>",
	Short: i18n.T(i18n.MsgWatchTradeShort),
	Args:  RequireMinArgs(1, i18n.T(i18n.ErrCandleMarketRequired)),
	Example: `  upbit watch trade KRW-BTC
  upbit watch trade KRW-BTC KRW-ETH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: "trade", Codes: codes}}
		return runPublicStream(cmd, sub, formatTrade, func() tea.Model { return tui.NewTradeModel() })
	},
}

// --- watch candle ---

var watchCandleCmd = &cobra.Command{
	Use:   "candle <market...>",
	Short: i18n.T(i18n.MsgWatchCandleShort),
	Args:  RequireMinArgs(1, i18n.T(i18n.ErrCandleMarketRequired)),
	Example: `  upbit watch candle KRW-BTC -i 1m
  upbit watch candle KRW-BTC KRW-ETH -i 1s`,
	RunE: func(cmd *cobra.Command, args []string) error {
		interval, _ := cmd.Flags().GetString("interval")
		candleType := "candle." + interval
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: candleType, Codes: codes}}

		// 초기 캔들 프리로드 (REST API → TUI에 주입)
		modelFn := func() tea.Model {
			market := codes[0]
			initial := preloadCandles(cmd.Context(), market, interval, 50)
			if len(initial) > 0 {
				return tui.NewCandleModelWithData(market, interval, initial)
			}
			return tui.NewCandleModel()
		}

		return runPublicStream(cmd, sub, formatCandle, modelFn)
	},
}

// --- watch my-order ---

var watchMyOrderCmd = &cobra.Command{
	Use:   "my-order [market...]",
	Short: i18n.T(i18n.MsgWatchMyOrderShort),
	Args:  cobra.ArbitraryArgs,
	Example: `  upbit watch my-order
  upbit watch my-order KRW-BTC`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: "myOrder", Codes: codes}}
		return runPrivateStream(cmd, sub, formatMyOrder, nil)
	},
}

// --- watch my-asset ---

var watchMyAssetCmd = &cobra.Command{
	Use:   "my-asset",
	Short: i18n.T(i18n.MsgWatchMyAssetShort),
	Args:  cobra.NoArgs,
	Example: `  upbit watch my-asset`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// myAsset은 codes를 지원하지 않음
		sub := []ws.SubscriptionType{{Type: "myAsset"}}
		return runPrivateStream(cmd, sub, formatMyAsset, nil)
	},
}

func init() {
	watchCandleCmd.Flags().StringP("interval", "i", "1m",
		i18n.T(i18n.FlagWatchIntervalUsage))

	watchCmd.AddCommand(watchTickerCmd)
	watchCmd.AddCommand(watchOrderbookCmd)
	watchCmd.AddCommand(watchTradeCmd)
	watchCmd.AddCommand(watchCandleCmd)
	watchCmd.AddCommand(watchMyOrderCmd)
	watchCmd.AddCommand(watchMyAssetCmd)
}

// --- 공통 스트림 실행 함수 ---

type lineFormatter func(data []byte) string
type tuiModelFactory func() tea.Model

// runPublicStream Public WebSocket 스트림 실행
func runPublicStream(cmd *cobra.Command, subs []ws.SubscriptionType, fmtFn lineFormatter, modelFn tuiModelFactory) error {
	// 마켓 유효성 검증
	if err := validateMarkets(cmd.Context(), subs); err != nil {
		return err
	}
	client := ws.NewWSClient(ws.PublicURL)
	return runStream(cmd, client, subs, fmtFn, modelFn)
}

// validateMarkets 구독 대상 마켓이 실제 존재하는지 확인 (캐시 활용)
func validateMarkets(ctx context.Context, subs []ws.SubscriptionType) error {
	// 구독에서 마켓 코드 추출
	var codes []string
	for _, sub := range subs {
		codes = append(codes, sub.Codes...)
	}
	if len(codes) == 0 {
		return nil
	}

	validMarkets := getValidMarkets(ctx)
	if validMarkets == nil {
		return nil // 검증 불가 시 진행
	}

	for _, code := range codes {
		if !validMarkets[code] {
			return fmt.Errorf("%s", i18n.Tf(i18n.ErrMarketNotFound, code))
		}
	}
	return nil
}

// getValidMarkets 유효한 마켓 목록 반환 (캐시 우선, 없으면 API 호출 후 캐시)
func getValidMarkets(ctx context.Context) map[string]bool {
	// 1. 캐시 확인 (TTL 1시간)
	mc, err := cache.NewMarketCache(1 * time.Hour)
	if err == nil {
		if cached := mc.Get(); cached != nil {
			return cached
		}
	}

	// 2. API 호출
	apiClient := GetClient()
	qc := quotation.NewQuotationClient(apiClient)
	markets, err := qc.GetMarkets(ctx)
	if err != nil {
		return nil
	}

	// 3. 캐시 저장
	validMarkets := make(map[string]bool, len(markets))
	marketCodes := make([]string, 0, len(markets))
	for _, m := range markets {
		validMarkets[m.Market] = true
		marketCodes = append(marketCodes, m.Market)
	}
	if mc != nil {
		_ = mc.Set(marketCodes)
	}

	return validMarkets
}

// runPrivateStream Private WebSocket 스트림 실행 (인증 필요)
func runPrivateStream(cmd *cobra.Command, subs []ws.SubscriptionType, fmtFn lineFormatter, modelFn tuiModelFactory) error {
	// 인증 확인
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrConfigLoad), err)
	}
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return fmt.Errorf("%s", i18n.T(i18n.ErrAuthRequired))
	}

	client := ws.NewWSClient(ws.PrivateURL, ws.WithAuth(cfg.AccessKey, cfg.SecretKey))
	return runStream(cmd, client, subs, fmtFn, modelFn)
}

// runStream WebSocket 스트림 메인 루프
func runStream(cmd *cobra.Command, client *ws.WSClient, subs []ws.SubscriptionType, fmtFn lineFormatter, modelFn tuiModelFactory) error {
	ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 연결
	if err := client.Connect(ctx); err != nil {
		return err
	}
	defer client.Close()

	// 구독 메시지 생성 및 전송
	subMsg, err := ws.BuildSubscribeMessage(subs)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrSubscribeBuild), err)
	}
	if err := client.Subscribe(subMsg); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrSubscribeSend), err)
	}

	// TTY + TUI 모델이 있으면 bubbletea TUI 실행
	if output.IsTTY() && modelFn != nil {
		return runStreamTUI(ctx, stop, client, modelFn())
	}

	// SIGINT/SIGTERM 시 ReadMessage 블로킹을 즉시 해제하기 위해
	// 별도 goroutine에서 ctx 취소를 감지하여 연결을 닫음
	go func() {
		<-ctx.Done()
		client.Close()
	}()

	isTTY := output.IsTTY()

	for {
		_, data, err := client.ReadMessageWithReconnect(ctx)
		if err != nil {
			// 정상 종료 (Ctrl+C 등)
			if ctx.Err() != nil {
				return nil
			}
			if gorillaWS.IsCloseError(err, gorillaWS.CloseNormalClosure, gorillaWS.CloseGoingAway) {
				return nil
			}
			return fmt.Errorf("%s: %w", i18n.T(i18n.ErrMessageReceive), err)
		}

		// status 메시지 무시 ({"status":"UP"})
		if isStatusMessage(data) {
			continue
		}

		// 에러 메시지 처리 — 서버 에러 시 즉시 종료
		if errMsg := parseErrorMessage(data); errMsg != "" {
			return fmt.Errorf("%s", i18n.Tf(i18n.ErrServerError, errMsg))
		}

		if isTTY {
			// tty (TUI 모델 없는 스트림): 한 줄 테이블 형태
			line := fmtFn(data)
			if line != "" {
				fmt.Println(line)
			}
		} else {
			// non-tty: JSON Lines
			fmt.Println(string(data))
		}
	}
}

// runStreamTUI bubbletea TUI로 WebSocket 스트림 표시
func runStreamTUI(ctx context.Context, stop context.CancelFunc, client *ws.WSClient, model tea.Model) error {
	p := tea.NewProgram(model, tea.WithAltScreen())

	// WebSocket -> TUI 브릿지
	go func() {
		defer func() {
			stop()
			client.Close()
		}()
		for {
			_, data, err := client.ReadMessageWithReconnect(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				if gorillaWS.IsCloseError(err, gorillaWS.CloseNormalClosure, gorillaWS.CloseGoingAway) {
					return
				}
				p.Send(tui.ErrMsg{Err: err})
				return
			}

			// status 메시지 무시
			if isStatusMessage(data) {
				continue
			}

			// 에러 메시지 처리
			if errMsg := parseErrorMessage(data); errMsg != "" {
				p.Send(tui.ServerErrMsg{Message: errMsg})
				return
			}

			p.Send(tui.WsMsg{Data: data})
		}
	}()

	// ctx 취소 시 TUI 종료
	go func() {
		<-ctx.Done()
		p.Quit()
	}()

	_, err := p.Run()
	return err
}

// --- 유틸 함수 ---

// upperCodes 마켓 코드 대문자 변환
func upperCodes(args []string) []string {
	codes := make([]string, len(args))
	for i, a := range args {
		codes[i] = strings.ToUpper(a)
	}
	return codes
}

// isStatusMessage {"status":"UP"} 메시지 여부
func isStatusMessage(data []byte) bool {
	var msg struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(data, &msg); err == nil && msg.Status != "" {
		return true
	}
	return false
}

// parseErrorMessage 에러 응답 파싱
func parseErrorMessage(data []byte) string {
	var msg struct {
		Error struct {
			Name    string `json:"name"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(data, &msg); err == nil && msg.Error.Name != "" {
		return fmt.Sprintf("[%s] %s", msg.Error.Name, msg.Error.Message)
	}
	return ""
}

// --- TTY 포맷터 ---

func formatTicker(data []byte) string {
	var t ws.TickerStream
	if err := json.Unmarshal(data, &t); err != nil {
		return ""
	}
	arrow := changeArrow(t.Change)
	return fmt.Sprintf("%-12s %s %14s  %s%.2f%%  %s: %.4f",
		t.Code, arrow, smartPrice(t.TradePrice),
		signPrefix(t.SignedChangeRate), t.SignedChangeRate*100,
		i18n.T(i18n.WatchVolume),
		t.AccTradeVolume24h,
	)
}

func formatOrderbook(data []byte) string {
	var o ws.OrderbookStream
	if err := json.Unmarshal(data, &o); err != nil {
		return ""
	}
	if len(o.OrderbookUnits) == 0 {
		return ""
	}
	best := o.OrderbookUnits[0]
	spread := best.AskPrice - best.BidPrice
	return fmt.Sprintf("%-12s  %s: %s (%.4f)  %s: %s (%.4f)  %s: %s",
		o.Code,
		i18n.T(i18n.WatchSell),
		smartPrice(best.AskPrice), best.AskSize,
		i18n.T(i18n.WatchBuy),
		smartPrice(best.BidPrice), best.BidSize,
		i18n.T(i18n.WatchSpread),
		smartPrice(spread),
	)
}

func formatTrade(data []byte) string {
	var t ws.TradeStream
	if err := json.Unmarshal(data, &t); err != nil {
		return ""
	}
	side := i18n.T(i18n.WatchBuy)
	if t.AskBid == "ASK" {
		side = i18n.T(i18n.WatchSell)
	}
	ts := time.UnixMilli(t.TradeTimestamp).In(kstLoc).Format("15:04:05")
	return fmt.Sprintf("%-12s  %s  %14s  %s: %.8f  %s",
		t.Code, side, smartPrice(t.TradePrice),
		i18n.T(i18n.WatchQty), t.TradeVolume, ts,
	)
}

func formatCandle(data []byte) string {
	var c ws.CandleStream
	if err := json.Unmarshal(data, &c); err != nil {
		return ""
	}
	return fmt.Sprintf("%-12s  %s  %s:%s %s:%s %s:%s %s:%s  %s: %.4f",
		c.Code, c.CandleDateTimeKST,
		i18n.T(i18n.WatchOpen), smartPrice(c.OpeningPrice),
		i18n.T(i18n.WatchHigh), smartPrice(c.HighPrice),
		i18n.T(i18n.WatchLow), smartPrice(c.LowPrice),
		i18n.T(i18n.WatchClose), smartPrice(c.TradePrice),
		i18n.T(i18n.WatchVolume), c.CandleAccTradeVolume,
	)
}

// preloadCandles REST API로 초기 캔들을 미리 로드
func preloadCandles(ctx context.Context, market, interval string, count int) []tui.CandleData {
	client := GetClient()
	qc := quotation.NewQuotationClient(client)
	candles, err := qc.GetCandles(ctx, market, interval, count)
	if err != nil || len(candles) == 0 {
		return nil
	}

	// API는 최신→오래된 순 → 뒤집기
	result := make([]tui.CandleData, len(candles))
	for i, c := range candles {
		result[len(candles)-1-i] = tui.CandleDataFromOHLCV(
			c.OpeningPrice, c.HighPrice, c.LowPrice, c.TradePrice,
			c.CandleAccTradeVolume, c.CandleDateTimeKst,
		)
	}
	return result
}

func formatMyOrder(data []byte) string {
	var o ws.MyOrderStream
	if err := json.Unmarshal(data, &o); err != nil {
		return ""
	}
	side := i18n.T(i18n.WatchBuy)
	if o.AskBid == "ASK" {
		side = i18n.T(i18n.WatchSell)
	}
	ts := time.UnixMilli(o.Timestamp).In(kstLoc).Format("15:04:05")
	return fmt.Sprintf("%-12s  %s  %s  %s  %s: %.8f  %s: %s  %s",
		o.Code, side, o.OrderType,
		smartPrice(o.Price),
		i18n.T(i18n.WatchQty), o.Volume,
		i18n.T(i18n.WatchState), o.State, ts,
	)
}

func formatMyAsset(data []byte) string {
	var a ws.MyAssetStream
	if err := json.Unmarshal(data, &a); err != nil {
		return ""
	}
	ts := time.UnixMilli(a.Timestamp).In(kstLoc).Format("15:04:05")
	parts := make([]string, 0, len(a.Assets))
	for _, asset := range a.Assets {
		parts = append(parts, fmt.Sprintf("%s: %.8f (%s: %.8f)",
			asset.Currency, asset.Balance,
			i18n.T(i18n.WatchLocking), asset.Locked))
	}
	return fmt.Sprintf("[%s] %s", ts, strings.Join(parts, " | "))
}

// changeArrow 변동 방향 화살표
func changeArrow(change string) string {
	switch change {
	case "RISE":
		return "▲"
	case "FALL":
		return "▼"
	default:
		return "-"
	}
}

// signPrefix 부호 접두사
func signPrefix(rate float64) string {
	if rate > 0 {
		return "+"
	}
	return ""
}

// smartPrice 가격을 크기에 맞게 포맷
// KRW 마켓 (큰 숫자): 106,000,000
// BTC 마켓 (소수점): 0.00002045
// kstLoc KST 시간대
var kstLoc = time.FixedZone("KST", 9*60*60)

func smartPrice(price float64) string {
	if price == 0 {
		return "0"
	}
	abs := price
	if abs < 0 {
		abs = -abs
	}
	switch {
	case abs >= 100:
		return fmt.Sprintf("%.0f", price)
	case abs >= 1:
		return fmt.Sprintf("%.2f", price)
	case abs >= 0.01:
		return fmt.Sprintf("%.4f", price)
	default:
		return fmt.Sprintf("%.8f", price)
	}
}
