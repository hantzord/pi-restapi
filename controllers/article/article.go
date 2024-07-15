package article

import (
	"capstone/controllers/article/request"
	"capstone/controllers/article/response"
	articleUseCase "capstone/entities/article"
	"capstone/utilities"
	"capstone/utilities/base"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ArticleController struct {
	articleUseCase articleUseCase.ArticleUseCaseInterface
}

func NewArticleController(articleUseCase articleUseCase.ArticleUseCaseInterface) *ArticleController {
	return &ArticleController{
		articleUseCase: articleUseCase,
	}
}

func (controller *ArticleController) CreateArticle(c echo.Context) error {
	newArticle := new(request.CreateArticleRequest)
	if err := c.Bind(newArticle); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse("Invalid request body"))
	}

	token := c.Request().Header.Get("Authorization")
	userId, err := utilities.GetUserIdFromToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, base.NewErrorResponse("Invalid token"))
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse("Invalid image format"))
	}

	imageURL, err := utilities.UploadImage(file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, base.NewErrorResponse("Failed to upload image"))
	}

	wordCount := utilities.CountWords(newArticle.Content)
	readTime := utilities.EstimateReadingTime(wordCount)

	articleEntity := articleUseCase.Article{
		Title:       newArticle.Title,
		Content:     newArticle.Content,
		ImageUrl:    imageURL,
		DoctorID:    uint(userId),
		ViewCount:   0,
		ReadingTime: readTime,
	}

	createdArticle, err := controller.articleUseCase.CreateArticle(&articleEntity, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	articleResponse := response.ArticleCreatedResponse{
		ID:          createdArticle.ID,
		DoctorID:    createdArticle.DoctorID,
		Title:       createdArticle.Title,
		Content:     createdArticle.Content,
		Date:        createdArticle.Date,
		ImageUrl:    createdArticle.ImageUrl,
		ReadingTime: createdArticle.ReadingTime,
		Doctor: response.DoctorInfoResponse{
			ID:   createdArticle.Doctor.ID,
			Name: createdArticle.Doctor.Name,
		},
	}

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Article created successfully", articleResponse))
}
func (controller *ArticleController) GetAllArticle(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	search := c.QueryParam("search")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	articles, err := controller.articleUseCase.GetAllArticle(*metadata, userId, search)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	articleResponse := make([]response.ArticleListResponse, 0, len(articles))
	for _, article := range articles {
		articleResponse = append(articleResponse, article.ToResponse())
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get All Articles", articleResponse))
}

func (controller *ArticleController) GetArticleById(c echo.Context) error {
	strId := c.Param("id")
	articleId, _ := strconv.Atoi(strId)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	article, err := controller.articleUseCase.GetArticleById(articleId, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	articleResp := response.ArticleCreatedResponse{
		ID:          article.ID,
		DoctorID:    article.DoctorID,
		Title:       article.Title,
		Content:     article.Content,
		Date:        article.Date,
		ImageUrl:    article.ImageUrl,
		IsLiked:     article.IsLiked,
		ReadingTime: article.ReadingTime,
		Doctor: response.DoctorInfoResponse{
			ID:   article.Doctor.ID,
			Name: article.Doctor.Name,
		},
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Article By Id", articleResp))
}

func (controller *ArticleController) GetLikedArticle(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	articles, err := controller.articleUseCase.GetLikedArticle(*metadata, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	articleResponse := make([]response.ArticleCreatedResponse, len(articles))
	for i, article := range articles {
		articleResponse[i] = response.ArticleCreatedResponse{
			ID:          article.ID,
			DoctorID:    article.DoctorID,
			Title:       article.Title,
			Content:     article.Content,
			Date:        article.Date,
			ImageUrl:    article.ImageUrl,
			IsLiked:     article.IsLiked,
			ReadingTime: article.ReadingTime,
			Doctor: response.DoctorInfoResponse{
				ID:   article.Doctor.ID,
				Name: article.Doctor.Name,
			},
		}
	}

	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get Liked Articles", metadata, articleResponse))
}

func (controller *ArticleController) LikeArticle(c echo.Context) error {
	var req request.ArticleLike

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse("Invalid request body"))
	}

	token := c.Request().Header.Get("Authorization")
	userId, err := utilities.GetUserIdFromToken(token)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, base.NewErrorResponse("Invalid token"))
	}

	err = controller.articleUseCase.LikeArticle(req.ArticleID, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success Like Article", nil))
}

func (controller *ArticleController) UnlikeArticle(c echo.Context) error {
	var req request.ArticleLike

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse("Invalid request body"))
	}

	token := c.Request().Header.Get("Authorization")
	userId, err := utilities.GetUserIdFromToken(token)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, base.NewErrorResponse("Invalid token"))
	}

	err = controller.articleUseCase.UnlikeArticle(req.ArticleID, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Unlike Article", nil))
}

