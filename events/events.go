package events

import (
	"context"
	"hamkorbank/config"
	"hamkorbank/events/trigger_listener_service"
	"hamkorbank/pkg/logger"
	"hamkorbank/pkg/rabbitmq"
	"hamkorbank/pkg/requests"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PubSubServer struct {
	cfg      config.Config
	rabbitmq rabbitmq.RabbitMQI
	log      logger.LoggerI
}

func NewEvents(cfg config.Config, log logger.LoggerI, ch *amqp.Channel) (*PubSubServer, error) {
	rabbit, err := rabbitmq.NewRabbitMQ(cfg, ch)
	if err != nil {
		return nil, err
	}

	initPublishers(rabbit)

	return &PubSubServer{
		cfg:      cfg,
		log:      log,
		rabbitmq: rabbit,
	}, nil
}

func (s *PubSubServer) InitServices(ctx context.Context, cfg config.Config, httpClient requests.HttpRequestI) {
	triggerListenerService := trigger_listener_service.NewTriggerListenerService(s.log, s.rabbitmq, cfg, httpClient)
	triggerListenerService.RegisterConsumers()
	s.rabbitmq.RunConsumers(ctx)
}

func initPublishers(rabbit rabbitmq.RabbitMQI) {
	_ = rabbit.AddPublisher("logger")
}
