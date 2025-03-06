package unit_test

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

func TestGetByIDEvent(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test_e2e case
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

	// Function to clear the venue collection before each test_e2e case
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

	// Function to clear the partner collection before each test_e2e case
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

	// Function to clear the event collection before each test_e2e case
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

	ev := event_repository.NewEventRepository(database, "event")
	err = ev.CreateOne(context.Background(), mockEventInput)

	// Define test_e2e cases
	tests := []struct {
		name      string
		inputID   primitive.ObjectID
		expectErr bool
	}{
		{
			name:      "success",
			inputID:   mockEventInput.ID,
			expectErr: false,
		},
		{
			name:      "error_invalid_id",
			inputID:   primitive.NilObjectID,
			expectErr: true,
		},
	}

	// Execute test_e2e cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = ev.GetByID(context.Background(), tt.inputID)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
