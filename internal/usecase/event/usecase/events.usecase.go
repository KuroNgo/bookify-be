package event_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	eventrepository "bookify/internal/repository/event/repository"
	eventtyperepository "bookify/internal/repository/event_type/repository"
	organizationrepository "bookify/internal/repository/organization/repository"
	userrepository "bookify/internal/repository/user/repository"
	venuerepository "bookify/internal/repository/venue/repository"
	"bookify/pkg/interface/cloudinary/utils/images"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/helper"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriven "go.mongodb.org/mongo-driver/mongo"
	"mime/multipart"
	"strconv"
	"time"
)

type IEventUseCase interface {
	GetByID(ctx context.Context, id string) (domain.Event, error)
	GetByTitle(ctx context.Context, title string) (domain.Event, error)
	GetByUserID(ctx context.Context, currentUser string) (domain.Event, error)
	GetByOrganizationIDAndStartTime(ctx context.Context, currentUser string, startTime string) ([]domain.Event, error)
	GetByStartTime(ctx context.Context, startTime string) ([]domain.Event, error)
	GetByStartTimePagination(ctx context.Context, startTime string, page string) ([]domain.Event, int64, int, error)
	GetAll(ctx context.Context) ([]domain.Event, error)
	GetAllPagination(ctx context.Context, page string) ([]domain.Event, int64, int, error)
	CreateOne(ctx context.Context, event *domain.EventInput) error
	CreateOneAsync(ctx context.Context, event *domain.EventInput) error
	UpdateOne(ctx context.Context, id string, event *domain.EventInput) error
	UpdateImage(ctx context.Context, id string, file *multipart.FileHeader) error
	DeleteOne(ctx context.Context, eventID string) error
}

type eventUseCase struct {
	database               *config.Database
	contextTimeout         time.Duration
	eventRepository        eventrepository.IEventRepository
	organizationRepository organizationrepository.IOrganizationRepository
	eventTypeRepository    eventtyperepository.IEventTypeRepository
	venueRepository        venuerepository.IVenueRepository
	userRepository         userrepository.IUserRepository
	client                 *mongodriven.Client
	cache                  *ristretto.Cache[string, domain.Event]
	caches                 *ristretto.Cache[string, []domain.Event]
}

// NewCache Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEvent() (*ristretto.Cache[string, domain.Event], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.Event]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCacheEvent Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEvents() (*ristretto.Cache[string, []domain.Event], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.Event]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewEventUseCase(database *config.Database, contextTimeout time.Duration, eventRepository eventrepository.IEventRepository,
	organizationRepository organizationrepository.IOrganizationRepository, eventTypeRepository eventtyperepository.IEventTypeRepository,
	venueRepository venuerepository.IVenueRepository, userRepository userrepository.IUserRepository, client *mongodriven.Client) IEventUseCase {
	cache, err := NewCacheEvent()
	if err != nil {
		panic(err)
	}

	caches, err := NewCacheEvents()
	if err != nil {
		panic(err)
	}
	return &eventUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout, eventRepository: eventRepository,
		organizationRepository: organizationRepository, eventTypeRepository: eventTypeRepository, venueRepository: venueRepository,
		userRepository: userRepository, client: client}
}

func (e eventUseCase) GetByID(ctx context.Context, id string) (domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(id)
	if found {
		return value, nil
	}

	eventID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Event{}, err
	}

	data, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return domain.Event{}, err
	}

	e.cache.Set(id, data, 1)
	e.cache.Wait()

	return data, nil
}

func (e eventUseCase) GetByTitle(ctx context.Context, title string) (domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(title)
	if found {
		return value, nil
	}

	if title == "" {
		return domain.Event{}, errors.New(constants.MsgInvalidInput)
	}

	data, err := e.eventRepository.GetByTitle(ctx, title)
	if err != nil {
		return domain.Event{}, err
	}

	e.cache.Set(title, data, 1)
	e.cache.Wait()

	return data, nil
}

func (e eventUseCase) GetByUserID(ctx context.Context, currentUser string) (domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(currentUser)
	if found {
		return value, nil
	}

	currentUserID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return domain.Event{}, err
	}

	organizationID, err := e.organizationRepository.GetByUserID(ctx, currentUserID)
	if err != nil {
		return domain.Event{}, err
	}

	data, err := e.eventRepository.GetByOrganizationID(ctx, organizationID.ID)
	if err != nil {
		return domain.Event{}, err
	}

	e.cache.Set(currentUser, data, 1)
	e.cache.Wait()

	return data, nil
}

