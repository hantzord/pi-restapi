package routes

import (
	"capstone/controllers/article"
	"capstone/controllers/chat"
	"capstone/controllers/chatbot"
	"capstone/controllers/complaint"
	"capstone/controllers/consultation"
	"capstone/controllers/doctor"
	"capstone/controllers/forum"
	"capstone/controllers/mood"
	"capstone/controllers/music"
	"capstone/controllers/notification"
	"capstone/controllers/otp"
	"capstone/controllers/post"
	"capstone/controllers/rating"
	"capstone/controllers/story"
	"capstone/controllers/transaction"
	"capstone/controllers/user"
	myMiddleware "capstone/middlewares"
	"capstone/utilities/base"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type RouteController struct {
	userController         *user.UserController
	doctorController       *doctor.DoctorController
	consultationController *consultation.ConsultationController
	storyController        *story.StoryController
	complaintController    *complaint.ComplaintController
	transactionController  *transaction.TransactionController
	musicController        *music.MusicController
	ratingController       *rating.RatingController
	moodController         *mood.MoodController
	forumController        *forum.ForumController
	postController         *post.PostController
	chatbotController      *chatbot.ChatbotController
	articleController      *article.ArticleController
	chatController         *chat.ChatController
	otpController          *otp.OtpController
	notificationController *notification.NotificationController
}

func NewRoute(
	userController *user.UserController,
	doctorController *doctor.DoctorController,
	consultationController *consultation.ConsultationController,
	storyContoller *story.StoryController,
	complaintController *complaint.ComplaintController,
	transactionController *transaction.TransactionController,
	musicController *music.MusicController,
	ratingController *rating.RatingController,
	moodController *mood.MoodController,
	forumController *forum.ForumController,
	postController *post.PostController,
	chatbotController *chatbot.ChatbotController,
	articleController *article.ArticleController,
	chatController *chat.ChatController,
	otpController *otp.OtpController,
	notification *notification.NotificationController) *RouteController {
	return &RouteController{
		userController:         userController,
		doctorController:       doctorController,
		consultationController: consultationController,
		storyController:        storyContoller,
		complaintController:    complaintController,
		transactionController:  transactionController,
		musicController:        musicController,
		ratingController:       ratingController,
		moodController:         moodController,
		forumController:        forumController,
		postController:         postController,
		chatbotController:      chatbotController,
		articleController:      articleController,
		chatController:         chatController,
		otpController:          otpController,
		notificationController: notification,
	}
}

func (r *RouteController) InitRoute(e *echo.Echo) {
	myMiddleware.LogMiddleware(e)

	e.HTTPErrorHandler = base.ErrorHandler
	e.Use(myMiddleware.CORSMiddleware())

	e.POST("/v1/payment-callback", r.transactionController.CallbackTransaction)

	// chatbot
	e.GET("/v1/users/chatbots/customer-service", r.chatbotController.ChatbotCS)        //customer service chatbot
	e.GET("/v1/users/chatbots/mental-health", r.chatbotController.ChatbotMentalHealth) //mental health chatbot
	e.GET("/v1/doctors/chatbots/treatment", r.chatbotController.ChatbotTreatment)      //Chatbot Treatment

	// Users
	userAuth := e.Group("/v1/users")
	userAuth.POST("/register", r.userController.Register)                //Register User
	userAuth.POST("/login", r.userController.Login)                      //Login User
	userAuth.PUT("/reset-password", r.userController.ResetPassword)      //Reset Password
	userAuth.PUT("/profiles", r.userController.UpdateProfile)            //Update Profile
	userAuth.PUT("/profiles/password", r.userController.ChangePassword)  // Change Password
	userAuth.PUT("/profiles/change-email", r.userController.ChangeEmail) // Change Email (Save new email to pending email)

	userAuth.GET("/auth/google/login", r.userController.GoogleLogin)
	userAuth.GET("/auth/google/callback", r.userController.GoogleCallback)

	userAuth.GET("/auth/facebook/login", r.userController.FacebookLogin)
	userAuth.GET("/auth/facebook/callback", r.userController.FacebookCallback)

	// OTP
	userAuth.POST("/otp/send", r.otpController.SendOtp)                             // Send OTP
	userAuth.POST("/otp/verify/forgot-password", r.otpController.VerifyOtp)         // Verify OTP Forgot Password
	userAuth.POST("/otp/verify/register", r.otpController.VerifyOTPRegister)        // Verify OTP Register

	userRoute := userAuth.Group("/")
	userRoute.Use(echojwt.JWT([]byte(os.Getenv("SECRET_JWT"))))
	// Doctor
	userRoute.GET("doctors/:id", r.doctorController.GetByID)         //Get Doctor By ID
	userRoute.GET("doctors", r.doctorController.GetAll)              //Get All Doctor
	userRoute.GET("doctors/available", r.doctorController.GetActive) //Get All Active Doctor
	userRoute.GET("doctors/search", r.doctorController.SearchDoctor) //Search Doctor

	// Consultation
	userRoute.POST("consultations", r.consultationController.CreateConsultation)     //Create Consultation
	userRoute.GET("consultations/:id", r.consultationController.GetConsultationByID) //Get Consultation By ID
	userRoute.GET("consultations", r.consultationController.GetAllConsultation)      //Get All Consultation

	// Inspirational Stories
	userRoute.GET("stories", r.storyController.GetAllStories)         //Get All Stories
	userRoute.GET("stories/:id", r.storyController.GetStoryById)      //Get Story By ID
	userRoute.GET("stories/liked", r.storyController.GetLikedStories) //Get Liked Stories
	userRoute.POST("stories/like", r.storyController.LikeStory)       //Like Story
	userRoute.DELETE("stories/like", r.storyController.UnlikeStory)   //Unlike Story

	// Music
	userRoute.GET("musics", r.musicController.GetAllMusics)         //Get All Music
	userRoute.GET("musics/:id", r.musicController.GetMusicByID)     //Get Music By ID
	userRoute.GET("musics/liked", r.musicController.GetLikedMusics) //Get Liked Music
	userRoute.POST("musics/like", r.musicController.LikeMusic)      //Like Music
	userRoute.DELETE("musics/like", r.musicController.UnlikeMusic)  //Unlike Music

	// Complaint
	userRoute.POST("complaint", r.complaintController.Create) // Create Complaint

	// Transaction
	userRoute.POST("payments/gateway", r.transactionController.InsertWithBuiltIn)               // Create Transaction
	userRoute.GET("transaction/:id", r.transactionController.FindByID)                          // Get Transaction By ID
	userRoute.GET("transaction/consultation/:id", r.transactionController.FindByConsultationID) // Get Transaction By Consultation ID
	userRoute.GET("transactions", r.transactionController.FindAll)                              // Get All Transaction
	userRoute.POST("payments/bank-transfer", r.transactionController.BankTransfer)              // Bank Transfer
	userRoute.POST("payments/e-wallet", r.transactionController.EWallet)                        // E-Wallet

	// Rating
	userRoute.POST("feedbacks", r.ratingController.SendFeedback) // Create Rating

	// Mood
	userRoute.POST("moods", r.moodController.CreateMood)     // Create Mood
	userRoute.GET("moods", r.moodController.GetAllMoods)     // Get All Moods
	userRoute.GET("moods/:id", r.moodController.GetMoodById) // Get Mood By ID

	// Forum
	userRoute.POST("forums/join", r.forumController.JoinForum)                       // Join Forum
	userRoute.DELETE("forums/:id", r.forumController.LeaveForum)                     // Leave Forum
	userRoute.GET("forums", r.forumController.GetJoinedForum)                        // Get All Forum
	userRoute.GET("forums/recommendation", r.forumController.GetRecommendationForum) // Get Recommendation Forum
	userRoute.GET("forums/:id", r.forumController.GetForumById)                      // Get Forum By ID

	// Posts
	userRoute.GET("forums/:forumId/posts", r.postController.GetAllPostsByForumId)   // Get All Posts By Forum ID
	userRoute.GET("posts/:id", r.postController.GetPostById)                        // Get Post By ID
	userRoute.POST("posts", r.postController.SendPost)                              // Create Post
	userRoute.POST("posts/like", r.postController.LikePost)                         // Like Post
	userRoute.DELETE("posts/like", r.postController.UnlikePost)                     // Unlike Post
	userRoute.POST("comments", r.postController.SendComment)                        // Create Comment
	userRoute.GET("posts/:postId/comments", r.postController.GetAllCommentByPostId) // Get All Comment By Post ID

	// Article
	userRoute.GET("articles", r.articleController.GetAllArticle) // Get All Article
	userRoute.GET("articles/:id", r.articleController.GetArticleById)
	userRoute.GET("articles/liked", r.articleController.GetLikedArticle)
	userRoute.POST("articles/like", r.articleController.LikeArticle)
	userRoute.DELETE("articles/like", r.articleController.UnlikeArticle)

	// Notification
	userRoute.GET("notifications", r.notificationController.GetAllUserNotification) // Get User Notification
	userRoute.PUT("notifications/:notificationID", r.notificationController.UpdateToReadConsultationUser)
	userRoute.DELETE("notifications/:notificationID", r.notificationController.DeleteToReadConsultationUser)

	// Profiles
	userRoute.GET("profiles", r.userController.GetDetailedProfile) // Get Profile By User ID

	// Consultation Notes
	userRoute.GET("consultation-notes/consultation/:chatId", r.consultationController.GetConsultationNotesByID) // Get Consultation Note By ID

	// Points
	userRoute.GET("points", r.userController.GetPointsByUserId) // Get Points

	// Chat
	userRoute.GET("chats", r.chatController.GetAllChatByUserId) // Get All Chat

	// Chat Messages
	userRoute.POST("chats/messages", r.chatController.SendMessage)           // Send Message
	userRoute.GET("chats/:chatId/messages", r.chatController.GetAllMessages) // Get All Message

	// OTP
	userRoute.POST("otp/send/change-email", r.otpController.SendOTPChangeEmail) // Send OTP Change Email
	userRoute.POST("otp/verify/change-email", r.otpController.VerifyOTPChangeEmail) // Verify OTP Change Email

	doctorAuth := e.Group("/v1/doctors")

	doctorAuth.POST("/register", r.doctorController.Register)                  //Register Doctor
	doctorAuth.POST("/login", r.doctorController.Login)                        //Login Doctor
	doctorAuth.GET("/auth/google/login", r.doctorController.GoogleLogin)       // Google Login
	doctorAuth.GET("/auth/google/callback", r.doctorController.GoogleCallback) // Google Callback

	doctorAuth.GET("/auth/facebook/login", r.doctorController.FacebookLogin)
	doctorAuth.GET("/auth/facebook/callback", r.doctorController.FacebookCallback)

	doctorRoute := doctorAuth.Group("/")
	doctorRoute.Use(echojwt.JWT([]byte(os.Getenv("SECRET_JWT"))))

	// articles
	doctorRoute.POST("articles", r.articleController.CreateArticle)                           // Post Article
	doctorRoute.GET("articles", r.articleController.GetAllArticleByDoctorId)                  // Get All Article
	doctorRoute.GET("articles/:id", r.articleController.GetArticleByIdForDoctor)              // Get Article By ID
	doctorRoute.GET("articles/count", r.articleController.CountArticleByDoctorId)             // Count Article By Doctor ID
	doctorRoute.GET("articles/like/count", r.articleController.CountArticleLikesByDoctorId)   // Count Article Likes By Doctor ID
	doctorRoute.GET("articles/view/count", r.articleController.CountArticleViewByDoctorId)    // Count Article View Count By Doctor ID
	doctorRoute.GET("articles/view/month/count", r.articleController.CountArticleViewByMonth) // Count Article View Count By Month
	doctorRoute.PUT("articles/:id", r.articleController.EditArticle)                          // Update Article
	doctorRoute.DELETE("articles/:id", r.articleController.DeleteArticle)                     // Delete Article

	// musics
	doctorRoute.POST("musics", r.musicController.PostMusic)                               // Post Music
	doctorRoute.GET("musics", r.musicController.GetAllMusicsByDoctorId)                   // Get All Music By Doctor ID
	doctorRoute.GET("musics/:id", r.musicController.GetMusicByIdForDoctor)                // Get Music By ID
	doctorRoute.GET("musics/count", r.musicController.CountMusicByDoctorId)               // Count Music By Doctor ID
	doctorRoute.GET("musics/like/count", r.musicController.CountMusicLikesByDoctorId)     // Count Music Likes By Doctor ID
	doctorRoute.GET("musics/view/count", r.musicController.CountMusicViewCountByDoctorId) // Count Music View Count By Doctor ID
	doctorRoute.PUT("musics/:id", r.musicController.EditMusic)                            // Update Music
	doctorRoute.DELETE("musics/:id", r.musicController.DeleteMusic)                       // Delete Music
	doctorRoute.GET("musics/view/month/count", r.musicController.CountMusicViewByMonth)   // Count Music View By Month

	// Inspirational Stories
	doctorRoute.POST("stories", r.storyController.PostStory)                             // Post Story
	doctorRoute.GET("stories", r.storyController.GetAllStoriesByDoctorId)                // Get All Story By Doctor ID
	doctorRoute.GET("stories/:id", r.storyController.GetStoryByIdForDoctor)              // Get Story By ID
	doctorRoute.GET("stories/count", r.storyController.CountStoriesByDoctorId)           // Count Stories By Doctor ID
	doctorRoute.GET("stories/like/count", r.storyController.CountStoryLikesByDoctorId)   // Count Stories Likes By Doctor ID
	doctorRoute.GET("stories/view/month/count", r.storyController.CountStoryViewByMonth) // Count Stories View Count By Month
	doctorRoute.GET("stories/view/count", r.storyController.CountStoryViewByDoctorId)    // Count Stories View Count By Doctor ID
	doctorRoute.PUT("stories/:id", r.storyController.EditStory)                          // Update Story
	doctorRoute.DELETE("stories/:id", r.storyController.DeleteStory)                     // Delete Story

	// Consultation
	doctorRoute.GET("consultations", r.consultationController.GetAllDoctorConsultation) //Get All Consultation
	doctorRoute.GET("consultations/count", r.consultationController.CountConsultation)
	doctorRoute.GET("consultations/today/count", r.consultationController.CountConsultationToday)
	doctorRoute.GET("consultations/:id", r.consultationController.GetDoctorConsultationByID) //Get Consultation By ID
	doctorRoute.PUT("consultations/:id", r.consultationController.UpdateStatusConsultation)  // Update Status Consultation

	doctorRoute.GET("stories/view/count", r.storyController.CountStoryViewByDoctorId) // Count Stories View Count By Doctor ID
	doctorRoute.PUT("stories/:id", r.storyController.EditStory)                       // Update Story
	doctorRoute.DELETE("stories/:id", r.storyController.DeleteStory)                  // Delete Story

	// consultation notes
	doctorRoute.POST("consultation-notes", r.consultationController.CreateConsultationNotes) // Post Consultation Note

	// Rating
	doctorRoute.GET("feedbacks", r.ratingController.GetAllFeedbacks) // Get All Feedbacks
	doctorRoute.GET("ratings", r.ratingController.GetSummaryRating)  // Get Summary Rating

	// Forum
	doctorRoute.POST("forums", r.forumController.CreateForum)                             // Create Forum
	doctorRoute.GET("forums", r.forumController.GetAllForumsByDoctorId)                   // Get All Forum By Doctor ID
	doctorRoute.PUT("forums/:id", r.forumController.UpdateForum)                          // Update Forum
	doctorRoute.DELETE("forums/:id", r.forumController.DeleteForum)                       // Delete Forum
	doctorRoute.GET("forums/:id", r.forumController.GetForumById)                         // Get Forum By ID
	doctorRoute.GET("forums/:forumId/members", r.forumController.GetForumMemberByForumId) // Get Members By Forum ID

	// Post
	doctorRoute.GET("forums/:forumId/posts", r.postController.GetAllPostsByForumId)   // Get All Posts By Forum ID
	doctorRoute.GET("posts/:postId/comments", r.postController.GetAllCommentByPostId) // Get All Comment By Post ID

	// Chat
	doctorRoute.GET("chats", r.chatController.GetAllChatByDoctorId) // Get All Chat By Doctor ID

	// Chat Message
	doctorRoute.POST("chats/messages", r.chatController.SendMessageDoctor)     // Send Message
	doctorRoute.GET("chats/:chatId/messages", r.chatController.GetAllMessages) // Get All Message

	// Transaction
	doctorRoute.GET("transactions", r.transactionController.FindAllByDoctorID)                // Get All Transaction By Doctor ID
	doctorRoute.GET("transaction/:id", r.transactionController.FindByID)                      // Get Transaction By ID
	doctorRoute.GET("transactions/count", r.transactionController.CountTransactionByDoctorID) // Count Transaction By Doctor ID
	doctorRoute.DELETE("transactions/:id", r.transactionController.DeleteTransaction)         // Delete Transaction

	// Patient
	doctorRoute.GET("patients", r.complaintController.GetAllComplaint) // Get All Patient
	doctorRoute.GET("patients/:id", r.complaintController.GetConsultationByComplaintID)
	doctorRoute.GET("patients/search", r.complaintController.SearchComplaintByPatientName)

	// Notification
	doctorRoute.GET("notifications", r.notificationController.GetAllDoctorNotification) // Get Doctor Notification
	doctorRoute.PUT("notifications/:notificationID", r.notificationController.UpdateToReadConsultationDoctor)
	doctorRoute.DELETE("notifications/:notificationID", r.notificationController.DeleteToReadConsultationDoctor)

	// Profiles
	doctorRoute.PUT("profiles", r.doctorController.UpdateDoctorProfile)  // Update Doctor Profile
	doctorRoute.GET("profiles", r.doctorController.GetDetailProfile)    // Get Detail Profile
}
