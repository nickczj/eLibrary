package helper

import (
	"eLibrary/dto"
	"eLibrary/internal/elibErr"
	"eLibrary/model"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func HandleGetBookResponse(c *gin.Context, book model.BookDetail, err error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, FailedAPIResponse("book not found", nil))
		} else {
			log.Error("Error while processing service response", "error", err)
			c.JSON(http.StatusInternalServerError, FailedAPIResponse("operation failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, model.SuccessBookResponse{Data: CreateBookDetailDTO(book)})
}

func HandleLoanResponse(c *gin.Context, loan model.LoanDetail, err error, successMessage string) {
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, FailedAPIResponse("loan not found", err.Error()))
		case errors.Is(err, elibErr.BookNotFound):
			c.JSON(http.StatusBadRequest, FailedAPIResponse("book not found", err.Error()))
		case errors.Is(err, elibErr.UserNotFound):
			c.JSON(http.StatusBadRequest, FailedAPIResponse("user not found", err.Error()))
		case errors.Is(err, elibErr.LoanAlreadyExists):
			c.JSON(http.StatusConflict, FailedAPIResponse("loan already exists", err.Error()))
		case errors.Is(err, elibErr.NoLoanFound):
			c.JSON(http.StatusConflict, FailedAPIResponse("loan not found", err.Error()))
		default:
			log.Error("Error while processing service response", "error", err)
			c.JSON(http.StatusInternalServerError, FailedAPIResponse("operation failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, SuccessLoanResponse(successMessage, CreateLoanDetailDTO(loan)))
}

func HandleCreateBookResponse(c *gin.Context, book model.BookDetail, err error, successMessage string) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, FailedAPIResponse("unable to create book", err.Error()))
	}
	c.JSON(http.StatusOK, SuccessCreateBookResponse(successMessage, CreateBookDetailDTO(book)))
}

func HandleCreateUserResponse(c *gin.Context, user model.User, err error, successMessage string) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, FailedAPIResponse("unable to create user", err.Error()))
	}

	c.JSON(http.StatusOK, SuccessCreateUserResponse(successMessage, CreateUserDTO(user)))
}

func FailedAPIResponse(message string, error interface{}) model.FailedResponse {
	r := model.FailedResponse{
		APIResponse: model.APIResponse{
			Status:  "error",
			Message: message,
		},
		ErrorDetails: error,
	}

	return r
}

func SuccessCreateBookResponse(message string, data dto.BookDetail) model.SuccessCreateBookResponse {
	r := model.SuccessCreateBookResponse{
		APIResponse: model.APIResponse{
			Status:  "success",
			Message: message,
		},
		Data: data,
	}

	return r
}

func SuccessLoanResponse(message string, loan dto.Loan) model.SuccessLoanResponse {
	r := model.SuccessLoanResponse{
		APIResponse: model.APIResponse{
			Status:  "success",
			Message: message,
		},
		Data: loan,
	}

	return r
}

func SuccessCreateUserResponse(message string, user dto.User) model.SuccessUserResponse {
	r := model.SuccessUserResponse{
		APIResponse: model.APIResponse{
			Status:  "success",
			Message: message,
		},
		Data: user,
	}

	return r
}
