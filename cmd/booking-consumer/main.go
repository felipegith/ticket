package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"

	"ticket/config"
	booking "ticket/ticket/application/use-cases/booking"
	brokerbooking "ticket/ticket/infrastructure/broker/booking"
	"ticket/ticket/infrastructure/cache"
	"ticket/ticket/infrastructure/persistence"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.PGCONECTIONSTRING)
	if err != nil {
		log.Fatalf("erro ao abrir conexão: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.REDISCONFIG.Host, cfg.REDISCONFIG.Port),
		Password: cfg.REDISCONFIG.Password,
		DB:       cfg.REDISCONFIG.DB,
	})
	defer redisClient.Close()
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("erro ao conectar no redis: %v", err)
	}

	ticketRepository := persistence.NewTicketRepository(db)
	bookingRepository := persistence.NewBookingRepository(db)
	eventCache := cache.NewRedisCache(redisClient)
	bookingProcess := booking.NewBookingProcess(ticketRepository, bookingRepository, eventCache)

	consumer, err := brokerbooking.NewRabbitMQBookingConsumer(cfg.RABBITMQURL)
	if err != nil {
		log.Fatalf("erro ao criar booking consumer: %v", err)
	}
	defer consumer.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("booking consumer iniciado")

	if err := consumer.Consume(ctx, bookingProcess); err != nil {
		log.Fatalf("erro no booking consumer: %v", err)
	}

	log.Println("booking consumer encerrado")
}
