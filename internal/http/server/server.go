package server

import (
	"github.com/gin-gonic/gin"

	"am-kafka-project/internal/http/handlers"
	"am-kafka-project/internal/kafka/producer"
	"am-kafka-project/internal/middleware"
)

// StartHTTPServer - function that starts the HTTP server.
// It takes a producer service as an argument and returns an error.
func StartHTTPServer(producer producer.ProducerService) (err error) {
	router := gin.New() // Create a new router
	if gin.IsDebugging() {
		router.Use(middleware.LoggerMiddleware()) // If debugging is enabled, use the logger middleware
	} else {
		router.Use(middleware.JsonLoggerMiddleware()) // If debugging is not enabled, use the JSON logger middleware
	}
	router.Use(gin.Recovery()) // Use the recovery middleware

	api := router.Group("/api")                            // Create a new group for the API
	apiv1 := api.Group("/v1")                              // Create a new group for the API version 1
	apiv1.POST("/am_alerts", handlers.NewAlerts(producer)) // Handle POST requests to /am_alerts with the NewAlerts handler

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		}) // Handle GET requests to /health with a JSON response
	})

	return router.Run() // Start the router and return any errors
}
