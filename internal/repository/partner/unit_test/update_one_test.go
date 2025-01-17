package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	partner_repository "bookify/internal/repository/partner/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUpdateOnePartner(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)
	// Function to clear the event collection before each test case
	clearEventTypeCollection := func() {
		err := database.Collection("partner").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear partner collection: %v", err)
		}
	}

	clearEventTypeCollection()
	mockPartner := &domain.Partner{
		ID:    primitive.NewObjectID(),
		Name:  "kuro",
		Email: "kuro@gmail.com",
		Phone: "0329245971",
	}

	par := partner_repository.NewPartnerRepository(database, "partner")
	err := par.CreateOne(context.Background(), mockPartner)
	assert.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		mockPartnerUpdate := &domain.Partner{
			ID:    mockPartner.ID,
			Name:  "andrew",
			Email: "kuro@gmail.com",
			Phone: "0329245971",
		}
		err = par.UpdateOne(context.Background(), mockPartnerUpdate)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockPartnerUpdateNil := &domain.Partner{
			ID:    primitive.NilObjectID,
			Name:  "andrew",
			Email: "kuro@gmail.com",
			Phone: "0329245971",
		}
		err = par.UpdateOne(context.Background(), mockPartnerUpdateNil)
		assert.Error(t, err)
	})
}
