package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"server/supabase"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/azoghal/succotash/webpage/metrics"
	"github.com/azoghal/succotash/webpage/server"
	"github.com/azoghal/succotash/webpage/utils"
)

const (
	supabaseKeyEnvVar    = "SUPABASE_API_KEY"
	supabaseApiUrlEnvVar = "SUPABASE_API_URL"
	dbUrlEnvVar          = "DATABASE_URL"
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

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Initialize zerolog logger
	logger := zerolog.New(stderr).With().Timestamp().Logger()

	apiKey := getenv(supabaseKeyEnvVar)
	apiUrl := getenv(supabaseApiUrlEnvVar)
	dbUrl := getenv(dbUrlEnvVar)

	config := server.ServerConfig{
		ApiKey: apiKey,
		ApiUrl: apiUrl,
	}

	restDbGetter := supabase.NewRestDBClientFactory(config.ApiUrl, config.ApiKey)

	pgDbGetter := supabase.NewPGClientFactory(dbUrl)

	server := server.NewServer(config, logger, restDbGetter, pgDbGetter)

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
	err = eg.Wait()
	if err != nil {
		log.Err(err).Msg("Error during operation")
	}

	<-egCtx.Done()

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, PanickingGetErr, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func PanickingGetErr(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("env %s not found", key))
	}
	return val
}
