package cloudinary

import (
	"bookify/internal/config"
	"github.com/cloudinary/cloudinary-go/v2"
)

func SetupCloudinary(env *config.Database) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(env.CloudinaryCloudName, env.CloudinaryAPIKey, env.CloudinaryAPISecret)
	if err != nil {
		return nil, err
	}

	return cld, nil
}
