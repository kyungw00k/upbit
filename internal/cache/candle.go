package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/kyungw00k/upbit/internal/types"
)

// CandleRow 캐시 DB 행 구조체
type CandleRow struct {
	DateTime string  `db:"datetime"`
	Open     float64 `db:"open"`
	High     float64 `db:"high"`
	Low      float64 `db:"low"`
	Close    float64 `db:"close"`
	Volume   float64 `db:"volume"`
	AccPrice float64 `db:"acc_price"`
}

// CandleCache SQLite 기반 캔들 캐시
type CandleCache struct {
	db *sqlx.DB
}

const createTableSQL = `
CREATE TABLE IF NOT EXISTS candles (
    market    TEXT NOT NULL,
    timeframe TEXT NOT NULL,
    datetime  TEXT NOT NULL,
    open      REAL NOT NULL,
    high      REAL NOT NULL,
    low       REAL NOT NULL,
    close     REAL NOT NULL,
    volume    REAL NOT NULL,
    acc_price REAL,
    PRIMARY KEY (market, timeframe, datetime)
);
`

// CandleCacheDir 캐시 DB가 저장될 디렉토리 경로 반환
// config.go의 configPath()와 동일한 디렉토리 결정 방식 사용
func CandleCacheDir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "upbit", "cache"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home 디렉토리 확인 실패: %w", err)
	}

	// ~/.config/upbit가 존재하면 그 아래 cache 사용
	xdgDefault := filepath.Join(home, ".config", "upbit")
	if _, err := os.Stat(xdgDefault); err == nil {
		return filepath.Join(xdgDefault, "cache"), nil
	}

	// 폴백: ~/.upbit/cache
	return filepath.Join(home, ".upbit", "cache"), nil
}

// NewCandleCache 캔들 캐시 생성 (DB 열기 + 테이블 생성)
func NewCandleCache() (*CandleCache, error) {
	dir, err := CandleCacheDir()
	if err != nil {
		return nil, fmt.Errorf("캐시 디렉토리 경로 확인 실패: %w", err)
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("캐시 디렉토리 생성 실패 (%s): %w", dir, err)
	}

	dbPath := filepath.Join(dir, "candles.db")
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("캐시 DB 열기 실패: %w", err)
	}

	// WAL 모드 활성화 (동시성 향상)
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("WAL 모드 설정 실패: %w", err)
	}

	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("캐시 테이블 생성 실패: %w", err)
	}

	return &CandleCache{db: db}, nil
}

// Close DB 연결 닫기
func (c *CandleCache) Close() error {
	return c.db.Close()
}

// GetRange 캐시된 데이터의 가장 오래된/최신 datetime 반환
func (c *CandleCache) GetRange(market, timeframe string) (oldest, newest string, err error) {
	row := c.db.QueryRow(
		`SELECT MIN(datetime), MAX(datetime) FROM candles WHERE market = ? AND timeframe = ?`,
		market, timeframe,
	)
	var minDT, maxDT *string
	if err := row.Scan(&minDT, &maxDT); err != nil {
		return "", "", err
	}
	if minDT == nil || maxDT == nil {
		return "", "", nil
	}
	return *minDT, *maxDT, nil
}

// Query 범위 조회
// from, to: ISO 8601 문자열 (빈 문자열이면 제한 없음)
func (c *CandleCache) Query(market, timeframe, from, to string, asc bool) ([]CandleRow, error) {
	query := `SELECT datetime, open, high, low, close, volume, acc_price FROM candles WHERE market = ? AND timeframe = ?`
	args := []interface{}{market, timeframe}

	if from != "" {
		query += ` AND datetime >= ?`
		args = append(args, from)
	}
	if to != "" {
		query += ` AND datetime <= ?`
		args = append(args, to)
	}

	if asc {
		query += ` ORDER BY datetime ASC`
	} else {
		query += ` ORDER BY datetime DESC`
	}

	var rows []CandleRow
	if err := c.db.Select(&rows, query, args...); err != nil {
		return nil, fmt.Errorf("캐시 조회 실패: %w", err)
	}
	return rows, nil
}

// Save 닫힌 캔들 저장 (INSERT OR IGNORE — 이미 존재하면 무시)
func (c *CandleCache) Save(market, timeframe string, candles []CandleRow) error {
	if len(candles) == 0 {
		return nil
	}

	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("트랜잭션 시작 실패: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO candles (market, timeframe, datetime, open, high, low, close, volume, acc_price) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("INSERT 준비 실패: %w", err)
	}
	defer stmt.Close()

	for _, row := range candles {
		if _, err := stmt.Exec(market, timeframe, row.DateTime, row.Open, row.High, row.Low, row.Close, row.Volume, row.AccPrice); err != nil {
			return fmt.Errorf("INSERT 실행 실패: %w", err)
		}
	}

	return tx.Commit()
}

// UpdateLast 마지막(현재) 캔들 업데이트 (INSERT OR REPLACE)
func (c *CandleCache) UpdateLast(market, timeframe string, candle CandleRow) error {
	_, err := c.db.Exec(
		`INSERT OR REPLACE INTO candles (market, timeframe, datetime, open, high, low, close, volume, acc_price) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		market, timeframe, candle.DateTime, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume, candle.AccPrice,
	)
	if err != nil {
		return fmt.Errorf("캔들 업데이트 실패: %w", err)
	}
	return nil
}

// Clear 전체 캐시 삭제
func (c *CandleCache) Clear() error {
	if _, err := c.db.Exec(`DELETE FROM candles`); err != nil {
		return fmt.Errorf("캐시 삭제 실패: %w", err)
	}
	return nil
}

// ClearMarket 특정 마켓/타임프레임 캐시 삭제
func (c *CandleCache) ClearMarket(market, timeframe string) error {
	if _, err := c.db.Exec(`DELETE FROM candles WHERE market = ? AND timeframe = ?`, market, timeframe); err != nil {
		return fmt.Errorf("캐시 삭제 실패: %w", err)
	}
	return nil
}

// CandleToRow types.Candle을 CandleRow로 변환
func CandleToRow(c types.Candle) CandleRow {
	return CandleRow{
		DateTime: c.CandleDateTimeKst,
		Open:     c.OpeningPrice,
		High:     c.HighPrice,
		Low:      c.LowPrice,
		Close:    c.TradePrice,
		Volume:   c.CandleAccTradeVolume,
		AccPrice: c.CandleAccTradePrice,
	}
}

// RowToCandle CandleRow를 types.Candle로 변환
func RowToCandle(market, timeframe string, r CandleRow) types.Candle {
	return types.Candle{
		Market:               market,
		CandleDateTimeKst:    r.DateTime,
		OpeningPrice:         r.Open,
		HighPrice:            r.High,
		LowPrice:             r.Low,
		TradePrice:           r.Close,
		CandleAccTradeVolume: r.Volume,
		CandleAccTradePrice:  r.AccPrice,
	}
}

// CandlesToRows types.Candle 슬라이스를 CandleRow 슬라이스로 변환
func CandlesToRows(candles []types.Candle) []CandleRow {
	rows := make([]CandleRow, len(candles))
	for i, c := range candles {
		rows[i] = CandleToRow(c)
	}
	return rows
}

// RowsToCandles CandleRow 슬라이스를 types.Candle 슬라이스로 변환
func RowsToCandles(market, timeframe string, rows []CandleRow) []types.Candle {
	candles := make([]types.Candle, len(rows))
	for i, r := range rows {
		candles[i] = RowToCandle(market, timeframe, r)
	}
	return candles
}
