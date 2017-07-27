package listener

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/azdlowry/simpleudpbeat/config"
	"github.com/pquerna/ffjson/ffjson"
)

type Listener struct {
	config             config.Config
	logEntriesRecieved chan common.MapStr
	logEntriesError    chan bool
}

func NewListener(cfg config.Config) *Listener {
	listener := &Listener{
		config: cfg,
	}
	return listener
}

func (listener *Listener) Start(logEntriesRecieved chan common.MapStr, logEntriesError chan bool) {

	listener.logEntriesRecieved = logEntriesRecieved
	listener.logEntriesError = logEntriesError

	address := fmt.Sprintf("%s:%d", listener.config.Address, listener.config.Port)

	listener.startUDP(listener.config.Protocol, address)
}

func (listener *Listener) startUDP(proto string, address string) {
	l, err := net.ListenPacket(proto, address)

	if err != nil {
		logp.Err("Error listening on % socket via %s: %v", listener.config.Protocol, address, err.Error())
		listener.logEntriesError <- true
		return
	}
	defer l.Close()

	logp.Info("Now listening for logs via %s on %s", listener.config.Protocol, address)
	buffer := make([]byte, listener.config.MaxMsgSize)

	for {
		length, _, err := l.ReadFrom(buffer)
		if err != nil {
			logp.Err("Error reading from buffer: %v", err.Error())
			continue
		}
		go listener.processMessage(buffer, length)
	}
}

func (listener *Listener) Shutdown() {
	close(listener.logEntriesError)
	close(listener.logEntriesRecieved)
}

func (listener *Listener) processMessage(buffer []byte, length int) {

	if length == 0 {
		return
	}

	logData := strings.TrimSpace(string(buffer[:length]))
	if logData == "" {
		logp.Err("Event is empty")
		return
	}
	event := common.MapStr{}

	if err := ffjson.Unmarshal([]byte(logData), &event); err != nil {
		logp.Err("Could not parse JSON: %v", err)
		event["message"] = logData
		event["tags"] = []string{"_protologbeat_json_parse_failure"}
	}

	if _, ok := event["type"]; !ok {
		event["type"] = listener.config.DefaultEsLogType
	}

	event["@timestamp"] = common.Time(time.Now())

	listener.logEntriesRecieved <- event
}
