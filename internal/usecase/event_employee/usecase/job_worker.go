package event_employee_usecase

import (
	cronjob "bookify/pkg/shared/schedules"
	"context"
)

type ICronjobEventEmployee interface {
	StartSchedulesSendQuestOfEmployeeInform(cs *cronjob.CronScheduler) error
	StopSchedulerDeadlineInform(cs *cronjob.CronScheduler) error
}

// SendQuestOfEmployeeInform send mail employee with task detail
func (e *eventEmployeeUseCase) SendQuestOfEmployeeInform(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

// StartSchedulesSendQuestOfEmployeeInform : schedules of SendQuestOfEmployeeInform
func (e *eventEmployeeUseCase) StartSchedulesSendQuestOfEmployeeInform(cs *cronjob.CronScheduler) error {
	oneMinute := cs.GenerateCronExpression(0, 0, 0, 1, 0)
	cs.AddCronJob("sendQuestOfEmployeeInform", oneMinute, e.SendQuestOfEmployeeInform)
	err := cs.RemoveJob("sendQuestOfEmployeeInform")
	if err != nil {
		return err
	}
	return nil
}

// DeadlineInform send mail employee with deadline of task detail
func (e *eventEmployeeUseCase) DeadlineInform(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

// StopSchedulerDeadlineInform : schedules of DeadlineInform
func (e *eventEmployeeUseCase) StopSchedulerDeadlineInform(cs *cronjob.CronScheduler) error {
	err := cs.RemoveJob("updateRemainingLeaveDays")
	if err != nil {
		return err
	}
	return nil
}
