package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionEmployeePartner = "employee_partner"
)

type EventPartner struct {
	EventID   primitive.ObjectID `bson:"event_id" json:"event_id"`
	PartnerID primitive.ObjectID `bson:"partner_id" json:"partner_id"`
	Role      string             `bson:"role" json:"role"`
}
