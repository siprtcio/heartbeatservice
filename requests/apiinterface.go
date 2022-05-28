package request

type HeartBeatPrivRequest struct {
	AccountID     string  `json:"accountid"`
	CallID        string  `json:"callid"`
	Rate          float64 `json:"rate"`
	Pulse         int64   `json:"pulse"`
	Duration      int64   `json:"duration"`
	MediaServerIP string  `json:"serverip"`
	Region        string  `json:"region"`
	MsLb          string  `json:"mslb"`
	ParentCallSid string  `json:"parent_call_sid"`
}
