package cli

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/internal/cache"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
	"github.com/kyungw00k/upbit/types"
)

var candleColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrCandleTime), Key: "candle_date_time_kst", Format: "datetime"},
	{Header: i18n.T(i18n.HdrOpen), Key: "opening_price", Format: "number"},
	{Header: i18n.T(i18n.HdrHigh), Key: "high_price", Format: "number"},
	{Header: i18n.T(i18n.HdrLow), Key: "low_price", Format: "number"},
	{Header: i18n.T(i18n.HdrClose), Key: "trade_price", Format: "number"},
	{Header: i18n.T(i18n.HdrCandleVolume), Key: "candle_acc_trade_volume", Format: "number"},
}

var candleCmd = &cobra.Command{
	Use:     "candle [market]",
	Short:   i18n.T(i18n.MsgCandleShort),
	GroupID: "quotation",
	Args:    RequireArgs(1, i18n.T(i18n.ErrCandleMarketRequired)),
	Example: `  upbit candle KRW-BTC                        # 일봉 200개
  upbit candle KRW-BTC -i 1m -c 50          # 1분봉 50개
  upbit candle KRW-BTC -i 1w                # 주봉
  upbit candle KRW-BTC -i 1M                # 월봉
  upbit candle KRW-BTC -o json              # JSON 출력
  upbit candle KRW-BTC --from 2025-01-01    # 2025-01-01부터 현재까지
  upbit candle KRW-BTC --from 2025-01-01 -c 500  # 2025-01-01부터 500개
  upbit candle KRW-BTC --desc               # 최신순 정렬
  upbit candle KRW-BTC --no-cache           # 캐시 무시`,
	RunE: runCandle,
}

func init() {
	candleCmd.Flags().StringP("interval", "i", "1d", i18n.T(i18n.FlagIntervalUsage))
	candleCmd.Flags().IntP("count", "c", 200, i18n.T(i18n.FlagCountUsage))
	candleCmd.Flags().String("from", "", i18n.T(i18n.FlagFromUsage))
	candleCmd.Flags().Bool("asc", true, i18n.T(i18n.FlagAscUsage))
	candleCmd.Flags().Bool("desc", false, i18n.T(i18n.FlagDescUsage))
	candleCmd.Flags().Bool("no-cache", false, i18n.T(i18n.FlagNoCacheUsage))
	AddForceFlag(candleCmd)
	rootCmd.AddCommand(candleCmd)
}

