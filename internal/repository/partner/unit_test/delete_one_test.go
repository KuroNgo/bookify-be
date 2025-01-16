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

func TestDeleteOnePartner(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

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
		err = ur.DeleteOne(context.Background(), mockPartner.ID)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		err = ur.DeleteOne(context.Background(), primitive.NilObjectID)
		assert.Error(t, err)
	})
}
