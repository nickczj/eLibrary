package elibErr

import "errors"

var BookNotFound = errors.New("A book with this title doesn't exist")
var UserNotFound = errors.New("A user with this id doesn't exist")
var LoanAlreadyExists = errors.New("A loan for this book already exists")
var NoLoanFound = errors.New("there is no existing loan for this book")
