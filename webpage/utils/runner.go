package utils

import (
	"context"
	"net/http"
	"time"
)

func RunServer(ctx context.Context, server *http.Server) error {

	var shutdownCtx context.Context
	var cancel context.CancelFunc

	errCh := make(chan error, 1)

	go func() {
		// log.Info().Msg("running HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		{
			// log.Info().Msg("Shutting down HTTP server due to context cancellation")
			shutdownCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
		}
	}

	return server.Shutdown(shutdownCtx)
}
