package user_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	"bookify/internal/repository/user/repository"
	"bookify/pkg/interface/cloudinary/utils/images"
	google_oauth2 "bookify/pkg/interface/oauth2/google"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/helper"
	"bookify/pkg/shared/mail/handles"
	"bookify/pkg/shared/password"
	"bookify/pkg/shared/token"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo_driven "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"mime/multipart"
	"sync"
	"time"
)

type IUserUseCase interface {
	FetchMany(ctx context.Context) ([]domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByIDForCheckCookie(ctx context.Context, accessToken string) (domain.User, error)
	GetByID(ctx context.Context, idUser string) (domain.User, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (domain.User, error)

	UpdateOne(ctx context.Context, userID string, input *domain.InputUser, file *multipart.FileHeader) error
	UpdateUserInfoOne(ctx context.Context, userID string, input *domain.UpdateUserInfo, file *multipart.FileHeader) error
	UpdateVerify(ctx context.Context, id string, input *domain.InputUser) error
	UpdateImage(ctx context.Context, id string, file *multipart.FileHeader) error
	UpdateSocialMedia(ctx context.Context, userID string, userSocial *domain.UpdateSocialMedia) error
	UpdateProfile(ctx context.Context, userID string, userProfile *domain.UpdateUserSettings, file *multipart.FileHeader) (string, error)
	UpdateProfileNotImage(ctx context.Context, userID string, userProfile *domain.UpdateUserSettings) error

	SignUp(ctx context.Context, input *domain.SignupUser) error
	LoginUser(ctx context.Context, signIn *domain.SignIn) (domain.OutputLogin, error)
	LoginGoogle(ctx context.Context, code string) (*domain.User, *domain.OutputLoginGoogle, error)
	DeleteOne(ctx context.Context, idUser string) error
	RefreshToken(ctx context.Context, refreshToken string) (*domain.OutputLogin, error)

	ForgetPassword(ctx context.Context, email string) error
	UpdateVerifyForChangePassword(ctx context.Context, verificationCode string) error
	UpdatePassword(ctx context.Context, id string, input *domain.ChangePasswordInput) error
}

type userUseCase struct {
	database       *config.Database
	contextTimeout time.Duration
	userRepository repository.IUserRepository
	mu             *sync.Mutex
	cache          *ristretto.Cache[string, domain.User]
	caches         *ristretto.Cache[string, []domain.User]
	client         *mongo_driven.Client
}

// NewCache Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCache() (*ristretto.Cache[string, domain.User], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.User]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewCaches() (*ristretto.Cache[string, []domain.User], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.User]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewUserUseCase(database *config.Database, contextTimeout time.Duration, userRepository repository.IUserRepository,
	client *mongo_driven.Client) IUserUseCase {
	cache, err := NewCache()
	if err != nil {
		panic(err)
	}

	caches, err := NewCaches()
	if err != nil {
		panic(err)
	}
	return &userUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout, userRepository: userRepository, client: client}
}

func (u *userUseCase) FetchMany(ctx context.Context) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := u.caches.Get("users")
	if found {
		return value, nil
	}

	userData, err := u.userRepository.FetchMany(ctx)
	if err != nil {
		return nil, err
	}

	var outputs []domain.User
	outputs = make([]domain.User, 0, len(userData))
	for _, user := range userData {
		output := domain.User{
			ID:          user.ID,
			Email:       user.Email,
			Phone:       user.Phone,
			FullName:    user.FullName,
			Gender:      user.Gender,
			Vocation:    user.Vocation,
			Address:     user.Address,
			City:        user.City,
			Region:      user.Region,
			DateOfBirth: user.DateOfBirth,
			AssetURL:    user.AssetURL,
			AvatarURL:   user.AvatarURL,
			Role:        user.Role,
		}

		outputs = append(outputs, output)
	}

	u.caches.Set("users", outputs, 1)
	// wait for value to pass through buffers
	u.caches.Wait()

	return outputs, nil
}

