package server

import (
	"fmt"
	"server/supabase"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/azoghal/succotash/webpage/log"
)

type ServerConfig struct {
	ApiKey string
	ApiUrl string
}

func NewServer(config ServerConfig, logger zerolog.Logger, getDbClient supabase.RestDBClientFactory) *gin.Engine {
	router := gin.New()                       // Create a new Gin router without default middleware
	router.Use(gin.Recovery())                // Add default recovery middleware
	router.Use(log.ZerologMiddleware(logger)) // Add custom zerolog logging middleware
	addRoutes(router, config, logger, getDbClient)
	return router
}

func addRoutes(
	router *gin.Engine,
	config ServerConfig,
	logger zerolog.Logger,
	getDbConn supabase.RestDBClientFactory,
) {
	router.Use(static.Serve("/", static.LocalFile("./webpage/dist", true)))

	// rpcs := router.Group("/r", func(c *gin.Context) { c.AbortWithStatus(http.StatusNotImplemented) })

	restApi := router.Group("/api/v1")
	{
		restApi.GET("/", apiHandler(config, logger))
		testGroup := restApi.Group("/test")
		{
			testGroup.GET("bob", testHandler(config, logger, getDbConn))
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

func testHandler(config ServerConfig, logger zerolog.Logger, dbClientGetter supabase.RestDBClientFactory) func(c *gin.Context) {
	return func(c *gin.Context) {
		client := dbClientGetter(config.ApiUrl, config.ApiKey)
		events, err := client.GetTestEvents()
		if err != nil {
			c.JSON(500, fmt.Sprintf("Failed to get events: %v", err))
		}

		for _, v := range events {
			logger.Info().Str("content", v.Content).Msg("an event")
		}
	}
}
