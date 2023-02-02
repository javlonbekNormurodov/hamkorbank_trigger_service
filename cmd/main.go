package main

import (
	"context"
	"fmt"
	"hamkorbank/config"
	"hamkorbank/events"
	"hamkorbank/pkg/logger"
	"hamkorbank/pkg/requests"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	cfg := config.Load()
	log := logger.NewLogger(cfg.LogLevel, "trigger_listener_service")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer ch.Close()

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic(err)
	}

	httpClient := requests.NewHttpClient("", 10)

	pubSubServer, err := events.NewEvents(cfg, log, ch)
	if err != nil {
		log.Error("error on event server")
	}

	ctx := context.Background()
	pubSubServer.InitServices(ctx, cfg, httpClient) // it should run forever if there is any consumer
}
