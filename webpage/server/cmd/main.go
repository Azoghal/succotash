package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	// TODO move
	"server/spapi"
	"server/supabase"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/sync/errgroup"

	"github.com/azoghal/succotash/webpage/metrics"
	"github.com/azoghal/succotash/webpage/server"
	"github.com/azoghal/succotash/webpage/utils"
)

const (
	dbUrlEnvVar               = "DATABASE_URL"
	spotifyClientIdEnvVar     = "SPOTIFY_CLIENT_ID"
	spotifyClientSecretEnvVar = "SPOTIFY_CLIENT_SECRET"
	webHostEnvVar             = "WEB_HOST"
	webPortEnvVar             = "WEB_PORT"
	ipv4HostEnvVar            = "WEB_HOST_IPV4"
)

func run(
	ctx context.Context,
	args []string,
	getEnv func(string) string,
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

	webHost := getEnv(webHostEnvVar)
	ipv4Host := getEnv(ipv4HostEnvVar)
	webPort := getEnv(webPortEnvVar)
	dbUrl := getEnv(dbUrlEnvVar)

	config := server.ServerConfig{
		WebUrl:  fmt.Sprintf("%s:%s", webHost, webPort),
		Ipv4Url: fmt.Sprintf("%s:%s", ipv4Host, webPort),
		DbUrl:   dbUrl,
	}

	pgDbGetter := supabase.NewPGClientFactory(dbUrl)

	// Load env for spotify stuff
	spotifyClientId := getEnv(spotifyClientIdEnvVar)
	spotifyClientSecret := getEnv(spotifyClientSecretEnvVar)

	// TODO get url from env
	spotifyAuth := spotifyauth.New(
		spotifyauth.WithClientID(spotifyClientId),
		spotifyauth.WithClientSecret(spotifyClientSecret),
		spotifyauth.WithRedirectURL(fmt.Sprintf("%s/spotify/callback", config.Ipv4Url)),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying),
	)

	server := server.NewServer(config, logger, pgDbGetter, spapi.NewSpotifyApiClient, spotifyAuth)

	addr := ":6789"
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
	if err := run(ctx, os.Args, panickingGetErr, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func panickingGetErr(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("env %s not found", key))
	}
	return val
}
