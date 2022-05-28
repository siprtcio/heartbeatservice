package main

import (
	"github.com/siprtcio/heartbeatservice/logger"
	"github.com/siprtcio/heartbeatservice/server"
)

func main() {
	logger.InitLogger()
	server.Init()
}
