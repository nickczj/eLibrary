package main

import (
	"eLibrary/database"
	"eLibrary/routes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//go:generate go env -w GO111MODULE=on
//go:generate go mod tidy
//go:generate go mod download

func main() {
	if !gin.IsDebugging() {
		log.SetLevel(log.InfoLevel)
	}

	database.Init()
	r := routes.SetupRouter()

	err := r.Run(":3000")
	if err != nil {
		log.Error("Error running app: ", err)
		return
	}
}
