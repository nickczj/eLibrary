package helper

import (
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
			c.JSON(http.StatusNotFound, gin.H{"bad request": "book not found"})
		} else {
			log.Error("Error while processing service response", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "operation failed", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"book": CreateBookDetailDTO(book)})
}

func HandleLoanResponse(c *gin.Context, loan model.LoanDetail, err error, successMessage string) {
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "loan not found", "details": err.Error()})
		case errors.Is(err, elibErr.BookNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "book not found", "details": err.Error()})
		case errors.Is(err, elibErr.UserNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found", "details": err.Error()})
		case errors.Is(err, elibErr.LoanAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "loan already exists", "details": err.Error()})
		case errors.Is(err, elibErr.NoLoanFound):
			c.JSON(http.StatusConflict, gin.H{"error": "loan not found", "details": err.Error()})
		default:
			log.Error("Error while processing service response", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "operation failed", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": successMessage,
		"loan":    CreateLoanDetailDTO(loan),
	})
}

func HandleCreateBookResponse(c *gin.Context, book model.BookDetail, err error, successMessage string) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create book", "details:": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": successMessage,
		"loan":    CreateBookDetailDTO(book),
	})
}

func HandleCreateUserResponse(c *gin.Context, user model.User, err error, successMessage string) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create user", "details:": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": successMessage,
		"loan":    CreateUserDTO(user),
	})
}
