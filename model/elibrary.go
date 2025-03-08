package model

import (
	"gorm.io/gorm"
	"time"
)

type BookDetail struct {
	gorm.Model
	Title           string `json:"title" gorm:"unique;not null" validate:"required"`
	Author          string `json:"author" gorm:"not null" validate:"required"`
	ISBN            string `json:"isbn" gorm:"not null" validate:"required"`
	AvailableCopies int    `json:"available_copies" gorm:"not null" validate:"required"`
}

type LoanDetail struct {
	gorm.Model
	BookID         uint       `json:"book_id" gorm:"not null"`
	BookDetail     BookDetail `gorm:"foreignkey:BookID"`
	UserID         uint       `json:"user_id" gorm:"not null"`
	User           User       `gorm:"foreignkey:UserID"`
	NameOfBorrower string     `json:"name_of_borrower"`
	LoanDate       time.Time  `json:"loan_date"`
	ReturnDate     time.Time  `json:"return_date"`
	IsReturned     bool       `json:"is_returned"`
}

type LoanRequest struct {
	Title  string `json:"title" validate:"required"`
	UserId int    `json:"user_id" validate:"required"`
}

type User struct {
	gorm.Model
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}
