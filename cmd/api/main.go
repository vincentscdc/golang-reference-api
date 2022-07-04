package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golangreferenceapi/internal/api"
	"golangreferenceapi/internal/api/configuration"

	"github.com/rs/zerolog/log"
)

// @title bnpl API
// @description This is the service used for buy now pay later.
// @version

// @contact.name Vincent
// @contact.email vincent.serpoul@crypto.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	if err := run(); err != nil {
		log.Error().Err(err).Msg("api failed with an error")

		os.Exit(1)
	}

	log.Info().Msg("exited")
}

func run() error {
	ctx := context.Background()

	// configuration
	currEnv := "local"
	if e := os.Getenv("APP_ENVIRONMENT"); e != "" {
		currEnv = e
	}

	configPath := "./config/api"

	cfg, err := configuration.GetConfig(configPath, currEnv)
	if err != nil {
		if errors.Is(err, configuration.MissingBaseConfigError{}) {
			return fmt.Errorf("GetConfig failed: %w", err)
		}

		log.Info().Err(err).Msg("GetConfig")
	}

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// api server
	apiSrv, errN := api.NewAPI(ctx, &cfg)
	if errN != nil {
		return fmt.Errorf("GetConfig failed: %w", errN)
	}

	shutdown, err := apiSrv.Start(ctx)
	if err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	<-sig

	shutdown()

	return nil
}
