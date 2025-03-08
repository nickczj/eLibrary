package handlers

import (
	_ "eLibrary/cmd/docs"
	_ "eLibrary/dto"
	"eLibrary/internal/helper"
	"eLibrary/internal/service"
	"eLibrary/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetBook godoc
//
//	@Summary		Retrieve Book Details
//	@Description	get book by title
//	@Produce		json
//	@Param			title	path		string	true	"Book Title"
//	@Success		200		{object}	model.SuccessBookResponse
//	@Failure		400		{object}	model.FailedResponse
//	@Failure		404		{object}	model.FailedResponse
//	@Failure		500		{object}	model.FailedResponse
//	@Router			/elibrary/v1/book/{title} [get]
func GetBook(c *gin.Context) {
	title := c.Param("title")
	if !helper.IsValidBookTitle(title) {
		c.JSON(http.StatusBadRequest, helper.FailedAPIResponse("invalid book title provided", nil))
		return
	}

	book, err := service.GetBook(title)
	helper.HandleGetBookResponse(c, book, err)
}

// BorrowBook godoc
//
//	@Summary		Borrow a book - with a 4-week loan
//	@Description	borrow a book by title and user id
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.LoanRequest	true	"Request Body"
//	@Success		200		{object}	model.SuccessBookResponse
//	@Failure		400		{object}	model.FailedResponse
//	@Failure		404		{object}	model.FailedResponse
//	@Failure		500		{object}	model.FailedResponse
//	@Router			/elibrary/v1/borrow [post]
func BorrowBook(c *gin.Context) {
	var loanRequest model.LoanRequest

	if !helper.HandleAndValidateRequest(c, &loanRequest) {
		return
	}

	loan, err := service.BorrowBook(loanRequest)
	helper.HandleLoanResponse(c, loan, err, "book borrowed successfully")
}

// ExtendBook godoc
//
//	@Summary		Extend a loan - by 3-weeks
//	@Description	extend a loan by book title and user id
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.LoanRequest	true	"Request Body"
//	@Success		200		{object}	model.SuccessBookResponse
//	@Failure		400		{object}	model.FailedResponse
//	@Failure		404		{object}	model.FailedResponse
//	@Failure		500		{object}	model.FailedResponse
//	@Router			/elibrary/v1/extend [post]
func ExtendBook(c *gin.Context) {
	var loanRequest model.LoanRequest

	if !helper.HandleAndValidateRequest(c, &loanRequest) {
		return
	}

	loan, err := service.ExtendBook(loanRequest)
	helper.HandleLoanResponse(c, loan, err, "loan extended successfully")
}

// ReturnBook godoc
//
//	@Summary		Return a book
//	@Description	return a book by title and user id
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.LoanRequest	true	"Request Body"
//	@Success		200		{object}	model.SuccessBookResponse
//	@Failure		400		{object}	model.FailedResponse
//	@Failure		404		{object}	model.FailedResponse
//	@Failure		500		{object}	model.FailedResponse
//	@Router			/elibrary/v1/return [post]
func ReturnBook(c *gin.Context) {
	var loanRequest model.LoanRequest

	if !helper.HandleAndValidateRequest(c, &loanRequest) {
		return
	}

	loan, err := service.ReturnBook(loanRequest)
	helper.HandleLoanResponse(c, loan, err, "book returned successfully")
}

// CreateBook godoc
//
//	@Summary		Create a book
//	@Description	create a book in the database
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.BookDetail	true	"Request Body"
//	@Success		200		{object}	model.SuccessBookResponse
//	@Failure		400		{object}	model.FailedResponse
//	@Failure		404		{object}	model.FailedResponse
//	@Failure		500		{object}	model.FailedResponse
//	@Router			/elibrary/v1/create-book [post]
func CreateBook(c *gin.Context) {
	var bookDetail model.BookDetail

	if !helper.ValidateBookCreationRequest(c, bookDetail) {
		return
	}

	book, err := service.CreateBook(bookDetail)
	helper.HandleCreateBookResponse(c, book, err, "book created successfully")
}

// CreateUser godoc
//
//	@Summary		Create a user
//	@Description	create a user in the eLibrary system
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.User	true	"Request Body"
//	@Success		200		{object}	model.SuccessBookResponse
//	@Failure		400		{object}	model.FailedResponse
//	@Failure		404		{object}	model.FailedResponse
//	@Failure		500		{object}	model.FailedResponse
//	@Router			/elibrary/v1/create-user [post]
func CreateUser(c *gin.Context) {
	user := model.User{}

	if !helper.ValidateUserCreationRequest(c, user) {
		return
	}

	user, err := service.CreateUser(user.FirstName, user.LastName, user.Username, user.Email)
	helper.HandleCreateUserResponse(c, user, err, "user created successfully")
}
