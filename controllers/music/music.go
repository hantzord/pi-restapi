package music

import (
	"capstone/controllers/music/request"
	"capstone/controllers/music/response"
	musicEntities "capstone/entities/music"
	"capstone/utilities"
	"capstone/utilities/base"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MusicController struct {
	musicUseCase musicEntities.UseCaseInterface
}

func NewMusicController(musicUseCase musicEntities.UseCaseInterface) *MusicController {
	return &MusicController{
		musicUseCase: musicUseCase,
	}
}

func (musicController *MusicController) GetAllMusics(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	searchParam := c.QueryParam("search")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	musics, err := musicController.musicUseCase.GetAllMusics(*metadata, userId, searchParam)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	musicResp := make([]response.MusicResponse, len(musics))

	for i, music := range musics {
		musicResp[i] = response.MusicResponse{
			Id:        music.Id,
			Title:     music.Title,
			Singer:    music.Singer,
			MusicUrl:  music.MusicUrl,
			ImageUrl:  music.ImageUrl,
			IsLiked:   music.IsLiked,
		}
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get All Musics", metadata, musicResp))
}

func (musicController *MusicController) GetAllMusicsByDoctorId(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	sortParam := c.QueryParam("sort")
	orderParam := c.QueryParam("order")
	searchParam := c.QueryParam("search")

	metadata := utilities.GetFullMetadata(pageParam, limitParam, sortParam, orderParam, searchParam)

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	musics, err := musicController.musicUseCase.GetAllMusicsByDoctorId(*metadata, doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	musicResp := make([]response.MusicGetDoctorResponse, len(musics))

	for i, music := range musics {
		musicResp[i] = response.MusicGetDoctorResponse{
			Id:        music.Id,
			Title:     music.Title,
			Singer:    music.Singer,
			MusicUrl:  music.MusicUrl,
			ImageUrl:  music.ImageUrl,
			ViewCount: music.ViewCount,
		}
	}

	return c.JSON(http.StatusOK, base.NewMetadataFullSuccessResponse("Success Get All Musics By Doctor Id", metadata, musicResp))
}

func (musicController *MusicController) GetMusicByID(c echo.Context) error {
	strId := c.Param("id")
	id, _ := strconv.Atoi(strId)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	music, err := musicController.musicUseCase.GetMusicById(id, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	musicResp := response.MusicResponse{
		Id:        music.Id,
		Title:     music.Title,
		Singer:    music.Singer,
		MusicUrl:  music.MusicUrl,
		ImageUrl:  music.ImageUrl,
		IsLiked:   music.IsLiked,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Music By Id", musicResp))
}

func (musicController *MusicController) GetLikedMusics(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	musics, err := musicController.musicUseCase.GetLikedMusics(*metadata, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	musicResp := make([]response.MusicResponse, len(musics))

	for i, music := range musics {
		musicResp[i] = response.MusicResponse{
			Id:        music.Id,
			Title:     music.Title,
			Singer:    music.Singer,
			MusicUrl:  music.MusicUrl,
			ImageUrl:  music.ImageUrl,
			IsLiked:   music.IsLiked,
		}
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get Liked Musics", metadata, musicResp))
}

func (musicController *MusicController) LikeMusic(c echo.Context) error {
	var req request.MusicLikeRequest

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	err = musicController.musicUseCase.LikeMusic(req.MusicId, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success Like Music", nil))
}

func (musicController *MusicController) UnlikeMusic(c echo.Context) error {
	var req request.MusicLikeRequest

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	err = musicController.musicUseCase.UnlikeMusic(req.MusicId, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Unlike Music", nil))
}

func (musicController *MusicController) CountMusicByDoctorId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	counter, err := musicController.musicUseCase.CountMusicByDoctorId(doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	counterResp := response.MusicCounter{
		Count: counter,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Music By Doctor Id", counterResp))
}

func (musicController *MusicController) CountMusicLikesByDoctorId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	counter, err := musicController.musicUseCase.CountMusicLikesByDoctorId(doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	counterResp := response.MusicCounter{
		Count: counter,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Music Liked By Doctor Id", counterResp))
}

func (musicController *MusicController) CountMusicViewCountByDoctorId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	counter, err := musicController.musicUseCase.CountMusicViewCountByDoctorId(doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	counterResp := response.MusicCounter{
		Count: counter,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Music Viewed By Doctor Id", counterResp))
}

func (musicController *MusicController) CountMusicViewByMonth(c echo.Context) error {
	startMonth := c.QueryParam("start_month")
	endMonth := c.QueryParam("end_month")

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	res, err := musicController.musicUseCase.CountMusicViewByMonth(doctorId, startMonth, endMonth)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Music View By Month", res))
}

func (musicController *MusicController) PostMusic(c echo.Context) error {
	var req request.MusicPostRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	fileImage, _ := c.FormFile("image")

	fileMusic, _ := c.FormFile("music")

	var musicEnt musicEntities.Music
	musicEnt.DoctorId = uint(doctorId)
	musicEnt.Title = req.Title
	musicEnt.Singer = req.Singer

	music, err := musicController.musicUseCase.PostMusic(musicEnt, fileImage, fileMusic)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.MusicGetDoctorResponse
	resp.Id = music.Id
	resp.Title = music.Title
	resp.Singer = music.Singer
	resp.MusicUrl = music.MusicUrl
	resp.ImageUrl = music.ImageUrl
	resp.ViewCount = music.ViewCount

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success Post Music", resp))
}

func (musicController *MusicController) GetMusicByIdForDoctor(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	music, err := musicController.musicUseCase.GetMusicByIdForDoctor(id)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.MusicGetDoctorResponse
	resp.Id = music.Id
	resp.Title = music.Title
	resp.Singer = music.Singer
	resp.MusicUrl = music.MusicUrl
	resp.ImageUrl = music.ImageUrl
	resp.ViewCount = music.ViewCount

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Music By Id", resp))
}

func (musicController *MusicController) EditMusic(c echo.Context) error {
	musicId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	var req request.StoryEditRequest
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	file, _ := c.FormFile("image")

	musicEnt := musicEntities.Music{
		Id:       uint(musicId),
		Title:    req.Title,
		Singer:   req.Singer,
	}

	music, err := musicController.musicUseCase.EditMusic(musicEnt, file)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.MusicGetDoctorResponse
	resp.Id = music.Id
	resp.Title = music.Title
	resp.Singer = music.Singer
	resp.MusicUrl = music.MusicUrl
	resp.ImageUrl = music.ImageUrl
	resp.ViewCount = music.ViewCount

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Edit Music", resp))
}

func (musicController *MusicController) DeleteMusic(c echo.Context) error {
	musicId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	err = musicController.musicUseCase.DeleteMusic(musicId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Delete Music", nil))
}