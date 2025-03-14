package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

type RequestNewToken struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RequestAdd struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type RequestEmployee struct {
	NameOrganization string `json:"nameOrg"`
	PhoneNumber      string `json:"phoneNumber"`
	Email            string `json:"email"`
	INN              string `json:"inn"`
}

const secretKEY = "ISP-7-21-borodinna"

func main() {
	cfg := config.MustLoad()
	storage, err := InitStorage(cfg)
	if err != nil {
		log.Fatalln("Произошла ошибка в инициализации бд: ", err.Error())
	}
	// Только для деплоя
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "api/v1"

	router.GET("/token/check", GetTimeToken(storage))

	router.GET("/all/vacs", GetAllVacancy(storage))

	router.GET("/vac", GetVacancy(storage))

	router.GET("/vac/:id", GetVacancyByID(storage))
	router.GET("/emp/:id", GetEmployerByID(storage))

	router.GET("/emp/vacs/:id", GetVacancyByEmployer(storage))

	router.POST("/vac", PostVacancy(storage))
	router.POST("/emp", PostEmployer(storage))

	router.POST("/user", PostUser(storage))

	router.POST("/user/otklik", AuthMiddleWare(), PostResponseOnVacancy(storage))

	router.GET("/user/otkliks/:id", AuthMiddleWare(), GetAllUserResponse(storage))

	router.GET("/auth/user", GetTokenForUser(storage))

	router.GET("/auth/test", AuthMiddleWare(), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "OK!",
			"auth":   "some text!",
		})
	})

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

// @Summary Получение списка всех откликов для пользователя
// @Description Возвращает список всех откликов для определенного пользователя по его ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param UID path int true "ID пользователя"
// @Success 200 {object} map[string]string "Возвращает статус и массив откликов. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки!"
// @Failure 404 {object} map[string]string "Возвращает ошибку, если не удалось преобразовать передаваемый параметр (ID) через URL."
// @Router /user/otkliks/{id} [get]
func GetAllUserResponse(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "Err",
				"error":  err.Error(),
			})
			return
		}
		result, err := storage.GetAllResponse(id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":  "OK!",
			"otkliks": result,
		})

	}
}

// @Summary Создание отклика на вакансию
// @Description Создает отклик на вакансию при помощи ID пользователя и вакансии. Статус отклика автоматически присваевается "Ожидание"
// @Tags vacancy
// @Accept  json
// @Produce  json
// @Param IDs body map[string]string true "ID пользователя и вакансии, на которую нужно добавить отклик"
// @Success 200 {int} respID "Возвращает ID отклика. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!"
// @Failure 400 {object} map[string]string "Возвращает ошибку, если не удалось распарсить request body. К ответу прикрепляется ID, который получил сервер, а также где именно произошла ошибка."
// @Router /user/otklik [post]
func PostResponseOnVacancy(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type RequestResponse struct {
			UID       int `json:"UID"`
			VacancyID int `json:"vac_id"`
		}
		var body RequestResponse
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  "Error in parse body! Please check our body in request!",
			})
			return
		}

		err := storage.CheckVacancyExist(body.VacancyID)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status":    "Err",
				"info":      "Error in vacancy part",
				"vacancyID": body.VacancyID,
				"error":     err.Error(),
			})
			return
		}

		err = storage.CheckUserExist(body.UID)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Error in user part",
				"UID":    body.UID,
				"error":  err.Error(),
			})
			return
		}
		respID, err := storage.MakeResponse(body.UID, body.VacancyID)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Error in 'MakeResponse' part",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":     "OK!",
			"responseID": respID,
		})

	}
}

