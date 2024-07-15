package configs

import (
	"capstone/repositories/mysql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default value")
	}
}

func InitConfigMySQL() mysql.Config {
	return mysql.Config{
		DBName: os.Getenv("DBName"),
		DBUser: os.Getenv("DBUser"),
		DBPass: os.Getenv("DBPass"),
		DBHost: os.Getenv("DBHost"),
		DBPort: os.Getenv("DBPort"),
	}
}

func GetGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "https://dev-capstone.practiceproject.tech/v1/users/auth/google/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func GetGoogleOAuthConfigDoctor() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID_DOCTOR"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET_DOCTOR"),
		RedirectURL:  "https://dev-capstone.practiceproject.tech/v1/doctors/auth/google/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func GetFacebookOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		RedirectURL:  "https://dev-capstone.practiceproject.tech/v1/users/auth/facebook/callback",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
}

func GetFacebookOAuthConfigDoctor() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("FACEBOOK_CLIENT_ID_DOCTOR"),
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET_DOCTOR"),
		RedirectURL:  "https://dev-capstone.practiceproject.tech/v1/doctors/auth/facebook/callback",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
}

func InitConfigJWT() string {
	return os.Getenv("SECRET_JWT")
}

func InitConfigCloudinary() string {
	return os.Getenv("CLOUDINARY_URL")
}

func InitConfigKeyChatbot() string {
	return os.Getenv("AI_KEY")
}

func InitConfigMyEmail() string {
	return os.Getenv("MY_EMAIL")
}

func InitConfigAppPassword() string {
	return os.Getenv("APP_PASSWORD")
}
