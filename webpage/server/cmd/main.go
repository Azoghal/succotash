package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/azoghal/succotash/webpage/metrics"
	"github.com/azoghal/succotash/webpage/server"
	"github.com/azoghal/succotash/webpage/utils"
)

const (
	apiKeyEnvVar = "SAM_API_KEY"
	apiUrlEnvVar = "API_URL"
)

func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	godotenv.Load()

	// Initialize zerolog logger
	logger := zerolog.New(stderr).With().Timestamp().Logger()

	apiKey := getenv(apiKeyEnvVar)
	apiUrl := getenv(apiUrlEnvVar)

	config := server.ServerConfig{
		ApiKey: apiKey,
		ApiUrl: apiUrl,
	}

	dbGetter := func() error {
		return nil
	}

	server := server.NewServer(config, logger, dbGetter)

	addr := ":6789" // Define the server address
	srv := &http.Server{
		Addr:    addr,
		Handler: server,
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		metricsAddr := ":9090" // Define the metrics server address
		logger.Info().Msgf("Starting metrics server on %s", metricsAddr)
		return metrics.RunMetricsServer(egCtx, metricsAddr)
	})
	eg.Go(func() error {
		logger.Info().Msgf("Starting server on %s", addr)
		return utils.RunServer(egCtx, srv)
	})
	err := eg.Wait()
	if err != nil {
		log.Err(err).Msg("Error during operation")
	}

	<-egCtx.Done()

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
