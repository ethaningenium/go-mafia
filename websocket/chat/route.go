package chat

import "github.com/gin-gonic/gin"

func RunWSRoutes(r *gin.Engine) {
	r.GET("/ws", func(c *gin.Context) {
		handleConnections(c)
	})
}
