package service

import (
	"eLibrary/global"
	"eLibrary/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func GetBook(title string) (book model.BookDetail, err error) {
	err = global.Database.Where("title = ?", title).First(&book).Error
	return book, err
}

func BorrowBook(request model.LoanRequest) (loan model.LoanDetail, err error) {
	book := model.BookDetail{}
	bookErr := global.Database.Where("title = ? AND available_copies > 0", request.Title).First(&book).Error
	if bookErr != nil {
		return loan, bookErr
	}

	user := model.User{}
	userErr := global.Database.Where("id = ?", request.UserId).First(&user).Error
	if userErr != nil && errors.Is(userErr, gorm.ErrRecordNotFound) {
		return loan, errors.New("user not found")
	} else if userErr != nil {
		return loan, userErr
	}

	loan = model.LoanDetail{}
	loanErr := global.Database.Where("book_id = ? AND user_id = ? AND is_returned = ?", book.ID, user.ID, false).First(&loan).Error
	if loanErr == nil {
		return loan, errors.New("a loan for this book exists")
	}

	book.AvailableCopies = book.AvailableCopies - 1
	global.Database.Save(&book)

	loan = model.LoanDetail{
		BookDetail:     book,
		User:           user,
		NameOfBorrower: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, 28),
		IsReturned:     false,
	}

	result := global.Database.Create(&loan)
	return loan, result.Error
}

func ExtendBook(userId int, title string) (loan model.LoanDetail, err error) {
	if loan, err = findLoan(userId, title); err != nil {
		return loan, err
	}

	loan.ReturnDate = loan.ReturnDate.AddDate(0, 0, 21)
	result := global.Database.Save(&loan)

	return loan, result.Error
}

func ReturnBook(userId int, title string) (loan model.LoanDetail, err error) {
	loan = model.LoanDetail{}

	err = global.Database.Preload("User").Preload("BookDetail").
		Joins("JOIN book_details ON book_details.id = loan_details.book_id").
		Where("book_details.title = ? AND loan_details.user_id = ? AND is_returned = ?", title, userId, false).
		First(&loan).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return loan, errors.New("there is no loan for this book")
	}

	loan.IsReturned = true
	result := global.Database.Save(&loan)
	if result.Error != nil {
		return loan, result.Error
	}
	
	loan.BookDetail.AvailableCopies = loan.BookDetail.AvailableCopies + 1
	saveBook := global.Database.Save(&loan.BookDetail)

	return loan, saveBook.Error
}

func CreateBook(detail model.BookDetail) (book model.BookDetail, err error) {
	book = model.BookDetail{
		Title:           detail.Title,
		Author:          detail.Author,
		ISBN:            detail.ISBN,
		AvailableCopies: detail.AvailableCopies,
	}

	result := global.Database.Create(&book)
	return book, result.Error
}

func CreateUser(firstName string, lastName string, username string, email string) (user model.User, err error) {
	user = model.User{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
	}

	result := global.Database.Create(&user)
	return user, result.Error
}

func findLoan(userId int, title string) (loan model.LoanDetail, err error) {
	loan = model.LoanDetail{}

	err = global.Database.Preload("User").Preload("BookDetail").
		Joins("JOIN book_details ON book_details.id = loan_details.book_id").
		Where("book_details.title = ? AND loan_details.user_id = ?", title, userId).
		First(&loan).Error

	return loan, err
}
