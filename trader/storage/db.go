package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyungw00k/upbit/trader/indicator"
	_ "modernc.org/sqlite"
)

func Open(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS candles (
			market   TEXT NOT NULL,
			interval TEXT NOT NULL,
			time_kst TEXT NOT NULL,
			open     REAL NOT NULL,
			high     REAL NOT NULL,
			low      REAL NOT NULL,
			close    REAL NOT NULL,
			volume   REAL NOT NULL,
			PRIMARY KEY (market, interval, time_kst)
		);
		CREATE INDEX IF NOT EXISTS idx_candles_range
			ON candles(market, interval, time_kst);
	`)
	return err
}

func InsertCandles(db *sql.DB, market, interval string, candles []indicator.CandleData) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO candles(market, interval, time_kst, open, high, low, close, volume)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	inserted := 0
	for _, c := range candles {
		res, err := stmt.Exec(market, interval, c.Time, c.Open, c.High, c.Low, c.Close, c.Volume)
		if err != nil {
			return 0, fmt.Errorf("insert %s: %w", c.Time, err)
		}
		n, _ := res.RowsAffected()
		inserted += int(n)
	}
	return inserted, tx.Commit()
}

func LoadCandles(db *sql.DB, market, interval, from, to string) ([]indicator.CandleData, error) {
	query := `SELECT time_kst, open, high, low, close, volume
			  FROM candles WHERE market = ? AND interval = ?`
	args := []any{market, interval}

	if from != "" {
		query += ` AND time_kst >= ?`
		args = append(args, from)
	}
	if to != "" {
		query += ` AND time_kst <= ?`
		args = append(args, to)
	}
	query += ` ORDER BY time_kst ASC`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []indicator.CandleData
	for rows.Next() {
		var c indicator.CandleData
		if err := rows.Scan(&c.Time, &c.Open, &c.High, &c.Low, &c.Close, &c.Volume); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func CandleCount(db *sql.DB, market, interval string) (int, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM candles WHERE market = ? AND interval = ?`,
		market, interval).Scan(&count)
	return count, err
}

func TimeRange(db *sql.DB, market, interval string) (min, max string, err error) {
	err = db.QueryRow(`SELECT MIN(time_kst), MAX(time_kst) FROM candles WHERE market = ? AND interval = ?`,
		market, interval).Scan(&min, &max)
	return
}
