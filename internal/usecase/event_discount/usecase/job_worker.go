package event_discount_usecase

import (
	cronjob "bookify/pkg/shared/cron"
	"bookify/pkg/shared/helper"
	"bookify/pkg/shared/mail/handles"
	"context"
	"log"
	"time"
)

type IJobWorkerEventDiscount interface {
	SendDiscountForApplicableUsersIfTheyHaveWishlist(ctx context.Context) error
	JobWorkerSendDiscountForApplicableUsersIfTheyHaveWishlist(cs *cronjob.CronScheduler) error
}

func (e *eventDiscountUseCase) SendDiscountForApplicableUsersIfTheyHaveWishlist(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Lấy danh sách tất cả users có wishlist
	users, err := e.wishlistRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		userID := user.ID

		wishlistData, err := e.wishlistRepository.GetByUserID(ctx, userID)
		if err != nil {
			continue // Bỏ qua user lỗi, tiếp tục xử lý user tiếp theo
		}

		eventDiscountData, err := e.eventDiscountRepository.GetByUserIDInApplicableAndEventID(ctx, userID, wishlistData.EventID)
		if err != nil || helper.IsZeroValue(eventDiscountData) {
			continue
		}

		userData, err := e.userRepository.GetByID(ctx, userID)
		if err != nil {
			continue
		}

		emailData := handles.EmailData{
			FullName:   userData.FullName,
			Subject:    "Exclusive Discount Just for You!",
			ExpireDate: eventDiscountData.EndDate.Format(time.UnixDate),
		}

		// Gửi email giảm giá
		go func(email string, data handles.EmailData) {
			err := handles.SendEmail(&data, email, "receive_exclusive.discount.html")
			if err != nil {
				// Ghi log lỗi thay vì return
				log.Println("Failed to send discount email:", err)
			}
		}(userData.Email, emailData)
	}

	return nil
}

func (e *eventDiscountUseCase) JobWorkerSendDiscountForApplicableUsersIfTheyHaveWishlist(cs *cronjob.CronScheduler) error {
	cronExpression := cs.GenerateCronExpression(0, 0, 12, 1, 0)

	cs.AddCronJob("sendDiscountEmails", cronExpression, func(ctx context.Context) error {
		err := e.SendDiscountForApplicableUsersIfTheyHaveWishlist(ctx)
		if err != nil {
			log.Println("Job send discount emails failed:", err)
		}
		return err // Quan trọng: Trả về lỗi để job có thể xử lý nếu cần
	})

	return nil
}
