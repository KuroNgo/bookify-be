package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionOrganization = "organization"
)

type Organization struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Name          string             `bson:"name" json:"name"`
	ContactPerson string             `bson:"contact_person" json:"contact_person"`
	Email         string             `bson:"email" json:"email"`
	Phone         string             `bson:"phone" json:"phone"`
}

type OrganizationInput struct {
	Name          string `bson:"name" json:"name"`
	ContactPerson string `bson:"contact_person" json:"contact_person"`
	Email         string `bson:"email" json:"email"`
	Phone         string `bson:"phone" json:"phone"`
}
