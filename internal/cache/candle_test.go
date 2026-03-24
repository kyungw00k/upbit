package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/kyungw00k/upbit/types"
)

// newTestCache 테스트용 인메모리 캐시 생성
func newTestCache(t *testing.T) *CandleCache {
	t.Helper()

	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("테스트 DB 열기 실패: %v", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		t.Fatalf("WAL 모드 설정 실패: %v", err)
	}

	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		t.Fatalf("테이블 생성 실패: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	return &CandleCache{db: db}
}

func TestNewCandleCache(t *testing.T) {
	// 임시 디렉토리 사용
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cc, err := NewCandleCache()
	if err != nil {
		t.Fatalf("NewCandleCache 실패: %v", err)
	}
	defer cc.Close()

	// DB 파일이 생성되었는지 확인
	dbPath := filepath.Join(tmpDir, "upbit", "cache", "candles.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("DB 파일이 생성되지 않음: %s", dbPath)
	}
}

func TestSaveAndQuery(t *testing.T) {
	cc := newTestCache(t)

	rows := []CandleRow{
		{DateTime: "2024-01-01T09:00:00", Open: 100, High: 110, Low: 90, Close: 105, Volume: 1000, AccPrice: 100000},
		{DateTime: "2024-01-01T10:00:00", Open: 105, High: 115, Low: 95, Close: 110, Volume: 1200, AccPrice: 120000},
		{DateTime: "2024-01-01T11:00:00", Open: 110, High: 120, Low: 100, Close: 115, Volume: 1100, AccPrice: 110000},
	}

	err := cc.Save("KRW-BTC", "1h", rows)
	if err != nil {
		t.Fatalf("Save 실패: %v", err)
	}

	// 전체 조회
	result, err := cc.Query("KRW-BTC", "1h", "", "", true)
	if err != nil {
		t.Fatalf("Query 실패: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("3개 행 기대, 실제: %d", len(result))
	}

	// 범위 조회
	result, err = cc.Query("KRW-BTC", "1h", "2024-01-01T10:00:00", "2024-01-01T11:00:00", true)
	if err != nil {
		t.Fatalf("범위 Query 실패: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("2개 행 기대, 실제: %d", len(result))
	}

	// 값 확인
	if result[0].Open != 105 {
		t.Errorf("첫 행 Open 기대: 105, 실제: %f", result[0].Open)
	}

	// DESC 조회
	result, err = cc.Query("KRW-BTC", "1h", "", "", false)
	if err != nil {
		t.Fatalf("DESC Query 실패: %v", err)
	}
	if result[0].DateTime != "2024-01-01T11:00:00" {
		t.Errorf("DESC 정렬 시 가장 최신이 첫 번째여야 함, 실제: %s", result[0].DateTime)
	}
}

func TestSave_EmptySlice(t *testing.T) {
	cc := newTestCache(t)

	err := cc.Save("KRW-BTC", "1h", nil)
	if err != nil {
		t.Fatalf("빈 슬라이스 Save는 에러 없이 반환되어야 함: %v", err)
	}
}

func TestSave_DuplicateIgnored(t *testing.T) {
	cc := newTestCache(t)

	row := CandleRow{DateTime: "2024-01-01T09:00:00", Open: 100, High: 110, Low: 90, Close: 105, Volume: 1000, AccPrice: 100000}

	err := cc.Save("KRW-BTC", "1h", []CandleRow{row})
	if err != nil {
		t.Fatalf("첫 Save 실패: %v", err)
	}

	// 동일 키로 다시 저장 (INSERT OR IGNORE)
	row.Close = 999
	err = cc.Save("KRW-BTC", "1h", []CandleRow{row})
	if err != nil {
		t.Fatalf("중복 Save 실패: %v", err)
	}

	result, _ := cc.Query("KRW-BTC", "1h", "", "", true)
	if len(result) != 1 {
		t.Errorf("중복 무시 후 1개 행 기대, 실제: %d", len(result))
	}
	// INSERT OR IGNORE이므로 원래 값 유지
	if result[0].Close != 105 {
		t.Errorf("중복 무시 시 원래 Close 유지 기대: 105, 실제: %f", result[0].Close)
	}
}

func TestGetRange(t *testing.T) {
	cc := newTestCache(t)

	// 데이터 없을 때
	oldest, newest, err := cc.GetRange("KRW-BTC", "1h")
	if err != nil {
		t.Fatalf("GetRange 실패: %v", err)
	}
	if oldest != "" || newest != "" {
		t.Errorf("데이터 없을 때 빈 문자열 기대, oldest=%s, newest=%s", oldest, newest)
	}

	// 데이터 추가 후
	rows := []CandleRow{
		{DateTime: "2024-01-01T09:00:00", Open: 100, High: 110, Low: 90, Close: 105, Volume: 1000, AccPrice: 100000},
		{DateTime: "2024-01-01T11:00:00", Open: 110, High: 120, Low: 100, Close: 115, Volume: 1100, AccPrice: 110000},
		{DateTime: "2024-01-01T10:00:00", Open: 105, High: 115, Low: 95, Close: 110, Volume: 1200, AccPrice: 120000},
	}
	cc.Save("KRW-BTC", "1h", rows)

	oldest, newest, err = cc.GetRange("KRW-BTC", "1h")
	if err != nil {
		t.Fatalf("GetRange 실패: %v", err)
	}
	if oldest != "2024-01-01T09:00:00" {
		t.Errorf("oldest 기대: 2024-01-01T09:00:00, 실제: %s", oldest)
	}
	if newest != "2024-01-01T11:00:00" {
		t.Errorf("newest 기대: 2024-01-01T11:00:00, 실제: %s", newest)
	}
}

func TestClear(t *testing.T) {
	cc := newTestCache(t)

	rows := []CandleRow{
		{DateTime: "2024-01-01T09:00:00", Open: 100, High: 110, Low: 90, Close: 105, Volume: 1000, AccPrice: 100000},
	}
	cc.Save("KRW-BTC", "1h", rows)
	cc.Save("KRW-ETH", "1d", rows)

	err := cc.Clear()
	if err != nil {
		t.Fatalf("Clear 실패: %v", err)
	}

	result, _ := cc.Query("KRW-BTC", "1h", "", "", true)
	if len(result) != 0 {
		t.Errorf("Clear 후 0개 행 기대, 실제: %d", len(result))
	}

	result, _ = cc.Query("KRW-ETH", "1d", "", "", true)
	if len(result) != 0 {
		t.Errorf("Clear 후 0개 행 기대, 실제: %d", len(result))
	}
}

func TestClearMarket(t *testing.T) {
	cc := newTestCache(t)

	rows := []CandleRow{
		{DateTime: "2024-01-01T09:00:00", Open: 100, High: 110, Low: 90, Close: 105, Volume: 1000, AccPrice: 100000},
	}
	cc.Save("KRW-BTC", "1h", rows)
	cc.Save("KRW-ETH", "1h", rows)

	err := cc.ClearMarket("KRW-BTC", "1h")
	if err != nil {
		t.Fatalf("ClearMarket 실패: %v", err)
	}

	result, _ := cc.Query("KRW-BTC", "1h", "", "", true)
	if len(result) != 0 {
		t.Errorf("ClearMarket 후 BTC 0개 행 기대, 실제: %d", len(result))
	}

	// ETH는 영향 없음
	result, _ = cc.Query("KRW-ETH", "1h", "", "", true)
	if len(result) != 1 {
		t.Errorf("ClearMarket 후 ETH 1개 행 유지 기대, 실제: %d", len(result))
	}
}

func TestCandleToRow(t *testing.T) {
	candle := types.Candle{
		Market:               "KRW-BTC",
		CandleDateTimeKst:    "2024-01-01T18:00:00",
		OpeningPrice:         50000000,
		HighPrice:            51000000,
		LowPrice:             49000000,
		TradePrice:           50500000,
		CandleAccTradeVolume: 100.5,
		CandleAccTradePrice:  5000000000,
	}

	row := CandleToRow(candle)

	if row.DateTime != "2024-01-01T18:00:00" {
		t.Errorf("DateTime 기대: 2024-01-01T18:00:00, 실제: %s", row.DateTime)
	}
	if row.Open != 50000000 {
		t.Errorf("Open 기대: 50000000, 실제: %f", row.Open)
	}
	if row.High != 51000000 {
		t.Errorf("High 기대: 51000000, 실제: %f", row.High)
	}
	if row.Low != 49000000 {
		t.Errorf("Low 기대: 49000000, 실제: %f", row.Low)
	}
	if row.Close != 50500000 {
		t.Errorf("Close 기대: 50500000, 실제: %f", row.Close)
	}
	if row.Volume != 100.5 {
		t.Errorf("Volume 기대: 100.5, 실제: %f", row.Volume)
	}
	if row.AccPrice != 5000000000 {
		t.Errorf("AccPrice 기대: 5000000000, 실제: %f", row.AccPrice)
	}
}

func TestRowToCandle(t *testing.T) {
	row := CandleRow{
		DateTime: "2024-01-01T18:00:00",
		Open:     50000000,
		High:     51000000,
		Low:      49000000,
		Close:    50500000,
		Volume:   100.5,
		AccPrice: 5000000000,
	}

	candle := RowToCandle("KRW-BTC", "1h", row)

	if candle.Market != "KRW-BTC" {
		t.Errorf("Market 기대: KRW-BTC, 실제: %s", candle.Market)
	}
	if candle.CandleDateTimeKst != "2024-01-01T18:00:00" {
		t.Errorf("CandleDateTimeKst 기대: 2024-01-01T18:00:00, 실제: %s", candle.CandleDateTimeKst)
	}
	if candle.OpeningPrice != 50000000 {
		t.Errorf("OpeningPrice 기대: 50000000, 실제: %f", candle.OpeningPrice)
	}
	if candle.TradePrice != 50500000 {
		t.Errorf("TradePrice 기대: 50500000, 실제: %f", candle.TradePrice)
	}
}

func TestCandleToRow_RowToCandle_Roundtrip(t *testing.T) {
	original := types.Candle{
		Market:               "KRW-ETH",
		CandleDateTimeKst:    "2024-06-15T12:00:00",
		OpeningPrice:         4500000,
		HighPrice:            4600000,
		LowPrice:             4400000,
		TradePrice:           4550000,
		CandleAccTradeVolume: 250.75,
		CandleAccTradePrice:  1125000000,
	}

	row := CandleToRow(original)
	roundtripped := RowToCandle("KRW-ETH", "1d", row)

	if roundtripped.Market != original.Market {
		t.Errorf("Market 불일치: %s vs %s", roundtripped.Market, original.Market)
	}
	if roundtripped.OpeningPrice != original.OpeningPrice {
		t.Errorf("OpeningPrice 불일치: %f vs %f", roundtripped.OpeningPrice, original.OpeningPrice)
	}
	if roundtripped.TradePrice != original.TradePrice {
		t.Errorf("TradePrice 불일치: %f vs %f", roundtripped.TradePrice, original.TradePrice)
	}
	if roundtripped.CandleAccTradeVolume != original.CandleAccTradeVolume {
		t.Errorf("Volume 불일치: %f vs %f", roundtripped.CandleAccTradeVolume, original.CandleAccTradeVolume)
	}
}

func TestUpdateLast(t *testing.T) {
	cc := newTestCache(t)

	row1 := CandleRow{DateTime: "2024-01-01T09:00:00", Open: 100, High: 110, Low: 90, Close: 105, Volume: 1000, AccPrice: 100000}
	cc.Save("KRW-BTC", "1h", []CandleRow{row1})

	// 같은 키로 UpdateLast (INSERT OR REPLACE)
	row2 := CandleRow{DateTime: "2024-01-01T09:00:00", Open: 100, High: 120, Low: 85, Close: 115, Volume: 1500, AccPrice: 150000}
	err := cc.UpdateLast("KRW-BTC", "1h", row2)
	if err != nil {
		t.Fatalf("UpdateLast 실패: %v", err)
	}

	result, _ := cc.Query("KRW-BTC", "1h", "", "", true)
	if len(result) != 1 {
		t.Fatalf("1개 행 기대, 실제: %d", len(result))
	}
	if result[0].Close != 115 {
		t.Errorf("UpdateLast 후 Close 기대: 115, 실제: %f", result[0].Close)
	}
	if result[0].High != 120 {
		t.Errorf("UpdateLast 후 High 기대: 120, 실제: %f", result[0].High)
	}
}
