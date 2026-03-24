package cli

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/kyungw00k/upbit/api"
	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/internal/i18n"
)

// decimalPlaces tickSize 문자열의 소수점 자릿수를 반환
// 예: "0.00000001" → 8, "1000" → 0, "0.01" → 2
func decimalPlaces(s string) int {
	if i := strings.Index(s, "."); i >= 0 {
		return len(s) - i - 1
	}
	return 0
}

// adjustPrice 호가 단위에 맞게 가격을 자동 보정
// side: "bid" -> 내림 (매수자에게 유리), "ask" -> 올림 (매도자에게 유리)
func adjustPrice(ctx context.Context, client *api.Client, market string, price string, side string) (adjusted string, wasAdjusted bool, err error) {
	qc := quotation.NewQuotationClient(client)
	tickSizes, err := qc.GetTickSizes(ctx, []string{market})
	if err != nil {
		return "", false, fmt.Errorf("%s: %w", i18n.T(i18n.ErrTickSizeFetch), err)
	}
	if len(tickSizes) == 0 {
		return "", false, fmt.Errorf("%s", i18n.Tf(i18n.ErrTickSizeEmpty, market))
	}

	tickSize, err := strconv.ParseFloat(tickSizes[0].TickSize, 64)
	if err != nil {
		return "", false, fmt.Errorf("%s: %w", i18n.T(i18n.ErrTickSizeParse), err)
	}
	if tickSize == 0 {
		return "", false, fmt.Errorf("%s", i18n.Tf(i18n.ErrTickSizeZero, market))
	}

	priceVal, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", false, fmt.Errorf("%s: %w", i18n.T(i18n.ErrPriceParse), err)
	}

	remainder := math.Mod(priceVal, tickSize)
	if remainder == 0 {
		return price, false, nil
	}

	var adjustedVal float64
	switch side {
	case "bid":
		// 매수: 내림 (매수자에게 유리)
		adjustedVal = math.Floor(priceVal/tickSize) * tickSize
	case "ask":
		// 매도: 올림 (매도자에게 유리)
		adjustedVal = math.Ceil(priceVal/tickSize) * tickSize
	default:
		return "", false, fmt.Errorf("%s", i18n.Tf(i18n.ErrUnknownSide, side))
	}

	// Q-2: tickSize의 소수점 자릿수에 따라 정밀도를 동적으로 설정
	// BTC/USDT 마켓 등에서 tick_size=0.00000001일 때 소수점이 잘리는 문제 수정
	prec := decimalPlaces(tickSizes[0].TickSize)
	adjustedStr := strconv.FormatFloat(adjustedVal, 'f', prec, 64)
	return adjustedStr, true, nil
}
