package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/kyungw00k/upbit/api"
	"github.com/kyungw00k/upbit/trader/agent"
	"github.com/kyungw00k/upbit/trader/agents"
	"github.com/kyungw00k/upbit/trader/config"
	traderCore "github.com/kyungw00k/upbit/trader"
	"github.com/kyungw00k/upbit/trader/sentiment"
	"github.com/kyungw00k/upbit/trader/types"

	_ "github.com/kyungw00k/upbit/trader/migrations"
)

func main() {
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	registerHooks(app)
	registerRoutes(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func registerHooks(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// Start daemon in background after server is ready.
		go startDaemon(app)

		return se.Next()
	})
}

func startDaemon(app core.App) {
	cfg, err := config.Load("trader/config/config.yaml")
	if err != nil {
		log.Printf("[daemon] config load failed, using defaults: %v", err)
		cfg = config.DefaultConfig()
	}

	interval, err := cfg.Daemon.GetDecisionInterval()
	if err != nil {
		interval = 4 * time.Hour
	}

	// Create agents.
	keltner := agents.NewKeltnerAgent()
	agentList := []traderCore.Agent{keltner}

	// Add strategy agent if enabled.
	if cfg.Strategy.Enabled {
		apiKey := cfg.API.GLMAPIKey
		baseURL := cfg.Strategy.BaseURL
		if baseURL == "" {
			baseURL = cfg.API.GLMBaseURL
		}
		model := cfg.Strategy.Model
		if model == "" {
			model = cfg.API.GLMModel
		}
		timeout := 15 * time.Second
		if cfg.Strategy.Timeout != "" {
			if d, err := time.ParseDuration(cfg.Strategy.Timeout); err == nil {
				timeout = d
			}
		}

		if apiKey != "" {
			strategyAgent := agent.NewStrategyAgent(apiKey, baseURL, model, timeout)
			agentList = append(agentList, strategyAgent)
			log.Printf("[daemon] strategy agent enabled (model=%s, timeout=%v)", model, timeout)
		} else {
			log.Printf("[daemon] strategy agent disabled: no API key")
		}
	}

	apiClient := api.NewClient(cfg.API.UpbitAccessKey, cfg.API.UpbitSecretKey)
	riskGate := traderCore.NewRiskGate(cfg.Risk)
	executor := traderCore.NewPaperExecutor()
	pipeline := traderCore.NewPipeline(agentList, riskGate, executor, cfg)

	log.Printf("[daemon] starting with %d agents, interval=%v", len(agentList), interval)

	// Wait until the next candle close boundary, then run at every interval.
	// For 1h candles: run at exactly :00 of each hour.
	// For 4h candles: run at 00:00, 04:00, 08:00, etc.
	now := time.Now()
	nextBoundary := now.Truncate(interval).Add(interval)
	initialDelay := time.Until(nextBoundary)

	log.Printf("[daemon] first run at %s (waiting %v)", nextBoundary.Format("2006-01-02 15:04:05"), initialDelay)

	// Run once immediately if we're close enough to a boundary (< 30s).
	if initialDelay < 30*time.Second {
		log.Printf("[daemon] running initial cycle (near boundary)")
		runCycle(app, pipeline, cfg, apiClient)
	}

	timer := time.NewTimer(initialDelay)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			runCycle(app, pipeline, cfg, apiClient)
			// Schedule next run at the next boundary.
			nextBoundary = time.Now().Truncate(interval).Add(interval)
			timer.Reset(time.Until(nextBoundary))
		}
	}
}

func runCycle(app core.App, pipeline *traderCore.Pipeline, cfg *config.Config, apiClient *api.Client) {
	ctx := context.Background()
	start := time.Now()
	log.Printf("[daemon] cycle started")

	// Step 1: Collect indicators.
	traderCore.CollectAndCalculate(app, apiClient, cfg.Trading.Markets, "60m")

	// Step 2: Get latest indicators from PocketBase.
	marketData := loadLatestIndicators(app, cfg.Trading.Markets[0])

	// Step 3: Get sentiment.
	var sentimentData *types.SentimentData
	sd, err := sentiment.FetchLatestFearGreedIndex(ctx)
	if err != nil {
		log.Printf("[daemon] sentiment fetch failed: %v", err)
	} else {
		sentimentData = &types.SentimentData{
			FearGreedIndex:    sd.Value,
			FearGreedCategory: sd.Classification,
			Source:            "alternative.me",
			Timestamp:         sd.Timestamp,
		}
	}

	// Step 4: Run pipeline.
	result, err := pipeline.Run(ctx, marketData, sentimentData)
	if err != nil {
		log.Printf("[daemon] pipeline error: %v", err)
	}

	// Step 5: Save pipeline run to PocketBase.
	savePipelineRun(app, start, result, err)

	log.Printf("[daemon] cycle completed in %v", time.Since(start))
}

func loadLatestIndicators(app core.App, market string) map[string]any {
	col, err := app.FindCollectionByNameOrId("indicators")
	if err != nil {
		return map[string]any{}
	}

	records, err := app.FindRecordsByFilter(col, "market={:market}", "-candle_time", 1, 0,
		map[string]any{"market": market})
	if err != nil || len(records) == 0 {
		return map[string]any{}
	}

	r := records[0]
	data := map[string]any{
		"close":         r.GetFloat("close"),
		"rsi":           r.GetFloat("rsi_14"),
		"atr":           r.GetFloat("atr_14"),
		"volume_ratio":  r.GetFloat("vol_ratio"),
		"trend":         r.GetString("trend"),
		"ema": map[string]any{
			"20": r.GetFloat("ema_20"),
			"50": r.GetFloat("ema_50"),
		},
		"bollinger_bands": map[string]any{
			"upper":  r.GetFloat("bb_upper"),
			"middle": r.GetFloat("bb_middle"),
			"lower":  r.GetFloat("bb_lower"),
		},
		"macd": map[string]any{
			"line":     r.GetFloat("macd_line"),
			"signal":   r.GetFloat("macd_signal"),
			"histogram": r.GetFloat("macd_hist"),
		},
	}

	return data
}

