package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os/signal"
	"strings"
	"syscall"
	"time"

	gorillaWS "github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/quotation"
	ws "github.com/kyungw00k/upbit/internal/api/websocket"
	"github.com/kyungw00k/upbit/internal/cache"
	"github.com/kyungw00k/upbit/internal/config"
	"github.com/kyungw00k/upbit/internal/output"
)

// --- watch ticker ---

var watchTickerCmd = &cobra.Command{
	Use:   "ticker <market...>",
	Short: "현재가 실시간 스트림",
	Args:  RequireMinArgs(1, "마켓 코드를 지정하세요 (예: KRW-BTC)"),
	Example: `  upbit watch ticker KRW-BTC
  upbit watch ticker KRW-BTC KRW-ETH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: "ticker", Codes: codes}}
		return runPublicStream(cmd, sub, formatTicker)
	},
}

// --- watch orderbook ---

var watchOrderbookCmd = &cobra.Command{
	Use:   "orderbook <market...>",
	Short: "호가 실시간 스트림",
	Args:  RequireMinArgs(1, "마켓 코드를 지정하세요 (예: KRW-BTC)"),
	Example: `  upbit watch orderbook KRW-BTC
  upbit watch orderbook KRW-BTC KRW-ETH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: "orderbook", Codes: codes}}
		return runPublicStream(cmd, sub, formatOrderbook)
	},
}

// --- watch trade ---

