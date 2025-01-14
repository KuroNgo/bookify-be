package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionEmployee = "employee"
)

type Employee struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	FirstName      string             `bson:"first_name" json:"first_name"`
	LastName       string             `bson:"last_name" json:"last_name"`
	JobTitle       string             `bson:"job_title" json:"job_title"`
	Email          string             `bson:"email" json:"email"`
}
