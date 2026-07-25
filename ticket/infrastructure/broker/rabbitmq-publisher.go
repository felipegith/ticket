package broker

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"ticket/ticket/domain/entities"
)

type RabbitMQPublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	brokerConfig, err := NewRabbitMQConfig(url, true)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar configuração do rabbitmq: %w", err)
	}

	if err := DeclareEventsTopology(brokerConfig.channel); err != nil {
		brokerConfig.connection.Close()
		return nil, err
	}

	return &RabbitMQPublisher{
		connection: brokerConfig.connection,
		channel:    brokerConfig.channel,
	}, nil
}

func (p *RabbitMQPublisher) Handle(ctx context.Context, message *entities.Outbox) error {
	confirmation, err := p.channel.PublishWithDeferredConfirmWithContext(
		ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    message.Id,
			Type:         message.Type,
			Timestamp:    time.Now(),
			Body:         []byte(message.Content),
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
		return fmt.Errorf("mensagem %s recebeu NACK do broker", message.Id)
	}

	return nil
}

func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.connection != nil {
		return p.connection.Close()
	}
	return nil
}
