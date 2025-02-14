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

	invalidate := true
	deleteParams := uploader.DestroyParams{
		PublicID:     assetID,
		Type:         "upload",
		ResourceType: "image",
		Invalidate:   &invalidate,
	}

	result, err := cld.Upload.Destroy(ctx, deleteParams)
	if err != nil {
		return "", err
	}

	return result.Result, nil
}
