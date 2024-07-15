package transaction

import (
	"capstone/constants"
	"capstone/controllers/transaction/request"
	"capstone/controllers/transaction/response"
	midtransEntities "capstone/entities/midtrans"
	transactionEntities "capstone/entities/transaction"
	"capstone/utilities"
	"capstone/utilities/base"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TransactionController struct {
	transactionUseCase transactionEntities.TransactionUseCase
	midtransUseCase    midtransEntities.MidtransUseCase
	validator          *validator.Validate
}

func NewTransactionController(transactionUseCase transactionEntities.TransactionUseCase, midtransUseCase midtransEntities.MidtransUseCase, validator *validator.Validate) *TransactionController {
	return &TransactionController{
		transactionUseCase: transactionUseCase,
		midtransUseCase:    midtransUseCase,
		validator:          validator,
	}
}

func (controller *TransactionController) InsertWithBuiltIn(c echo.Context) error {
	var transactionRequest request.TransactionRequest
	if err := c.Bind(&transactionRequest); err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	if err := controller.validator.Struct(transactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(constants.ErrBadRequest.Error()))
	}

	userId, err := utilities.GetUserIdFromToken(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, base.NewErrorResponse(constants.ErrUnauthorized.Error()))
	}

	transactionResponse, err := controller.transactionUseCase.InsertWithBuiltInInterface(transactionRequest.ToEntities(), transactionRequest.UsePoint, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusCreated, transactionResponse.ToUserResponse())
}

func (controller *TransactionController) FindByID(c echo.Context) error {
	id := c.Param("id")
	transactionResponse, err := controller.transactionUseCase.FindByID(id)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, transactionResponse.ToUserResponse())
}

func (controller *TransactionController) FindByConsultationID(c echo.Context) error {
	c.Param("id")
	transactionResponse, err := controller.transactionUseCase.FindByConsultationID(1)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, transactionResponse.ToUserResponse())
}

func (controller *TransactionController) FindAll(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	status := c.QueryParam("status")

	metadata := utilities.GetMetadata(pageParam, limitParam)
	token := c.Request().Header.Get("Authorization")
	userId, err := utilities.GetUserIdFromToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	transactions, err := controller.transactionUseCase.FindAllByUserID(metadata, uint(userId), status)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var transactionResponses []response.UserTransactionResponse
	for _, transaction := range *transactions {
		transactionResponses = append(transactionResponses, *transaction.ToUserResponse())
	}
	return c.JSON(http.StatusOK, transactionResponses)
}

func (controller *TransactionController) BankTransfer(c echo.Context) error {
	var transaction request.TransactionRequest
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}
	bankName := c.QueryParam("bank")
	transaction.Bank = bankName
	transaction.PaymentType = constants.BankTransfer
	transactionRequest := transaction.ToEntities()
	userId, err := utilities.GetUserIdFromToken(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, base.NewErrorResponse(constants.ErrUnauthorized.Error()))
	}

	transactionResponse, err := controller.transactionUseCase.InsertWithBuiltInInterface(transactionRequest, transaction.UsePoint, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Transaction created", transactionResponse.ToUserResponse()))
}

func (controller *TransactionController) EWallet(c echo.Context) error {
	var transaction request.TransactionRequest
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}
	transaction.PaymentType = constants.GoPay
	transactionRequest := transaction.ToEntities()
	userId, err := utilities.GetUserIdFromToken(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, base.NewErrorResponse(constants.ErrUnauthorized.Error()))
	}

	transactionResponse, err := controller.transactionUseCase.InsertWithBuiltInInterface(transactionRequest, transaction.UsePoint, userId)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Transaction created", transactionResponse.ToUserResponse()))
}

func (controller *TransactionController) CallbackTransaction(c echo.Context) error {
	var transactionCallback response.TransactionCallback
	if err := c.Bind(&transactionCallback); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}
	statusCode, err := controller.midtransUseCase.VerifyPayment(transactionCallback.OrderID)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	fmt.Println(transactionCallback.TransactionID)
	transaction, err := controller.transactionUseCase.ConfirmedPayment(transactionCallback.OrderID, statusCode)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, base.NewSuccessResponse("Payment confirmed", transaction.ToUserResponse()))
}

func (controller *TransactionController) CountTransactionByDoctorID(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	doctorId, err := utilities.GetUserIdFromToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}
	transactionResponse, err := controller.transactionUseCase.CountTransactionByDoctorID(uint(doctorId))
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, base.NewSuccessResponse("Count transaction success", transactionResponse))
}

func (controller *TransactionController) FindAllByDoctorID(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	query := c.QueryParam("status")

	metadata := utilities.GetMetadata(pageParam, limitParam)
	token := c.Request().Header.Get("Authorization")
	doctorId, err := utilities.GetUserIdFromToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	transactions, err := controller.transactionUseCase.FindAllByDoctorID(metadata, uint(doctorId), query)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var transactionResponses []response.DoctorTransactionResponse
	for _, transaction := range *transactions {
		transactionResponses = append(transactionResponses, *transaction.ToDoctorResponse())
	}
	return c.JSON(http.StatusOK, base.NewSuccessResponse("Get all transaction success", transactionResponses))
}

func (controller *TransactionController) DeleteTransaction(c echo.Context) error {
	strID := c.Param("id")
	err := controller.transactionUseCase.Delete(strID)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, base.NewSuccessResponse("Transaction deleted", nil))
}
