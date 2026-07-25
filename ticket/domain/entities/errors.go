package entities

import "errors"

var ErrSeatTaken = errors.New("seat already taken")

var ErrEventNotFound = errors.New("event not found")
