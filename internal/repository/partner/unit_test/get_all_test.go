package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	partnerrepository "bookify/internal/repository/partner/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindAllPartner(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Mock data
	mockPartner := &domain.Partner{
		ID:    primitive.NewObjectID(),
		Name:  "kuro",
		Email: "kuro@gmail.com",
		Phone: "0329245971",
	}

	par := partnerrepository.NewPartnerRepository(database, "partner")
	err := par.CreateOne(context.Background(), mockPartner)
	assert.Nil(t, err)

	// Define test_e2e cases
	tests := []struct {
		name        string
		expectedErr bool
		description string
	}{
		{
			name:        "success_find_all_partners",
			expectedErr: false,
			description: "Should successfully fetch all partners",
		},
	}

	// Execute test_e2e cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := par.GetAll(context.Background())

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
