package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionEventType = "event_type"
)

type EventType struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type EventTypeInput struct {
	Name string `bson:"name" json:"name"`
}
