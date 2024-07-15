package forum

import (
	"capstone/controllers/forum/request"
	"capstone/controllers/forum/response"
	forumEntities "capstone/entities/forum"
	"capstone/utilities"
	"capstone/utilities/base"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ForumController struct {
	forumUseCase forumEntities.UseCaseInterface
}

func NewForumController(forumUseCase forumEntities.UseCaseInterface) *ForumController {
	return &ForumController{
		forumUseCase: forumUseCase,
	}
}

func (forumController *ForumController) JoinForum(c echo.Context) error {
	var req request.ForumJoinRequest

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	err = forumController.forumUseCase.JoinForum(req.ForumID, uint(userId))
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success Join Forum", nil))
}

func (forumController *ForumController) GetJoinedForum(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	search := c.QueryParam("search")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	forums, err := forumController.forumUseCase.GetJoinedForum(uint(userId), *metadata, search)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp []response.ForumJoinedResponse
	for _, forum := range forums {
		resp = append(resp, response.ForumJoinedResponse{
			ForumID:         forum.ID,
			Name:            forum.Name,
			ImageUrl:        forum.ImageUrl,
			NumberOfMembers: forum.NumberOfMembers,
		})

		for _, user := range forum.User {
			resp[len(resp)-1].User = append(resp[len(resp)-1].User, response.UserJoined{
				UserID:   uint(user.Id),
				ProfilePicture: user.ProfilePicture,
			})
		}
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get Joined Forum", metadata, resp))
}

func (forumController *ForumController) GetRecommendationForum(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	search := c.QueryParam("search")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	forums, err := forumController.forumUseCase.GetRecommendationForum(uint(userId), *metadata, search)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp []response.ForumRecommendationResponse
	for _, forum := range forums {
		resp = append(resp, response.ForumRecommendationResponse{
			ForumID:         forum.ID,
			Name:            forum.Name,
			ImageUrl:        forum.ImageUrl,
			NumberOfMembers: forum.NumberOfMembers,
		})
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get Recommendation Forum", metadata, resp))
}

func (forumController *ForumController) GetForumById(c echo.Context) error {
	forumId := c.Param("id")
	forumIdInt, _ := strconv.Atoi(forumId)

	forum, err := forumController.forumUseCase.GetForumById(uint(forumIdInt))
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.ForumDetailResponse
	resp.ForumID = forum.ID
	resp.Name = forum.Name
	resp.Description = forum.Description
	resp.ImageUrl = forum.ImageUrl

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Forum By Id", resp))
}

func (forumController *ForumController) LeaveForum(c echo.Context) error {
	forumId := c.Param("id")
	forumIdInt, _ := strconv.Atoi(forumId)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	err := forumController.forumUseCase.LeaveForum(uint(forumIdInt), uint(userId))
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Leave Forum", nil))
}

func (forumController *ForumController) CreateForum(c echo.Context) error {
	var req request.ForumCreateRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	file, _ := c.FormFile("image")

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	var forumEnt forumEntities.Forum
	forumEnt.Name = req.Name
	forumEnt.Description = req.Description
	forumEnt.DoctorID = uint(doctorId)

	forum, err := forumController.forumUseCase.CreateForum(forumEnt, file)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.ForumDetailResponse
	resp.ForumID = forum.ID
	resp.Name = forum.Name
	resp.Description = forum.Description
	resp.ImageUrl = forum.ImageUrl

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success Create Forum", resp))
}

func (forumController *ForumController) GetAllForumsByDoctorId(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	searchParam := c.QueryParam("search")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	forums, err := forumController.forumUseCase.GetAllForumsByDoctorId(uint(doctorId), *metadata, searchParam)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp []response.ForumGetDoctorResponse
	for _, forum := range forums {
		resp = append(resp, response.ForumGetDoctorResponse{
			ID:              forum.ID,
			Name:            forum.Name,
			ImageUrl:        forum.ImageUrl,
			NumberOfMembers: forum.NumberOfMembers,
		})
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get Forum By Doctor Id", metadata, resp))
}

func (forumController *ForumController) UpdateForum(c echo.Context) error {
	var req request.ForumCreateRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	forumId := c.Param("id")
	forumIdInt, _ := strconv.Atoi(forumId)

	forum := forumEntities.Forum{
		ID:          uint(forumIdInt),
		Name:        req.Name,
		Description: req.Description,
	}

	file, _ := c.FormFile("image")

	forum, err = forumController.forumUseCase.UpdateForum(forum, file)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.ForumDetailResponse
	resp.ForumID = forum.ID
	resp.Name = forum.Name
	resp.Description = forum.Description
	resp.ImageUrl = forum.ImageUrl

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Update Forum", resp))
}

func (forumController *ForumController) DeleteForum(c echo.Context) error {
	forumId := c.Param("id")
	forumIdInt, _ := strconv.Atoi(forumId)

	err := forumController.forumUseCase.DeleteForum(uint(forumIdInt))
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Delete Forum", nil))
}

func (forumController *ForumController) GetForumMemberByForumId(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	
	metadata := utilities.GetMetadata(pageParam, limitParam)

	forumId := c.Param("forumId")
	forumIdInt, _ := strconv.Atoi(forumId)

	users, err := forumController.forumUseCase.GetForumMemberByForumId(uint(forumIdInt), *metadata)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp []response.ForumMemberResponse
	for _, user := range users {
		resp = append(resp, response.ForumMemberResponse{
			ID:       uint(user.Id),
			Username: user.Username,
			Name:     user.Name,
			ImageUrl: user.ProfilePicture,
		})
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get Forum Member", metadata, resp))
}