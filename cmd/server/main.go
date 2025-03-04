package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "main.go/cmd/server/docs"
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

type RequestEmployee struct {
	NameOrganization string `json:"nameOrg"`
	PhoneNumber      string `json:"phoneNumber"`
	Email            string `json:"email"`
	Geography        string `json:"geography"`
	About            string `json:"about"`
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
func main() {
	cfg := config.MustLoad()
	storage, err := InitStorage(cfg)
	if err != nil {
		log.Fatalln("Произошла ошибка в инициализации бд")
	}
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "api/v1"
	router.GET("/vacs", GetVacancy(storage))

	router.GET("/vac/:id", GetVacancyByID(storage))
	router.GET("/emp/:id", GetEmployerByID(storage))

	router.GET("/emp/vacs/:id", GetVacancyByEmployer(storage))

	router.POST("/vac", PostVacancy(storage))
	router.POST("/emp", PostEmployer(storage))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run("localhost:4252")
}

func InitStorage(cfg *config.Config) (*sqlite.Storage, error) {
	_, err := sqlite.CreateVacancyTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateVacancy Table")
	}
	_, err = sqlite.CreateResponeVacTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateEmployee Table")
	}
	_, err = sqlite.CreateEmployeeTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateEmployee Table")
	}
	_, err = sqlite.CreateTableUser(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateEmployee Table")
	}
	_, err = sqlite.CreateStatusTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateEmployee Table")
	}
	_, err = sqlite.CreateExperienceTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateEmployee Table")
	}
	storage, err := sqlite.CreateResumeTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateVacancy Table")
	}

	return storage, nil
}

// @Success 200 {string} GetVacancyByEmployer
// @Router /emp/vacs/id [get]
func GetVacancyByEmployer(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Error",
				"info":   "Error in get id's from URL parametr! PLS check ur id",
			})
			return
		}
		result, err := storage.GetAllVacsForEmployee(id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"status": "Error",
				"info":   err.Error(),
			})
			return
		}
		ctx.JSON(200, result)
	}
}

// @Success 200 {string} PostEmployer
// @Router /emp [post]
func PostEmployer(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RequestEmployee
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, "error in parse body! Please check our body in request!")
			return
		}

		id, err := storage.AddEmployee(0, req.NameOrganization, req.PhoneNumber, req.Email, req.Geography, req.About)
		if err != nil {
			ctx.JSON(200, "Error in method AddEmployee")
			return
		}
		ctx.JSON(200, gin.H{
			"emp_id": id,
			"status": "OK",
		})

	}

}

// @Success 200 {string} PostVacancy
// @Router /vac [post]
func PostVacancy(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body Vacancy_Body
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, "error in parse body! Please check our body in request!")
			return
		}
		vac_id, emp_limit, err := storage.AddVacancy(body.Emp_ID, body.Vac_Name, body.Price, body.Location, body.Experience)
		if err != nil {
			ctx.JSON(200, "ERROR IN GET ALL VACANCY in SQLITE")
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

// @Success 200 {string} GetEmployerByID
// @Router /emp/:id [get]
func GetEmployerByID(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Error",
				"info":   "Error in get id's from URL parametr! PLS check ur id",
			})
			return
		}

		res, err := storage.GetEmployee(id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"status": "Error",
				"info":   "Произошла какая-то ошибка в методе. Напишите об этом разработчику",
			})
			return
		}
		ctx.JSON(200, res)

	}

}

// @Success 200 {string} GetVacancyByID
// @Router /vac/:id [get]
func GetVacancyByID(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Error",
				"info":   "Error in get id's from URL parametr! PLS check ur id",
			})
			return
		}
		response, err := storage.VacancyByID(id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"status": "Error",
				"info":   "Произошла какая-то ошибка в методе. Напишите об этом разработчику",
			})
			return
		}
		ctx.JSON(200, response)

	}

}

// @Success 200 {string} GetVacancy
// @Router /vac [get]
func GetVacancy(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := storage.GetAllVacancy()
		if err != nil {
			ctx.JSON(200, "ERROR IN GET ALL VACANCY in SQLITE")
		}
		ctx.JSON(200, response)
	}

}
