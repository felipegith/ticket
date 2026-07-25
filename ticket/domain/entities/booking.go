package entities

import "github.com/google/uuid"

type Booking struct {
	Id       string
	TicketId string
	UserId   string
}

func NewBooking(name, ticketId, userId string) *Booking {
	id := uuid.NewString()
	return &Booking{
		Id:       id,
		TicketId: ticketId,
		UserId:   userId,
	}
}
