package doctor

import (
	"capstone/constants"
	"capstone/controllers/doctor/request"
	"capstone/controllers/doctor/response"
	doctorUseCase "capstone/entities/doctor"
	"capstone/utilities"
	"capstone/utilities/base"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

type DoctorController struct {
	doctorUseCase doctorUseCase.DoctorUseCaseInterface
	validator     *validator.Validate
}

func NewDoctorController(doctorUseCase doctorUseCase.DoctorUseCaseInterface, validator *validator.Validate) *DoctorController {
	return &DoctorController{
		doctorUseCase: doctorUseCase,
		validator:     validator,
	}
}

func (controller *DoctorController) Register(c echo.Context) error {
	var doctorFromRequest request.DoctorRegisterRequest
	c.Bind(&doctorFromRequest)
	if err := controller.validator.Struct(doctorFromRequest); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(constants.ErrBadRequest.Error()))
	}

	imageFromRequest, err := c.FormFile("profile_picture")
	doctorFromRequest.ProfilePicture = imageFromRequest

	doctorRequest, err := doctorFromRequest.ToDoctorEntities()
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	doctorResult, err := controller.doctorUseCase.Register(doctorRequest)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	doctorResponse := doctorResult.ToResponse()
	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Register", doctorResponse))
}

func (controller *DoctorController) Login(c echo.Context) error {
	var doctorFromRequest request.DoctorLoginRequest
	c.Bind(&doctorFromRequest)

	doctorRequest := doctorFromRequest.ToDoctorLoginEntities()
	doctorResult, err := controller.doctorUseCase.Login(doctorRequest)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	doctorResponse := doctorResult.ToResponse()
	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Login", doctorResponse))
}

func (controller *DoctorController) GetByID(c echo.Context) error {
	strDoctorID := c.Param("id")
	doctorID, err := strconv.Atoi(strDoctorID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}
	doctorResult, err := controller.doctorUseCase.GetDoctorByID(doctorID)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	doctorResponse := doctorResult.ToDoctorResponse()
	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Doctor By ID", doctorResponse))
}

func (controller *DoctorController) GetAll(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	doctorResult, err := controller.doctorUseCase.GetAllDoctor(metadata)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var doctorResponse []response.DoctorResponse
	for _, doctor := range *doctorResult {
		doctorResponse = append(doctorResponse, *doctor.ToDoctorResponse())
	}
	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get All Doctor", metadata, doctorResponse))
}

func (controller *DoctorController) GetActive(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	doctorResult, err := controller.doctorUseCase.GetActiveDoctor(metadata)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var doctorResponse []response.DoctorResponse
	for _, doctor := range *doctorResult {
		doctorResponse = append(doctorResponse, *doctor.ToDoctorResponse())
	}
	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Get Active Doctor", metadata, doctorResponse))
}

