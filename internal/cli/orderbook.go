package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderbookCmd = &cobra.Command{
	Use:        "orderbook [market...]",
	Short:      i18n.T(i18n.MsgOrderbookShort),
	SuggestFor: []string{"ob", "book"},
	GroupID:    "quotation",
	Args:    RequireMinArgs(1, i18n.T(i18n.ErrOrderbookMarket)),
	Example: `  upbit orderbook KRW-BTC               # 비트코인 호가
  upbit orderbook KRW-BTC KRW-ETH      # 복수 마켓
  upbit orderbook KRW-BTC -o json      # JSON 출력`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		qc := quotation.NewQuotationClient(client)

		orderbooks, err := qc.GetOrderbooks(cmd.Context(), args)
		if err != nil {
			return err
		}

		// 테이블 출력일 때는 호가 단위를 펼쳐서 표시
		formatter := GetFormatter()
		if _, ok := formatter.(*output.TableFormatter); ok || isTableOutput() {
			return formatOrderbookTable(orderbooks)
		}

		return formatter.Format(orderbooks)
	},
}

// isTableOutput 현재 출력이 테이블 모드인지 확인
func isTableOutput() bool {
	switch flagOutput {
	case "table":
		return true
	case "auto", "":
		return output.IsTTY()
	}
	return false
}

// formatOrderbookTable 호가 정보를 커스텀 테이블로 출력
func formatOrderbookTable(orderbooks interface{}) error {
	// orderbooks를 map으로 변환하여 접근
	type orderbookUnit struct {
		AskPrice float64
		BidPrice float64
		AskSize  float64
		BidSize  float64
	}
	type orderbookData struct {
		Market       string
		TotalAskSize float64
		TotalBidSize float64
		Units        []orderbookUnit
	}

	// JSON을 통한 변환
	m := output.ToMapSlice(orderbooks)
	if m == nil {
		return nil
	}

	w := os.Stdout

	// padRight 헬퍼: 표시 폭 기준으로 오른쪽 패딩
	padRight := func(s string, width int) string {
		sw := runewidth.StringWidth(s)
		if sw >= width {
			return s
		}
		return s + strings.Repeat(" ", width-sw)
	}

	for idx, ob := range m {
		if idx > 0 {
			fmt.Fprintln(w)
		}

		market, _ := ob["market"].(string)
		totalAskSize, _ := ob["total_ask_size"].(float64)
		totalBidSize, _ := ob["total_bid_size"].(float64)

		fmt.Fprint(w, i18n.Tf(i18n.MsgOrderbookMarket, market))
		fmt.Fprint(w, i18n.Tf(i18n.MsgOrderbookTotalSizes,
			output.FormatNumberPublic(totalAskSize),
			output.FormatNumberPublic(totalBidSize)))
		fmt.Fprintln(w)

		headers := []string{
			i18n.T(i18n.HdrAskSize),
			i18n.T(i18n.HdrAskPrice),
			i18n.T(i18n.HdrBidPrice),
			i18n.T(i18n.HdrBidSize),
		}

		units, _ := ob["orderbook_units"].([]interface{})

		// 모든 행 데이터 미리 계산
		type unitRow struct{ askSize, askPrice, bidPrice, bidSize string }
		rows := make([]unitRow, 0, len(units))
		for _, u := range units {
			unit, ok := u.(map[string]interface{})
			if !ok {
				continue
			}
			askPrice, _ := unit["ask_price"].(float64)
			bidPrice, _ := unit["bid_price"].(float64)
			askSize, _ := unit["ask_size"].(float64)
			bidSize, _ := unit["bid_size"].(float64)
			rows = append(rows, unitRow{
				askSize:  output.FormatNumberPublic(askSize),
				askPrice: output.FormatNumberPublic(askPrice),
				bidPrice: output.FormatNumberPublic(bidPrice),
				bidSize:  output.FormatNumberPublic(bidSize),
			})
		}

		// 각 컬럼의 최대 표시 폭 계산
		colWidths := make([]int, 4)
		for j, h := range headers {
			colWidths[j] = runewidth.StringWidth(h)
		}
		for _, r := range rows {
			vals := []string{r.askSize, r.askPrice, r.bidPrice, r.bidSize}
			for j, v := range vals {
				if w2 := runewidth.StringWidth(v); w2 > colWidths[j] {
					colWidths[j] = w2
				}
			}
		}

		// 헤더 출력
		headerParts := make([]string, 4)
		for j, h := range headers {
			headerParts[j] = padRight(h, colWidths[j])
		}
		fmt.Fprintln(w, strings.Join(headerParts, "  "))

		// 데이터 행 출력
		for _, r := range rows {
			vals := []string{r.askSize, r.askPrice, r.bidPrice, r.bidSize}
			parts := make([]string, 4)
			for j, v := range vals {
				parts[j] = padRight(v, colWidths[j])
			}
			fmt.Fprintln(w, strings.Join(parts, "  "))
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(orderbookCmd)
}
