package event_discount_route

import (
	event_discount_controller "bookify/internal/api/controller/event_discount"
	"bookify/internal/config"
	"bookify/internal/domain"
	event_discount_repository "bookify/internal/repository/event_discount/repository"
	userrepository "bookify/internal/repository/user/repository"
	event_discount_usecase "bookify/internal/usecase/event_discount/usecase"
	cronjob "bookify/pkg/shared/cron"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEventDiscountRouter(env *config.Database, cs *cronjob.CronScheduler, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ev := event_discount_repository.NewEventDiscountRepository(db, domain.CollectionEventDiscount)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)

	eventDiscount := &event_discount_controller.EventDiscountController{
		EventDiscountUseCase: event_discount_usecase.NewEventDiscountUseCase(env, cs, timeout, ev, ur),
		Database:             env,
	}

	router := group.Group("/event-discounts")
	router.POST("/create", eventDiscount.CreateOne)
	router.PUT("/update", eventDiscount.UpdateOne)
	router.POST("/delete", eventDiscount.DeleteOne)
}
