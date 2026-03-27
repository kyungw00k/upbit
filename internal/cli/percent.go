package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/kyungw00k/upbit/api"
	"github.com/kyungw00k/upbit/api/exchange"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/types"
)

// isPercent checks if a value string ends with "%"
func isPercent(s string) bool {
	return strings.HasSuffix(s, "%")
}

// parsePercent parses "50%" into 0.5, "100%" into 1.0, etc.
func parsePercent(s string) (float64, error) {
	trimmed := strings.TrimSuffix(s, "%")
	val, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i18n.T(i18n.ErrPercentParse), err)
	}
	if val <= 0 || val > 100 {
		return 0, fmt.Errorf("%s", i18n.T(i18n.ErrPercentRange))
	}
	return val / 100, nil
}

// resolvePercentOrder resolves percentage values to actual amounts using GetOrderChance.
// For buy (side="bid"):
//   - volume%: available_krw * ratio / (1 + fee) / price = volume
//   - total%:  available_krw * ratio / (1 + fee) = total
//
// For sell (side="ask"):
//   - volume%: available_coin * ratio = volume
func resolvePercentOrder(ctx context.Context, client *api.Client, market, side, price, volume, total string) (resolvedPrice, resolvedVolume, resolvedTotal string, err error) {
	resolvedPrice = price
	resolvedVolume = volume
	resolvedTotal = total

	needResolve := isPercent(volume) || isPercent(total)
	if !needResolve {
		return
	}

	ec := exchange.NewExchangeClient(client)
	chance, err := ec.GetOrderChance(ctx, market)
	if err != nil {
		return "", "", "", fmt.Errorf("%s: %w", i18n.T(i18n.ErrChanceFetch), err)
	}

	switch side {
	case "bid":
		err = resolveBuyPercent(chance, price, volume, total, &resolvedPrice, &resolvedVolume, &resolvedTotal)
	case "ask":
		err = resolveSellPercent(chance, volume, &resolvedVolume)
	}
	return
}

func resolveBuyPercent(chance *types.OrderChance, price, volume, total string, resolvedPrice, resolvedVolume, resolvedTotal *string) error {
	available := float64(chance.BidAccount.Balance)
	fee := float64(chance.BidFee)

	if isPercent(total) {
		// 시장가 매수: total = available * ratio / (1 + fee)
		ratio, err := parsePercent(total)
		if err != nil {
			return err
		}
		actualTotal := available * ratio / (1 + fee)
		*resolvedTotal = smartPrice(actualTotal)
		return nil
	}

	if isPercent(volume) {
		// 지정가 매수: volume = available * ratio / (1 + fee) / price
		ratio, err := parsePercent(volume)
		if err != nil {
			return err
		}
		if price == "" {
			return fmt.Errorf("%s", i18n.T(i18n.ErrPercentBuyNeedsPrice))
		}
		priceVal, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.ErrPriceParse), err)
		}
		if priceVal <= 0 {
			return fmt.Errorf("%s", i18n.T(i18n.ErrPriceParse))
		}
		actualVolume := available * ratio / (1 + fee) / priceVal
		*resolvedVolume = strconv.FormatFloat(actualVolume, 'f', 8, 64)
		return nil
	}

	return nil
}

func resolveSellPercent(chance *types.OrderChance, volume string, resolvedVolume *string) error {
	if !isPercent(volume) {
		return nil
	}

	available := float64(chance.AskAccount.Balance)
	ratio, err := parsePercent(volume)
	if err != nil {
		return err
	}
	actualVolume := available * ratio
	*resolvedVolume = strconv.FormatFloat(actualVolume, 'f', 8, 64)
	return nil
}
