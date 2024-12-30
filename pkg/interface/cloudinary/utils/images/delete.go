package images

import (
	"bookify/internal/config"
	"bookify/pkg/interface/cloudinary"
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func DeleteToCloudinary(assetID string, env *config.Database) (string, error) {
	ctx := context.Background()
	cld, err := cloudinary.SetupCloudinary(env)
	if err != nil {
		return "", err
	}

	deleteParams := uploader.DestroyParams{
		PublicID: assetID,
	}

	result, err := cld.Upload.Destroy(ctx, deleteParams)
	if err != nil {
		return "", err
	}

	return result.Result, nil
}
