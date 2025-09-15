package spapi

import (
	"context"
	"net/http"

	"github.com/zmb3/spotify/v2"
)

type SpotifyApiClient interface {
	GetCurrentPlayingName(ctx context.Context) (string, error)
}

type spotifyApiClient struct {
	client *spotify.Client
}

var _ SpotifyApiClient = (*spotifyApiClient)(nil)

type SpotifyApiClientFactory func(ctx context.Context) (SpotifyApiClient, error)

func NewSpotifyApiClient(ctx context.Context) (SpotifyApiClient, error) {
	client := spotify.New(http.DefaultClient)

	return &spotifyApiClient{
		client,
	}, nil
}

func (c *spotifyApiClient) GetCurrentPlayingName(ctx context.Context) (string, error) {

	currentlyPlaying, err := c.client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return "", err
	}

	if currentlyPlaying.Item == nil {
		return "", nil
	}
	trackName := currentlyPlaying.Item.Name

	return trackName, nil
}
