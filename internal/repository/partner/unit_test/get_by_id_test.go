package unit_test

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	partner_repository "bookify/internal/repository/partner/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindByIDPartner(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test case
	clearPartnerCollection := func() {
		err := database.Collection("partner").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear partner collection: %v", err)
		}
	}

	clearPartnerCollection()
	mockPartner := &domain.Partner{
		ID:    primitive.NewObjectID(),
		Name:  "kuro",
		Email: "kuro@gmail.com",
		Phone: "0329245971",
	}
	ur := partner_repository.NewPartnerRepository(database, "partner")
	err := ur.CreateOne(context.Background(), mockPartner)
	assert.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		_, err = ur.GetByID(context.Background(), mockPartner.ID)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		_, err = ur.GetByID(context.Background(), primitive.NilObjectID)
		assert.Error(t, err)
	})
}
