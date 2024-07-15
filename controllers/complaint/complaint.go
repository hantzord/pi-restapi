package complaint

import (
	"capstone/constants"
	"capstone/controllers/complaint/request"
	"capstone/controllers/complaint/response"
	response2 "capstone/controllers/consultation/response"
	complaintUseCase "capstone/entities/complaint"
	consultationEntities "capstone/entities/consultation"
	"capstone/utilities"
	"capstone/utilities/base"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ComplaintController struct {
	complaintUseCase    complaintUseCase.ComplaintUseCase
	consultationUseCase consultationEntities.ConsultationUseCase
	validator           *validator.Validate
}

func NewComplaintController(complaint complaintUseCase.ComplaintUseCase, consultationUseCase consultationEntities.ConsultationUseCase, validator *validator.Validate) *ComplaintController {
	return &ComplaintController{
		complaintUseCase:    complaint,
		consultationUseCase: consultationUseCase,
		validator:           validator,
	}
}

func (controller *ComplaintController) Create(c echo.Context) error {
	var complaintRequest request.ComplaintRequest
	c.Bind(&complaintRequest)
	if err := controller.validator.Struct(complaintRequest); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(constants.ErrBadRequest.Error()))
	}
	complaint, err := controller.complaintUseCase.Create(complaintRequest.ToEntities())
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Complaint Created", complaint.ToResponse()))
}

func (controller *ComplaintController) GetAllByDoctorID(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	metadata := utilities.GetMetadata(pageParam, limitParam)

	token := c.Request().Header.Get("Authorization")
	doctorID, err := utilities.GetUserIdFromToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}
	complaintsResponse, err := controller.complaintUseCase.GetAllByDoctorID(metadata, doctorID)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var complaints []response.ComplaintResponse
	for _, value := range *complaintsResponse {
		complaints = append(complaints, *value.ToResponse())
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Complaints Retrieved", complaints))
}

func (controller *ComplaintController) GetByComplaintID(c echo.Context) error {
	strComplaintID := c.Param("id")
	complaintID, err := strconv.Atoi(strComplaintID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	consultation, err := controller.consultationUseCase.GetConsultationByComplaintID(complaintID)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Complaint Retrieved", consultation.ToDoctorResponse()))
}

func (controller *ComplaintController) SearchComplaintByPatientName(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	name := c.QueryParam("name")

	metadata := utilities.GetMetadata(pageParam, limitParam)
	token := c.Request().Header.Get("Authorization")
	doctorID, err := utilities.GetUserIdFromToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	complaints, err := controller.consultationUseCase.SearchConsultationByComplaintName(metadata, doctorID, name)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var complaintsResponse []response2.ConsultationDoctorResponse
	for _, value := range *complaints {
		complaintsResponse = append(complaintsResponse, *value.ToDoctorResponse())
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Complaint Retrieved", complaintsResponse))
}

func (controller *ComplaintController) GetConsultationByComplaintID(c echo.Context) error {
	strComplaintID := c.Param("id")
	complaintID, err := strconv.Atoi(strComplaintID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	consultation, err := controller.consultationUseCase.GetConsultationByComplaintID(complaintID)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Consultation Retrieved", consultation.ToDoctorResponse()))
}

func (controller *ComplaintController) GetAllComplaint(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	metadata := utilities.GetMetadata(pageParam, limitParam)

	doctorID, err := utilities.GetUserIdFromToken(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	complaints, err := controller.consultationUseCase.GetDoctorConsultationByComplaint(metadata, doctorID)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var complaintsResponse []response2.ConsultationDoctorResponse
	for _, value := range *complaints {
		complaintsResponse = append(complaintsResponse, *value.ToDoctorResponse())
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Complaint Retrieved", complaintsResponse))
}