func (u *userUseCase) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := u.cache.Get(email)
	if found {
		return value, nil
	}

	userData, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	output := domain.User{
		ID:          userData.ID,
		Email:       userData.Email,
		Phone:       userData.Phone,
		FullName:    userData.FullName,
		Gender:      userData.Gender,
		Vocation:    userData.Vocation,
		Address:     userData.Address,
		City:        userData.City,
		Region:      userData.Region,
		DateOfBirth: userData.DateOfBirth,
		AssetURL:    userData.AssetURL,
		AvatarURL:   userData.AvatarURL,
		Role:        userData.Role,
	}

	u.cache.Set(email, userData, 1)
	// wait for value to pass through buffers
	u.cache.Wait()

	return output, nil
}

func (u *userUseCase) GetByIDForCheckCookie(ctx context.Context, accessToken string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	sub, err := token.ValidateToken(accessToken, u.database.AccessTokenPublicKey)
	if err != nil {
		return domain.User{}, err
	}

	// get value from cache
	value, found := u.cache.Get(fmt.Sprint(sub))
	if found {
		return value, nil
	}

	userID, err := primitive.ObjectIDFromHex(fmt.Sprint(sub))
	if err != nil {
		return domain.User{}, err
	}

	userData, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}

	u.cache.Set(fmt.Sprint(sub), userData, 1)
	// wait for value to pass through buffers
	u.cache.Wait()

	return userData, nil
}

func (u *userUseCase) GetByID(ctx context.Context, idUser string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := u.cache.Get(idUser)
	if found {
		return value, nil
	}

	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return domain.User{}, err
	}

	userData, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}

	u.cache.Set(idUser, userData, 1)
	// wait for value to pass through buffers
	u.cache.Wait()

	return userData, nil
}

func (u *userUseCase) GetByVerificationCode(ctx context.Context, verificationCode string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := u.cache.Get(verificationCode)
	if found {
		return value, nil
	}

	user, err := u.userRepository.GetByVerificationCode(ctx, verificationCode)
	if err != nil {
		return domain.User{}, err
	}

	updUser := domain.User{
		ID:        user.ID,
		Verified:  true,
		UpdatedAt: time.Now(),
	}

	// Update User in Database
	err = u.userRepository.UpdateVerificationCode(ctx, &updUser)
	if err != nil {
		return domain.User{}, err
	}

	response := domain.User{
		ID:          user.ID,
		Email:       user.Email,
		Phone:       user.Phone,
		FullName:    user.FullName,
		Gender:      user.Gender,
		Vocation:    user.Vocation,
		Address:     user.Address,
		City:        user.City,
		Region:      user.Region,
		DateOfBirth: user.DateOfBirth,
		AssetURL:    user.AssetURL,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
	}

	u.cache.Set(verificationCode, response, 1)
	// wait for value to pass through buffers
	u.cache.Wait()

	return response, nil
}

func (u *userUseCase) UpdateOne(ctx context.Context, userID string, input *domain.InputUser, file *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	userData, err := u.userRepository.GetByID(ctx, idUser)
	if err != nil {
		return err
	}

	if file == nil {
		user := domain.User{
			ID: userData.ID,
			//Username:  input.Username,
			UpdatedAt: time.Now(),
		}

		err = u.userRepository.UpdateOne(ctx, &user)
		if err != nil {
			return err
		}

		return nil
	}

	if !helper.IsImage(file.Filename) {
		return err
	}

	f, err := file.Open()
	if err != nil {
		return err
	}
	defer func(f multipart.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)

	user := domain.User{
		ID: idUser,
		//Username:  input.Username,
		UpdatedAt: time.Now(),
	}

	err = u.userRepository.UpdateOne(ctx, &user)
	if err != nil {
		return err
	}
	u.cache.Clear()

	return nil
}

func (u *userUseCase) UpdateSocialMedia(ctx context.Context, userID string, userSocial *domain.UpdateSocialMedia) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	userData, err := u.userRepository.GetByID(ctx, idUser)
	if err != nil {
		return err
	}

	user := domain.User{
		ID:                             userData.ID,
		FacebookSc:                     userSocial.FacebookSc,
		InstagramSc:                    userSocial.InstagramSc,
		LinkedInSc:                     userSocial.LinkedInSc,
		YoutubeSc:                      userSocial.YoutubeSc,
		EnableAutomaticSharingOfEvents: userSocial.EnableAutomaticSharingOfEvents,
		EnableSharingOn:                userSocial.EnableSharingOn,
	}

	err = u.userRepository.UpdateSocialMedia(ctx, &user)
	if err != nil {
		return err
	}
	u.cache.Clear()

	return nil
}

