package user

import "context"

type User struct {
	Id             int
	Name           string
	Username       string
	Email          string
	Password       string
	Address        string
	Bio            string
	PhoneNumber    string
	Gender         string
	Age            int
	ProfilePicture string
	Token          string
	IsOauth        bool
	Points         int
	IsActive       bool
	PendingEmail   string
}

type RepositoryInterface interface {
	Register(user *User) (User, int64, error)
	Login(user *User) (User, error)
	Create(email string, picture string, name string, username string) (User, error)
	OauthFindByEmail(email string) (User, int, error)
	GetPointsByUserId(id int) (int, error)
	ResetPassword(email string, password string) error
	UpdatePointsByUserID(id int, point int) error
	UpdateUserProfile(user *User) (User, error)
	ChangePassword(userId int, oldPassword, newPassword string) error
	ChangeEmail(userId int, email string) error
	GetDetailedProfile(id int) (User, error)
}

type UseCaseInterface interface {
	Register(user *User) (User, error)
	Login(user *User) (User, error)
	HandleGoogleLogin() string
	HandleGoogleCallback(ctx context.Context, code string) (User, error)
	GetPointsByUserId(id int) (int, error)
	ResetPassword(email string, password string) error
	HandleFacebookLogin() string
	HandleFacebookCallback(ctx context.Context, code string) (User, error)
	UpdateSuccessPointByUserID(id int, pointSpend int) error
	UpdateFailedPointByUserID(id int, pointSpend int) error
	UpdateUserProfile(user *User) (User, error)
	ChangePassword(userId int, oldPassword, newPassword string) error
	ChangeEmail(userId int, newEmail string) error
	GetDetailedProfile(id int) (User, error)
}
