package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/siprtcio/heartbeatservice/adapters"
	"github.com/siprtcio/heartbeatservice/logger"
	"github.com/siprtcio/heartbeatservice/models"
	request "github.com/siprtcio/heartbeatservice/requests"
	util "github.com/siprtcio/heartbeatservice/utils"

	"github.com/astaxie/beego"
	"github.com/nsqio/go-nsq"
)

type HeartBeatService struct {
	nsqHeartBeatServe *NSQServices
	cs                *adapters.CallState
}

var balanceServiceUrl = beego.AppConfig.String("balanceservice_url")

func (hbs *HeartBeatService) handleHeartBeatMessage() nsq.HandlerFunc {
	return nsq.HandlerFunc(func(m *nsq.Message) error {
		if len(m.Body) == 0 {
			return errors.New("body is blank re-enqueue message")
		}
		heartbeat := request.HeartBeatPrivRequest{}
		err := json.Unmarshal(m.Body, &heartbeat)

		if err != nil {
			logger.Error("handleHeartBeatMessage", logger.LogFields{
				"authid":    heartbeat.AccountID,
				"callid":    heartbeat.CallID,
				"heartbeat": heartbeat,
			})
		}

		logger.Debug("handleHeartBeatMessage", logger.LogFields{
			"authid":    heartbeat.AccountID,
			"callid":    heartbeat.CallID,
			"heartbeat": heartbeat,
		})

		// set call route in redis common instance - used by webapi to delete calls.
		hbs.cs.SetCallRouteTimeout(heartbeat.CallID, heartbeat.MsLb)

		if heartbeat.Rate == 0 {
			return err
		}

		bsReq := &models.BalanceServiceRequest{Balance: -heartbeat.Rate}

		url := fmt.Sprintf("%s/v1/Accounts/%s/Balances", balanceServiceUrl, heartbeat.AccountID)
		client := resty.New()
		resp, err := client.R().
			EnableTrace().
			SetHeader("Accept", "application/json").
			SetBody(bsReq).
			Patch(url)

		if err != nil || resp.StatusCode() != 200 {
			logger.Error("handleHeartBeatMessage", logger.LogFields{
				"authid":  heartbeat.AccountID,
				"callid":  heartbeat.CallID,
				"message": "balance amount fetch failed",
				"error":   err,
			})
			_, _ = hbs.deleteCall(heartbeat.AccountID, heartbeat.CallID)
		}
		return err
	})
}

func (hbs *HeartBeatService) SubscribeHeartBeatReceiver() {
	handler := hbs.handleHeartBeatMessage()
	_ = hbs.nsqHeartBeatServe.nsqSubscribe(beego.AppConfig.String("nsqlookupd_url"), "heartbeat_billing", "heartbeat_billingchannel", handler)
}

func (hbs *HeartBeatService) PubishHeartBeat(message []byte) {
	_ = hbs.nsqHeartBeatServe.nsqPublish("heartbeat_billing", message)
}

func (hbs *HeartBeatService) deleteCall(authID string, callUUID string) (string, error) {
	callState := new(request.HeartBeatPrivRequest)

	callStatKey := fmt.Sprintf("callstate:api::%s", callUUID)

	if callStats, err := hbs.cs.GetCallState(callStatKey); err == nil {

		_ = json.Unmarshal(callStats, callState)

		go util.RelayDeleteCallToRegionGW(callState.MsLb, authID, callUUID, callState.MediaServerIP)

		hbs.cs.DeleteCallRoute(callStatKey)
	}
	return "call deleted", nil
}

func (hbs *HeartBeatService) setCallState(callUUID string, callState request.HeartBeatPrivRequest) {

	callStatKey := fmt.Sprintf("callstate:api::%s", callUUID)
	route, err := hbs.cs.GetCallState(callStatKey)

	if err == nil && len(route) < 100 {
		return
	}

	callStateByte, _ := json.Marshal(callState)
	_ = hbs.cs.SetCallState(callStatKey, callStateByte)

	//we need to insert call state
	return
}

func (hbs *HeartBeatService) InitializeHeartBeatService() {
	hbs.nsqHeartBeatServe = new(NSQServices)
	hbs.cs = new(adapters.CallState)
	_ = hbs.cs.InitCallState()
	_ = hbs.nsqHeartBeatServe.initNSQServices(beego.AppConfig.String("nsqd_url"), 10, 1)
	go hbs.SubscribeHeartBeatReceiver()
}
