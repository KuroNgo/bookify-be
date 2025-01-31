package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUser = "user"
)

type User struct {
	ID                             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email                          string             `bson:"email" json:"email"`
	PasswordHash                   string             `bson:"password_hash" json:"password_hash"` // Hash of the password
	Phone                          string             `json:"phone" bson:"phone"`
	FullName                       string             `bson:"full_name" json:"full_name"`
	Gender                         string             `bson:"gender" json:"gender"`
	Vocation                       string             `bson:"vocation" json:"vocation"`
	Address                        string             `bson:"address" json:"address"`
	City                           string             `bson:"city" json:"city"`
	Region                         string             `bson:"region" json:"region"`
	DateOfBirth                    time.Time          `bson:"date_of_birth" json:"date_of_birth"`
	AssetURL                       string             `bson:"asset_url"  json:"asset_url"`
	AvatarURL                      string             `bson:"avatar_url"  json:"avatar_url"`
	Verified                       bool               `bson:"verify"   json:"verify"`
	VerificationCode               string             `bson:"verification_code" json:"verification_code"`
	Provider                       string             `bson:"provider" json:"provider"`
	Role                           string             `bson:"role" json:"role"` // Example: "Admin", "Manager", "Employee"
	FacebookSc                     string             `bson:"facebook_sc" json:"facebook_sc"`
	InstagramSc                    string             `bson:"instagram_sc" json:"instagram_sc"`
	LinkedInSc                     string             `bson:"linked_in_sc" json:"linked_in_sc"`
	YoutubeSc                      string             `bson:"youtube_sc" json:"youtube_sc"`
	ShowInterest                   bool               `bson:"show_interest" json:"show_interest"`
	SocialMedia                    bool               `bson:"social_media" json:"social_media"`
	EnableAutomaticSharingOfEvents bool               `bson:"enable_automatic_sharing_of_events" json:"enable_automatic_sharing_of_events"`
	EnableSharingOn                []string           `bson:"enable_sharing_on" json:"enable_sharing_on"`
	CreatedAt                      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt                      time.Time          `bson:"updated_at" json:"updated_at"`
}

type InputUser struct {
	ID               string    `bson:"_id" json:"id,omitempty"`
	Gender           string    `bson:"gender" json:"gender"`
	Vocation         string    `bson:"vocation" json:"vocation"`
	Address          string    `bson:"address" json:"address"`
	City             string    `bson:"city" json:"city"`
	Region           string    `bson:"region" json:"region"`
	DateOfBirth      time.Time `bson:"date_of_birth" json:"date_of_birth"`
	FullName         string    `bson:"full_name" json:"full_name"`
	PasswordHash     string    `bson:"password" json:"password"` // Hash of the password
	AvatarURL        string    `bson:"avatar_url"  json:"avatar_url"`
	AssetURL         string    `bson:"asset_url"  json:"asset_url"`
	Email            string    `bson:"email" json:"email"`
	Phone            string    `json:"phone" bson:"phone"`
	Verified         bool      `bson:"verify"   json:"verify"`
	VerificationCode string    `bson:"verification_code" json:"verification_code"`
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

type UpdateUserInfo struct {
	ID                             string   `bson:"_id" json:"id,omitempty"`
	Gender                         string   `bson:"gender" json:"gender"`
	Vocation                       string   `bson:"vocation" json:"vocation"`
	Address                        string   `bson:"address" json:"address"`
	City                           string   `bson:"city" json:"city"`
	Region                         string   `bson:"region" json:"region"`
	DateOfBirth                    string   `bson:"date_of_birth" json:"date_of_birth"`
	FullName                       string   `bson:"full_name" json:"full_name"`
	AvatarURL                      string   `bson:"avatar_url"  json:"avatar_url"`
	AssetURL                       string   `bson:"asset_url"  json:"asset_url"`
	FacebookSc                     string   `bson:"facebook_sc" json:"facebook_sc"`
	InstagramSc                    string   `bson:"instagram_sc" json:"instagram_sc"`
	LinkedInSc                     string   `bson:"linked_in_sc" json:"linked_in_sc"`
	YoutubeSc                      string   `bson:"youtube_sc" json:"youtube_sc"`
	ShowInterest                   bool     `bson:"show_interest" json:"show_interest"`
	SocialMedia                    bool     `bson:"social_media" json:"social_media"`
	EnableAutomaticSharingOfEvents bool     `bson:"enable_automatic_sharing_of_events" json:"enable_automatic_sharing_of_events"`
	EnableSharingOn                []string `bson:"enable_sharing_on" json:"enable_sharing_on"`
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

type UpdateSocialMedia struct {
	ID                             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FacebookSc                     string             `bson:"facebook_sc" json:"facebook_sc"`
	InstagramSc                    string             `bson:"instagram_sc" json:"instagram_sc"`
	LinkedInSc                     string             `bson:"linked_in_sc" json:"linked_in_sc"`
	YoutubeSc                      string             `bson:"youtube_sc" json:"youtube_sc"`
	EnableAutomaticSharingOfEvents bool               `bson:"enable_automatic_sharing_of_events" json:"enable_automatic_sharing_of_events"`
	EnableSharingOn                []string           `bson:"enable_sharing_on" json:"enable_sharing_on"`
}

type UpdateUserSettings struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Gender       string             `bson:"gender" json:"gender"`
	Phone        string             `json:"phone" bson:"phone"`
	Vocation     string             `bson:"vocation" json:"vocation"`
	Address      string             `bson:"address" json:"address"`
	City         string             `bson:"city" json:"city"`
	Region       string             `bson:"region" json:"region"`
	DateOfBirth  string             `bson:"date_of_birth" json:"date_of_birth"`
	FullName     string             `bson:"full_name" json:"full_name"`
	ShowInterest bool               `bson:"show_interest" json:"show_interest"`
	SocialMedia  bool               `bson:"social_media" json:"social_media"`
}

type OutputLogin struct {
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	IsLogged     string `bson:"is_logged" json:"is_logged"`
}

type OutputLoginGoogle struct {
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	IsLogged     string `bson:"is_logged" json:"is_logged"`
	SignedToken  string `bson:"signed_token" json:"signed_token"`
}
