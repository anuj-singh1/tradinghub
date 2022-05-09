package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tradingdata/internal/config"
)

// Routes with some middleware
func Routes(router *gin.Engine, global config.GlobalInstance) {
	router.Use(configMiddleware(global))
	router.GET("/status", ping)
	router.NoRoute(notFound)
	// cron route middleware
	router.GET("/login", login)
	router.GET("/get_authcode_url", getAuthCodeUrl)
}

// not found route
func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  http.StatusNotFound,
		"message": "Route Not Found",
	})
}