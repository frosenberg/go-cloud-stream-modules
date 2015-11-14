package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/frosenberg/go-cloud-stream/api"
	"github.com/frosenberg/go-cloud-stream/stream"
)

func bridge(input <-chan *api.Message, output chan<- *api.Message) {
	log.Infoln("bridge-processor started")
	for {
		msg := <-input
		output<- msg
	}
}

func main() {
	stream.Init();
	stream.RunProcessor(bridge)
	stream.Cleanup()
}
