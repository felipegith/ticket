package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	usecase "ticket/ticket/application/use-cases/booking"
	"ticket/ticket/infrastructure/broker"
)

type Processor interface {
	Process(ctx context.Context, message *usecase.BookingInputModel) error
}

type RabbitMQBookingConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQBookingConsumer(url string) (*RabbitMQBookingConsumer, error) {

	config, err := broker.NewRabbitMQConfig(url, false)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar configuração do rabbitmq: %w", err)
	}

	if err := broker.DeclareBookingTopology(config.Channel()); err != nil {
		config.Connection().Close()
		return nil, err
	}

	return &RabbitMQBookingConsumer{
		connection: config.Connection(),
		channel:    config.Channel(),
	}, nil
}

func (c *RabbitMQBookingConsumer) Consume(ctx context.Context, processor Processor) error {
	if err := c.channel.Qos(10, 0, false); err != nil {
		return fmt.Errorf("erro ao configurar qos: %w", err)
	}

	deliveries, err := c.channel.ConsumeWithContext(
		ctx,
		broker.QueueBookingPayed,
		"booking-processor",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("erro ao registrar consumer: %w", err)
	}

	log.Println("booking consumer escutando a fila booking.payed...")

	for delivery := range deliveries {
		var message usecase.BookingInputModel
		if err := json.Unmarshal(delivery.Body, &message); err != nil {

			log.Printf("[booking] body inválido, descartando: %v", err)
			_ = delivery.Nack(false, false)
			continue
		}

		if err := processor.Process(ctx, &message); err != nil {
			log.Printf("[booking] erro ao processar %s: %v — devolvendo à fila", message.TicketId, err)
			_ = delivery.Nack(false, true)
			continue
		}

		if err := delivery.Ack(false); err != nil {
			log.Printf("[booking] erro ao dar ack: %v", err)
		}
	}

	return nil
}

func (c *RabbitMQBookingConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.connection != nil {
		c.connection.Close()
	}
}
