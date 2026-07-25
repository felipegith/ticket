package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"

	"ticket/config"
	booking "ticket/ticket/application/use-cases/booking"
	event "ticket/ticket/application/use-cases/event"
	ticket "ticket/ticket/application/use-cases/ticket"
	"ticket/ticket/infrastructure/broker"
	brokerbooking "ticket/ticket/infrastructure/broker/booking"
	"ticket/ticket/infrastructure/cache"
	"ticket/ticket/infrastructure/esearch"
	"ticket/ticket/infrastructure/persistence"
	"ticket/ticket/infrastructure/worker"
	"ticket/ticket/presentation/controller"
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

	eventRepository := persistence.NewEventRepository(db)
	ticketRepository := persistence.NewTicketRepository(db)
	eventCache := cache.NewRedisCache(redisClient)
	eleastiSearchClient, err := esearch.NewElasticsearchClient(cfg.ELASTICSEARCHURL)
	if err != nil {
		log.Fatalf("erro ao criar cliente Elasticsearch: %v", err)
	}

	createEvent := event.NewEvent(eventRepository)
	getEvents := event.NewGetEvents(eventRepository, eventCache, 60*time.Second)
	getEventsWithKeyWords := event.NewGetEventsKeyword(eleastiSearchClient)
	getEventById := event.NewGetEventById(eventRepository, eventCache)
	createTicket := ticket.NewTicket(ticketRepository)
	unavalilable := booking.NewUnavailable(eventCache)

	bookingPublisher, err := brokerbooking.NewRabbitMQBookingPayedPublisher(cfg.RABBITMQURL)
	if err != nil {
		log.Fatalf("erro ao criar publisher de booking: %v", err)
	}
	defer bookingPublisher.Close()
	bookingPayed := booking.NewBookingPayed(bookingPublisher)

	inputController := controller.NewController(controller.Dependencies{
		Create:                createEvent,
		GetEvents:             getEvents,
		GetEventsWithKeyWords: getEventsWithKeyWords,
		GetEventById:          getEventById,
		CreateTicket:          createTicket,
		UnavailableSeat:       unavalilable,
		BookingPayed:          bookingPayed,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	publisher, err := broker.NewRabbitMQPublisher(cfg.RABBITMQURL)
	if err != nil {
		log.Fatalf("erro ao criar publisher do rabbitmq: %v", err)
	}
	defer publisher.Close()

	outboxRepository := persistence.NewOutboxRepository(db)
	outboxWorker := worker.NewOutboxWorker(
		outboxRepository,
		publisher,
		1*time.Second,
		100,
	)
	go outboxWorker.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /events", inputController.CreateEvent)
	mux.HandleFunc("GET /events", inputController.GetEvents)
	mux.HandleFunc("GET /events-keyword", inputController.GetEventsWithKeyWords)
	mux.HandleFunc("GET /events/{id}", inputController.GetEventById)
	mux.HandleFunc("POST /tickets", inputController.CreateTicket)
	mux.HandleFunc("POST /booking-unavailable", inputController.UnavailableSeat)
	mux.HandleFunc("POST /booking-payed", inputController.BookingPayed)

	log.Println("servidor rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("erro ao subir o servidor: %v", err)
	}
}