func (u *userUseCase) UpdateProfile(ctx context.Context, userID string, userProfile *domain.UpdateUserSettings, file *multipart.FileHeader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate input
	if err := validate_data.ValidateUser4(userProfile); err != nil {
		return "", err
	}

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	userData, err := u.userRepository.GetByID(ctx, idUser)
	if err != nil {
		return "", err
	}

	// Nếu không có file, chỉ cập nhật thông tin người dùng
	if file == nil {
		parseDateOfBirth, err := time.Parse(time.RFC3339, userProfile.DateOfBirth)
		if err != nil {
			return "", errors.New(constants.MsgInvalidInput)
		}

		user := domain.User{
			ID:           userData.ID,
			Gender:       userProfile.Gender,
			Vocation:     userProfile.Vocation,
			Address:      userProfile.Address,
			City:         userProfile.City,
			Region:       userProfile.Region,
			DateOfBirth:  parseDateOfBirth,
			FullName:     userProfile.FullName,
			ShowInterest: userProfile.ShowInterest,
			SocialMedia:  userProfile.SocialMedia,
		}

		if err := u.userRepository.UpdateProfile(ctx, &user); err != nil {
			return "", err
		}
		return "", nil
	}

	// Kiểm tra xem file có phải ảnh hợp lệ không
	if !helper.IsImage(file.Filename) {
		return "", errors.New("invalid image format")
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Xóa ảnh cũ nếu có
	if userData.AssetURL != "" {
		result, err := images.DeleteToCloudinary(userData.AvatarURL, u.database)
		if err != nil {
			log.Println("Failed to delete old image from Cloudinary:", err)
		} else {
			log.Println("Deleted old image from Cloudinary successfully" + result)
		}
	}

	// Upload ảnh mới lên Cloudinary
	imageURL, err := images.UploadImageToCloudinary(f, file.Filename, u.database.CloudinaryUploadFolderUser, u.database)
	if err != nil {
		log.Println("Failed to upload image to Cloudinary:", err)
		return "", err
	}

	// Parse ngày sinh
	parseDateOfBirth, err := time.Parse(time.RFC3339, userProfile.DateOfBirth)
	if err != nil {
		return "", errors.New(constants.MsgInvalidInput)
	}

	user := domain.User{
		ID:           userData.ID,
		Gender:       userProfile.Gender,
		Vocation:     userProfile.Vocation,
		Address:      userProfile.Address,
		City:         userProfile.City,
		Region:       userProfile.Region,
		DateOfBirth:  parseDateOfBirth,
		FullName:     userProfile.FullName,
		AvatarURL:    imageURL.ImageURL,
		AssetURL:     imageURL.AssetID,
		ShowInterest: userProfile.ShowInterest,
		SocialMedia:  userProfile.SocialMedia,
	}

	// Cập nhật thông tin người dùng
	if err = u.userRepository.UpdateProfile(ctx, &user); err != nil {
		log.Println("Failed to update user profile, rolling back image upload:", err)
		_, _ = images.DeleteToCloudinary(imageURL.ImageURL, u.database) // Xóa ảnh mới nếu có lỗi
		return "", err
	}

	u.cache.Clear()

	return imageURL.ImageURL, nil
}

func (u *userUseCase) UpdateProfileNotImage(ctx context.Context, userID string, userProfile *domain.UpdateUserSettings) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if err := validate_data.ValidateUser4(userProfile); err != nil {
		return err

	}

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	userData, err := u.userRepository.GetByID(ctx, idUser)
	if err != nil {
		return err

	}

	parseDateOfBirth, err := time.Parse(time.RFC3339, userProfile.DateOfBirth)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	user := domain.User{
		ID:           userData.ID,
		Gender:       userProfile.Gender,
		Vocation:     userProfile.Vocation,
		Address:      userProfile.Address,
		City:         userProfile.City,
		Region:       userProfile.Region,
		DateOfBirth:  parseDateOfBirth,
		FullName:     userProfile.FullName,
		ShowInterest: userProfile.ShowInterest,
		SocialMedia:  userProfile.SocialMedia,
	}

	err = u.userRepository.UpdateProfileNotImage(ctx, &user)
	if err != nil {
		return err
	}
	u.cache.Clear()

	return nil
}

