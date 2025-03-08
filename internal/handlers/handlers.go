package handlers

import (
	"eLibrary/internal/helper"
	"eLibrary/internal/service"
	"eLibrary/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBook(c *gin.Context) {
	title := c.Param("title")
	if !helper.IsValidBookTitle(title) {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "invalid book title provided"})
		return
	}

	book, err := service.GetBook(title)
	helper.HandleGetBookResponse(c, book, err)
}

func BorrowBook(c *gin.Context) {
	var loanRequest model.LoanRequest

	if !helper.HandleAndValidateRequest(c, &loanRequest) {
		return
	}

	loan, err := service.BorrowBook(loanRequest)
	helper.HandleLoanResponse(c, loan, err, "book borrowed successfully")
}

func ExtendBook(c *gin.Context) {
	var loanRequest model.LoanRequest

	if !helper.HandleAndValidateRequest(c, &loanRequest) {
		return
	}

	loan, err := service.ExtendBook(loanRequest)
	helper.HandleLoanResponse(c, loan, err, "loan extended successfully")
}

func ReturnBook(c *gin.Context) {
	var loanRequest model.LoanRequest

	if !helper.HandleAndValidateRequest(c, &loanRequest) {
		return
	}

	loan, err := service.ReturnBook(loanRequest)
	helper.HandleLoanResponse(c, loan, err, "book returned successfully")
}

func CreateBook(c *gin.Context) {
	var bookDetail model.BookDetail

	if !helper.ValidateBookCreationRequest(c, bookDetail) {
		return
	}

	book, err := service.CreateBook(bookDetail)
	helper.HandleCreateBookResponse(c, book, err, "book created successfully")
}

func CreateUser(c *gin.Context) {
	user := model.User{}

	if !helper.ValidateUserCreationRequest(c, user) {
		return
	}

	user, err := service.CreateUser(user.FirstName, user.LastName, user.Username, user.Email)
	helper.HandleCreateUserResponse(c, user, err, "user created successfully")
}
