package cli

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/kyungw00k/upbit/internal/api"
	"github.com/kyungw00k/upbit/internal/api/quotation"
)

// adjustPrice 호가 단위에 맞게 가격을 자동 보정
// side: "bid" -> 내림 (매수자에게 유리), "ask" -> 올림 (매도자에게 유리)
func adjustPrice(ctx context.Context, client *api.Client, market string, price string, side string) (adjusted string, wasAdjusted bool, err error) {
	qc := quotation.NewQuotationClient(client)
	tickSizes, err := qc.GetTickSizes(ctx, []string{market})
	if err != nil {
		return "", false, fmt.Errorf("호가 단위 조회 실패: %w", err)
	}
	if len(tickSizes) == 0 {
		return "", false, fmt.Errorf("호가 단위 정보 없음: %s", market)
	}

	tickSize, err := strconv.ParseFloat(tickSizes[0].TickSize, 64)
	if err != nil {
		return "", false, fmt.Errorf("호가 단위 파싱 실패: %w", err)
	}
	if tickSize == 0 {
		return "", false, fmt.Errorf("호가 단위가 0: %s", market)
	}

	priceVal, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", false, fmt.Errorf("가격 파싱 실패: %w", err)
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
		return "", false, fmt.Errorf("알 수 없는 side: %s", side)
	}

	// 정수 문자열로 반환 (소수점 없이)
	adjustedStr := strconv.FormatFloat(adjustedVal, 'f', 0, 64)
	return adjustedStr, true, nil
}
