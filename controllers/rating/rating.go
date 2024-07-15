package rating

import (
	"capstone/controllers/rating/request"
	"capstone/controllers/rating/response"
	ratingEntities "capstone/entities/rating"
	"capstone/utilities"
	"capstone/utilities/base"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RatingController struct {
	ratingUseCase ratingEntities.UseCaseInterface
}

func NewRatingController(ratingUseCase ratingEntities.UseCaseInterface) *RatingController {
	return &RatingController{
		ratingUseCase: ratingUseCase,
	}
}

func (ratingController *RatingController) SendFeedback(c echo.Context) error {
	var feedBackReq request.SendFeedbackRequest
	c.Bind(&feedBackReq)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	ratingEnt := ratingEntities.Rating{
		DoctorId: feedBackReq.DoctorId,
		UserId:   uint(userId),
		Rate:     feedBackReq.Rate,
		Message:  feedBackReq.Message,
	}

	result, err := ratingController.ratingUseCase.SendFeedback(ratingEnt)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	responseResult := response.RatingCreateResponse{
		Id:       result.Id,
		UserId:   result.UserId,
		DoctorId: result.DoctorId,
		Rate:     result.Rate,
		Message:  result.Message,
	}

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success Send Feedback", responseResult))
}

func (ratingController *RatingController) GetAllFeedbacks(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	result, err := ratingController.ratingUseCase.GetAllFeedbacks(*metadata, uint(doctorId))
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	responseResult := make([]response.RatingGetResponse, len(result))
	for i, rating := range result {
		responseResult[i] = response.RatingGetResponse{
			Id:      rating.Id,
			User:    response.UserRatingGetResponse{
				Id: rating.UserId, 
				Name: rating.User.Name, 
				Username: rating.User.Username, 
				ImageUrl: rating.User.ProfilePicture,
			},
			Rate:    rating.Rate,
			Message: rating.Message,
			Date:    rating.Date,
		}
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get All Feedbacks", metadata, responseResult))
}

func (ratingController *RatingController) GetSummaryRating(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	result, err := ratingController.ratingUseCase.GetSummaryRating(uint(doctorId))
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := response.RatingSummaryResponse{
		OneStarCount: result.OneStarCount,
		TwoStarCount: result.TwoStarCount,
		ThreeStarCount: result.ThreeStarCount,
		FourStarCount: result.FourStarCount,
		FiveStarCount: result.FiveStarCount,
		Average: result.Average,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Summary Rating", resp))
}