package images

import (
	"bookify/internal/config"
	"bookify/pkg/interface/cloudinary"
	"bookify/pkg/interface/cloudinary/models"
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"mime/multipart"
)

func UploadImageToCloudinary(file multipart.File, filePath string, folder string, env *config.Database) (models.UploadImage, error) {
	ctx := context.Background()
	cld, err := cloudinary.SetupCloudinary(env)
	if err != nil {
		return models.UploadImage{}, err
	}

	uploadParams := uploader.UploadParams{
		PublicID: filePath,
		Folder:   folder,
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return models.UploadImage{}, err
	}

	resultRes := models.UploadImage{
		ImageURL: result.SecureURL,
		AssetID:  result.AssetID,
	}
	return resultRes, nil
}
