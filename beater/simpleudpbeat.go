package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/azdlowry/simpleudpbeat/config"
	"github.com/azdlowry/simpleudpbeat/listener"
)

type Simpleudpbeat struct {
	done     chan struct{}
	config   config.Config
	client   publisher.Client
	listener *listener.Listener
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Simpleudpbeat{
		done:     make(chan struct{}),
		config:   config,
		listener: listener.NewListener(config),
	}
	return bt, nil
}

func (bt *Simpleudpbeat) Run(b *beat.Beat) error {
	logp.Info("simpleudpbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()

	logEntriesRecieved := make(chan common.MapStr, 100000)
	logEntriesErrors := make(chan bool, 1)

	go func(logs chan common.MapStr, errs chan bool) {
		bt.listener.Start(logs, errs)
	}(logEntriesRecieved, logEntriesErrors)

	var event common.MapStr

	for {
		select {
		case <-bt.done:
			return nil
		case <-logEntriesErrors:
			return nil
		case event = <-logEntriesRecieved:
			if event == nil {
				return nil
			}
			bt.client.PublishEvent(event)
			logp.Info("Event sent")
		}
	}
}

func (bt *Simpleudpbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