func (e eventUseCase) GetByOrganizationIDAndStartTime(ctx context.Context, currentUser string, startTime string) ([]domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.caches.Get(currentUser + startTime)
	if found {
		return value, nil
	}

	currentUserID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return nil, err
	}

	organizationID, err := e.organizationRepository.GetByUserID(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	// Parse thời gian từ chuỗi (ISO 8601)
	parseStartTime, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, errors.New(constants.MsgInvalidInput)
	}

	// Chuyển về đầu ngày (00:00:00 UTC)
	startOfDay := time.Date(parseStartTime.Year(), parseStartTime.Month(), parseStartTime.Day(), 0, 0, 0, 0, time.UTC)

	data, err := e.eventRepository.GetByOrganizationIDAndStartTime(ctx, organizationID.ID, startOfDay)
	if err != nil {
		return nil, err
	}

	e.caches.Set(currentUser+startTime, data, 1)
	e.caches.Wait()

	return data, nil
}

func (e eventUseCase) GetByStartTime(ctx context.Context, startTime string) ([]domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.caches.Get(startTime)
	if found {
		return value, nil
	}

	// Parse thời gian từ chuỗi (ISO 8601)
	parseStartTime, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, errors.New(constants.MsgInvalidInput)
	}

	// Chuyển về đầu ngày (00:00:00 UTC)
	startOfDay := time.Date(parseStartTime.Year(), parseStartTime.Month(), parseStartTime.Day(), 0, 0, 0, 0, time.UTC)

	data, err := e.eventRepository.GetByStartTime(ctx, startOfDay)
	if err != nil {
		return nil, err
	}

	e.caches.Set(startTime, data, 1)
	e.caches.Wait()

	return data, nil
}

func (e eventUseCase) GetByStartTimePagination(ctx context.Context, startTime string, page string) ([]domain.Event, int64, int, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	//value, found := e.caches.Get(startTime + page)
	//if found {
	//	return value, nil
	//}

	layout := "20/1/2025"
	parseStartTime, err := time.Parse(layout, startTime)
	if err != nil {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}

	pageChoose, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, 0, err
	}

	if pageChoose < 1 {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}

	data, totalPage, pageCurrent, err := e.eventRepository.GetByStartTimePagination(ctx, parseStartTime, page)
	if err != nil {
		return nil, 0, 0, err
	}

	//e.caches.Set(startTime, data, 1)
	//e.caches.Wait()

	return data, totalPage, pageCurrent, nil
}

func (e eventUseCase) GetAll(ctx context.Context) ([]domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e eventUseCase) GetAllPagination(ctx context.Context, page string) ([]domain.Event, int64, int, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	pageChoose, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, 0, err
	}

	if pageChoose < 1 {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}

	data, totalPage, pageCurrent, err := e.eventRepository.GetAllPagination(ctx, page)
	if err != nil {
		return nil, 0, 0, err
	}

	return data, totalPage, pageCurrent, nil
}

