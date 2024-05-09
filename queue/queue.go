package queue

import (
	"context"
	"encoding/json"
	"github.com/KirillShikunov/fruit-core/event"
	"github.com/KirillShikunov/fruit-core/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Manager struct {
	channel *rabbitmq.Channel
}

func (m *Manager) Post(ctx context.Context, event event.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return m.channel.Publish(ctx, event.Name(), amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (m *Manager) Listen(ctx context.Context, name string, handler func(e event.Event)) error {
	return m.channel.Listen(ctx, name, func(delivery amqp.Delivery) {
		e, err := event.CreateEvent(name)
		if err != nil {
			log.Fatalf("Error creating event: %s", err)
			return
		}

		err = e.Deserialize(delivery.Body)
		if err != nil {
			log.Fatalf("Error deserializing message: %s", err)
			return
		}
		handler(e)
	})
}

func NewManager(channel *rabbitmq.Channel) *Manager {
	return &Manager{channel}
}
