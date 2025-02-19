package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

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
	Status         string             `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
	WhoCreated     string             `bson:"who_created_at" json:"who_created_at"`
}

type EmployeeInput struct {
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	FirstName      string             `bson:"first_name" json:"first_name"`
	LastName       string             `bson:"last_name" json:"last_name"`
	JobTitle       string             `bson:"job_title" json:"job_title"`
	Email          string             `bson:"email" json:"email"`
}

type EmployeeResponse struct {
	Employee []Employee `json:"employee"`
}
