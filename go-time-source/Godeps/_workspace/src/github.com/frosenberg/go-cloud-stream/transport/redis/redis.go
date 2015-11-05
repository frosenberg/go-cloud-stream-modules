package redis

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/frosenberg/go-cloud-stream/api"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/pubsub"
	"strings"
	"regexp"
)

//
// Basic Redis transport information
//
type RedisTransport struct {
	api.Transport
	Address string

	// Timeout for blocking receives in seconds
	Timeout int

	MaxConnections int
	pool           *pool.Pool
}

// Creates a new RedisTransport instance with
// sensible default values.
func NewRedisTransport(address string, inputBinding string, outputBinding string) *RedisTransport {

	// set some reasonable defaults
	if address == "" || address == ":6379" {
		address = "localhost:6379"
	} else  { // check if it has a port

		match, _ := regexp.MatchString("^.+:\\d+$", address)
		if !match {
			log.Debugf("Appending default redis port :6379 to %s", address)
			address = address + ":6379"
		}

	}
	if inputBinding == "" {
		inputBinding = "input"
	}
	if outputBinding == "" {
		outputBinding = "output"
	}

	transport := &RedisTransport{
		Transport: api.Transport{InputBinding: prefix(inputBinding),
			OutputBinding: prefix(outputBinding)},
		Address:        address,
		Timeout:        1,  // TODO parameterize via CLI?
		MaxConnections: 10, // TODO parameterize via CLI?
	}
	return transport
}

func (t *RedisTransport) Connect() (err error) {
	log.Debugln("Connecting to Redis server: ", t.Address)

	// create redis pool
	pool, err := pool.New("tcp", t.Address, t.MaxConnections)
	if err != nil {
		return err
	}
	t.pool = pool

	// do a ping to ensure we are connected
	conn, err := t.pool.Get()
	if err != nil {
		log.Debugln("Cannot get connection from Redis pool.")
		return err
	}

	resp := conn.Cmd("PING")
	if resp.Err != nil {
		log.Debugln("Cannot while pinging Redis.")
		return resp.Err
	}
	defer t.pool.Put(conn)

	return nil
}

// Disconnects from the Redis transport. It does not fai
// if you are not connected.
func (t *RedisTransport) Disconnect() {
	log.Debugln("Disconnecting from Redis: ", t.Address)

	// nothing to do really
}

func (t *RedisTransport) Send(m *api.Message) (err error) {
	conn, _ := t.pool.Get()
	defer t.pool.Put(conn)

	if t.isOutputTopicSemantics() {
		resp := conn.Cmd("PUBLISH", t.OutputBinding, m.ToRawByteArray())
		log.Debugln("resp (publish): ", resp)
		if resp.Err != nil {
			log.Errorf("Cannot PUBLISH on queue '%v': %v", t.OutputBinding, err)
		} else {
			log.Debugf("Published '%s' to topic '%s'\n", m.Content, t.OutputBinding)
		}
	} else {
		resp := conn.Cmd("RPUSH", t.OutputBinding, m.ToRawByteArray())
		if resp.Err != nil {
			log.Errorf("Cannot RPUSH on queue '%v': %v", t.OutputBinding, err)
		} else {
			log.Debugf("Pushed '%s' to queue '%s'\n", m.Content, t.OutputBinding)
		}

	}
	return err
}

func (t *RedisTransport) Receive() <-chan api.Message {
	out := make(chan api.Message)

	if t.isInputTopicSemantics() { // topic processing

		go func() {
			conn, _ := t.pool.Get()
			defer t.pool.Put(conn)
			psc := pubsub.NewSubClient(conn)
			psc.Subscribe(t.InputBinding)
			defer psc.Unsubscribe(t.InputBinding)

			for {
				resp := psc.Receive()
				//log.Debugln("after: ", resp)

				if resp.Err != nil {
					out <- *api.NewMessageFromRawBytes([]byte(resp.Err.Error()))
				} else {
					out <- *api.NewMessageFromRawBytes([]byte(resp.Message))
				}
			}
		}()

	} else { // queue processing

		go func() {
			conn, _ := t.pool.Get()
			defer t.pool.Put(conn)

			for {
				content, err := conn.Cmd("BRPOP", t.InputBinding, 0).List()
				if err != nil {
					log.Errorf("Cannot RPOP on '%v': %v", t.InputBinding, err)
				} else {
					//log.Debugln(content)
					out <- *api.NewMessageFromRawBytes([]byte(content[1]))
				}
			}
		}()
	}
	return out
}

func (t *RedisTransport) isOutputTopicSemantics() bool {
	return strings.HasPrefix(t.OutputBinding, "topic.")
}

func (t *RedisTransport) isInputTopicSemantics() bool {
	return strings.HasPrefix(t.InputBinding, "topic.")
}

// Set the prefix of a binding correctly as it is
// expected by the underlying transformer.

func prefix(binding string) string {
	if strings.HasPrefix(binding, "topic:") {
		return strings.Replace(binding, "topic:", "topic.", 1)
	}

	if strings.HasPrefix(binding, "queue:") {
		return strings.Replace(binding, "queue:", "queue.", 1)
	} else {
		return fmt.Sprintf("queue.%s", binding)
	}
}
