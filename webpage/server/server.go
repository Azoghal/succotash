package server

import (
	"context"
	"fmt"
	"net/http"
	"server/supabase"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/rs/zerolog"

	"github.com/azoghal/succotash/webpage/log"
)

type ServerConfig struct {
	ApiKey string
	ApiUrl string
}

func NewServer(config ServerConfig, logger zerolog.Logger, getRestDbClient supabase.DBClientFactory, getPgDbClient supabase.DBClientFactory) *gin.Engine {
	router := gin.New()                       // Create a new Gin router without default middleware
	router.Use(gin.Recovery())                // Add default recovery middleware
	router.Use(log.ZerologMiddleware(logger)) // Add custom zerolog logging middleware
	addRoutes(router, config, logger, getRestDbClient, getPgDbClient)
	return router
}

func addRoutes(
	router *gin.Engine,
	config ServerConfig,
	logger zerolog.Logger,
	getRestDbClient supabase.DBClientFactory,
	getPgClient supabase.DBClientFactory,
) {
	// redirect anything at / to the landing page
	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/p/landing") })

	router.Static("assets", "./webpage/dist/assets")

	// pages are all served under /p prefix
	webpages := router.Group("/p")
	webpages.GET("/*filepath", func(c *gin.Context) {
		c.File("./webpage/dist/index.html")
	})

	// RPC endpoints
	// rpcs := router.Group("/r", func(c *gin.Context) { c.AbortWithStatus(http.StatusNotImplemented) })

	// API endpoints
	restApi := router.Group("/api/v1", debugMiddleware(logger))
	{
		restApi.GET("/", apiHandler(config, logger))
		testGroup := restApi.Group("/test")
		{
			testGroup.GET("bob", testHandler(config, logger, getRestDbClient))
			testGroup.GET("bill", testHandler(config, logger, getPgClient))
		}
	}

	// Add more routes here as needed
	router.NoRoute(func(c *gin.Context) {
		logger.Warn().Msg("404 Not Found")
		c.JSON(404, gin.H{"error": "Not Found"})
	})
}

func apiHandler(config ServerConfig, logger zerolog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger.Info().Msg("apiHandler called")
		logger.Warn().Msg("404 Not Found")
		c.JSON(404, gin.H{"error": "API endpoint not Found"})
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

func debugMiddleware(logger zerolog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("supasession")
		if err != nil {
			logger.Err(err).Msg("failed to get cookie")
			return
		}
		tokenString := cookie.Value
		logger.Info().Str("supabasesession token", tokenString).Send()

		keyLocation := "https://ocdegtteilykjvohsxrl.supabase.co/auth/v1/.well-known/jwks.json"
		keySet, err := jwk.Fetch(context.Background(), keyLocation)
		if err != nil {
			logger.Err(err).Msg("failed to get key set")
			return
		}

		parsedVerified, err := jwt.Parse([]byte(tokenString), jwt.WithVerify(true), jwt.WithKeySet(keySet), jwt.WithValidate(true))
		if err != nil {
			fmt.Printf("failed to parse JWT: %s\n", err)
			return
		}

		var email string
		err = parsedVerified.Get("email", &email)
		if err != nil {
			logger.Err(err).Msg("failed to get email claim")
			return
		}

		logger.Info().Str("email", email).Msg("woohoo")
	}
}