func (controller *ArticleController) GetArticleByIdForDoctor(c echo.Context) error {
	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	article, err := controller.articleUseCase.GetArticleByIdForDoctor(articleId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.ArticleGetDoctorResponse
	resp.ID = article.ID
	resp.DoctorID = article.DoctorID
	resp.Title = article.Title
	resp.Content = article.Content
	resp.Date = article.Date
	resp.ImageUrl = article.ImageUrl
	resp.ViewCount = article.ViewCount
	resp.ReadingTime = article.ReadingTime

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Article For Doctor", resp))
}

func (controller *ArticleController) GetAllArticleByDoctorId(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	sortParam := c.QueryParam("sort")
	orderParam := c.QueryParam("order")
	searchParam := c.QueryParam("search")

	metadata := utilities.GetFullMetadata(pageParam, limitParam, sortParam, orderParam, searchParam)

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	articles, err := controller.articleUseCase.GetAllArticleByDoctorId(*metadata, doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	articleResp := make([]response.ArticleGetDoctorResponse, len(articles))

	for i, article := range articles {
		articleResp[i] = response.ArticleGetDoctorResponse{
			ID:          article.ID,
			DoctorID:    article.DoctorID,
			Title:       article.Title,
			Content:     article.Content,
			Date:        article.Date,
			ImageUrl:    article.ImageUrl,
			ViewCount:   article.ViewCount,
			ReadingTime: article.ReadingTime,
		}
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get All Article By Doctor Id", articleResp))
}

func (controller *ArticleController) CountArticleByDoctorId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	count, err := controller.articleUseCase.CountArticleByDoctorId(doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := response.ArticleCounter{
		Count: count,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Article By Doctor Id", resp))
}

func (controller *ArticleController) CountArticleLikesByDoctorId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	count, err := controller.articleUseCase.CountArticleLikesByDoctorId(doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := response.ArticleCounter{
		Count: count,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Article Likes By Doctor Id", resp))
}

func (controller *ArticleController) CountArticleViewByDoctorId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	count, err := controller.articleUseCase.CountArticleViewByDoctorId(doctorId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := response.ArticleCounter{
		Count: count,
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Article View By Doctor Id", resp))
}

func (controller *ArticleController) CountArticleViewByMonth(c echo.Context) error {
	startMonth := c.QueryParam("start_month")
	endMonth := c.QueryParam("end_month")

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	res, err := controller.articleUseCase.CountArticleViewByMonth(doctorId, startMonth, endMonth)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Count Article View By Month", res))
}

func (controller *ArticleController) EditArticle(c echo.Context) error {
	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	var req request.UpdateArticleRequest
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	file, _ := c.FormFile("image")

	articleEnt := articleUseCase.Article{
		ID:      uint(articleId),
		Title:   req.Title,
		Content: req.Content,
	}

	article, err := controller.articleUseCase.EditArticle(articleEnt, file)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.ArticleEditResponse
	resp.ID = article.ID
	resp.Title = article.Title
	resp.Content = article.Content
	resp.Date = article.Date
	resp.ImageUrl = article.ImageUrl
	resp.ViewCount = article.ViewCount
	resp.ReadingTime = article.ReadingTime

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Edit Article", resp))
}

func (controller *ArticleController) DeleteArticle(c echo.Context) error {
	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	err = controller.articleUseCase.DeleteArticle(articleId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Delete Article", nil))
}
