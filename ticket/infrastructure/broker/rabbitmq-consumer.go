package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const retryBackoff = 2 * time.Second

type Indexer interface {
	IndexDocument(ctx context.Context, index, documentID string, document []byte) error
}

type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQConsumer(url string) (*RabbitMQConsumer, error) {
	brokerConfig, err := NewRabbitMQConfig(url, false)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar configuração do rabbitmq: %w", err)
	}

	if err := DeclareEventsTopology(brokerConfig.channel); err != nil {
		brokerConfig.connection.Close()
		return nil, err
	}

	return &RabbitMQConsumer{
		connection: brokerConfig.connection,
		channel:    brokerConfig.channel,
	}, nil
}

func (c *RabbitMQConsumer) Consume(ctx context.Context, indexer Indexer) error {

	if err := c.channel.Qos(10, 0, false); err != nil {
		return fmt.Errorf("erro ao configurar qos: %w", err)
	}

	deliveries, err := c.channel.ConsumeWithContext(
		ctx,
		queueName,
		"es-indexer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("erro ao registrar consumer: %w", err)
	}

	log.Println("consumer escutando a fila, aguardando mensagens...")

	for delivery := range deliveries {

		var payload struct {
			Id string `json:"Id"`
		}
		if err := json.Unmarshal(delivery.Body, &payload); err != nil {

			log.Printf("[consumer] body inválido (poison), enviando para DLQ: %v", err)
			_ = delivery.Nack(false, false)
			continue
		}

		if err := indexer.IndexDocument(ctx, eventsIndex, payload.Id, delivery.Body); err != nil {

			log.Printf("[consumer] erro ao indexar evento %s: %v — retentando em %s", payload.Id, err, retryBackoff)
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(retryBackoff):
			}
			if nackErr := delivery.Nack(false, true); nackErr != nil {
				log.Printf("[consumer] erro ao dar nack: %v", nackErr)
			}
			continue
		}

		log.Printf("[consumer] evento %s indexado no ES", payload.Id)
		if err := delivery.Ack(false); err != nil {
			log.Printf("[consumer] erro ao dar ack: %v", err)
		}
	}

	return nil
}

func (c *RabbitMQConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.connection != nil {
		c.connection.Close()
	}
}
