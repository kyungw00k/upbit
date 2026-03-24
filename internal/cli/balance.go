package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api"
	"github.com/kyungw00k/upbit/api/exchange"
	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
	"github.com/kyungw00k/upbit/types"
)

var balanceColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrCurrency), Key: "currency"},
	{Header: i18n.T(i18n.HdrBalance), Key: "balance", Format: "number"},
	{Header: i18n.T(i18n.HdrLocked), Key: "locked", Format: "number"},
	{Header: i18n.T(i18n.HdrAvgBuyPrice), Key: "avg_buy_price", Format: "number"},
	{Header: i18n.T(i18n.HdrEvalKRW), Key: "eval_krw", Format: "number"},
}

// balanceWithEval 평가금액 포함 잔고
type balanceWithEval struct {
	Currency    string      `json:"currency"`
	Balance     types.Float64 `json:"balance"`
	Locked      types.Float64 `json:"locked"`
	AvgBuyPrice types.Float64 `json:"avg_buy_price"`
	EvalKRW     string      `json:"eval_krw"`
}

var balanceCmd = &cobra.Command{
	Use:     "balance [currency]",
	Short:   i18n.T(i18n.MsgBalanceShort),
	GroupID: "trading",
	Args:    cobra.MaximumNArgs(1),
	Example: `  upbit balance              # 전체 잔고
  upbit balance KRW          # KRW 잔고만
  upbit balance BTC          # BTC 잔고만
  upbit balance -o json      # JSON 출력`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		ec := exchange.NewExchangeClient(client)

		accounts, err := ec.GetAccounts(cmd.Context())
		if err != nil {
			return err
		}

		// currency 필터
		if len(args) == 1 {
			currency := strings.ToUpper(args[0])
			var filtered []types.Account
			for _, a := range accounts {
				if a.Currency == currency {
					filtered = append(filtered, a)
				}
			}
			if len(filtered) == 0 {
				return fmt.Errorf("%s", i18n.Tf(i18n.ErrBalanceNotFound, currency))
			}
			accounts = filtered
		}

		if emptyMessage(accounts, i18n.T(i18n.MsgBalanceEmpty)) {
			return nil
		}

		// 평가금액 계산
		enriched := enrichBalance(cmd.Context(), client, accounts)

		formatter := GetFormatterWithColumns(balanceColumns)
		return formatter.Format(enriched)
	},
}

// enrichBalance 잔고에 평가금액(KRW) 추가
func enrichBalance(ctx context.Context, client *api.Client, accounts []types.Account) []balanceWithEval {
	// 유효한 마켓만 필터하여 현재가 조회
	validMarkets := getValidMarkets(ctx)

	var markets []string
	for _, a := range accounts {
		if a.Currency != "KRW" {
			market := "KRW-" + a.Currency
			if validMarkets == nil || validMarkets[market] {
				markets = append(markets, market)
			}
		}
	}

	// 현재가 조회 (실패해도 평가금액 없이 진행)
	priceMap := make(map[string]float64)
	if len(markets) > 0 {
		qc := quotation.NewQuotationClient(client)
		tickers, err := qc.GetTickers(ctx, markets)
		if err == nil {
			for _, t := range tickers {
				parts := strings.SplitN(t.Market, "-", 2)
				if len(parts) == 2 {
					priceMap[parts[1]] = t.TradePrice
				}
			}
		}
	}

	result := make([]balanceWithEval, len(accounts))
	for i, a := range accounts {
		result[i] = balanceWithEval{
			Currency:    a.Currency,
			Balance:     a.Balance,
			Locked:      a.Locked,
			AvgBuyPrice: a.AvgBuyPrice,
		}

		bal := float64(a.Balance) + float64(a.Locked)
		if a.Currency == "KRW" {
			result[i].EvalKRW = smartPrice(bal)
		} else if price, ok := priceMap[a.Currency]; ok && bal > 0 {
			result[i].EvalKRW = smartPrice(bal * price)
		} else {
			result[i].EvalKRW = "-"
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
