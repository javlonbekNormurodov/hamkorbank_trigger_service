package rabbitmq

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	name string
}

func (r RabbitMQ) AddPublisher(name string) error {
	if _, ok := r.publishers[name]; ok {
		panic(errors.New("consumer with the same name already exists: " + name))
	}

	err := r.channel.ExchangeDeclare(
		"logs",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		panic(err)
	}

	r.publishers[name] = &Publisher{
		name: name,
	}

	return nil
}

func (r RabbitMQ) Publish(ctx context.Context, name string, data []byte) error {

	err := r.channel.PublishWithContext(ctx,
		"logs", // exchange
		name,   // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})

	if err != nil {
		panic(err)
	}

	return err
}