func savePipelineRun(app core.App, start time.Time, result *types.PipelineResult, pipelineErr error) {
	col, err := app.FindCollectionByNameOrId("pipeline_runs")
	if err != nil {
		return
	}

	record := core.NewRecord(col)
	record.Set("started_at", start.UTC().Format(time.RFC3339))
	record.Set("finished_at", time.Now().UTC().Format(time.RFC3339))
	record.Set("duration_ms", time.Since(start).Milliseconds())
	record.Set("status", "completed")
	if pipelineErr != nil {
		record.Set("error", pipelineErr.Error())
		record.Set("status", "error")
	}
	app.Save(record)
}

func registerRoutes(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// Dashboard status API.
		se.Router.GET("/api/trader/status", func(e *core.RequestEvent) error {
			// Get latest pipeline run.
			var lastRun map[string]any
			if col, err := app.FindCollectionByNameOrId("pipeline_runs"); err == nil {
				records, _ := app.FindRecordsByFilter(col, "1=1", "-started_at", 1, 0)
				if len(records) > 0 {
					r := records[0]
					lastRun = map[string]any{
						"started_at":  r.GetString("started_at"),
						"finished_at": r.GetString("finished_at"),
						"status":      r.GetString("status"),
						"duration_ms": r.GetInt("duration_ms"),
					}
				}
			}

			// Get open positions.
			var positions []map[string]any
			if col, err := app.FindCollectionByNameOrId("positions"); err == nil {
				records, _ := app.FindRecordsByFilter(col, "1=1", "-created", 10, 0)
				for _, r := range records {
					positions = append(positions, map[string]any{
						"market":        r.GetString("market"),
						"side":          r.GetString("side"),
						"entry_price":   r.GetFloat("entry_price"),
						"current_price": r.GetFloat("current_price"),
						"volume":        r.GetFloat("volume"),
						"unrealized_pnl": r.GetFloat("unrealized_pnl"),
						"stop_loss":     r.GetFloat("stop_loss"),
						"strategy_source": r.GetString("strategy_source"),
					})
				}
			}

			// Get recent agent reports (LLM reasoning).
			var agentReports []map[string]any
			if col, err := app.FindCollectionByNameOrId("agent_reports"); err == nil {
				records, _ := app.FindRecordsByFilter(col, "1=1", "-cycle_time", 20, 0)
				for _, r := range records {
					report := map[string]any{
						"cycle_time":   r.GetString("cycle_time"),
						"agent_id":     r.GetString("agent_id"),
						"signal":       r.GetString("signal"),
						"confidence":   r.GetFloat("confidence"),
						"reasoning":    r.GetString("reasoning"),
						"market_regime": r.GetString("market_regime"),
						"position_size": r.GetFloat("position_size"),
					}
					// Parse strategy_response JSON for LLM reasoning details.
					if sr := r.GetString("strategy_response"); sr != "" {
						var parsed map[string]any
						if json.Unmarshal([]byte(sr), &parsed) == nil {
							report["strategy"] = parsed
						}
					}
					agentReports = append(agentReports, report)
				}
			}

			return e.JSON(200, map[string]any{
				"status":        "running",
				"last_run":      lastRun,
				"positions":     positions,
				"agent_reports": agentReports,
			})
		})

		// Recent trades API.
		se.Router.GET("/api/trader/trades", func(e *core.RequestEvent) error {
			var trades []map[string]any
			if col, err := app.FindCollectionByNameOrId("trades"); err == nil {
				records, _ := app.FindRecordsByFilter(col, "1=1", "-created", 50, 0)
				for _, r := range records {
					trades = append(trades, map[string]any{
						"market":  r.GetString("market"),
						"side":    r.GetString("side"),
						"price":   r.GetFloat("price"),
						"volume":  r.GetFloat("volume"),
						"fee":     r.GetFloat("fee"),
						"pnl":     r.GetFloat("pnl"),
						"mode":    r.GetString("mode"),
						"created": r.GetString("created"),
					})
				}
			}
			return e.JSON(200, map[string]any{"trades": trades})
		})

		// Daily summary API.
		se.Router.GET("/api/trader/summary", func(e *core.RequestEvent) error {
			var summaries []map[string]any
			if col, err := app.FindCollectionByNameOrId("daily_summary"); err == nil {
				records, _ := app.FindRecordsByFilter(col, "1=1", "-date", 30, 0)
				for _, r := range records {
					summaries = append(summaries, map[string]any{
						"date":         r.GetString("date"),
						"total_pnl":    r.GetFloat("total_pnl"),
						"total_pnl_pct": r.GetFloat("total_pnl_pct"),
						"trade_count":  r.GetInt("trade_count"),
						"win_rate":     r.GetFloat("win_rate"),
						"max_drawdown": r.GetFloat("max_drawdown"),
					})
				}
			}
			return e.JSON(200, map[string]any{"summaries": summaries})
		})

		return se.Next()
	})
}
