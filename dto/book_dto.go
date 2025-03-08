package dto

type BookDetail struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	AvailableCopies int    `json:"available_copies"`
}
