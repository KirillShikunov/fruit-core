package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQ struct {
	ch *amqp.Channel
}

func (r *RabbitMQ) getQueue(name string) (amqp.Queue, error) {
	return r.ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Publish(ctx context.Context, name string, publishing amqp.Publishing) error {
	q, err := r.getQueue(name)
	if err != nil {
		return err
	}

	err = r.ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		publishing,
	)

	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

	return nil
}

func (r *RabbitMQ) Listen(ctx context.Context, name string, handler func(amqp.Delivery)) error {
	messages, err := r.ch.ConsumeWithContext(
		ctx,
		name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}
	go func() {
		for d := range messages {
			handler(d)
		}
	}()
	return nil
}

func NewRabbitMQ(ch *amqp.Channel) *RabbitMQ {
	return &RabbitMQ{ch}
}
