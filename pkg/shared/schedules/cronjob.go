package cronjob

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

type CronScheduler struct {
	c        *cron.Cron
	entryMap map[string]cron.EntryID
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		c:        cron.New(),
		entryMap: make(map[string]cron.EntryID),
	}
}

func (cs *CronScheduler) Start() {
	// Bắt đầu schedules scheduler
	log.Println("Starting schedules scheduler...")
	cs.c.Start()
}

func (cs *CronScheduler) AddCronJob(name, cronExpression string, taskFunc func(ctx context.Context) error) cron.EntryID {
	entryID, err := cs.c.AddFunc(cronExpression, func() {
		err := taskFunc(context.Background())
		if err != nil {
			log.Printf("Error in schedules job: %v", err)
		}
	})

	if err != nil {
		log.Printf("Error adding schedules job: %v", err)
	}

	cs.entryMap[name] = entryID
	return entryID
}

func (cs *CronScheduler) RemoveJob(name string) error {
	if entryID, exists := cs.entryMap[name]; exists {
		cs.c.Remove(entryID)
		delete(cs.entryMap, name)
		log.Printf("Removed job: %s", name)
		return nil
	}
	return fmt.Errorf("job '%s' not found", name)
}

func (cs *CronScheduler) GetJobCount() int {
	return len(cs.entryMap)
}

func (cs *CronScheduler) GenerateCronExpression(day, month, hour, minute, dayOfWeek int) string {
	if day < 1 {
		day = 1
	}
	if month < 1 {
		month = 1
	}

	log.Printf("Implemented job worker")
	return fmt.Sprintf("%d %d %d %d %d", minute, hour, day, month, dayOfWeek)
}
