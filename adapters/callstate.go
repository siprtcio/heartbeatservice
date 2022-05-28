package adapters

type CallStateI interface {
	InitCallState() error
	GetCallState(string) ([]byte, error)
	SetCallState(string, string) error
	SetCallRoute(string, string) error
	SetCallRouteTimeout(string, string) error
	DelCallState(key string) error
	GetCallRoute(string) string
	IncrementHeartBeatCounter(string) string
	DeleteHeartBeatCounter(string)
}
