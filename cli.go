package main

import (
	"gopkg.in/op/go-logging.v1"
	"os"
	"time"
)

var log = &logging.Logger{}
var backend = logging.NewLogBackend(os.Stdout, "", 0)

func main() {
	backend1Leveled := logging.AddModuleLevel(backend)
	backend1Leveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(backend1Leveled)

	serverConfig := defaultServerConfig
	serverConfig.portRangeStart = 4000
	serverConfig.portRangeEnd = 4100
	serverConfig.networks = []network{udp4}
	newServer(&serverConfig).start()
	clientConfig := clientConfig{
		host:            "127.0.0.1",
		portRangeStart:  serverConfig.portRangeStart,
		portRangeEnd:    serverConfig.portRangeEnd,
		networks:        serverConfig.networks,
		waitTime:        2 * time.Second,
		parallelisation: 1,
	}
	fails, _ := run(&clientConfig)
	if len(fails) == 0 {
		log.Info("All ports have passed the test")
	}
	for k, v := range fails {
		log.Infof("Network %s fails: %v\n", k, v)
	}
}