func (e eventUseCase) CreateOne(ctx context.Context, event *domain.EventInput) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Check validate data
	if err := validate_data.ValidateEventInput(event); err != nil {
		return err
	}

	// Check exist data
	eventTypeData, err := e.eventTypeRepository.GetByName(ctx, event.EventTypeName)
	if err != nil {
		return err
	}

	organizationID, err := primitive.ObjectIDFromHex(event.OrganizationID)
	if err != nil {
		return err
	}

	organizationData, err := e.organizationRepository.GetByID(ctx, organizationID)
	if err != nil {
		return err
	}
	if organizationData.ID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	//venueData, err := e.venueRepository.GetByID(ctx, event.VenueID)
	//if err != nil {
	//	return err
	//}
	//if venueData.ID == primitive.NilObjectID {
	//	return errors.New(constants.MsgInvalidInput)
	//}

	parseStartTime, err := time.Parse(time.RFC3339, event.StartTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	parseEndTime, err := time.Parse(time.RFC3339, event.EndTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	eventInput := &domain.Event{
		ID:                primitive.NewObjectID(),
		EventTypeID:       eventTypeData.ID,
		OrganizationID:    organizationID,
		Title:             event.Title,
		ShortDescription:  event.ShortDescription,
		Description:       event.Description,
		ImageURL:          event.ImageURL,
		AssetURL:          event.AssetURL,
		StartTime:         parseStartTime,
		EndTime:           parseEndTime,
		Mode:              event.Mode,
		EstimatedAttendee: event.EstimatedAttendee,
		ActualAttendee:    event.ActualAttendee,
		TotalExpenditure:  event.TotalExpenditure,
	}

	err = e.eventRepository.CreateOne(ctx, eventInput)
	if err != nil {
		return err
	}

	return nil
}

func (e eventUseCase) CreateOneAsync(ctx context.Context, event *domain.EventInput) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Validate event and venue inputs
	if err := validate_data.ValidateEventInput(event); err != nil {
		return err
	}

	// Start MongoDB session for transaction
	session, err := e.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionCtx mongodriven.SessionContext) (interface{}, error) {
		// Create venue
		venueInput := &domain.Venue{
			ID:          primitive.NewObjectID(),
			Capacity:    event.Capacity,
			AddressLine: event.AddressLine,
			City:        event.City,
			Country:     event.Country,
			EventMode:   event.EventMode,
			LinkAttend:  event.LinkAttend,
			FromAttend:  event.FromAttend,
		}
		if err = e.venueRepository.CreateOne(sessionCtx, venueInput); err != nil {
			return nil, err
		}

		organizationID, err := primitive.ObjectIDFromHex(event.OrganizationID)
		if err != nil {
			return nil, err
		}

		parseStartTime, err := time.Parse(time.RFC3339, event.StartTime)
		if err != nil {
			return nil, errors.New(constants.MsgInvalidInput)
		}

		parseEndTime, err := time.Parse(time.RFC3339, event.EndTime)
		if err != nil {
			return nil, errors.New(constants.MsgInvalidInput)
		}

		eventTypeData, err := e.eventTypeRepository.GetByName(ctx, event.EventTypeName)
		if err != nil {
			return nil, errors.New(constants.MsgInvalidInput)
		}

		// Create event
		eventInput := &domain.Event{
			ID:                primitive.NewObjectID(),
			EventTypeID:       eventTypeData.ID,
			VenueID:           venueInput.ID,
			OrganizationID:    organizationID,
			Title:             event.Title,
			ShortDescription:  event.ShortDescription,
			Description:       event.Description,
			ImageURL:          event.ImageURL,
			AssetURL:          event.AssetURL,
			StartTime:         parseStartTime,
			EndTime:           parseEndTime,
			Mode:              event.Mode,
			EstimatedAttendee: event.EstimatedAttendee,
			ActualAttendee:    event.ActualAttendee,
			TotalExpenditure:  event.TotalExpenditure,
			Tags:              event.Tags,
		}
		if err := e.eventRepository.CreateOne(sessionCtx, eventInput); err != nil {
			return nil, err
		}

		return nil, nil // Successfully created event and venue
	}

	// Run transaction
	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}

	return nil // Transaction successful, no need to commit explicitly
}

func (e eventUseCase) UpdateOne(ctx context.Context, id string, event *domain.EventInput) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	eventTypeData, err := e.eventTypeRepository.GetByName(ctx, event.EventTypeName)
	if err != nil {
		return err
	}

	organizationID, err := primitive.ObjectIDFromHex(event.OrganizationID)
	if err != nil {
		return err
	}

	parseStartTime, err := time.Parse(time.RFC3339, event.StartTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	parseEndTime, err := time.Parse(time.RFC3339, event.EndTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	eventInput := &domain.Event{
		ID:                eventID,
		EventTypeID:       eventTypeData.ID,
		OrganizationID:    organizationID,
		Title:             event.Title,
		Description:       event.Description,
		ImageURL:          event.ImageURL,
		AssetURL:          event.AssetURL,
		StartTime:         parseStartTime,
		EndTime:           parseEndTime,
		Mode:              event.Mode,
		EstimatedAttendee: event.EstimatedAttendee,
		ActualAttendee:    event.ActualAttendee,
		TotalExpenditure:  event.TotalExpenditure,
		Tags:              event.Tags,
	}

	err = e.eventRepository.UpdateOne(ctx, eventInput)
	if err != nil {
		return err
	}

	return nil
}

func (e eventUseCase) UpdateImage(ctx context.Context, id string, file *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	session, err := e.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionCtx mongodriven.SessionContext) (interface{}, error) {
		if file == nil {
			return nil, errors.New("images not nil")
		}
		// Kiểm tra xem file có phải là hình ảnh không
		if !helper.IsImage(file.Filename) {
			return nil, err
		}

		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()

		imageURL, err := images.UploadImageToCloudinary(f, file.Filename, e.database.CloudinaryUploadFolderUser, e.database)
		if err != nil {
			return nil, err
		}

		// Đảm bảo xóa ảnh trên Cloudinary nếu xảy ra lỗi sau khi tải lên thành công
		defer func() {
			if err != nil {
				_, _ = images.DeleteToCloudinary(imageURL.AssetID, e.database)
			}
		}()

		eventInput := &domain.Event{
			ID:       eventID,
			ImageURL: imageURL.ImageURL,
			AssetURL: imageURL.AssetID,
		}

		err = e.eventRepository.UpdateImage(ctx, eventInput)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Run the transaction
	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}

	return session.CommitTransaction(ctx)
}

func (e eventUseCase) DeleteOne(ctx context.Context, eventID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return err
	}

	err = e.eventRepository.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
