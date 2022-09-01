package model

import "time"

type EventType int16

const (
	EventType_ProductCreated = EventType(1)
	EventType_ProductUpdated = EventType(2)
	EventType_ProductDeleted = EventType(3)
)

type EventStatus int16

const (
	EventStatus_Created = EventStatus(1)
	EventStatus_Locked  = EventStatus(2)
)

type ProductEvent struct {
	ID          uint64
	EventType   EventType
	EventStatus EventStatus
	OccuredAt   time.Time
	Data        Product
}
