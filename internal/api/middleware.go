package api

import (
	"github.com/gin-gonic/gin"
	"tradingdata/internal/config"
)

func configMiddleware(global config.GlobalInstance) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.GIN_ENV_GLOBAL_INSTANCE, global)
		c.Next()
	}
}