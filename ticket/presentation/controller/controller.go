package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	booking "ticket/ticket/application/use-cases/booking"
	event "ticket/ticket/application/use-cases/event"
	ticket "ticket/ticket/application/use-cases/ticket"
	"ticket/ticket/domain/entities"
	"time"
)

type Dependencies struct {
	Create                *event.CreateEvent
	GetEvents             *event.GetEvents
	GetEventsWithKeyWords *event.GetEventsKeyword
	GetEventById          *event.GetEventById
	CreateTicket          *ticket.CreateTicket
	UnavailableSeat       *booking.Unavailable
	BookingPayed          *booking.BookingPayed
}

type UseCase struct {
	create                *event.CreateEvent
	getEvents             *event.GetEvents
	getEventsWithKeyWords *event.GetEventsKeyword
	getEventById          *event.GetEventById
	createTicket          *ticket.CreateTicket
	unavailableSeat       *booking.Unavailable
	bookingPayed          *booking.BookingPayed
}

func NewController(deps Dependencies) *UseCase {
	return &UseCase{
		create:                deps.Create,
		getEvents:             deps.GetEvents,
		getEventsWithKeyWords: deps.GetEventsWithKeyWords,
		getEventById:          deps.GetEventById,
		createTicket:          deps.CreateTicket,
		unavailableSeat:       deps.UnavailableSeat,
		bookingPayed:          deps.BookingPayed,
	}
}

type createEventRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Date        time.Time `json:"date"`
}

type createTicketRequest struct {
	Price   float64 `json:"price"`
	Status  string  `json:"status"`
	Seat    string  `json:"seat"`
	UserId  string  `json:"userId"`
	EventId string  `json:"eventId"`
}

type createEventResponse struct {
	Id string `json:"id"`
}
type createTicketResponse struct {
	Id string `json:"id"`
}

type unavailableSeatResponse struct {
	TicketId string `json:"ticketId"`
	UserId   string `json:"userId"`
	Seat     string `json:"seat"`
}

type bookingPayedRequest struct {
	TicketId string `json:"ticketId"`
	Seat     string `json:"seat"`
	UserId   string `json:"userId"`
}

func (c *UseCase) CreateEvent(writter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var input createEventRequest

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		http.Error(writter, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := c.create.Execute(request.Context(), event.EventInputModel{
		Name:        input.Name,
		Description: input.Description,
		Status:      input.Status,
		Date:        input.Date,
	})

	if err != nil {
		log.Printf("erro ao criar evento: %v", err)
		http.Error(writter, "falha ao criar evento", http.StatusInternalServerError)
		return
	}

	writter.Header().Set("Content-Type", "application/json")
	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(createEventResponse{Id: id})
}

func (c *UseCase) GetEvents(writter http.ResponseWriter, request *http.Request) {
	events, err := c.getEvents.Execute(request.Context())
	if err != nil {
		http.Error(writter, "falha ao buscar eventos", http.StatusInternalServerError)
		return
	}

	writter.Header().Set("Content-Type", "application/json")
	writter.WriteHeader(http.StatusOK)
	json.NewEncoder(writter).Encode(events)
}

func (c *UseCase) GetEventsWithKeyWords(writter http.ResponseWriter, request *http.Request) {
	keyword := request.URL.Query().Get("keyword")
	events, err := c.getEventsWithKeyWords.Execute(request.Context(), "events", keyword)
	if err != nil {
		http.Error(writter, "falha ao buscar eventos com palavras chave", http.StatusInternalServerError)
		return
	}

	writter.Header().Set("Content-Type", "application/json")
	writter.WriteHeader(http.StatusOK)
	json.NewEncoder(writter).Encode(events)
}

func (c *UseCase) GetEventById(writter http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	event, err := c.getEventById.Execute(request.Context(), id)
	if err != nil {
		if errors.Is(err, entities.ErrEventNotFound) {
			http.Error(writter, "evento não encontrado", http.StatusNotFound)
			return
		}
		log.Printf("erro ao buscar evento por id: %v", err)
		http.Error(writter, "falha ao buscar evento", http.StatusInternalServerError)
		return
	}

	writter.Header().Set("Content-Type", "application/json")
	writter.WriteHeader(http.StatusOK)
	json.NewEncoder(writter).Encode(event)
}

func (c *UseCase) CreateTicket(writter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var input createTicketRequest

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		http.Error(writter, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := c.createTicket.Execute(request.Context(), ticket.TicketInputModel{
		Seat:    input.Seat,
		UserId:  input.UserId,
		Status:  input.Status,
		EventId: input.EventId,
		Price:   input.Price,
	})

	if err != nil {
		log.Printf("erro ao criar ticket: %v", err)
		if errors.Is(err, entities.ErrSeatTaken) {
			http.Error(writter, "assento já ocupado", http.StatusConflict)
			return
		}
		http.Error(writter, "falha ao criar ticket", http.StatusInternalServerError)
		return
	}

	writter.Header().Set("Content-Type", "application/json")
	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(createTicketResponse{Id: response})
}
func (c *UseCase) UnavailableSeat(writter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var input unavailableSeatResponse

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		http.Error(writter, err.Error(), http.StatusBadRequest)
		return
	}
	err := c.unavailableSeat.Execute(request.Context(), booking.UnavailableInputModel{
		Seat:     input.Seat,
		UserId:   input.UserId,
		TicketId: input.TicketId,
	})

	if err != nil {
		http.Error(writter, err.Error(), http.StatusBadRequest)
		return
	}
}

func (c *UseCase) BookingPayed(writter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var input bookingPayedRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		http.Error(writter, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.bookingPayed.Execute(request.Context(), booking.BookingInputModel{
		TicketId: input.TicketId,
		Seat:     input.Seat,
		UserId:   input.UserId,
	}); err != nil {
		log.Printf("erro ao processar pagamento: %v", err)
		http.Error(writter, "falha ao processar pagamento", http.StatusInternalServerError)
		return
	}

	writter.WriteHeader(http.StatusAccepted)
}
