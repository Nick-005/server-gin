package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/config"
	"main.go/internal/storage/sqlite"
)

type Vacancy_Body struct {
	Emp_ID     int    `json:"emp_id"`
	Vac_Name   string `json:"vac_name"`
	Price      int    `json:"price"`
	Location   string `json:"location"`
	Experience string `json:"exp"`
}

func main() {
	cfg := config.MustLoad()
	storage, err := InitStorage(cfg)
	if err != nil {
		log.Fatalln("Произошла ошибка в инициализации бд")
	}
	router := gin.Default()

	router.GET("/vac", GetVacancy(storage))
	router.POST("/vac", PostVacancy(storage))
	router.Run("localhost:4252")
}

func InitStorage(cfg *config.Config) (*sqlite.Storage, error) {
	_, err := sqlite.CreateVacancyTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateVacancy Table")
	}
	_, err = sqlite.CreateEmployeeTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateEmployee Table")
	}
	storage, err := sqlite.CreateTableUser(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateVacancy Table")
	}
	return storage, nil
}

func PostVacancy(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body Vacancy_Body
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, "error in parse body! Please check our body in request!")
			return
		}
		vac_id, emp_limit, err := storage.AddVacancy(body.Emp_ID, body.Vac_Name, body.Price, body.Location, body.Experience)
		if err != nil {
			ctx.JSON(404, "ERROR IN GET ALL VACANCY in SQLITE")
		}
		if vac_id == -1 {
			ctx.JSON(200, gin.H{
				"vacancyID": vac_id,
				"status":    "Error",
				"message":   "Employee has a limit",
				"Emp_limit": emp_limit,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Emp_limit": emp_limit,
			"vacancyID": vac_id,
			"status":    "Success",
		})
	}

}

func GetVacancy(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

}
