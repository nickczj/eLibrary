package helper

import (
	"eLibrary/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"
)

var validate = validator.New()

func ValidateBookCreationRequest(c *gin.Context, bookDetail model.BookDetail) bool {
	if err := c.ShouldBindJSON(&bookDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return false
	}

	if err := ValidateBook(bookDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
	}

	return true
}

func HandleAndValidateRequest(c *gin.Context, loanRequest *model.LoanRequest) bool {
	if err := c.ShouldBindJSON(loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return false
	}

	if err := validateLoanRequest(*loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return false
	}

	return true
}

func ValidateUserCreationRequest(c *gin.Context, user model.User) bool {
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return false
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
	}

	return true
}

func validateLoanRequest(req model.LoanRequest) error {
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("invalid request structure: %w", err)
	}

	if !IsValidBookTitle(req.Title) {
		return fmt.Errorf("invalid book title provided")
	}

	return nil
}

func ValidateBook(book model.BookDetail) error {
	if err := validate.Struct(book); err != nil {
		return fmt.Errorf("invalid request structure: %w", err)
	}

	if !IsValidBookTitle(book.Title) {
		return fmt.Errorf("invalid book title provided")
	}

	return nil
}

func IsValidBookTitle(title string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9\s.,'":;!?-]+$`)
	return re.MatchString(title)
}