var watchTradeCmd = &cobra.Command{
	Use:   "trade <market...>",
	Short: "체결 실시간 스트림",
	Args:  RequireMinArgs(1, "마켓 코드를 지정하세요 (예: KRW-BTC)"),
	Example: `  upbit watch trade KRW-BTC
  upbit watch trade KRW-BTC KRW-ETH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: "trade", Codes: codes}}
		return runPublicStream(cmd, sub, formatTrade)
	},
}

// --- watch candle ---

var watchCandleCmd = &cobra.Command{
	Use:   "candle <market...>",
	Short: "캔들 실시간 스트림",
	Args:  RequireMinArgs(1, "마켓 코드를 지정하세요 (예: KRW-BTC)"),
	Example: `  upbit watch candle KRW-BTC -i 1m
  upbit watch candle KRW-BTC KRW-ETH -i 1s`,
	RunE: func(cmd *cobra.Command, args []string) error {
		interval, _ := cmd.Flags().GetString("interval")
		candleType := "candle." + interval
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: candleType, Codes: codes}}
		return runPublicStream(cmd, sub, formatCandle)
	},
}

// --- watch my-order ---

var watchMyOrderCmd = &cobra.Command{
	Use:   "my-order [market...]",
	Short: "내 주문 실시간 스트림 (인증 필요)",
	Args:  cobra.ArbitraryArgs,
	Example: `  upbit watch my-order
  upbit watch my-order KRW-BTC`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codes := upperCodes(args)
		sub := []ws.SubscriptionType{{Type: "myOrder", Codes: codes}}
		return runPrivateStream(cmd, sub, formatMyOrder)
	},
}

// --- watch my-asset ---

var watchMyAssetCmd = &cobra.Command{
	Use:   "my-asset",
	Short: "내 자산 실시간 스트림 (인증 필요)",
	Args:  cobra.NoArgs,
	Example: `  upbit watch my-asset`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// myAsset은 codes를 지원하지 않음
		sub := []ws.SubscriptionType{{Type: "myAsset"}}
		return runPrivateStream(cmd, sub, formatMyAsset)
	},
}

func init() {
	watchCandleCmd.Flags().StringP("interval", "i", "1m",
		"캔들 단위 (1s, 1m, 3m, 5m, 10m, 15m, 30m, 60m, 240m)")

	watchCmd.AddCommand(watchTickerCmd)
	watchCmd.AddCommand(watchOrderbookCmd)
	watchCmd.AddCommand(watchTradeCmd)
	watchCmd.AddCommand(watchCandleCmd)
	watchCmd.AddCommand(watchMyOrderCmd)
	watchCmd.AddCommand(watchMyAssetCmd)
}

// --- 공통 스트림 실행 함수 ---

type lineFormatter func(data []byte) string

// runPublicStream Public WebSocket 스트림 실행
func runPublicStream(cmd *cobra.Command, subs []ws.SubscriptionType, fmtFn lineFormatter) error {
	// 마켓 유효성 검증
	if err := validateMarkets(cmd.Context(), subs); err != nil {
		return err
	}
	client := ws.NewWSClient(ws.PublicURL)
	return runStream(cmd, client, subs, fmtFn)
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
			return fmt.Errorf("존재하지 않는 마켓: %s\n\n유효한 마켓 목록은 'upbit market' 명령으로 확인하세요", code)
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
func runPrivateStream(cmd *cobra.Command, subs []ws.SubscriptionType, fmtFn lineFormatter) error {
	// 인증 확인
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("설정 로드 실패: %w", err)
	}
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return fmt.Errorf("인증이 필요합니다: ACCESS_KEY 및 SECRET_KEY를 설정하세요")
	}

	client := ws.NewWSClient(ws.PrivateURL, ws.WithAuth(cfg.AccessKey, cfg.SecretKey))
	return runStream(cmd, client, subs, fmtFn)
}

// runStream WebSocket 스트림 메인 루프
func runStream(cmd *cobra.Command, client *ws.WSClient, subs []ws.SubscriptionType, fmtFn lineFormatter) error {
	ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 연결
	if err := client.Connect(ctx); err != nil {
		return err
	}
	defer client.Close()

	// SIGINT/SIGTERM 시 ReadMessage 블로킹을 즉시 해제하기 위해
	// 별도 goroutine에서 ctx 취소를 감지하여 연결을 닫음
	go func() {
		<-ctx.Done()
		client.Close()
	}()

	// 구독 메시지 생성 및 전송
	subMsg, err := ws.BuildSubscribeMessage(subs)
	if err != nil {
		return fmt.Errorf("구독 메시지 생성 실패: %w", err)
	}
	if err := client.Subscribe(subMsg); err != nil {
		return fmt.Errorf("구독 전송 실패: %w", err)
	}

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
			return fmt.Errorf("메시지 수신 오류: %w", err)
		}

		// status 메시지 무시 ({"status":"UP"})
		if isStatusMessage(data) {
			continue
		}

		// 에러 메시지 처리 — 서버 에러 시 즉시 종료
		if errMsg := parseErrorMessage(data); errMsg != "" {
			return fmt.Errorf("서버 에러: %s", errMsg)
		}

		if isTTY {
			// tty: 한 줄 테이블 형태
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
	return fmt.Sprintf("%-12s %s %14s  %s%.2f%%  거래량: %.4f",
		t.Code, arrow, smartPrice(t.TradePrice),
		signPrefix(t.SignedChangeRate), t.SignedChangeRate*100,
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
	return fmt.Sprintf("%-12s  매도: %s (%.4f)  매수: %s (%.4f)  스프레드: %s",
		o.Code,
		smartPrice(best.AskPrice), best.AskSize,
		smartPrice(best.BidPrice), best.BidSize,
		smartPrice(spread),
	)
}

func formatTrade(data []byte) string {
	var t ws.TradeStream
	if err := json.Unmarshal(data, &t); err != nil {
		return ""
	}
	side := "매수"
	if t.AskBid == "ASK" {
		side = "매도"
	}
	ts := time.UnixMilli(t.TradeTimestamp).In(kstLoc).Format("15:04:05")
	return fmt.Sprintf("%-12s  %s  %14s  수량: %.8f  %s",
		t.Code, side, smartPrice(t.TradePrice), t.TradeVolume, ts,
	)
}

func formatCandle(data []byte) string {
	var c ws.CandleStream
	if err := json.Unmarshal(data, &c); err != nil {
		return ""
	}
	return fmt.Sprintf("%-12s  %s  시:%s 고:%s 저:%s 종:%s  거래량: %.4f",
		c.Code, c.CandleDateTimeKST,
		smartPrice(c.OpeningPrice), smartPrice(c.HighPrice),
		smartPrice(c.LowPrice), smartPrice(c.TradePrice),
		c.CandleAccTradeVolume,
	)
}

func formatMyOrder(data []byte) string {
	var o ws.MyOrderStream
	if err := json.Unmarshal(data, &o); err != nil {
		return ""
	}
	side := "매수"
	if o.AskBid == "ASK" {
		side = "매도"
	}
	ts := time.UnixMilli(o.Timestamp).In(kstLoc).Format("15:04:05")
	return fmt.Sprintf("%-12s  %s  %s  %s  수량: %.8f  상태: %s  %s",
		o.Code, side, o.OrderType,
		smartPrice(o.Price), o.Volume,
		o.State, ts,
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
		parts = append(parts, fmt.Sprintf("%s: %.8f (주문중: %.8f)",
			asset.Currency, asset.Balance, asset.Locked))
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
