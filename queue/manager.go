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
	rabbitMQ *rabbitmq.RabbitMQ
}

func (b *Manager) Post(ctx context.Context, event event.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return b.rabbitMQ.Publish(ctx, event.Name(), amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (b *Manager) Listen(ctx context.Context, name string, handler func(e event.Event)) error {
	factory := event.NewFactory()
	return b.rabbitMQ.Listen(ctx, name, func(delivery amqp.Delivery) {
		e, err := factory.Create(name)
		if err != nil {
			log.Printf("Error creating event: %s", err)
			return
		}

		err = e.Deserialize(delivery.Body)
		if err != nil {
			log.Printf("Error deserializing message: %s", err)
			return
		}
		handler(e)
	})
}

func NewManager(rabbitMQ *rabbitmq.RabbitMQ) *Manager {
	return &Manager{rabbitMQ}
}
