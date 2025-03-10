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
	Emp_ID      int    `json:"emp_id"`
	Vac_Name    string `json:"vac_name"`
	Price       int    `json:"price"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Location    string `json:"location"`
	Experience  int    `json:"exp"`
	About       string `json:"about"`
	Is_visible  bool   `json:"is_visible"`
}

type RequestVac struct {
	Limit   int `json:"limit"`
	Last_id int `json:"last_id"`
}

type RequestAdd struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type RequestEmployee struct {
	NameOrganization string `json:"nameOrg"`
	PhoneNumber      string `json:"phoneNumber"`
	Email            string `json:"email"`
	INN              string `json:"inn"`
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
		log.Fatalln("Произошла ошибка в инициализации бд: ", err.Error())
	}
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "api/v1"

	router.GET("/all/vacs", GetAllVacancy(storage))

	router.GET("/vac", GetVacancy(storage))

	router.GET("/vac/:id", GetVacancyByID(storage))
	router.GET("/emp/:id", GetEmployerByID(storage))

	router.GET("/emp/vacs/:id", GetVacancyByEmployer(storage))

	router.POST("/vac", PostVacancy(storage))
	router.POST("/emp", PostEmployer(storage))

	router.POST("/user", PostUser(storage))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run("localhost:8089")
}

func InitStorage(cfg *config.Config) (*sqlite.Storage, error) {
	_, err := sqlite.CreateVacancyTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateVacancy Table")
	}
	_, err = sqlite.CreateResponeVacTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateResponeVacTable Table. %w", err)
	}
	_, err = sqlite.CreateEmployeeTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateEmployee Table. %w", err)
	}
	_, err = sqlite.CreateTableUser(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateTableUser Table. %w", err)
	}
	_, err = sqlite.CreateStatusTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateStatusTable Table. %w", err)
	}
	_, err = sqlite.CreateExperienceTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateExperienceTable Table. %w", err)
	}
	storage, err := sqlite.CreateResumeTable(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in CreateResumeTable Table. %w", err)
	}

	return storage, nil
}

// @Success 200 {string} PostUser
// @Router /user [post]
func PostUser(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body RequestAdd
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, "Error in parse body! Please check our body in request!")
			return
		}
		uid, err := storage.AddUser(body.Email, body.Password, body.Name, body.PhoneNumber)
		if err != nil {
			ctx.JSON(400, gin.H{
				"Status": "Error",
				"Info":   err.Error(),
			})
			return
		}
		token, err := storage.CreateAccessToken(body.Email)
		if err != nil {
			ctx.JSON(400, gin.H{
				"Status": "Error",
				"Info":   err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status": "OK",
			"Token":  token,
			"UID":    uid,
		})
	}
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
		// status = 1  == Заблокирован
		// status = 2  == Активен																  |
		// status = 3  == Требует активации.													  |
		// В данном случае, у нас по умолчанию будет работодателю требоваться активация аккаунта  V
		id, err := storage.AddEmployee(req.NameOrganization, req.PhoneNumber, req.Email, req.INN, 3)
		if err != nil {
			ctx.JSON(200, fmt.Errorf("error in add employer. Error is: %w", err).Error())
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
		vac_id, err := storage.AddVacancy(body.Emp_ID, body.Vac_Name, body.Price, body.Email, body.PhoneNumber, body.Location, body.Experience, body.About, body.Is_visible)
		if err != nil {
			ctx.JSON(200, fmt.Errorf("error in add vacancy! Error is: %w", err).Error())
			return
		}
		ctx.JSON(200, gin.H{
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
				"info":   fmt.Errorf("ошибка: %w", err).Error(),
			})
			return
		}
		ctx.JSON(200, response)

	}

}

// @Success 200 {string} GetVacancy
// @Router /vac [get]
func GetAllVacancy(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := storage.GetAllVacancy()
		if err != nil {
			ctx.JSON(200, fmt.Errorf("ERROR IN GET ALL VACANCY in SQLITE. %w", err).Error())
			return
		}
		ctx.JSON(200, response)
	}

}

// @Success 200 {string} GetVacancy
// @Router /vac [get]
func GetVacancy(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body RequestVac
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, "Error in parse body! Please check our body in request!")
			return
		}
		response, err := storage.VacancyByLimit(body.Limit, body.Last_id)
		if err != nil {
			ctx.JSON(200, fmt.Errorf("error in GET vacancies! %w", err).Error())
		}
		ctx.JSON(200, response)
	}

}
