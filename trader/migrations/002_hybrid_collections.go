package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// --- positions (current open positions) ---
		positions := core.NewBaseCollection("positions")
		positions.Fields.Add(&core.TextField{Name: "market"})
		positions.Fields.Add(&core.TextField{Name: "side"})
		positions.Fields.Add(&core.NumberField{Name: "entry_price"})
		positions.Fields.Add(&core.NumberField{Name: "volume"})
		positions.Fields.Add(&core.NumberField{Name: "cost"})
		positions.Fields.Add(&core.NumberField{Name: "current_price"})
		positions.Fields.Add(&core.NumberField{Name: "stop_loss"})
		positions.Fields.Add(&core.NumberField{Name: "trailing_stop"})
		positions.Fields.Add(&core.NumberField{Name: "unrealized_pnl"})
		positions.Fields.Add(&core.TextField{Name: "entry_time"})
		positions.Fields.Add(&core.TextField{Name: "strategy_source"})
		positions.Fields.Add(&core.JSONField{Name: "execution_params"})
		positions.AddIndex("idx_positions_market", false, "market", "")
		if err := app.Save(positions); err != nil {
			return err
		}

		// --- agent_reports (per-cycle agent outputs) ---
		agentReports := core.NewBaseCollection("agent_reports")
		agentReports.Fields.Add(&core.DateField{Name: "cycle_time"})
		agentReports.Fields.Add(&core.TextField{Name: "agent_id"})
		agentReports.Fields.Add(&core.TextField{Name: "market"})
		agentReports.Fields.Add(&core.TextField{Name: "signal"})
		agentReports.Fields.Add(&core.NumberField{Name: "confidence"})
		agentReports.Fields.Add(&core.JSONField{Name: "strategy_response"})
		agentReports.Fields.Add(&core.TextField{Name: "reasoning"})
		agentReports.Fields.Add(&core.TextField{Name: "market_regime"})
		agentReports.Fields.Add(&core.NumberField{Name: "position_size"})
		agentReports.Fields.Add(&core.JSONField{Name: "exit_conditions"})
		agentReports.AddIndex("idx_agent_reports_cycle", false, "cycle_time, agent_id", "")
		if err := app.Save(agentReports); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		for _, name := range []string{"positions", "agent_reports"} {
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
