package story

import (
	"capstone/controllers/story/request"
	"capstone/controllers/story/response"
	storyEntities "capstone/entities/story"
	"capstone/utilities"
	"capstone/utilities/base"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type StoryController struct {
	storyUseCase storyEntities.UseCaseInterface
}

func NewStoryController(storyUseCase storyEntities.UseCaseInterface) *StoryController {
	return &StoryController{
		storyUseCase: storyUseCase,
	}
}

func (storyController *StoryController) GetAllStories(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	search := c.QueryParam("search")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	stories, err := storyController.storyUseCase.GetAllStories(*metadata, userId, search)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	storiesResp := make([]response.StoriesGetAllResponse, len(stories))

	for i, story := range stories {
		storiesResp[i] = response.StoriesGetAllResponse{
			ID:        story.Id,
			Title:     story.Title,
			Content:   story.Content,
			Date:      story.Date,
			ImageUrl:  story.ImageUrl,
			IsLiked:   story.IsLiked,
			Doctor: response.DoctorGetAllResponse{
				ID:   story.Doctor.ID,
				Name: story.Doctor.Name,
			},
		}
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get All Stories", metadata, storiesResp))
}

func (storyController *StoryController) GetStoryById(c echo.Context) error {
	strId := c.Param("id")
	storyId, _ := strconv.Atoi(strId)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	story, err := storyController.storyUseCase.GetStoryById(storyId, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	storyResp := response.StoriesGetAllResponse{
		ID:        story.Id,
		Title:     story.Title,
		Content:   story.Content,
		Date:      story.Date,
		ImageUrl:  story.ImageUrl,
		IsLiked:   story.IsLiked,
		Doctor: response.DoctorGetAllResponse{
			ID:   story.Doctor.ID,
			Name: story.Doctor.Name,
		},
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Story By Id", storyResp))
}

func (storyController *StoryController) GetLikedStories(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	stories, err := storyController.storyUseCase.GetLikedStories(*metadata, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	storiesResp := make([]response.StoriesGetAllResponse, len(stories))

	for i, story := range stories {
		storiesResp[i] = response.StoriesGetAllResponse{
			ID:        story.Id,
			Title:     story.Title,
			Content:   story.Content,
			Date:      story.Date,
			ImageUrl:  story.ImageUrl,
			IsLiked:   story.IsLiked,
			Doctor: response.DoctorGetAllResponse{
				ID:   story.Doctor.ID,
				Name: story.Doctor.Name,
			},
		}
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get Liked Stories", metadata, storiesResp))
}

func (storyController *StoryController) LikeStory(c echo.Context) error {
	var req request.StoryLike

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	err = storyController.storyUseCase.LikeStory(req.StoryId, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success Like Story", nil))
}

func (storyController *StoryController) UnlikeStory(c echo.Context) error {
	var req request.StoryLike

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	err = storyController.storyUseCase.UnlikeStory(req.StoryId, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Unlike Story", nil))
}

func (storyController *StoryController) CountStoriesByDoctorId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	count, err := storyController.storyUseCase.CountStoriesByDoctorId(doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := response.StoriesCounter{
		Count: count,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Stories By Doctor Id", resp))
}

func (storyController *StoryController) CountStoryLikesByDoctorId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	count, err := storyController.storyUseCase.CountStoryLikesByDoctorId(doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := response.StoriesCounter{
		Count: count,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Story Likes By Doctor Id", resp))
}

func (storyController *StoryController) CountStoryViewByDoctorId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	count, err := storyController.storyUseCase.CountStoryViewByDoctorId(doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := response.StoriesCounter{
		Count: count,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Story View Count By Doctor Id", resp))
}

func (storyController *StoryController) CountStoryViewByMonth(c echo.Context) error {
	startMonth := c.QueryParam("start_month")
	endMonth := c.QueryParam("end_month")

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	res, err := storyController.storyUseCase.CountStoryViewByMonth(doctorId, startMonth, endMonth)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Story View Count By Month", res))
}

func (storyController *StoryController) PostStory(c echo.Context) error {
	var req request.StoryPostRequest

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	file, _ := c.FormFile("image")

	var storyEnt storyEntities.Story
	storyEnt.Title = req.Title
	storyEnt.Content = req.Content
	storyEnt.ImageUrl = file.Filename
	storyEnt.Date = time.Now()
	storyEnt.DoctorId = uint(doctorId)

	story, err := storyController.storyUseCase.PostStory(storyEnt, file)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.StoriesGetDoctorResponse
	resp.ID = int(story.Id)
	resp.Title = story.Title
	resp.Content = story.Content
	resp.Date = story.Date
	resp.ImageUrl = story.ImageUrl
	resp.ViewCount = story.ViewCount

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success Post Story", resp))
}

func (storyController *StoryController) GetStoryByIdForDoctor(c echo.Context) error {
	storyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	story, err := storyController.storyUseCase.GetStoryByIdForDoctor(storyId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.StoriesGetDoctorResponse
	resp.ID = int(story.Id)
	resp.Title = story.Title
	resp.Content = story.Content
	resp.Date = story.Date
	resp.ImageUrl = story.ImageUrl
	resp.ViewCount = story.ViewCount

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Story By Id For Doctor", resp))
}

func (storyController *StoryController) GetAllStoriesByDoctorId(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	sortParam := c.QueryParam("sort")
	orderParam := c.QueryParam("order")
	searchParam := c.QueryParam("search")

	metadata := utilities.GetFullMetadata(pageParam, limitParam, sortParam, orderParam, searchParam)

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	stories, err := storyController.storyUseCase.GetAllStoriesByDoctorId(*metadata, doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	storyResp := make([]response.StoriesGetDoctorResponse, len(stories))

	for i, story := range stories {
		storyResp[i] = response.StoriesGetDoctorResponse{
			ID:        int(story.Id),
			Title:     story.Title,
			Content:   story.Content,
			Date:      story.Date,
			ImageUrl:  story.ImageUrl,
			ViewCount: story.ViewCount,
		}
	}

	return c.JSON(http.StatusOK, base.NewMetadataFullSuccessResponse("Success Get All Stories By Doctor Id", metadata, storyResp))
}

func (storyController *StoryController) EditStory(c echo.Context) error {
	storyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	var req request.StoryPostRequest
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	file, _ := c.FormFile("image")

	storyEnt := storyEntities.Story{
		Id:       uint(storyId),
		Title:    req.Title,
		Content:  req.Content,
	}

	story, err := storyController.storyUseCase.EditStory(storyEnt, file)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.StoriesGetDoctorResponse
	resp.ID = int(story.Id)
	resp.Title = story.Title
	resp.Content = story.Content
	resp.Date = story.Date
	resp.ImageUrl = story.ImageUrl
	resp.ViewCount = story.ViewCount

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Edit Story", resp))
}

func (storyController *StoryController) DeleteStory(c echo.Context) error {
	storyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	err = storyController.storyUseCase.DeleteStory(storyId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Delete Story", nil))
}