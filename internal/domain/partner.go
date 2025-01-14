package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionPartner = "partner"
)

type Partner struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
	Phone string             `bson:"phone" json:"phone"`
}
