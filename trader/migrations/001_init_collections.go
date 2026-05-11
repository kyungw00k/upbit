package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// --- indicators ---
		indicators := core.NewBaseCollection("indicators")
		indicators.Fields.Add(&core.TextField{Name: "market"})
		indicators.Fields.Add(&core.TextField{Name: "interval"})
		indicators.Fields.Add(&core.DateField{Name: "candle_time"})
		indicators.Fields.Add(&core.NumberField{Name: "open"})
		indicators.Fields.Add(&core.NumberField{Name: "high"})
		indicators.Fields.Add(&core.NumberField{Name: "low"})
		indicators.Fields.Add(&core.NumberField{Name: "close"})
		indicators.Fields.Add(&core.NumberField{Name: "volume"})
		indicators.Fields.Add(&core.NumberField{Name: "rsi_14"})
		indicators.Fields.Add(&core.NumberField{Name: "macd_line"})
		indicators.Fields.Add(&core.NumberField{Name: "macd_signal"})
		indicators.Fields.Add(&core.NumberField{Name: "macd_hist"})
		indicators.Fields.Add(&core.NumberField{Name: "bb_upper"})
		indicators.Fields.Add(&core.NumberField{Name: "bb_middle"})
		indicators.Fields.Add(&core.NumberField{Name: "bb_lower"})
		indicators.Fields.Add(&core.NumberField{Name: "ema_20"})
		indicators.Fields.Add(&core.NumberField{Name: "ema_50"})
		indicators.Fields.Add(&core.NumberField{Name: "atr_14"})
		indicators.Fields.Add(&core.NumberField{Name: "vol_ratio"})
		indicators.Fields.Add(&core.TextField{Name: "trend"})
		indicators.AddIndex("idx_indicators_unique", true, "market, interval, candle_time", "")
		if err := app.Save(indicators); err != nil {
			return err
		}

		// --- sentiment ---
		sentiment := core.NewBaseCollection("sentiment")
		sentiment.Fields.Add(&core.TextField{Name: "source"})
		sentiment.Fields.Add(&core.NumberField{Name: "value"})
		sentiment.Fields.Add(&core.TextField{Name: "classification"})
		sentiment.Fields.Add(&core.DateField{Name: "timestamp"})
		sentiment.AddIndex("idx_sentiment_unique", true, "source, timestamp", "")
		if err := app.Save(sentiment); err != nil {
			return err
		}

		// --- signals ---
		signals := core.NewBaseCollection("signals")
		signals.Fields.Add(&core.TextField{Name: "market"})
		signals.Fields.Add(&core.TextField{Name: "signal"})
		signals.Fields.Add(&core.NumberField{Name: "confidence"})
		signals.Fields.Add(&core.NumberField{Name: "entry_price"})
		signals.Fields.Add(&core.NumberField{Name: "stop_loss"})
		signals.Fields.Add(&core.NumberField{Name: "take_profit"})
		signals.Fields.Add(&core.TextField{Name: "reasoning"})
		signals.Fields.Add(&core.TextField{Name: "model"})
		signals.Fields.Add(&core.JSONField{Name: "input_snapshot"})
		signals.Fields.Add(&core.TextField{Name: "source"})
		signals.Fields.Add(&core.TextField{Name: "timeframe"})
		signals.Fields.Add(&core.BoolField{Name: "executed"})
		if err := app.Save(signals); err != nil {
			return err
		}

		// --- trades ---
		trades := core.NewBaseCollection("trades")
		trades.Fields.Add(&core.TextField{Name: "market"})
		trades.Fields.Add(&core.TextField{Name: "side"})
		trades.Fields.Add(&core.TextField{Name: "ord_type"})
		trades.Fields.Add(&core.NumberField{Name: "price"})
		trades.Fields.Add(&core.NumberField{Name: "volume"})
		trades.Fields.Add(&core.NumberField{Name: "cost"})
		trades.Fields.Add(&core.NumberField{Name: "fee"})
		trades.Fields.Add(&core.TextField{Name: "order_uuid"})
		trades.Fields.Add(&core.TextField{Name: "status"})
		trades.Fields.Add(&core.TextField{Name: "mode"})
		trades.Fields.Add(&core.NumberField{Name: "pnl"})
		if err := app.Save(trades); err != nil {
			return err
		}

		// --- pipeline_runs ---
		pipelineRuns := core.NewBaseCollection("pipeline_runs")
		pipelineRuns.Fields.Add(&core.DateField{Name: "started_at"})
		pipelineRuns.Fields.Add(&core.DateField{Name: "finished_at"})
		pipelineRuns.Fields.Add(&core.NumberField{Name: "duration_ms"})
		pipelineRuns.Fields.Add(&core.JSONField{Name: "markets_analyzed"})
		pipelineRuns.Fields.Add(&core.NumberField{Name: "signals_generated"})
		pipelineRuns.Fields.Add(&core.NumberField{Name: "trades_placed"})
		pipelineRuns.Fields.Add(&core.TextField{Name: "status"})
		pipelineRuns.Fields.Add(&core.TextField{Name: "error"})
		pipelineRuns.Fields.Add(&core.TextField{Name: "mode"})
		if err := app.Save(pipelineRuns); err != nil {
			return err
		}

		// --- daily_summary ---
		dailySummary := core.NewBaseCollection("daily_summary")
		dailySummary.Fields.Add(&core.DateField{Name: "date"})
		dailySummary.Fields.Add(&core.NumberField{Name: "total_pnl"})
		dailySummary.Fields.Add(&core.NumberField{Name: "total_pnl_pct"})
		dailySummary.Fields.Add(&core.NumberField{Name: "trade_count"})
		dailySummary.Fields.Add(&core.NumberField{Name: "win_rate"})
		dailySummary.Fields.Add(&core.NumberField{Name: "max_drawdown"})
		dailySummary.Fields.Add(&core.JSONField{Name: "portfolio_snapshot"})
		dailySummary.AddIndex("idx_daily_summary_date", true, "date", "")
		if err := app.Save(dailySummary); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		for _, name := range []string{"indicators", "sentiment", "signals", "trades", "pipeline_runs", "daily_summary"} {
			col, err := app.FindCollectionByNameOrId(name)
			if err != nil {
				continue
			}
			if err := app.Delete(col); err != nil {
				return err
			}
		}
		return nil
	})
}
