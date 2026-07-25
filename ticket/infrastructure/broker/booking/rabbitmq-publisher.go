package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	usecase "ticket/ticket/application/use-cases/booking"
	"ticket/ticket/infrastructure/broker"
)

var _ usecase.Handler = (*RabbitMQBookingPayedPublisher)(nil)

type RabbitMQBookingPayedPublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQBookingPayedPublisher(url string) (*RabbitMQBookingPayedPublisher, error) {

	config, err := broker.NewRabbitMQConfig(url, true)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar configuração do rabbitmq: %w", err)
	}

	if err := broker.DeclareBookingTopology(config.Channel()); err != nil {
		config.Connection().Close()
		return nil, err
	}

	return &RabbitMQBookingPayedPublisher{
		connection: config.Connection(),
		channel:    config.Channel(),
	}, nil
}

func (p *RabbitMQBookingPayedPublisher) Handle(ctx context.Context, message *usecase.BookingInputModel) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("erro ao serializar booking: %w", err)
	}

	confirmation, err := p.channel.PublishWithDeferredConfirmWithContext(
		ctx,
		broker.ExchangeBookingPayed,
		broker.RoutingKeyBookingPayed,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    message.TicketId,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("erro ao publicar no rabbitmq: %w", err)
	}

	acked, err := confirmation.WaitContext(ctx)
	if err != nil {
		return fmt.Errorf("erro aguardando confirmação do broker: %w", err)
	}
	if !acked {
		return fmt.Errorf("booking do ticket %s recebeu NACK do broker", message.TicketId)
	}

	return nil
}

func (p *RabbitMQBookingPayedPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.connection != nil {
		return p.connection.Close()
	}
	return nil
}
