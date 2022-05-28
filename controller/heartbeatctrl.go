package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siprtcio/heartbeatservice/logger"
	"github.com/siprtcio/heartbeatservice/managers"
	request "github.com/siprtcio/heartbeatservice/requests"
	"github.com/siprtcio/heartbeatservice/services"
)

type HeartbeatController struct {
	heartbeatManager *managers.HeartbeatManager
}

func (hc *HeartbeatController) InitializeHeartbeatController(heartBeatService *services.HeartBeatService) {
	hc.heartbeatManager = new(managers.HeartbeatManager)
	hc.heartbeatManager.InitHeartbeatManager(heartBeatService)
}

func (hc HeartbeatController) ProcessHeartBeat(c *gin.Context) {
	accountID := c.Param("accountid")
	callID := c.Param("callid")
	logger.Debug("ProcessHeartBeat", logger.LogFields{"authid": accountID, "callid": callID})
	cr := request.HeartBeatPrivRequest{}
	if err := c.BindJSON(&cr); err == nil {
		hc.heartbeatManager.HeartBeatProcess(cr)
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "all is well!!"})
	} else {
		logger.Error("ProcessHeartBeat", logger.LogFields{"authid": accountID, "callid": callID, "request": cr})
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Incorrect JSON Request"})
	}
}
