package helper

import (
	"eLibrary/model"
	"fmt"
	"time"
)

func ConstructNewLoan(loan model.LoanDetail, book model.BookDetail, user model.User) model.LoanDetail {
	loan = model.LoanDetail{
		BookDetail:     book,
		User:           user,
		NameOfBorrower: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, 28),
		IsReturned:     false,
	}
	return loan
}
