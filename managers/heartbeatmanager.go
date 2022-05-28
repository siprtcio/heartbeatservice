package managers

import (
	"encoding/json"

	request "github.com/siprtcio/heartbeatservice/requests"
	"github.com/siprtcio/heartbeatservice/services"
)

type HeartbeatManager struct {
	heartBeatService *services.HeartBeatService
}

func (hm *HeartbeatManager) InitHeartbeatManager(heartBeatService *services.HeartBeatService) {
	hm.heartBeatService = heartBeatService
}

func (hm HeartbeatManager) HeartBeatProcess(hbpReq request.HeartBeatPrivRequest) {
	data, _ := json.Marshal(hbpReq)
	hm.heartBeatService.PubishHeartBeat(data)
}