func (u *userUseCase) UpdateUserInfoOne(ctx context.Context, userID string, input *domain.UpdateUserInfo, file *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	userData, err := u.userRepository.GetByID(ctx, idUser)
	if err != nil {
		return err
	}

	layout := "31/02/2002 16:06:00"
	parseDateOfBirth, err := time.Parse(layout, input.DateOfBirth)
	if err != nil {
		return err
	}

	if file == nil {
		user := domain.User{
			ID:          userData.ID,
			FullName:    input.FullName,
			Gender:      input.Gender,
			Vocation:    input.Vocation,
			Address:     input.Address,
			City:        input.City,
			Region:      input.Region,
			DateOfBirth: parseDateOfBirth,
			AvatarURL:   input.AvatarURL,
			AssetURL:    input.AssetURL,
			UpdatedAt:   time.Now(),
		}

		err = u.userRepository.UpdateOne(ctx, &user)
		if err != nil {
			return err
		}

		return nil
	}

	if !helper.IsImage(file.Filename) {
		return err
	}

	f, err := file.Open()
	if err != nil {
		return err
	}
	defer func(f multipart.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)

	user := domain.User{
		ID:          userData.ID,
		FullName:    input.FullName,
		Gender:      input.Gender,
		Vocation:    input.Vocation,
		Address:     input.Address,
		City:        input.City,
		Region:      input.Region,
		DateOfBirth: parseDateOfBirth,
		UpdatedAt:   time.Now(),
	}

	err = u.userRepository.UpdateOne(ctx, &user)
	if err != nil {
		return err
	}

	u.cache.Clear()

	return nil
}

func (u *userUseCase) UpdateVerify(ctx context.Context, id string, input *domain.InputUser) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	idUser, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err = validate_data.ValidateUser(input); err != nil {
		return err
	}

	user := domain.User{
		ID:               idUser,
		VerificationCode: input.VerificationCode,
		Verified:         false,
		UpdatedAt:        time.Now(),
	}

	return u.userRepository.UpdateVerify(ctx, &user)
}

func (u *userUseCase) UpdateImage(ctx context.Context, id string, file *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	idUser, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	session, err := u.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionCtx mongo_driven.SessionContext) (interface{}, error) {
		if file == nil {
			return nil, errors.New("images not nil")
		}

		// Kiểm tra xem file có phải là hình ảnh không
		if !helper.IsImage(file.Filename) {
			return nil, err
		}

		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()

		imageURL, err := images.UploadImageToCloudinary(f, file.Filename, u.database.CloudinaryUploadFolderUser, u.database)
		if err != nil {
			return nil, err
		}

		// Đảm bảo xóa ảnh trên Cloudinary nếu xảy ra lỗi sau khi tải lên thành công
		defer func() {
			if err != nil {
				_, _ = images.DeleteToCloudinary(imageURL.ImageURL, u.database)
			}
		}()

		user := domain.User{
			ID:        idUser,
			AvatarURL: imageURL.ImageURL,
			AssetURL:  imageURL.AssetID,
			UpdatedAt: time.Now(),
		}

		err = u.userRepository.UpdateImage(sessionCtx, &user)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Run the transaction
	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}

	u.cache.Clear()
	return session.CommitTransaction(ctx)
}

