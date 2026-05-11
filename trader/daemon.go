package trader

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Daemon runs the trading pipeline on a fixed interval.
type Daemon struct {
	pipeline *Pipeline
	interval time.Duration
}

// NewDaemon creates a new Daemon with the given pipeline and execution interval.
func NewDaemon(pipeline *Pipeline, interval time.Duration) *Daemon {
	return &Daemon{
		pipeline: pipeline,
		interval: interval,
	}
}

// Start begins the daemon loop. It executes the pipeline immediately once,
// then repeatedly at the configured interval until the context is cancelled
// or a termination signal is received.
func (d *Daemon) Start(ctx context.Context) error {
	// Set up signal handling for graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Execute immediately on start.
	log.Printf("[daemon] starting initial pipeline run")
	if err := d.runOnce(ctx); err != nil {
		log.Printf("[daemon] initial run failed: %v", err)
	}

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[daemon] context cancelled, shutting down")
			return nil
		case <-sigCh:
			log.Printf("[daemon] received termination signal, shutting down")
			return nil
		case <-ticker.C:
			log.Printf("[daemon] running scheduled pipeline execution")
			if err := d.runOnce(ctx); err != nil {
				log.Printf("[daemon] pipeline run failed: %v", err)
			}
		}
	}
}

// Stop performs a clean shutdown of the daemon.
func (d *Daemon) Stop() {
	log.Printf("[daemon] stop requested")
}

// runOnce executes a single pipeline run.
func (d *Daemon) runOnce(ctx context.Context) error {
	start := time.Now()
	result, err := d.pipeline.Run(ctx, map[string]any{}, nil)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	if result != nil {
		log.Printf("[daemon] pipeline completed in %v", elapsed)
	}
	return nil
}
