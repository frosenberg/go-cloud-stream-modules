package main

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"github.com/frosenberg/go-cloud-stream/api"
	"github.com/frosenberg/go-cloud-stream/stream"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"strings"
)

var (
	// CLI arguments
	url = kingpin.Flag("url", "URL of the target resource").Required().String()
	httpMethod = kingpin.Flag("httpMethod", "HTTP method, e.g., GET, POST, PUT, DELETE").Default("GET").String()
	body = kingpin.Flag("body", "HTTP body to send along with the request").String()
	headers = kingpin.Flag("headers", " a map of HTTP headers to send along with the request").StringMap()

	// TODO how do we deal with redirects?
	client = &http.Client{ /* CheckRedirect: redirectPolicyFunc,*/ }
)

func httpclient(ch api.InputOutputChannel) {
	log.Infoln("httpclient processor started")
	log.Debugf("--url: %s", *url)
	log.Debugf("--httpMethod: %s", *httpMethod)
	log.Debugf("--body: %s", *body)
	log.Debugf("--headers: %s", *headers)

	trimmedUrl := strings.Trim(*url, "'\"") // strip because we have seen weird things on some shells

	out := ch.Receive()
	for {
		msg := <-out
		if (msg.Content != nil) {

			req, err := http.NewRequest(*httpMethod, trimmedUrl, nil)
			if err != nil {
				log.Errorf("Error while creating request %s %s", *httpMethod, trimmedUrl)
				ch.Send(api.NewTextMessage([]byte(err.Error())))
			}

			// TODO support JSON path on body to dynamically create request

			// TODO support headers
			if len(*headers) > 0 {
				log.Warnf("Header no supported yet: %s", *headers)
			}

			resp, err := client.Do(req)
			if (err != nil) {
				log.Errorf("Error while invoking HTTP %s on %s", *httpMethod, trimmedUrl)
				// TODO send JSON error message downstream
				ch.Send(api.NewTextMessage([]byte(err.Error())))
			} else {
				body, _ := ioutil.ReadAll(resp.Body)
				ch.Send(api.NewTextMessage(body))
			}
		}
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	stream.Init()
	stream.RunProcessor(httpclient)
	stream.Cleanup()
}