// SignUp Giảm độ phức tạp của hàm bằng cách không thực hiện việc tải file mà để công việc đó vào update Image
// Chỉ thực hiện đúng nhiệm vụ đăng ký tài khoản người dùng
func (u *userUseCase) SignUp(ctx context.Context, input *domain.SignupUser) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	session, err := u.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionCtx mongo_driven.SessionContext) (interface{}, error) {
		if err = validate_data.ValidateUser2(input); err != nil {
			return nil, err
		}

		// Bên phía client sẽ phải so sánh password thêm một lần nữa đã đúng chưa
		if !helper.PasswordStrong(input.Password) {
			return nil, errors.New("password must have at least 8 characters including uppercase letters, lowercase letters and numbers")
		}

		// Băm mật khẩu
		hashedPassword, err := password.HashPassword(input.Password)
		if err != nil {
			return nil, err
		}

		newUser := &domain.User{
			ID:                             primitive.NewObjectID(),
			Email:                          input.Email,
			PasswordHash:                   hashedPassword,
			FullName:                       "",
			Gender:                         "",
			Phone:                          input.Phone,
			Address:                        "",
			City:                           "",
			Region:                         "",
			Vocation:                       "",
			DateOfBirth:                    time.Now(),
			AssetURL:                       "",
			AvatarURL:                      "",
			Verified:                       false,
			VerificationCode:               "",
			Provider:                       "inside",
			Role:                           "User",
			FacebookSc:                     "",
			InstagramSc:                    "",
			LinkedInSc:                     "",
			YoutubeSc:                      "",
			ShowInterest:                   false,
			SocialMedia:                    false,
			EnableAutomaticSharingOfEvents: false,
			EnableSharingOn:                []string{},
			CreatedAt:                      time.Now(),
			UpdatedAt:                      time.Now(),
		}

		err = u.userRepository.CreateOne(sessionCtx, newUser)
		if err != nil {
			return nil, err
		}

		// logic chỗ này tương lai sẽ đổi
		var code string
		code = randstr.Dec(6)

		updUser := domain.User{
			ID:               newUser.ID,
			VerificationCode: code,
			Verified:         false,
			UpdatedAt:        time.Now(),
		}

		// Update User in Database
		err = u.userRepository.UpdateVerificationCode(sessionCtx, &updUser)
		if err != nil {
			return nil, err
		}

		emailData := handles.EmailData{
			Code:     code,
			FullName: newUser.FullName,
			Subject:  "Your account verification code: " + code,
		}

		err = handles.SendEmail(&emailData, input.Email, "sign_up.user.html")
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Run the transaction
	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}

	return session.CommitTransaction(ctx)
}

func (u *userUseCase) LoginUser(ctx context.Context, signIn *domain.SignIn) (domain.OutputLogin, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userData, err := u.userRepository.GetByEmail(ctx, signIn.Email)
	if err != nil || userData.Verified == false {
		return domain.OutputLogin{}, err
	}

	err = password.VerifyPassword(userData.PasswordHash, signIn.Password)
	if err != nil {
		return domain.OutputLogin{}, err
	}

	accessToken, err := token.CreateToken(u.database.AccessTokenExpiresIn, userData.ID, u.database.AccessTokenPrivateKey)
	if err != nil {
		return domain.OutputLogin{}, err
	}

	refreshToken, err := token.CreateToken(u.database.RefreshTokenExpiresIn, userData.ID, u.database.RefreshTokenPrivateKey)
	if err != nil {
		return domain.OutputLogin{}, err
	}

	response := domain.OutputLogin{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsLogged:     "1",
	}

	return response, nil
}