// @Summary Узнать время актуальности 'Bearer Token' пользователя
// @Description Позволяет узнать текущее время и время, когда 'Bearer Token' пользователя перестанет быть актуальным. Время учитывается в системе UNIX. Зачем нужен этот endpoint? А я и не знаю... по приколу...
// @Tags token
// @Accept  json
// @Produce  json
// @Param UserToken header string true "Токен пользователя"
// @Success 200 {object} map[string]string "Возвращает время, когда 'Bearer Token' перестанет быть актуальным. Время учитывается в системе UNIX. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!"
// @Failure 401 {object} map[string]string "Возвращает ошибку, если не удалось распарсить header, который отвечает за токен или когда токен не соответствует 'Bearer Token'"
// @Router /token/check [get]
func GetTimeToken(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, gin.H{
				"status": "Err",
				"error":  "Authorization header is required"},
			)
			return
		}

		// Проверяем, что заголовок начинается с "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(401, gin.H{
				"status": "Err",
				"error":  "Invalid authorization format"},
			)
			return
		}

		// Извлекаем токен, удаляя "Bearer " из строки
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// fmt.Println(tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKEY), nil
		})
		// fmt.Println(token)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"error":  err.Error(),
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// fmt.Println("Expires at:", time.Unix(int64(claims["exp"].(float64)), 0))
			ctx.JSON(200, gin.H{
				"status":      "OK",
				"token":       "valid",
				"expiredTime": int64(claims["exp"].(float64)),
				"nowTime":     time.Now().Unix(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"error":  "something get wrong! Please write to nick-005",
			})
			return
		}

	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  "Authorization header is required"},
			)
			ctx.Abort()
			return
		}

		// Проверяем, что заголовок начинается с "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  "Invalid authorization format"},
			)
			ctx.Abort()
			return
		}

		// Извлекаем токен, удаляя "Bearer " из строки
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// fmt.Println(tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKEY), nil
		})
		// fmt.Println(token)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// @Summary Выдать новый токен пользователю
// @Description Позволяет выдать новый токен пользователю, если у него нету актуального 'Bearer Token' или был, но он уже не действителен.
// @Tags token
// @Accept  json
// @Produce  json
// @Param UserEmailNPassword body map[string]string true "Актуальные логин (почта) и пароль пользователя"
// @Success 200 {object} map[string]string "Возвращает актуальный и новый токен для пользователя. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!"
// @Failure 400 {object} map[string]string "Возвращает ошибку, если не удалось распарсить body, который отвечает за данные пользователя!"
// @Failure 401 {object} map[string]string "Возвращает ошибку, если не удалось найти пользователя в БД, который соответствовал бы данным, которые были получены сервером в результате этого запроса!"
// @Router /auth/user [get]
func GetTokenForUser(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body RequestNewToken
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, "Error in parse body! Please check our body in request!")
			return
		}

		data, err := storage.CheckPasswordNEmail(body.Email, body.Password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  err.Error(),
			})
			return
		}

		token, err := storage.CreateAccessToken(data.Email)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status": "OK",
			"token":  token,
		})
	}
}

// @Summary Создать нового пользователя
// @Description Позволяет добавить нового пользователя в систему, если пользователя с такими данными не существовало!
// @Tags user
// @Accept  json
// @Produce  json
// @Param UserData body map[string]string true "Данные пользователя. А именно: Почта (email), пароль (password), name (имя), номер телефона (phoneNumber)"
// @Success 200 {object} map[string]string "Возвращает актуальный токен для пользователя, а также ID пользователя. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!"
// @Failure 400 {object} map[string]string "Возвращает ошибку, если не удалось распарсить body, который отвечает за данные пользователя!"
// @Failure 401 {object} map[string]string "Возвращает ошибку, если не удалось добавить пользователя в БД, который соответствовал бы данным, которые были получены сервером в результате этого запроса или не удалось создать для него токен! Конкретная ошибка будет в результате запроса!"
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
			ctx.JSON(401, gin.H{
				"status": "Er",
				"error":  err.Error(),
			})
			return
		}
		token, err := storage.CreateAccessToken(body.Email)
		if err != nil {
			ctx.JSON(401, gin.H{
				"status": "Err",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status": "OK",
			"token":  token,
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
				"status": "Err",
				"error":  err.Error(),
			})
			return
		}
		result, err := storage.GetAllVacsForEmployee(id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"status": "Err",
				"error":  err.Error(),
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

// @Success 200 {string} GetAllVacancy
// @Router /all/vac [get]
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

// @Success 200 {string} GetLimitVacancy
// @Router /vac [get]
func GetVacancy(storage *sqlite.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body RequestVac
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  "Error in parse body! Please check our body in request!",
			})
			return
		}
		response, err := storage.VacancyByLimit(body.Limit, body.Last_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"error":  fmt.Errorf("error in GET vacancies! %w", err).Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":    "OK!",
			"vacancies": response,
		})
	}

}
