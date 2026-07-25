package entities

import "github.com/google/uuid"

type Ticket struct {
	Id      string
	EventId string
	Price   float64
	Status  string
	Seat    string
	UserId  string
}

func NewTicket(eventId string, price float64, status string, seat string, userId string) *Ticket {
	id := uuid.NewString()
	return &Ticket{
		Id:      id,
		EventId: eventId,
		Price:   price,
		Status:  status,
		Seat:    seat,
		UserId:  userId,
	}
}
