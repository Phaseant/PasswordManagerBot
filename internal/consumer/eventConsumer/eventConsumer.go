package eventconsumer

import (
	"log"
	"time"

	"github.com/Phaseant/PasswordManagerBot/internal/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() {
	log.Print("Bot started")
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("Consumer: failed to fetch events: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("Consumer: failed to handle events: %s", err.Error())

			continue
		}

	}
}

func (c Consumer) handleEvents(newEvents []events.Event) error {
	for _, e := range newEvents {
		if e.Type != events.Unknown {
			log.Printf("Consumer: got event: %v", e)
			if err := c.processor.Process(e); err != nil {
				log.Printf("Consumer: failed to process event: %s", err.Error())
			}
		}
		continue
	}

	return nil
}
