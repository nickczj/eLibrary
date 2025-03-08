package main

import (
	"eLibrary/database"
	"eLibrary/routes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate go env -w GO111MODULE=on
//go:generate go mod tidy
//go:generate go mod download

//	@title			eLibrary API
//	@version		1.0
//	@description	This is a sample eLibrary server.

//	@contact.name	Nick Chow
//	@contact.email	x@nickczj.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3000
//	@BasePath	/elibrary/v1

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	if !gin.IsDebugging() {
		log.SetLevel(log.InfoLevel)
	}

	database.Init()
	r := routes.SetupRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run(":3000")
	if err != nil {
		log.Error("Error running app: ", err)
		return
	}
}
