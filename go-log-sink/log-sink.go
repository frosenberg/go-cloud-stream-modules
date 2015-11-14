package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/frosenberg/go-cloud-stream/api"
	"github.com/frosenberg/go-cloud-stream/stream"
)

func logSink(input <-chan *api.Message) {
	log.Debugln("logSink started")
	for {
		s := <-input
		fmt.Println(string(s.Content))
	}
}

func main() {
	stream.Init()
	stream.RunSink(logSink)
	stream.Cleanup()
}
