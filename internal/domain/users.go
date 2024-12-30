package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUser = "user"
)

// User represents a user in the system.
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username         string             `bson:"username" json:"username"`
	PasswordHash     string             `bson:"password_hash" json:"password_hash"` // Hash of the password
	Email            string             `bson:"email" json:"email"`
	Phone            string             `json:"phone" bson:"phone"`
	AssetURL         string             `bson:"asset_url"  json:"asset_url"`
	AvatarURL        string             `bson:"avatar_url"  json:"avatar_url"`
	Verified         bool               `bson:"verify"   json:"verify"`
	VerificationCode string             `bson:"verification_code" json:"verification_code"`
	Provider         string             `bson:"provider" json:"provider"`
	Role             string             `bson:"role" json:"role"` // Example: "Admin", "Manager", "Employee"
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type InputUser struct {
	ID               string `bson:"_id" json:"id,omitempty"`
	Username         string `bson:"username" json:"username"`
	PasswordHash     string `bson:"password" json:"password"` // Hash of the password
	AvatarURL        string `bson:"avatar_url"  json:"avatar_url"`
	AssetURL         string `bson:"asset_url"  json:"asset_url"`
	Email            string `bson:"email" json:"email"`
	Phone            string `json:"phone" bson:"phone"`
	Verified         bool   `bson:"verify"   json:"verify"`
	VerificationCode string `bson:"verification_code" json:"verification_code"`
}

type SignupUser struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"` // Hash of the password
	Phone    string `json:"phone" bson:"phone"`
}

type SignIn struct {
	Email    string `bson:"email" json:"email" example:"admin@admin.com" `
	Password string `bson:"password" json:"password" example:"12345"`
}

type VerificationInput struct {
	VerificationCode string `json:"verification_code" binding:"required"`
}

type ChangePasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordCompare string `json:"password_compare" binding:"required"`
}

type ForgetPassword struct {
	Email string `json:"email" bson:"email"`
}

type OutputUser struct {
	User User `bson:"user" json:"user"`
}

type OutputLogin struct {
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	IsLogged     string `bson:"is_logged" json:"is_logged"`
}

type OutputLoginGoogle struct {
	RefreshToken string `bson:"refresh_token" json:"refresh_token""`
	AccessToken  string `bson:"access_token" json:"access_token"`
	IsLogged     string `bson:"is_logged" json:"is_logged"`
	SignedToken  string `bson:"signed_token" json:"signed_token"`
}
