package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionEventPartner = "event_partner"
)

type EventPartner struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	EventID   primitive.ObjectID `bson:"event_id" json:"event_id"`
	PartnerID primitive.ObjectID `bson:"partner_id" json:"partner_id"`
	Role      string             `bson:"role" json:"role"`
}

type EventPartnerInput struct {
	EventID   string `bson:"event_id" json:"event_id"`
	PartnerID string `bson:"partner_id" json:"partner_id"`
	Role      string `bson:"role" json:"role"`
}
