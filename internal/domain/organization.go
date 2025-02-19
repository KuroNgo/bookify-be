package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionOrganization = "organization"
)

type Organization struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name          string             `bson:"name" json:"name"`
	ContactPerson string             `bson:"contact_person" json:"contact_person"`
	Email         string             `bson:"email" json:"email"`
	Phone         string             `bson:"phone" json:"phone"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type OrganizationInput struct {
	Name          string             `bson:"name" json:"name"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	ContactPerson string             `bson:"contact_person" json:"contact_person"`
	Email         string             `bson:"email" json:"email"`
	Phone         string             `bson:"phone" json:"phone"`
}
