package helper

import (
	"eLibrary/dto"
	"eLibrary/model"
)

func CreateBookDetailDTO(book model.BookDetail) dto.BookDetail {
	bookDTO := dto.BookDetail{
		Title:           book.Title,
		Author:          book.Author,
		AvailableCopies: book.AvailableCopies,
	}

	return bookDTO
}

func CreateLoanDetailDTO(loan model.LoanDetail) dto.Loan {
	loanDTO := dto.Loan{
		Title:      loan.BookDetail.Title,
		Author:     loan.BookDetail.Author,
		ISBN:       loan.BookDetail.ISBN,
		LoanDate:   loan.LoanDate,
		ReturnDate: loan.ReturnDate,
	}

	return loanDTO
}

func CreateUserDTO(user model.User) dto.User {
	userDTO := dto.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.FirstName,
		Email:     user.Email,
	}

	return userDTO
}
