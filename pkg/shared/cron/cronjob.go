package cron

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

func init() {
	c := NewCronScheduler()
	c.Start()
}

type Scheduler struct {
	c        *cron.Cron
	entryMap map[string]cron.EntryID
}

func NewCronScheduler() *Scheduler {
	return &Scheduler{
		c:        cron.New(),
		entryMap: make(map[string]cron.EntryID),
	}
}

func (cs *Scheduler) Start() {
	// Bắt đầu cron scheduler
	log.Println("Starting cron scheduler...")
	cs.c.Start()
}

func (cs *Scheduler) AddCronJob(name, cronExpression string, taskFunc func(ctx context.Context) error) cron.EntryID {
	entryID, err := cs.c.AddFunc(cronExpression, func() {
		err := taskFunc(context.Background())
		if err != nil {
			log.Printf("Error in cron job: %v", err)
		}
	})

	if err != nil {
		log.Printf("Error adding cron job: %v", err)
	}

	cs.entryMap[name] = entryID
	return entryID
}

func (cs *Scheduler) RemoveJob(name string) error {
	if entryID, exists := cs.entryMap[name]; exists {
		cs.c.Remove(entryID)
		delete(cs.entryMap, name)
		log.Printf("Removed job: %s", name)
		return nil
	}
	return fmt.Errorf("job '%s' not found", name)
}

func (cs *Scheduler) GetJobCount() int {
	return len(cs.entryMap)
}

func (cs *Scheduler) GenerateCronExpression(day, month, hour, minute, dayOfWeek int) string {
	log.Printf("Implemented job worker")
	return fmt.Sprintf("%d %d %d %d %d", minute, hour, day, month, dayOfWeek)
}
