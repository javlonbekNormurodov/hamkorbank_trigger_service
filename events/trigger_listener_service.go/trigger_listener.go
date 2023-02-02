package trigger_listener_service

import (
	"hamkorbank/config"
	"hamkorbank/pkg/logger"
	"hamkorbank/pkg/rabbitmq"
	"hamkorbank/pkg/requests"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	RecordId string `json:"record_id"`
}

type triggerListener struct {
	log        logger.LoggerI
	rabbitmq   rabbitmq.RabbitMQI
	conn       *amqp.Connection
	httpClient requests.HttpRequestI
	cfg        config.Config
}

type Response struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        Phone  `json:"data"`
}

type Phone struct {
	ID        string `json:"id"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type NotFound struct {
	NotFound string `json:"not_found"`
}

func NewTriggerListenerService(log logger.LoggerI, rabbit rabbitmq.RabbitMQI, cfg config.Config, client requests.HttpRequestI) *triggerListener {
	return &triggerListener{
		log:        log,
		rabbitmq:   rabbit,
		httpClient: client,
		cfg:        cfg,
	}
}

func (t *triggerListener) RegisterConsumers() {
	_ = t.rabbitmq.AddConsumer(config.Consumer, t.Listen)
}
