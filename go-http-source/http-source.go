package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/frosenberg/go-cloud-stream/api"
	"github.com/frosenberg/go-cloud-stream/stream"
	"net/http"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// CLI arguments
	pathPattern = kingpin.Flag("pathPattern", "The request mapping path").Default("/messages").String()
)

func httpSource(ch api.OutputChannel) {
	log.Printf("http-source started on port %s", *stream.ServerPort)

	http.HandleFunc(*pathPattern, func (w http.ResponseWriter, r *http.Request) {
		err := ch.Send(api.NewMessageFromHttpRequest(r))
		if err != nil {
			log.Errorf("Error writing http message to transport: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%s", *stream.ServerPort), nil)
}

func main() {
	stream.Init()
	stream.RunSource(httpSource)
	stream.Cleanup()
}
