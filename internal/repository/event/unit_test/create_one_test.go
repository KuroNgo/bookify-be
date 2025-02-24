package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	event_repository "bookify/internal/repository/event/repository"
	event_type_repository "bookify/internal/repository/event_type/repository"
	organizationrepository "bookify/internal/repository/organization/repository"
	venuerepository "bookify/internal/repository/venue/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestCreateOneEvent(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test case
	clearEventTypeCollection := func() {
		err := database.Collection("event_type").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event_type collection: %v", err)
		}
	}
	clearEventTypeCollection()

	mockEventType := domain.EventType{
		ID:   primitive.NewObjectID(),
		Name: "music",
	}
	et := event_type_repository.NewEventTypeRepository(database, "event_type")
	err := et.CreateOne(context.Background(), mockEventType)
	assert.Nil(t, err)

	// Function to clear the venue collection before each test case
	clearVenueCollection := func() {
		err := database.Collection("venue").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear venue collection: %v", err)
		}
	}
	clearVenueCollection()

	mockVenue := &domain.Venue{
		ID:          primitive.NewObjectID(),
		Capacity:    100,
		AddressLine: "123 Main Street",
		City:        "New York",
		State:       "NY",
		Country:     "USA",
		PostalCode:  "10001",
		OnlineFlat:  false,
	}
	ve := venuerepository.NewVenueRepository(database, "venue")
	err = ve.CreateOne(context.Background(), mockVenue)
	assert.Nil(t, err)

	// Function to clear the partner collection before each test case
	clearOrganizationCollection := func() {
		err := database.Collection("organization").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear organization collection: %v", err)
		}
	}
	clearOrganizationCollection()

	// Mock data
	mockOrganization := &domain.Organization{
		ID:            primitive.NewObjectID(),
		Name:          "Tech Corp",
		ContactPerson: "John Doe",
		Email:         "john.doe@techcorp.com",
		Phone:         "0329245971",
	}
	or := organizationrepository.NewOrganizationRepository(database, "organization")
	err = or.CreateOne(context.Background(), mockOrganization)
	assert.Nil(t, err)

	// Function to clear the event collection before each test case
	clearEventCollection := func() {
		err := database.Collection("event").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event collection: %v", err)
		}
	}
	clearEventCollection()

	// Mock data
	mockEventInput := &domain.Event{
		ID:                primitive.NewObjectID(),
		EventTypeID:       mockEventType.ID,
		VenueID:           mockVenue.ID,
		OrganizationID:    mockOrganization.ID,
		Title:             "Music Concert 2025",
		Description:       "A live music concert featuring popular artists.",
		ImageURL:          "https://example.com/image.jpg",
		AssetURL:          "https://example.com/asset.mp4",
		StartTime:         time.Date(2025, time.March, 10, 19, 0, 0, 0, time.UTC),
		EndTime:           time.Date(2025, time.March, 10, 22, 0, 0, 0, time.UTC),
		Mode:              "Public",
		EstimatedAttendee: 500,
		ActualAttendee:    450,
		TotalExpenditure:  15000.50,
	}
	mockEventInputNil := &domain.Event{}

	// Define test cases
	tests := []struct {
		name        string
		input       *domain.Event
		expectedErr bool
		description string
	}{
		{
			name:        "success_create_event",
			input:       mockEventInput,
			expectedErr: false,
			description: "Should successfully create a event",
		},
		{
			name:        "error_create_event_with_nil",
			input:       mockEventInputNil,
			expectedErr: true,
			description: "Should return error when creating a event with nil fields",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := event_repository.NewEventRepository(database, "event")
			err = ev.CreateOne(context.Background(), tt.input)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
