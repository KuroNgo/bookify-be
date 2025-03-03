package activity_log_usecase

import (
	"bookify/pkg/shared/mail/handles"
	cronjob "bookify/pkg/shared/schedules"
	"context"
	"fmt"
	"log"
)

type IJobWorkerActivityLogs interface {
	SendInformForAdminToExpireTimeActivityLog30DaysStart(ctx context.Context) error
	SendInformForAdminToExpireTimeActivityLog3DaysStart(ctx context.Context) error
	RemoveActivityLog(ctx context.Context) error
	JobWorkerSendInformForAdminToExpireTimeActivityLog30DaysStart(cs *cronjob.CronScheduler) error
	JobWorkerSendInformForAdminToExpireTimeActivityLog3DaysStart(cs *cronjob.CronScheduler) error
	JobWorkerRemoveActivityLog(cs *cronjob.CronScheduler) error
	RemoveJobWorkerSendInformForAdminToExpireTimeActivityLog30DaysStart(cs *cronjob.CronScheduler) error
	RemoveJobWorkerSendInformForAdminToExpireTimeActivityLog3DaysStart(cs *cronjob.CronScheduler) error
	RemoveJobRemoveActivityLog(cs *cronjob.CronScheduler) error
}

func (a *activityUseCase) SendInformForAdminToExpireTimeActivityLog30DaysStart(ctx context.Context) error {
	emailData := handles.EmailData{
		FullName: "Administrator",
		Subject:  "[Bookify] - Log Expiration Notice",
		Email:    "hoaiphong01012002@gmail.com",
	}

	if err := handles.SendEmail(&emailData, emailData.Email, "expiring_log_30days.log.html"); err != nil {
		log.Printf("Failed to send email: %v", err)
	}

	return nil
}

func (a *activityUseCase) SendInformForAdminToExpireTimeActivityLog3DaysStart(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a *activityUseCase) RemoveActivityLog(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a *activityUseCase) JobWorkerSendInformForAdminToExpireTimeActivityLog30DaysStart(cs *cronjob.CronScheduler) error {
	cronExpression := cs.GenerateCronExpression(0, 0, 12, 1, 0)
	log.Println("=== Job is running ===") // Log để kiểm tra job có chạy không
	cs.AddCronJob("sendActivityExpiringEmails", cronExpression, func(ctx context.Context) error {
		err := a.SendInformForAdminToExpireTimeActivityLog30DaysStart(ctx)
		if err != nil {
			_ = a.RemoveJobWorkerSendInformForAdminToExpireTimeActivityLog30DaysStart(cs)
			log.Println("Job send discount emails failed:", err)
		}
		return err // Quan trọng: Trả về lỗi để job có thể xử lý nếu cần
	})

	fmt.Print("create job worker success")
	return nil
}

func (a *activityUseCase) JobWorkerSendInformForAdminToExpireTimeActivityLog3DaysStart(cs *cronjob.CronScheduler) error {
	//TODO implement me
	panic("implement me")
}

func (a *activityUseCase) JobWorkerRemoveActivityLog(cs *cronjob.CronScheduler) error {
	//TODO implement me
	panic("implement me")
}

func (a *activityUseCase) RemoveJobWorkerSendInformForAdminToExpireTimeActivityLog30DaysStart(cs *cronjob.CronScheduler) error {
	err := cs.RemoveJob("sendActivityExpiringEmails")
	if err != nil {
		return err
	}

	return nil
}

func (a *activityUseCase) RemoveJobWorkerSendInformForAdminToExpireTimeActivityLog3DaysStart(cs *cronjob.CronScheduler) error {
	//TODO implement me
	panic("implement me")
}

func (a *activityUseCase) RemoveJobRemoveActivityLog(cs *cronjob.CronScheduler) error {
	//TODO implement me
	panic("implement me")
}
