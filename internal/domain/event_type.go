package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionEventType = "event_type"
)

type EventType struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	EventTypeName string             `bson:"event_type_name" json:"event_type_name"`
}
