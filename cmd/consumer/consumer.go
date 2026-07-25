package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ticket/config"
	"ticket/ticket/infrastructure/broker"
	"ticket/ticket/infrastructure/esearch"
)

func main() {
	cfg := config.Load()

	esClient, err := esearch.NewElasticsearchClient(cfg.ELASTICSEARCHURL)
	if err != nil {
		log.Fatalf("erro ao criar cliente elasticsearch: %v", err)
	}

	consumer, err := broker.NewRabbitMQConsumer(cfg.RABBITMQURL)
	if err != nil {
		log.Fatalf("erro ao criar consumer: %v", err)
	}
	defer consumer.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("consumer iniciado")

	if err := consumer.Consume(ctx, esClient); err != nil {
		log.Fatalf("erro no consumer: %v", err)
	}

	log.Println("consumer encerrado")
}
