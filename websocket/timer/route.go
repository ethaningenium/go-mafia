package timer

import "github.com/gin-gonic/gin"

func RunGameRoutes(r *gin.Engine) {
	r.GET("/sub/timer", func(c *gin.Context) {
		handleConnections(c)
	})
}