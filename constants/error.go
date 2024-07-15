package constants

import "errors"

var ErrEmptyInputUser error = errors.New("fullname, username, email or password cannot be empty")

var ErrHashedPassword error = errors.New("error hashing password")

var ErrInsertDatabase error = errors.New("failed insert data in database")

var ErrUsernameAlreadyExist error = errors.New("username already exist")

var ErrEmailAlreadyExist error = errors.New("email already exist")

var ErrEmptyInputLogin error = errors.New("username or password cannot be empty")

var ErrUserNotFound error = errors.New("user not found")

var ErrDataNotFound error = errors.New("data not found")

var ErrInvalidToken error = errors.New("invalid token")

var ErrDataEmpty error = errors.New("data empty")

var ErrEmptyInputArticle error = errors.New("title or content cannot be empty")

var ErrServer error = errors.New("server error")

var ErrInvalidRate error = errors.New("rate must be between 1 and 5")

var ErrCloudinary error = errors.New("cloudinary url not found")

var ErrEmptyInputMood error = errors.New("mood type id or date cannot be empty")

var ErrUploadImage error = errors.New("failed upload image")

var ErrEmptyRangeDateMood error = errors.New("start date or end date cannot be empty")

var ErrInvalidStartDate error = errors.New("invalid format start date")

var ErrInvalidEndDate error = errors.New("invalid format end date")

var ErrStartDateGreater error = errors.New("start date must be less than end date")

var ErrAlreadyLiked error = errors.New("already liked")

var ErrEmptyInputForum error = errors.New("name or description cannot be empty")

var ErrEmptyInputPost error = errors.New("forum id or content cannot be empty")

var ErrEmptyInputLike error = errors.New("post id cannot be empty")

var ErrEmptyInputComment error = errors.New("post id or content cannot be empty")

var ErrExcange error = errors.New("failed excange")

var ErrNewServiceGoogle error = errors.New("failed new service google")

var ErrNewUserInfo error = errors.New("failed new user info")

var ErrInsertOAuth error = errors.New("failed insert oauth")

var ErrEmptyInputMusic error = errors.New("title or singer cannot be empty")

var ErrEmptyInputStory error = errors.New("title or content cannot be empty")

var ErrInputTime error = errors.New("invalid time format")

var ErrInvalidConsultationID error = errors.New("invalid consultation id")

var ErrEmptyCreateForum error = errors.New("name, description, or image cannot be empty")

var ErrEmptyChat error = errors.New("chat id or message cannot be empty")

var ErrEmptyInputEmailOTP error = errors.New("email cannot be empty")

var ErrEmptyInputVerifyOTP error = errors.New("email or code cannot be empty")

var ErrInvalidOTP error = errors.New("invalid otp")

var ErrExpiredOTP error = errors.New("expired otp")

var ErrEmptyResetPassword error = errors.New("new password cannot be empty")

var ErrDeleteDatabase error = errors.New("failed delete data in database")

var ErrInvalidPrice error = errors.New("price must be greater than 0")

var ErrUnauthorized error = errors.New("unauthorized")

var ErrPointSpend error = errors.New("point spend must be greater than 0")

var ErrInsufficientPoint error = errors.New("insufficient point")

var ErrConsultationAlreadyRejected error = errors.New("status already rejected")

var ErrInvalidCredentials error = errors.New("invalid credentials")

var ErrEmptyInputViewByMonth error = errors.New("start month or end month cannot be empty")

var ErrBadRequest error = errors.New("bad request")

var ErrEmptyNewEmail error = errors.New("new email cannot be empty")

var ErrLocationNotFound error = errors.New("location not found")

var ErrNotificationAlreadyRead error = errors.New("notification already read")

var ErrUpdateDatabase error = errors.New("failed update data in database")
