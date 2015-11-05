package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/frosenberg/go-cloud-stream/api"
	"github.com/frosenberg/go-cloud-stream/stream"
)

func logSink(ch api.InputChannel) {
	log.Debugln("logSink started")

	out := ch.Receive()
	for {
		s := <-out
		fmt.Println(s)
	}
}

func main() {
	stream.Init()
	stream.RunSink(logSink)
	stream.Cleanup()
}
