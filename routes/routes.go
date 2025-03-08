package routes

import (
	"eLibrary/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//r.Use(middleware.Logger())

	eLibrary := r.Group("/elibrary/v1")
	{
		eLibrary.GET("/book/:title", handlers.GetBook)
		eLibrary.POST("/borrow", handlers.BorrowBook)
		eLibrary.POST("/extend", handlers.ExtendBook)
		eLibrary.POST("/return", handlers.ReturnBook)
		eLibrary.POST("/create-book", handlers.CreateBook)
		eLibrary.POST("/create-user", handlers.CreateUser)
	}

	return r
}
