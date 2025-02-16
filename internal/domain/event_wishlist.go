package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEventWishlist = "event_wishlist"
)

type EventWishlist struct {
	ID        primitive.ObjectID `bson:"_id" json:"id" `
	EventID   primitive.ObjectID `bson:"event_id" json:"event_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Notes     string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
