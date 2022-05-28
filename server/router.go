package server

import (
	"github.com/gin-gonic/gin"
	"github.com/siprtcio/heartbeatservice/controller"
	"github.com/siprtcio/heartbeatservice/services"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "0",
		})
	})

	heartBeatService := new(services.HeartBeatService)
	heartBeatService.InitializeHeartBeatService()

	heartbeatController := new(controller.HeartbeatController)

	heartbeatController.InitializeHeartbeatController(heartBeatService)

	/* /priv/Mdr/Data */
	privPath := router.Group("/priv")
	{
		privPath.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		heartBeat := privPath.Group("HeartBeat")
		heartBeat.POST(":accountid/Heartbeat/:callid", heartbeatController.ProcessHeartBeat)
	}
	return router
}