func runCandle(cmd *cobra.Command, args []string) error {
	client := GetClient()
	qc := quotation.NewQuotationClient(client)
	formatter := GetFormatterWithColumns(candleColumns)

	market := args[0]
	interval, _ := cmd.Flags().GetString("interval")
	count, _ := cmd.Flags().GetInt("count")
	fromStr, _ := cmd.Flags().GetString("from")
	descFlag, _ := cmd.Flags().GetBool("desc")
	ascFlag, _ := cmd.Flags().GetBool("asc")
	noCache, _ := cmd.Flags().GetBool("no-cache")
	force := GetForce()

	// --desc가 명시적으로 설정되면 desc 우선
	ascending := ascFlag && !descFlag

	// 자동 페이지네이션 모드 판단: --from이 있거나 -c > 200
	usePagination := fromStr != "" || count > 200

	if !usePagination {
		// 기존 동작: 단일 API 호출, 캐시 안 씀
		candles, err := qc.GetCandles(cmd.Context(), market, interval, count)
		if err != nil {
			return err
		}
		// 기본 API 응답은 최신→오래된 순
		if ascending {
			sort.Slice(candles, func(i, j int) bool {
				return candles[i].CandleDateTimeKst < candles[j].CandleDateTimeKst
			})
		}
		return formatter.Format(candles)
	}

	// --- 자동 페이지네이션 모드 ---

	// from 파싱
	var fromTime time.Time
	if fromStr != "" {
		var err error
		fromTime, err = parseFrom(fromStr)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.ErrFromParse), err)
		}
	}

	// 예상 캔들 개수 / API 호출 횟수 계산 후 확인 프롬프트
	estimatedCount := count
	if fromStr != "" && count == 200 {
		// count를 명시하지 않은 경우 (기본값 200) → from부터 현재까지 전부
		estimatedCount = estimateCandleCount(fromTime, interval)
	} else if fromStr != "" && count > 200 {
		est := estimateCandleCount(fromTime, interval)
		if est < count {
			estimatedCount = est
		}
	}

	apiCalls := int(math.Ceil(float64(estimatedCount) / 200.0))
	if apiCalls >= 10 {
		msg := i18n.Tf(i18n.MsgCandleConfirm, estimatedCount, apiCalls)
		ok, err := output.Confirm(msg, force)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCancelled))
			return nil
		}
	}

	// from이 설정되었고 count가 기본값(200)이면 → count=0 (전부)
	fetchCount := count
	if fromStr != "" && count == 200 {
		fetchCount = 0
	}

	// from을 KST ISO 8601 형태로 변환 (API 및 캐시에서 사용)
	var fromKST string
	if fromStr != "" {
		fromKST = fromTime.In(types.KSTLoc).Format("2006-01-02T15:04:05")
	}

	var candles []types.Candle

	if noCache {
		// 비캐시 모드: 직접 API 호출
		result, err := qc.GetCandlesAll(cmd.Context(), market, interval, fromKST, fetchCount)
		if err != nil {
			return err
		}
		candles = result
	} else {
		// 캐시 모드
		result, err := fetchWithCache(cmd, qc, market, interval, fromKST, fetchCount)
		if err != nil {
			return err
		}
		candles = result
	}

	// 정렬
	if ascending {
		sort.Slice(candles, func(i, j int) bool {
			return candles[i].CandleDateTimeKst < candles[j].CandleDateTimeKst
		})
	} else {
		sort.Slice(candles, func(i, j int) bool {
			return candles[i].CandleDateTimeKst > candles[j].CandleDateTimeKst
		})
	}

	return formatter.Format(candles)
}

// fetchWithCache 캐시를 활용한 캔들 조회
func fetchWithCache(cmd *cobra.Command, qc *quotation.QuotationClient, market, interval, fromKST string, count int) ([]types.Candle, error) {
	cc, err := cache.NewCandleCache()
	if err != nil {
		// 캐시 생성 실패 시 직접 API 호출로 폴백
		fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCacheInitFailed), err)
		return qc.GetCandlesAll(cmd.Context(), market, interval, fromKST, count)
	}
	defer cc.Close()

	// 캐시에서 현재 범위 확인
	cachedOldest, cachedNewest, err := cc.GetRange(market, interval)
	if err != nil {
		fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCacheRangeError), err)
	}

	nowKST := time.Now().In(types.KSTLoc).Format("2006-01-02T15:04:05")

	if cachedOldest == "" {
		// 캐시 비어있음 → 전체 API 호출
		result, err := qc.GetCandlesAll(cmd.Context(), market, interval, fromKST, count)
		if err != nil {
			return nil, err
		}
		// 결과를 캐시에 저장
		saveToCacheAll(cc, market, interval, result)
		return result, nil
	}

	// 캐시에 데이터가 있는 경우: 부족한 구간 확인 + 보충
	if fromKST != "" && fromKST < cachedOldest {
		// 오래된 구간 보충
		olderCandles, err := qc.GetCandlesAll(cmd.Context(), market, interval, fromKST, 0)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T(i18n.ErrOlderRangeFetch), err)
		}
		// cachedOldest 이전 캔들만 필터링하여 저장
		var toSave []cache.CandleRow
		for _, c := range olderCandles {
			if c.CandleDateTimeKst < cachedOldest {
				toSave = append(toSave, cache.CandleToRow(c))
			}
		}
		if len(toSave) > 0 {
			if err := cc.Save(market, interval, toSave); err != nil {
				fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCacheSaveError), err)
			}
		}
	}

	if cachedNewest < nowKST {
		// 최신 구간 보충
		newerCandles, err := qc.GetCandlesAll(cmd.Context(), market, interval, cachedNewest, 0)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T(i18n.ErrNewerRangeFetch), err)
		}
		if len(newerCandles) > 0 {
			rows := cache.CandlesToRows(newerCandles)
			// 마지막 캔들은 아직 닫히지 않았을 수 있으므로 UpdateLast
			if len(rows) > 1 {
				if err := cc.Save(market, interval, rows[:len(rows)-1]); err != nil {
					fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCacheSaveError), err)
				}
			}
			if err := cc.UpdateLast(market, interval, rows[len(rows)-1]); err != nil {
				fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCacheUpdateError), err)
			}
		}
	}

	// 캐시에서 전체 범위 조회
	rows, err := cc.Query(market, interval, fromKST, nowKST, true)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.T(i18n.ErrCacheQuery), err)
	}

	result := cache.RowsToCandles(market, interval, rows)

	// count 제한 적용
	if count > 0 && len(result) > count {
		result = result[:count]
	}

	return result, nil
}