func (c *DoctorController) GoogleLogin(ctx echo.Context) error {
	url := c.doctorUseCase.HandleGoogleLogin()
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (c *DoctorController) GoogleCallback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	result, err := c.doctorUseCase.HandleGoogleCallback(ctx.Request().Context(), code)
	if err != nil {
		return ctx.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var res response.DoctorLoginAndRegisterResponse
	res.ID = result.ID
	res.Token = result.Token

	return ctx.JSON(http.StatusOK, base.NewSuccessResponse("Success Login Oauth", res))
}

func (c *DoctorController) FacebookLogin(ctx echo.Context) error {
	url := c.doctorUseCase.HandleFacebookLogin()
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (c *DoctorController) FacebookCallback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	result, err := c.doctorUseCase.HandleFacebookCallback(ctx.Request().Context(), code)
	if err != nil {
		return ctx.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var res response.DoctorLoginAndRegisterResponse
	res.ID = result.ID
	res.Token = result.Token

	return ctx.JSON(http.StatusOK, base.NewSuccessResponse("Success Login Oauth", res))
}

func (controller *DoctorController) SearchDoctor(c echo.Context) error {
	query := c.QueryParam("query")
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	metadata := utilities.GetMetadata(pageParam, limitParam)

	doctorResult, err := controller.doctorUseCase.SearchDoctor(query, metadata)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var doctorResponse []response.DoctorResponse
	for _, doctor := range *doctorResult {
		doctorResponse = append(doctorResponse, *doctor.ToDoctorResponse())
	}
	return c.JSON(http.StatusOK, base.NewMetadataSuccessResponse("Success Search Doctor", metadata, doctorResponse))
}

func (c *DoctorController) UpdateDoctorProfile(ctx echo.Context) error {
	var doctorFromRequest request.UpdateDoctorProfileRequest
	if err := ctx.Bind(&doctorFromRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	token := ctx.Request().Header.Get("Authorization")
	doctorID, err := utilities.GetUserIdFromToken(token)
	if err != nil {
		return ctx.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	imageURL, err := utilities.UploadImage(file)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, base.NewErrorResponse(err.Error()))
	}

	doctorEntities := doctorUseCase.Doctor{
		ID:               uint(doctorID),
		// Username:         doctorFromRequest.Username,
		Name:             doctorFromRequest.Name,
		// Address:          doctorFromRequest.Address,
		// PhoneNumber:      doctorFromRequest.PhoneNumber,
		Gender:           doctorFromRequest.Gender,
		ProfilePicture:   imageURL,
		// Experience:       doctorFromRequest.Experience,
		BachelorAlmamater: doctorFromRequest.BachelorAlmamater,
		// BachelorGraduationYear: doctorFromRequest.BachelorGraduationYear,
		MasterAlmamater:        doctorFromRequest.MasterAlmamater,
		// MasterGraduationYear:   doctorFromRequest.MasterGraduationYear,
		PracticeLocation: doctorFromRequest.PracticeLocation,
		// PracticeCity:     doctorFromRequest.PracticeCity,
		// Fee:              doctorFromRequest.Fee,
		Specialist:       doctorFromRequest.Specialist,
	}

	updatedDoctor, err := c.doctorUseCase.UpdateDoctorProfile(&doctorEntities)
	if err != nil {
		return ctx.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	doctorResponse := response.DoctorUpdateProfileResponse{
		ID:               updatedDoctor.ID,
		Username:         updatedDoctor.Username,
		Email:            updatedDoctor.Email,
		Name:             updatedDoctor.Name,
		Address:          updatedDoctor.Address,
		PhoneNumber:      updatedDoctor.PhoneNumber,
		Gender:           updatedDoctor.Gender,
		ProfilePicture:   updatedDoctor.ProfilePicture,
		Experience:       updatedDoctor.Experience,
		BachelorAlmamater: updatedDoctor.BachelorAlmamater,
		BachelorGraduationYear: updatedDoctor.BachelorGraduationYear,
		MasterAlmamater: updatedDoctor.MasterAlmamater,
		MasterGraduationYear: updatedDoctor.MasterGraduationYear,
		PracticeLocation: updatedDoctor.PracticeLocation,
		PracticeCity:     updatedDoctor.PracticeCity,
		Fee:              updatedDoctor.Fee,
		Specialist:       updatedDoctor.Specialist,
	}

	return ctx.JSON(http.StatusOK, base.NewSuccessResponse("Success Update Profile", doctorResponse))
}

func (c *DoctorController) GetDetailProfile(ctx echo.Context) error {
	token := ctx.Request().Header.Get("Authorization")
	doctorID, err := utilities.GetUserIdFromToken(token)
	if err != nil {
		return ctx.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	doctor, err := c.doctorUseCase.GetDetailProfile(uint(doctorID))
	if err != nil {
		return ctx.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	doctorResponse := response.DoctorResponse{
		ID:               doctor.ID,
		Username:         doctor.Username,
		Email:            doctor.Email,
		Name:             doctor.Name,
		Address:          doctor.Address,
		PhoneNumber:      doctor.PhoneNumber,
		Gender:           doctor.Gender,
		ProfilePicture:   doctor.ProfilePicture,
		Experience:       doctor.Experience,
		BachelorAlmamater: doctor.BachelorAlmamater,
		BachelorGraduationYear: doctor.BachelorGraduationYear,
		MasterAlmamater: doctor.MasterAlmamater,
		MasterGraduationYear: doctor.MasterGraduationYear,
		PracticeLocation: doctor.PracticeLocation,
		PracticeCity:     doctor.PracticeCity,
		Fee:              doctor.Fee,
		Specialist:       doctor.Specialist,
		Amount:           doctor.Amount,
		IsAvailable:      doctor.IsAvailable,
		Balance:          doctor.Balance,
		RatingPrecentage: doctor.RatingPrecentage,
	}

	return ctx.JSON(http.StatusOK, base.NewSuccessResponse("Success Get Detail Profile", doctorResponse))
}