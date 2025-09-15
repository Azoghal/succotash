package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"server/spapi"
	"server/supabase"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/rs/zerolog"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/azoghal/succotash/webpage/log"
)

const (
	actorDisplayNameKey = "actor-display-id"
	spotifyTokenKey     = "spotify-token-key"
)

// TODO move elsewhere
// keys for attaching stuff to context.Context
type actorDisplayName struct{}
type spotifyTokenKeyStruct struct{}

func GetSpotifyTokenFromCtx(ctx context.Context) (string, error) {
	val, ok := ctx.Value(spotifyTokenKeyStruct{}).(string)
	if !ok {
		return "", errors.New("no spotify token")
	}

	return val, nil
}

type ServerConfig struct {
	WebUrl  string
	Ipv4Url string
	DbUrl   string
}

func NewServer(config ServerConfig, logger zerolog.Logger, getPgDbClient supabase.DBClientFactory, spapiClient spapi.SpotifyApiClientFactory, spotifyAuth *spotifyauth.Authenticator) *gin.Engine {
	router := gin.New()                       // Create a new Gin router without default middleware
	router.Use(gin.Recovery())                // Add default recovery middleware
	router.Use(log.ZerologMiddleware(logger)) // Add custom zerolog logging middleware
	addRoutes(router, config, logger, getPgDbClient, spapiClient, spotifyAuth)
	return router
}

func addRoutes(
	router *gin.Engine,
	config ServerConfig,
	logger zerolog.Logger,
	getPgClient supabase.DBClientFactory,
	spapiClient spapi.SpotifyApiClientFactory,
	spotifyAuth *spotifyauth.Authenticator,
) {
	// redirect anything at / to the landing page
	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/p/landing") })

	router.Static("assets", "./webpage/dist/assets")

	// pages are all served under /p prefix
	webpages := router.Group("/p")
	webpages.GET("/*filepath", func(c *gin.Context) {
		c.File("./webpage/dist/index.html")
	})

	// Spotify auth endpoints
	spotifyLinking := router.Group("/spotify")
	{

		// TODO
		tempState := "temporarystatenotreallygood"

		// redirect to the login
		spotifyLinking.GET("/link", func(c *gin.Context) {
			spotifyLoginUrl := spotifyAuth.AuthURL(tempState)
			c.Redirect(http.StatusTemporaryRedirect, spotifyLoginUrl)
		})

		spotifyLinking.GET("/callback", func(c *gin.Context) {
			token, err := spotifyAuth.Token(c.Request.Context(), tempState, c.Request)
			if err != nil {
				logger.Err(err).Msg("failed to get token")
				c.AbortWithError(http.StatusForbidden, errors.New("failed to get token"))
				return
			}
			if reqState := c.Request.FormValue("state"); reqState != tempState {
				logger.Err(err).Str("expected", tempState).Str("received", reqState).Msg("state does not match")
				c.AbortWithError(http.StatusNotFound, errors.New("failed to get token"))
				return
			}

			logger.Info().Str("tok", token.AccessToken).Msg("user's token!!")

			c.Redirect(http.StatusTemporaryRedirect, "/p/landing")

		})
	}

	// RPC endpoints
	// rpcs := router.Group("/r", func(c *gin.Context) { c.AbortWithStatus(http.StatusNotImplemented) })

	// API endpoints
	restApi := router.Group("/api/v1/")
	{
		unauthedRestApi := restApi.Group("/unauthed")
		{
			unauthedRestApi.GET("/test", helloWorldHandler(config, logger))
		}
		authedRestApi := restApi.Group("/authed", authMiddleware(logger, getPgClient))
		{
			authedRestApi.GET("/", helloUserHandler(config, logger))
			testGroup := authedRestApi.Group("/test")
			{
				testGroup.GET("alice", helloUserHandler(config, logger))
				testGroup.GET("bob", testHandler(config, logger, getPgClient))
			}
			spotifyGroup := authedRestApi.Group("/spotify", spotifyMiddleware(logger))
			{
				spotifyGroup.GET("test", spotifyHandler(config, logger, spapiClient))
			}
		}
	}

	router.NoRoute(func(c *gin.Context) {
		logger.Warn().Msg("404 Not Found")
		c.JSON(404, gin.H{"error": "Not Found"})
	})
}

func helloWorldHandler(config ServerConfig, logger zerolog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger.Info().Msg("helloWorldHandler called")
		c.JSON(200, gin.H{"message": "Hello World"})
	}
}

func helloUserHandler(config ServerConfig, logger zerolog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger.Info().Msg("helloUserHandler called")
		name, ok := c.Get(actorDisplayNameKey)
		if !ok {
			c.JSON(500, gin.H{"error": "failed to get user name"})
			return
		}
		c.JSON(200, gin.H{"message": fmt.Sprintf("Hello %s", name)})
	}
}

