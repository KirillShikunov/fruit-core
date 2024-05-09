package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Channel struct {
	ch *amqp.Channel
}

func (r *Channel) getQueue(name string) (amqp.Queue, error) {
	return r.ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *Channel) Publish(ctx context.Context, name string, publishing amqp.Publishing) error {
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

func (r *Channel) Listen(ctx context.Context, name string, handler func(amqp.Delivery)) error {
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

func NewChannel(ch *amqp.Channel) *Channel {
	return &Channel{ch}
}
