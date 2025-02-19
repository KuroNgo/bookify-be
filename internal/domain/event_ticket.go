package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEventTicket = "event_ticket"
)

type EventTicket struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	EventID   primitive.ObjectID `bson:"event_id" json:"event_id"`
	Price     float64            `bson:"price" json:"price"`
	Quantity  float64            `bson:"quantity" json:"quantity"`
	Status    string             `bson:"status" json:"status"` // Trạng thái vé - automatic (available, sold_out, expired, canceled, pending, reserved)
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type EventTicketInput struct {
	EventID  primitive.ObjectID `bson:"event_id" json:"event_id"`
	Price    float64            `bson:"price" json:"price"`
	Quantity float64            `bson:"quantity" json:"quantity"`
	Status   string             `bson:"status" json:"status"` //
}
