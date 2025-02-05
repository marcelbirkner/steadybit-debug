package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/steadybit-debug/agent"
	"github.com/steadybit/steadybit-debug/config"
	"github.com/steadybit/steadybit-debug/output"
	"github.com/steadybit/steadybit-debug/platform"
	"os"
	"sync"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.GetConfig()
	output.AddOutputDirectory(&cfg)

	output.AddJsonOutput(output.AddJsonOutputOptions{
		Config:     &cfg,
		Content:    cfg,
		OutputPath: []string{"debugging_config.yaml"},
	})
	gatherInformation(&cfg)
	output.ZipOutputDirectory(&cfg)

	if !cfg.NoCleanup {
		err := os.RemoveAll(cfg.OutputPath)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to remove output directory '%s' after completion", cfg.OutputPath)
		}
	}
}

func gatherInformation(cfg *config.Config) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		platform.AddPlatformDebuggingInformation(cfg)
	}()

	go func() {
		defer wg.Done()
		agent.AddAgentDebuggingInformation(cfg)
	}()

	wg.Wait()
}
