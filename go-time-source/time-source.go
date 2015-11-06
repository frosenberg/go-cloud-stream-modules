package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/frosenberg/go-cloud-stream/api"
	"github.com/frosenberg/go-cloud-stream/stream"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

var (
	//	timeformat = kingpin.Flag("format", "The time format").Default("yyyy-MM-dd HH:mm:ss").String()
	period = kingpin.Flag("fixedDelay", "Time interval in seconds").Default("5").Int()
)

func timeSource(ch api.OutputChannel) {
	log.Println("timesource started")

	// TODO add timeformat support
	// TODO support initialDelay - delay before the first message (default: 1)
	// TODO support timeUnit - the time unit for the fixed and initial delays (default: "seconds")

	t := time.Tick(time.Duration(*period) * time.Second)
	for now := range t {
		ch.Send(api.NewTextMessage([]byte(fmt.Sprint(now))))
	}

}

func main() {
	stream.Init()
	stream.RunSource(timeSource)
	stream.Cleanup()
}
