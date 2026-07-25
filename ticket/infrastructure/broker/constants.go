package broker

const (
	exchangeName = "events"
	queueName    = "events.elasticsearch"
	routingKey   = "event.created"
	eventsIndex  = "events"

	deadLetterExchange = "events.dlx"
	deadLetterQueue    = "events.elasticsearch.dlq"

	ExchangeBookingPayed     = "booking"
	QueueBookingPayed        = "booking.payed"
	RoutingKeyBookingPayed   = "booking.payed"
	bookingIndexBookingPayed = "booking"
)
