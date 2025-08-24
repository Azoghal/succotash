package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/azoghal/succotash/webpage/log"
)

type ServerConfig struct {
	ApiKey string
	ApiUrl string
}

func NewServer(config ServerConfig, logger zerolog.Logger, getDbConn func() error) *gin.Engine {
	router := gin.New()                       // Create a new Gin router without default middleware
	router.Use(gin.Recovery())                // Add default recovery middleware
	router.Use(log.ZerologMiddleware(logger)) // Add custom zerolog logging middleware
	addRoutes(router, config, logger, getDbConn)
	return router
}

func addRoutes(
	router *gin.Engine,
	config ServerConfig,
	logger zerolog.Logger,
	getDbConn func() error,
) {
	api := router.Group("/api/v1")
	{
		api.GET("/", apiHandler(config, logger))
	}

	// beercamGroup := api.Group("/beercam")
	// {
	// 	beercamGroup.GET("/", beercamHandler(config, logger))
	// 	beercamGroup.GET("/person/:person_name", personHandler(config, logger, getDbConn))
	// }

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
