package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Connection struct {
	dsn string
}

func (c *Connection) GetChannel(ctx context.Context) (*amqp.Channel, error) {
	conn, err := amqp.Dial(c.dsn)

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	go func() {
		<-ctx.Done()

		err = ch.Close()
		if err != nil {
			log.Fatalf("Failed to close channel: %s", err)
		}

		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection: %s", err)
		}
	}()

	return ch, err
}

func NewConnection(dsn string) *Connection {
	return &Connection{dsn}
}
