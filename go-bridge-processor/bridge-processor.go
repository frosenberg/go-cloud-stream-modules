package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/frosenberg/go-cloud-stream/api"
	"github.com/frosenberg/go-cloud-stream/stream"
)

func bridge(ch api.InputOutputChannel) {
	log.Debugf("bridge processor started")
	out := ch.Receive()
	for {
		msg := <-out
		ch.Send(&msg)
	}
}

func main() {
	stream.Init();
	stream.RunProcessor(bridge)
	stream.Cleanup()
}
