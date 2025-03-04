package handlers

import (
	"eLibrary/internal/service"
	"eLibrary/model"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"regexp"
)

var validate = validator.New()

func GetBook(c *gin.Context) {
	title := c.Param("title")
	if !isValidBookTitle(title) {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "invalid book title provided"})
	} else if book, err := service.GetBook(title); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"bad request": "book not found"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to retrieve net worth"})
	} else {
		c.JSON(http.StatusOK, gin.H{"book": book})
	}
}

func BorrowBook(c *gin.Context) {
	var loanRequest model.LoanRequest
	if err := c.ShouldBindJSON(&loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "invalid request body", "details:": err.Error()})
	} else if err := validate.Struct(loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "validation failed", "details:": err.Error()})
	} else if loan, err := service.BorrowBook(loanRequest); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, gin.H{"bad request": "there are no more available books to borrow"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to borrow book", "details:": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"loan": loan})
	}
}

func ExtendBook(c *gin.Context) {
	var loanRequest model.LoanRequest
	if err := c.ShouldBindJSON(&loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "invalid request body", "details:": err.Error()})
	} else if err := validate.Struct(loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "validation failed", "details:": err.Error()})
	} else if !isValidBookTitle(loanRequest.Title) {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "invalid book title provided"})
	} else if loan, err := service.ExtendBook(loanRequest.UserId, loanRequest.Title); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, gin.H{"bad request": "loan not found", "details:": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"loan": loan})
	}
}

func ReturnBook(c *gin.Context) {
	var loanRequest model.LoanRequest
	if err := c.ShouldBindJSON(&loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "invalid request body", "details:": err.Error()})
	} else if err := validate.Struct(loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "validation failed", "details:": err.Error()})
	} else if !isValidBookTitle(loanRequest.Title) {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "invalid book title provided"})
	} else if loan, err := service.ReturnBook(loanRequest.UserId, loanRequest.Title); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, gin.H{"bad request": "loan not found", "details:": err.Error()})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal server error": "something went wrong", "details:": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"loan": loan})
	}
}

func CreateBook(c *gin.Context) {
	bookDetail := model.BookDetail{}
	if err := c.ShouldBindJSON(&bookDetail); err != nil {
	}
	book, err := service.CreateBook(bookDetail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create book", "details:": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func CreateUser(c *gin.Context) {
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
	}
	user, err := service.CreateUser(user.FirstName, user.LastName, user.Username, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create user", "details:": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func isValidBookTitle(title string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9\s.,'":;!?-]+$`)
	return re.MatchString(title)
}
