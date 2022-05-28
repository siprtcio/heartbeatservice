package utils

import (
	resty "github.com/go-resty/resty/v2"
	"github.com/siprtcio/heartbeatservice/logger"
)

func RelayDeleteCallToRegionGW(rGw string, acID string, callId string, mediaServerIp string) {
	apiURL := rGw + "/api/v1/account/" + acID + "/Call/" + callId
	logger.Debug("Realy delete call to voice call routing gateway",
		logger.LogFields{"authid": acID, "callid": callId, "gateway": rGw, "url": apiURL})
	client := resty.New()
	resp, err := client.R().
		SetHeader("X-MediaServer-IP", mediaServerIp).
		Delete(apiURL)
	if err != nil {
		logger.Error(err.Error(), logger.LogFields{"authid": acID, "callid": callId, "error": err.Error()})
	} else {
		if resp.StatusCode() != 200 {
			logger.Error("Http Response",
				logger.LogFields{"authid": acID, "callid": callId, "HttpResponse": resp.StatusCode()})
		}
	}
}