// saveToCacheAll API 결과 전체를 캐시에 저장
func saveToCacheAll(cc *cache.CandleCache, market, interval string, candles []types.Candle) {
	if len(candles) == 0 {
		return
	}
	rows := cache.CandlesToRows(candles)
	// 마지막 캔들(최신)은 UpdateLast, 나머지는 Save
	if len(rows) > 1 {
		if err := cc.Save(market, interval, rows[:len(rows)-1]); err != nil {
			fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCacheSaveError), err)
		}
	}
	if err := cc.UpdateLast(market, interval, rows[len(rows)-1]); err != nil {
		fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCacheUpdateError), err)
	}
}

// parseFrom 다양한 형식의 시각 문자열을 time.Time으로 파싱
func parseFrom(s string) (time.Time, error) {
	formats := []struct {
		layout string
		loc    *time.Location
	}{
		// RFC3339 (타임존 포함)
		{"2006-01-02T15:04:05Z07:00", nil},
		{"2006-01-02T15:04:05Z", time.UTC},
		// 타임존 없는 형식 → KST 가정
		{"2006-01-02T15:04:05", types.KSTLoc},
		{"2006-01-02T15:04", types.KSTLoc},
		// 날짜만 → KST 00:00:00
		{"2006-01-02", types.KSTLoc},
		{"2006/01/02", types.KSTLoc},
	}

	s = strings.TrimSpace(s)

	for _, f := range formats {
		if f.loc == nil {
			t, err := time.Parse(f.layout, s)
			if err == nil {
				return t, nil
			}
		} else {
			t, err := time.ParseInLocation(f.layout, s, f.loc)
			if err == nil {
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("%s", i18n.Tf(i18n.ErrUnrecognizedTime, s))
}

// estimateCandleCount from 시각부터 현재까지 예상 캔들 개수 계산
func estimateCandleCount(from time.Time, interval string) int {
	duration := time.Since(from)
	if duration <= 0 {
		return 0
	}

	switch interval {
	case "1s":
		return int(duration.Seconds())
	case "1m":
		return int(duration.Minutes())
	case "3m":
		return int(duration.Minutes() / 3)
	case "5m":
		return int(duration.Minutes() / 5)
	case "10m":
		return int(duration.Minutes() / 10)
	case "15m":
		return int(duration.Minutes() / 15)
	case "30m":
		return int(duration.Minutes() / 30)
	case "60m":
		return int(duration.Minutes() / 60)
	case "240m":
		return int(duration.Minutes() / 240)
	case "1d":
		return int(duration.Hours() / 24)
	case "1w":
		return int(duration.Hours() / (24 * 7))
	case "1M":
		return int(duration.Hours() / (24 * 30))
	case "1y":
		return int(duration.Hours() / (24 * 365))
	default:
		return int(duration.Hours() / 24)
	}
}