func testHandler(config ServerConfig, logger zerolog.Logger, dbClientGetter supabase.DBClientFactory) func(c *gin.Context) {
	return func(c *gin.Context) {
		client := dbClientGetter()
		events, err := client.GetTestEvents()
		if err != nil {
			c.JSON(500, fmt.Sprintf("Failed to get events: %v", err))
		}

		if len(events) > 0 {
			c.JSON(http.StatusOK, gin.H{"content": events[0].Content})
		} else {
			c.JSON(http.StatusOK, gin.H{"content": "nothing"})
		}

		for _, v := range events {
			logger.Info().Str("content", v.Content).Msg("an event")
		}
	}
}

func spotifyHandler(config ServerConfig, logger zerolog.Logger, spapiFactory spapi.SpotifyApiClientFactory) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, spotifyTokenKeyStruct{}, c.Value(spotifyTokenKey))

		spapiClient, err := spapiFactory(ctx)
		if err != nil {
			logger.Err(err).Msg("Failed to get spapi client")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		currentSongName, err := spapiClient.GetCurrentPlayingName(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to get current playing name")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"current_song": currentSongName})
	}
}

type UserMetadata struct {
	AvatarUrl     string `json:"avatar_url"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FullName      string `json:"full_name"`
	Iss           string `json:"iss"`
	Name          string `json:"name"`
	PhoneVerified bool   `json:"phone_verified"`
	Picture       string `json:"picture"`
	ProviderId    string `json:"provider_id"`
	Sub           string `json:"sub"`
}

// For clarity right now, we do it all every time
// probably we should lazily update the keyset if we ever look for a key that doesn't exist
func authMiddleware(logger zerolog.Logger, dbClientGetter supabase.DBClientFactory) func(c *gin.Context) {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("supasession")
		if err != nil {
			logger.Err(err).Msg("failed to get cookie")
			abortUnauthorised(c)
			return
		}
		tokenString := cookie.Value
		logger.Info().Str("supabasesession token", tokenString).Send()

		keyLocation := "https://ocdegtteilykjvohsxrl.supabase.co/auth/v1/.well-known/jwks.json"
		keySet, err := jwk.Fetch(context.Background(), keyLocation)
		if err != nil {
			logger.Err(err).Msg("failed to get key set")
			abortInternal(c)
			return
		}

		parsedVerified, err := jwt.Parse([]byte(tokenString), jwt.WithVerify(true), jwt.WithKeySet(keySet), jwt.WithValidate(true))
		if err != nil {
			fmt.Printf("failed to parse JWT: %s\n", err)
			abortUnauthorised(c)
			return
		}

		var email string
		err = parsedVerified.Get("email", &email)
		if err != nil {
			logger.Err(err).Msg("failed to get email claim")
			abortInternal(c)
			return
		}

		var userMetadata map[string]interface{}
		err = parsedVerified.Get("user_metadata", &userMetadata)
		if err != nil {
			logger.Err(err).Msg("failed to get user metadata claim")
			abortInternal(c)
			return
		}

		userNameA, ok := userMetadata["name"]
		if !ok {
			logger.Error().Msg("failed to get name")
			abortInternal(c)
			return
		}
		userName, ok := userNameA.(string)
		if !ok {
			logger.Error().Msg("failed to cast name")
			abortInternal(c)
			return
		}

		logger.Info().Str("user name", userName).Msg("got the user metadata properly!")

		var sessionId string
		err = parsedVerified.Get("session_id", &sessionId)
		if err != nil {
			logger.Err(err).Msg("failed to get sessionId claim")
			abortInternal(c)
			return
		}

		logger.Info().Str("email", email).Msg("woohoo")

		// Now look up the session id and make sure it's real
		ok, err = dbClientGetter().CheckSession(sessionId)
		if err != nil {
			logger.Err(err).Msg("failed to check session")
			abortInternal(c)
			return
		}

		if !ok {
			logger.Info().Msg("user not logged in")
			abortUnauthorised(c)
			return
		}

		// can set some key on the context if we fancied
		c.Set(actorDisplayNameKey, userName)
		logger.Info().Bool("session exists", ok).Msg("woohooo 2")
	}
}

// Gets the user's spotify token outta the cookie
func spotifyMiddleware(logger zerolog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("succotashspotify")
		if err != nil {
			logger.Err(err).Msg("failed to get spotify cookie")
			return
		}
		tokenString := cookie.Value
		logger.Info().Str("succotashspotify token", tokenString).Send()

		c.Set(spotifyTokenKey, tokenString)
	}
}

func abortUnauthorised(c *gin.Context) {
	c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
}

func abortInternal(c *gin.Context) {
	c.AbortWithError(http.StatusInternalServerError, errors.New("unexpected error"))
}
