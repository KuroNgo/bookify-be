package activity_log_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	activitylogrepository "bookify/internal/repository/activity_log/repository"
	userrepository "bookify/internal/repository/user/repository"
	"context"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IActivityUseCase interface {
	GetByID(ctx context.Context, id string) (domain.ActivityLog, error)
	GetByLevel(ctx context.Context, id string) (domain.ActivityLog, error)
	GetByUserID(ctx context.Context, id string) (domain.ActivityLog, error)
	GetAll(ctx context.Context) ([]domain.ActivityLog, error)
	CreateOne(ctx context.Context, activity *domain.ActivityLogInput, currentUser string) error
	DeleteOne(ctx context.Context, id string) error
}

type activityUseCase struct {
	database           *config.Database
	contextTimeout     time.Duration
	activityRepository activitylogrepository.IActivityLogRepository
	userRepository     userrepository.IUserRepository
	cache              *ristretto.Cache[string, domain.Employee]
	caches             *ristretto.Cache[string, []domain.Employee]
}

func NewActivityUseCase(database *config.Database, contextTimeout time.Duration, activityRepository activitylogrepository.IActivityLogRepository,
	userRepository userrepository.IUserRepository) IActivityUseCase {
	return &activityUseCase{database: database, contextTimeout: contextTimeout, activityRepository: activityRepository, userRepository: userRepository}
}

func (a *activityUseCase) GetByID(ctx context.Context, id string) (domain.ActivityLog, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ActivityLog{}, err
	}

	data, err := a.activityRepository.GetByID(ctx, employeeID)
	if err != nil {
		return domain.ActivityLog{}, err
	}

	return data, nil
}

func (a *activityUseCase) GetByLevel(ctx context.Context, id string) (domain.ActivityLog, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ActivityLog{}, err
	}

	data, err := a.activityRepository.GetByID(ctx, employeeID)
	if err != nil {
		return domain.ActivityLog{}, err
	}

	return data, nil
}

func (a *activityUseCase) GetByUserID(ctx context.Context, id string) (domain.ActivityLog, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ActivityLog{}, err
	}

	data, err := a.activityRepository.GetByID(ctx, employeeID)
	if err != nil {
		return domain.ActivityLog{}, err
	}

	return data, nil
}

func (a *activityUseCase) GetAll(ctx context.Context) ([]domain.ActivityLog, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	data, err := a.activityRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (a *activityUseCase) CreateOne(ctx context.Context, activity *domain.ActivityLogInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	employeeInput := &domain.ActivityLog{
		ID:           primitive.NewObjectID(),
		ClientIP:     activity.ClientIP,
		UserID:       userID,
		Level:        activity.Level,
		Method:       activity.Method,
		StatusCode:   activity.StatusCode,
		BodySize:     activity.BodySize,
		Path:         activity.Path,
		Latency:      activity.Latency,
		Error:        activity.Error,
		ActivityTime: activity.ActivityTime,
		ExpireAt:     activity.ExpireAt,
	}

	err = a.activityRepository.CreateOne(ctx, employeeInput)
	if err != nil {
		return err
	}

	return nil
}

func (a *activityUseCase) DeleteOne(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	activityID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = a.activityRepository.DeleteOne(ctx, activityID)
	if err != nil {
		return err
	}

	return nil
}
