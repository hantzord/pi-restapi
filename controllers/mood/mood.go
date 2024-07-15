package mood

import (
	"capstone/controllers/mood/request"
	"capstone/controllers/mood/response"
	moodEntities "capstone/entities/mood"
	"capstone/utilities"
	"capstone/utilities/base"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MoodController struct {
	moodUseCase moodEntities.UseCaseInterface
}

func NewMoodController(moodUseCase moodEntities.UseCaseInterface) *MoodController {
	return &MoodController{
		moodUseCase: moodUseCase,
	}
}

func (moodController *MoodController) CreateMood(c echo.Context) error {
	var req request.MoodCreate
	c.Bind(&req)

	file, _ := c.FormFile("image")

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	moodEnt := moodEntities.Mood{
		Message:    req.Message,
		MoodTypeId: req.MoodTypeId,
		Date:       req.Date,
	}

	moodEnt.UserId = uint(userId)

	result, err := moodController.moodUseCase.SendMood(file, moodEnt)

	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	response := response.MoodResponse{
		ID:       result.ID,
		MoodType: response.MoodTypeResponse{ID: result.MoodType.ID, Name: result.MoodType.Name},
		Date:     result.Date,
		Message:  result.Message,
		ImageUrl: result.ImageUrl,
	}

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success Create Mood", response))
}

func (moodController *MoodController) GetAllMoods(c echo.Context) error {
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	moods, err := moodController.moodUseCase.GetAllMoods(userId, startDateStr, endDateStr)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := make([]response.MoodGetAllResponse, len(moods))
	for i, mood := range moods {
		resp[i] = response.MoodGetAllResponse{
			ID:       mood.ID,
			Date:     mood.Date,
			MoodType: response.MoodTypeResponse{ID: mood.MoodType.ID, Name: mood.MoodType.Name},
		}
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get All Moods", resp))
}

func (moodController *MoodController) GetMoodById(c echo.Context) error {
	moodId := c.Param("id")
	moodIdInt, _ := strconv.Atoi(moodId)

	mood, err := moodController.moodUseCase.GetMoodById(moodIdInt)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := response.MoodResponse{
		ID:       mood.ID,
		MoodType: response.MoodTypeResponse{ID: mood.MoodType.ID, Name: mood.MoodType.Name},
		Date:     mood.Date,
		Message:  mood.Message,
		ImageUrl: mood.ImageUrl,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Mood By Id", resp))
}