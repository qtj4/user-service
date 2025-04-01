package event

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type EventPublisher struct {
	ch *amqp.Channel
}

func NewEventPublisher(ch *amqp.Channel) *EventPublisher {
	return &EventPublisher{ch: ch}
}

func (p *EventPublisher) Publish(eventType string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = p.ch.Publish(
		"",        // exchange
		eventType, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Published event: %s", eventType)
	return nil
}