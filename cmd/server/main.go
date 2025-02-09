package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"main.go/internal/config"
	"main.go/internal/storage/sqlite"
)

func main() {
	router := gin.Default()

	router.GET("/vac", GetVacancy)
	router.Run("0.0.0.0:4252")
}

func GetVacancy(ctx *gin.Context) {
	cfg := config.MustLoad()
	storage, err := sqlite.CreateVacancyTable(cfg.StoragePath)
	if err != nil {
		log.Fatal("error in CreateVacancy Table", err)
	}
	response, err := storage.GetAllVacancy()
	if err != nil {
		ctx.JSON(404, "ERROR IN GET ALL VACANCY in SQLITE")
	}
	ctx.JSON(200, response)
}
