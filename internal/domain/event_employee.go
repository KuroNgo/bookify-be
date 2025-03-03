package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEventEmployee = "event_employee"
)

type EventEmployee struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	EventID    primitive.ObjectID `bson:"event_id" json:"event_id"`
	EmployeeID primitive.ObjectID `bson:"employee_id" bson:"employee_id"`
	Task       []Task             `bson:"task" json:"task"`
}

type Task struct {
	TaskName       string    `bson:"task_name" json:"task_name"`
	ImportantLevel int       `bson:"important_level" json:"important_level"`
	StartDate      time.Time `bson:"start_date" json:"start_date"`
	Deadline       time.Time `bson:"deadline" json:"deadline"`
	TaskCompleted  bool      `bson:"task_completed" json:"task_completed"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
	WhoCreated     string    `bson:"who_created" json:"who_created"`
}

type EventEmployeeInput struct {
	EventID        primitive.ObjectID `bson:"event_id" json:"event_id"`
	EmployeeID     primitive.ObjectID `bson:"employee_id" bson:"employee_id"`
	Task           string             `bson:"task" json:"task"`
	TaskName       string             `bson:"task_name" json:"task_name"`
	ImportantLevel int                `bson:"important_level" json:"important_level"`
	StartDate      time.Time          `bson:"start_date" json:"start_date"`
	Deadline       time.Time          `bson:"deadline" json:"deadline"`
	TaskCompleted  bool               `bson:"task_completed" json:"task_completed"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
	WhoCreated     string             `bson:"who_created" json:"who_created"`
}

type EventEmployeeResponse struct {
	EventEmployee         EventEmployee         `json:"event_employee"`
	ResultEventUnComplete ResultEventUnComplete `json:"result_event_uncompleted"`
	ResultEventComplete   ResultEventComplete   `json:"result_event_complete"`
}

type ResultEventUnComplete struct {
	TotalTasks      int     `bson:"total_tasks"`
	IncompleteTasks int     `bson:"incomplete_tasks"`
	Percentage      float64 `bson:"percentage"`
}

type ResultEventComplete struct {
	TotalTasks    int     `bson:"total_tasks"`
	CompleteTasks int     `bson:"complete_tasks"`
	Percentage    float64 `bson:"percentage"`
}