func (u *userUseCase) LoginGoogle(ctx context.Context, code string) (*domain.User, *domain.OutputLoginGoogle, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	googleOauthConfig := &oauth2.Config{
		ClientID:     u.database.GoogleClientID,
		ClientSecret: u.database.GoogleClientSecret,
		RedirectURL:  u.database.GoogleOAuthRedirectUrl,
		Scopes:       []string{"profile", "email"}, // Adjust scopes as needed
		Endpoint:     google.Endpoint,
	}

	tokenData, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Error exchanging code: " + err.Error())
		return nil, nil, err
	}

	userInfo, err := google_oauth2.GetUserInfo(tokenData.AccessToken)
	if err != nil {
		fmt.Println("Error getting user info: " + err.Error())
		return nil, nil, err
	}

	// Giả sử userInfo là một map[string]interface{}
	email := userInfo["email"].(string)
	phone := userInfo["phone"].(string)
	fullName := userInfo["name"].(string)
	avatarURL := userInfo["picture"].(string)
	verifiedEmail := userInfo["verified_email"].(bool)

	user := &domain.User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Phone:     phone,
		FullName:  fullName,
		AvatarURL: avatarURL,
		Provider:  "google",
		Verified:  verifiedEmail,
		Role:      "guess",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updateUser, err := u.userRepository.UpsertOne(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	signedToken, err := google_oauth2.SignJWT(userInfo)
	if err != nil {
		fmt.Println("Error signing token: " + err.Error())
		return nil, nil, err
	}

	accessToken, err := token.CreateToken(u.database.AccessTokenExpiresIn, updateUser.ID, u.database.AccessTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := token.CreateToken(u.database.RefreshTokenExpiresIn, updateUser.ID, u.database.RefreshTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	output := &domain.User{
		ID:          user.ID,
		Email:       user.Email,
		Phone:       user.Phone,
		FullName:    user.FullName,
		Gender:      user.Gender,
		Vocation:    user.Vocation,
		Address:     user.Address,
		City:        user.City,
		Region:      user.Region,
		DateOfBirth: user.DateOfBirth,
		AssetURL:    user.AssetURL,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
	}

	output2 := &domain.OutputLoginGoogle{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsLogged:     "1",
		SignedToken:  signedToken,
	}

	return output, output2, nil
}

func (u *userUseCase) DeleteOne(ctx context.Context, idUser string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return err
	}

	err = u.userRepository.DeleteOne(ctx, userID)
	if err != nil {
		return err
	}

	u.cache.Clear()
	return nil
}

func (u *userUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.OutputLogin, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	sub, err := token.ValidateToken(refreshToken, u.database.RefreshTokenPublicKey)
	if err != nil {
		return nil, err
	}

	userID, err := primitive.ObjectIDFromHex(fmt.Sprint(sub))
	if err != nil {
		return nil, err
	}

	userData, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	accessToken, err := token.CreateToken(u.database.AccessTokenExpiresIn, userData.ID, u.database.AccessTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	refresh, err := token.CreateToken(u.database.RefreshTokenExpiresIn, userData.ID, u.database.RefreshTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	response := &domain.OutputLogin{
		AccessToken:  accessToken,
		RefreshToken: refresh,
		IsLogged:     "1",
	}

	return response, nil
}

func (u *userUseCase) ForgetPassword(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	session, err := u.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionCtx mongo_driven.SessionContext) (interface{}, error) {
		user, err := u.userRepository.GetByEmail(sessionCtx, email)
		if err != nil {
			return nil, err
		}

		var code string
		code = randstr.Dec(6)

		updUser := &domain.User{
			ID:       user.ID,
			Verified: true,
		}

		// Update User in Database
		err = u.userRepository.UpdateVerify(sessionCtx, updUser)
		if err != nil {
			return nil, err
		}

		emailData := handles.EmailData{
			Code:     code,
			FullName: user.FullName,
			Subject:  "Khôi phục mật khẩu",
		}

		err = handles.SendEmail(&emailData, user.Email, "user.forget_password.html")
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Run the transaction
	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}

	return session.CommitTransaction(ctx)
}

func (u *userUseCase) UpdateVerifyForChangePassword(ctx context.Context, verificationCode string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.GetByVerificationCode(ctx, verificationCode)
	if err != nil {
		return err
	}

	if user.Verified == false {
		return errors.New("verification code check failed")
	}

	updUser := domain.User{
		ID:       user.ID,
		Verified: true,
	}

	return u.userRepository.UpdateVerify(ctx, &updUser)
}

func (u *userUseCase) UpdatePassword(ctx context.Context, id string, input *domain.ChangePasswordInput) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if input.Password != input.PasswordCompare {
		return errors.New("the passwords provided do not match")
	}

	user, err := u.userRepository.GetByVerificationCode(ctx, id)
	if err != nil {
		return err
	}

	input.Password, err = password.HashPassword(input.Password)
	if err != nil {
		return err
	}

	updateUser := &domain.User{
		ID:           user.ID,
		PasswordHash: input.Password,
		UpdatedAt:    time.Now(),
	}

	return u.userRepository.UpdatePassword(ctx, updateUser)
}
