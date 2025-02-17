package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEventEmployee = "event_employee"
)

type EventEmployee struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	EventID       primitive.ObjectID `bson:"event_id" json:"event_id"`
	EmployeeID    primitive.ObjectID `bson:"employee_id" bson:"employee_id"`
	Task          string             `bson:"task" json:"task"`
	StartDate     time.Time          `bson:"start_date" json:"start_date"`
	Deadline      time.Time          `bson:"deadline" json:"deadline"`
	TaskCompleted bool               `bson:"task_completed" json:"task_completed"`
}

type EventEmployeeInput struct {
	EventID       primitive.ObjectID `bson:"event_id" json:"event_id"`
	EmployeeID    primitive.ObjectID `bson:"employee_id" bson:"employee_id"`
	Task          string             `bson:"task" json:"task"`
	StartDate     time.Time          `bson:"start_date" json:"start_date"`
	Deadline      time.Time          `bson:"deadline" json:"deadline"`
	TaskCompleted bool               `bson:"task_completed" json:"task_completed"`
}
