package main

import (
	"capstone/configs"
	articleController "capstone/controllers/article"
	chatController "capstone/controllers/chat"
	chatbotController "capstone/controllers/chatbot"
	complaintController "capstone/controllers/complaint"
	consultationController "capstone/controllers/consultation"
	doctorController "capstone/controllers/doctor"
	forumController "capstone/controllers/forum"
	moodController "capstone/controllers/mood"
	musicController "capstone/controllers/music"
	notificationController "capstone/controllers/notification"
	otpController "capstone/controllers/otp"
	postController "capstone/controllers/post"
	ratingController "capstone/controllers/rating"
	storyController "capstone/controllers/story"
	transactionController "capstone/controllers/transaction"
	userController "capstone/controllers/user"
	"capstone/repositories/mysql"
	articleRepositories "capstone/repositories/mysql/article"
	chatRepositories "capstone/repositories/mysql/chat"
	complaintRepositories "capstone/repositories/mysql/complaint"
	consultationRepositories "capstone/repositories/mysql/consultation"
	doctorRepositories "capstone/repositories/mysql/doctor"
	forumRepositories "capstone/repositories/mysql/forum"
	moodRepositories "capstone/repositories/mysql/mood"
	musicRepositories "capstone/repositories/mysql/music"
	notificationRepositories "capstone/repositories/mysql/notification"
	otpRepositories "capstone/repositories/mysql/otp"
	postRepositories "capstone/repositories/mysql/post"
	ratingRepositories "capstone/repositories/mysql/rating"
	storyRepositories "capstone/repositories/mysql/story"
	transactionRepositories "capstone/repositories/mysql/transaction"
	userRepositories "capstone/repositories/mysql/user"
	"capstone/routes"
	articleUseCase "capstone/usecases/article"
	chatUseCase "capstone/usecases/chat"
	chatbotUseCase "capstone/usecases/chatbot"
	complaintUseCase "capstone/usecases/complaint"
	consultationUseCase "capstone/usecases/consultation"
	"capstone/usecases/cronjob"
	doctorUseCase "capstone/usecases/doctor"
	forumUseCase "capstone/usecases/forum"
	midtransUseCase "capstone/usecases/midtrans"
	moodUseCase "capstone/usecases/mood"
	musicUseCase "capstone/usecases/music"
	notificationUseCase "capstone/usecases/notification"
	otpUseCase "capstone/usecases/otp"
	postUseCase "capstone/usecases/post"
	ratingUseCase "capstone/usecases/rating"
	storyUseCase "capstone/usecases/story"
	transactionUseCase "capstone/usecases/transaction"
	userUseCase "capstone/usecases/user"
	"log"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

func main() {
	configs.LoadEnv()
	db := mysql.ConnectDB(configs.InitConfigMySQL())
	midtransConfig := configs.MidtransConfig()
	validate := validator.New()
	oauthConfig := configs.GetGoogleOAuthConfig()
	oauthConfigDoctor := configs.GetGoogleOAuthConfigDoctor()
	oauthConfigFB := configs.GetFacebookOAuthConfig()
	oauthConfigFBDoctor := configs.GetFacebookOAuthConfigDoctor()

	gcron, err := gocron.NewScheduler()
	if err != nil {
		log.Print(err.Error())
	}

	userRepo := userRepositories.NewUserRepo(db)
	doctorRepo := doctorRepositories.NewDoctorRepo(db)
	consultationRepo := consultationRepositories.NewConsultationRepo(db)
	storyRepo := storyRepositories.NewStoryRepo(db)
	complaintRepo := complaintRepositories.NewComplaintRepo(db)
	transactionRepo := transactionRepositories.NewTransactionRepo(db)
	musicRepo := musicRepositories.NewMusicRepo(db)
	ratingRepo := ratingRepositories.NewRatingRepo(db)
	moodRepo := moodRepositories.NewMoodRepo(db)
	forumRepo := forumRepositories.NewForumRepo(db)
	postRepo := postRepositories.NewPostRepo(db)
	articleRepo := articleRepositories.NewArticleRepo(db)
	chatRepo := chatRepositories.NewChatRepo(db)
	otpRepo := otpRepositories.NewOtpRepo(db)
	notificationRepo := notificationRepositories.NewNotificationRepository(db)

	userUC := userUseCase.NewUserUseCase(userRepo, oauthConfig, oauthConfigFB)
	notificationUC := notificationUseCase.NewNotificationUseCase(notificationRepo)
	doctorUC := doctorUseCase.NewDoctorUseCase(doctorRepo, ratingRepo, oauthConfigDoctor, oauthConfigFBDoctor)
	consultationUC := consultationUseCase.NewConsultationUseCase(consultationRepo, transactionRepo, userUC, doctorRepo, notificationUC, validate, chatRepo)
	storyUC := storyUseCase.NewStoryUseCase(storyRepo)
	complaintUC := complaintUseCase.NewComplaintUseCase(complaintRepo, notificationUC, consultationUC)
	midtransUC := midtransUseCase.NewMidtransUseCase(midtransConfig)
	transactionUC := transactionUseCase.NewTransactionUseCase(transactionRepo, midtransUC, consultationRepo, doctorRepo, userUC, validate)
	musicUC := musicUseCase.NewMusicUseCase(musicRepo)
	ratingUC := ratingUseCase.NewRatingUseCase(ratingRepo)
	moodUC := moodUseCase.NewMoodUseCase(moodRepo)
	forumUC := forumUseCase.NewForumUseCase(forumRepo)
	postUC := postUseCase.NewPostUseCase(postRepo)
	chatbotUC := chatbotUseCase.NewChatbotUsecase()
	articleUC := articleUseCase.NewArticleUseCase(articleRepo)
	chatUC := chatUseCase.NewChatUseCase(chatRepo)
	otpUC := otpUseCase.NewOtpUseCase(otpRepo)
	cronjobUC := cronjob.NewCronJob(gcron, consultationRepo)

	userCont := userController.NewUserController(userUC)
	doctorCont := doctorController.NewDoctorController(doctorUC, validate)
	consultationCont := consultationController.NewConsultationController(consultationUC, validate)
	storyCont := storyController.NewStoryController(storyUC)
	complaintCont := complaintController.NewComplaintController(complaintUC, consultationUC, validate)
	transactionCont := transactionController.NewTransactionController(transactionUC, midtransUC, validate)
	musicCont := musicController.NewMusicController(musicUC)
	ratingCont := ratingController.NewRatingController(ratingUC)
	moodCont := moodController.NewMoodController(moodUC)
	forumCont := forumController.NewForumController(forumUC)
	postCont := postController.NewPostController(postUC)
	chatbotCont := chatbotController.NewChatbotController(chatbotUC)
	articleCont := articleController.NewArticleController(articleUC)
	chatCont := chatController.NewChatController(chatUC)
	otpCont := otpController.NewOtpController(otpUC)
	notificationCont := notificationController.NewNotificationController(notificationUC)

	route := routes.NewRoute(userCont, doctorCont, consultationCont, storyCont, complaintCont, transactionCont, musicCont, ratingCont, moodCont, forumCont, postCont, chatbotCont, articleCont, chatCont, otpCont, notificationCont)

	e := echo.New()
	route.InitRoute(e)
	cronjobUC.InitCronJob()

	go gcron.Start()

	e.Logger.Fatal(e.Start(":8080"))
}
