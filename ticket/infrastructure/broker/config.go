package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQConfig(url string, enableConfirms bool) (*RabbitMQConfig, error) {
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no rabbitmq: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("erro ao abrir canal: %w", err)
	}

	if enableConfirms {
		if err := channel.Confirm(false); err != nil {
			channel.Close()
			connection.Close()
			return nil, fmt.Errorf("erro ao ativar publisher confirms: %w", err)
		}
	}

	return &RabbitMQConfig{
		connection: connection,
		channel:    channel,
	}, nil
}

func (c *RabbitMQConfig) Channel() *amqp.Channel { return c.channel }

func (c *RabbitMQConfig) Connection() *amqp.Connection { return c.connection }

func DeclareEventsTopology(channel *amqp.Channel) error {
	if err := channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("erro ao declarar exchange: %w", err)
	}

	if err := channel.ExchangeDeclare(deadLetterExchange, "fanout", true, false, false, false, nil); err != nil {
		return fmt.Errorf("erro ao declarar dead-letter exchange: %w", err)
	}

	if _, err := channel.QueueDeclare(deadLetterQueue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("erro ao declarar dead-letter queue: %w", err)
	}

	if err := channel.QueueBind(deadLetterQueue, "", deadLetterExchange, false, nil); err != nil {
		return fmt.Errorf("erro ao fazer bind da dead-letter queue: %w", err)
	}

	mainQueueArgs := amqp.Table{"x-dead-letter-exchange": deadLetterExchange}
	if _, err := channel.QueueDeclare(queueName, true, false, false, false, mainQueueArgs); err != nil {
		return fmt.Errorf("erro ao declarar fila: %w", err)
	}

	if err := channel.QueueBind(queueName, routingKey, exchangeName, false, nil); err != nil {
		return fmt.Errorf("erro ao fazer bind da fila: %w", err)
	}

	return nil
}

func DeclareBookingTopology(channel *amqp.Channel) error {
	if err := channel.ExchangeDeclare(ExchangeBookingPayed, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("erro ao declarar exchange de booking: %w", err)
	}

	if _, err := channel.QueueDeclare(QueueBookingPayed, true, false, false, false, nil); err != nil {
		return fmt.Errorf("erro ao declarar fila de booking: %w", err)
	}

	if err := channel.QueueBind(QueueBookingPayed, RoutingKeyBookingPayed, ExchangeBookingPayed, false, nil); err != nil {
		return fmt.Errorf("erro ao fazer bind da fila de booking: %w", err)
	}

	return nil
}
