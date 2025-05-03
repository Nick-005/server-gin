package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx" // swagger embed files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"main.go/docs"
	s "main.go/internal/api/Struct"
	"main.go/internal/config"
	sqlp "main.go/internal/storage/postSQL"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.MustLoad()
	storage, err := sqlx.Connect("pgx", cfg.StoragePath)
	if err != nil {
		log.Fatalln("Произошла ошибка в инициализации бд: ", err.Error())
	}
	defer storage.Close()

	// var addStatus Status
	// addStatus.Name = "Без опыта"

	// query, args, err := psql.Insert("status").Columns("name").Values(addStatus.Name).ToSql()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Nice")
	// result, err := storage.Exec(query, args...)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(result, "NICE VERY NICE")

	// Только для деплоя
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/status", GetAllStatus(storage))
		apiV1.POST("/status", AddNewStatus(storage))

		apiV1.POST("/emp", PostNewEmployer(storage))
		apiV1.GET("/emp", GetAllEmployee(storage))

		apiV1.POST("/exp", PostNewExperience(storage))
		apiV1.GET("/exp", GetAllExperience(storage))

		apiV1.POST("/user", PostNewCandidate(storage))
		apiV1.GET("/user", GetAllCandidates(storage))

		apiV1.POST("/resume", PostNewResume(storage))
		apiV1.GET("/resume", GetResumeOfCandidates(storage))

		apiV1.POST("/vac", PostNewVacancy(storage))
		// apiV1.GET("/token/check", GetTimeToken(storage))

		// apiV1.GET("/all/vacs", GetAllVacancy(storage))

		// apiV1.GET("/vac", GetVacancy(storage))

		// apiV1.GET("/vacID", GetVacancyByID(storage))
		// apiV1.GET("/empID", GetEmployerByID(storage))

		// apiV1.GET("/emp/vacs", GetVacancyByEmployer(storage))

		// apiV1.POST("/vac", PostVacancy(storage))
		// apiV1.POST("/emp", PostEmployer(storage))

		// apiV1.POST("/user", PostUser(storage))

		// apiV1.POST("/user/otklik", AuthMiddleWare(), PostResponseOnVacancy(storage))

		// apiV1.GET("/user/otkliks/:id", AuthMiddleWare(), GetAllUserResponse(storage))

		// apiV1.POST("/auth/user", GetTokenForUser(storage))

		// apiV1.GET("/auth/test", AuthMiddleWare(), func(ctx *gin.Context) {
		// 	ctx.JSON(200, gin.H{
		// 		"status": "OK!",
		// 		"auth":   "some text!",
		// 	})
		// })

	}

	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8089")
}

type AllUserResponseOK struct {
	Status  string
	Otkliks string
}

type SimpleError struct {
	Status string
	Error  string
}

type InfoError struct {
	SimpleError
	Info string
}

func GetAllVacanciesByEmployee(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storag.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}

		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}()

		// var id int
		// id, err = strconv.Atoi(ctx.Query("user_id"))

	}
}

func PostNewVacancy(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storag.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}

		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}()
		var req s.ResponseVac
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}
		// vacancy, employer, experience, err := sqlp.PostNewVacancy(storag, req)
		// if err != nil {
		// 	ctx.JSON(200, gin.H{
		// 		"status": "Err",
		// 		"info":   "Ошибка в SQL файле",
		// 		"error":  err.Error(),
		// 	})
		// 	return
		// }

		// ctx.JSON(200, gin.H{
		// 	"status":          "Ok!",
		// 	"vacancy_info":    vacancy,
		// 	"employee_info":   employer,
		// 	"experience_info": experience,
		// })

		tx.Commit()
	}
}

func PostNewRespone(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storag.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}

		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}()
		var req s.RequestResponse
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}
	}
}

func GetResumeOfCandidates(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storag.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}

		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}()
		var id int
		id, err = strconv.Atoi(ctx.Query("user_id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Ошибка в получении данных из строки",
				"error":  err.Error(),
			})
			return
		}

		data, err := sqlp.GetAllResumeByCandidate(storag, id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status": "Ok!",
			"Info":   data,
		})
		tx.Commit()
	}
}

func GetAllCandidates(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storag.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}

		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}()

		data, err := sqlp.GetAllCandidates(storag)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":          "Ok!",
			"Candidates_Info": data,
		})

	}
}

func PostNewCandidate(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storag.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}

		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}()

		var req s.RequestCandidate
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}

		data, err := sqlp.PostNewCandidate(storag, req)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":         "Ok!",
			"condidate_Info": data,
		})
		tx.Commit()
	}
}

func PostNewResume(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storag.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}

		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}()

		var req s.RequestResume
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}

		err = sqlp.PostNewResume(storag, req)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status": "Ok!",
		})
		tx.Commit()

	}
}

func GetAllExperience(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storage.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
			return
		}
		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("Ошибка в откате транзакции: %v", err)
			}
		}()
		data, err := sqlp.GetAllExperience(storage)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":        "Ok!",
			"AllExperience": data,
		})

		tx.Commit()

	}
}

func PostNewExperience(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storage.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}
		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}()
		name := ctx.Query("name")
		err = sqlp.PostNewExperience(storage, name)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}
		// panic("hello")
		if err := tx.Commit(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в коммите транзакции",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status": "Ok!",
		})

	}
}

func PostNewEmployer(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storage.Begin()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}
		defer tx.Rollback()
		var req s.RequestEmployee
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}
		data, err := sqlp.PostNewEmployer(storage, req)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":    "OK!",
			"AllStatus": data,
		})
		tx.Commit()
	}
}

func GetAllEmployee(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storag.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
			return
		}
		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("Ошибка в откате транзакции: %v", err)
			}
		}()

		data, err := sqlp.GetAllEmployee(storag)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":        "Ok!",
			"Info_Employes": data,
		})

		tx.Commit()

	}
}

func AddNewStatus(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storage.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}
		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}()
		name := ctx.Query("name")
		err = sqlp.PostNewStatus(storage, name)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}
		// panic("hello")
		if err := tx.Commit(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в коммите транзакции",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status": "Ok!",
		})

	}
}

func GetAllStatus(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storage.Begin()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"status": "Err",
				"info":   "Ошибка в создании транзакции для БД",
				"error":  err.Error(),
			})
		}
		defer tx.Rollback()
		data, err := sqlp.GetAllStatus(storage)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":    "OK!",
			"AllStatus": data,
		})
		tx.Commit()
	}
}

// // @Summary Получение списка всех откликов для пользователя
// // @Description Возвращает список всех откликов для определенного пользователя по его ID
// // @Security ApiKeyAuth
// // @Tags user
// // @Produce  json
// // @Param id path int true "ID пользователя"
// // @Success 200 {object} AllUserResponseOK "Возвращает статус и массив откликов. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки!"
// // @Failure 404 {object} SimpleError "Возвращает ошибку, если не удалось преобразовать передаваемый параметр (ID) через URL."
// // @Router /user/otkliks/{id} [get]
// func GetAllUserResponse(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		id, err := strconv.Atoi(ctx.Param("id"))
// 		if err != nil {
// 			ctx.JSON(http.StatusNotFound, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		result, err := storage.GetAllResponse(id)
// 		if err != nil {
// 			ctx.JSON(200, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"status":  "OK!",
// 			"otkliks": result,
// 		})

// 	}
// }

// type ResponseOnVacancy struct {
// 	Status string
// 	RespID int64
// }

// type RequestResponse struct {
// 	UID       int `json:"UID"`
// 	VacancyID int `json:"vac_id"`
// }

// // @Summary Создание отклика на вакансию
// // @Description Создает отклик на вакансию при помощи ID пользователя и вакансии. Статус отклика автоматически присваевается "Ожидание"
// // @Tags vacancy
// // @Security ApiKeyAuth
// // @Accept  json
// // @Produce  json
// // @Param IDs body RequestResponse true "ID пользователя и вакансии, на которую нужно добавить отклик"
// // @Success 200 {integer} ResponseOnVacancy "Возвращает ID отклика. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!"
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить request body. К ответу прикрепляется ID, который получил сервер, а также где именно произошла ошибка."
// // @Router /user/otklik [post]
// func PostResponseOnVacancy(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {

// 		var body RequestResponse
// 		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Please check our body in request!",
// 				"error":  "Error in parse body!",
// 			})
// 			return
// 		}

// 		err := storage.CheckVacancyExist(body.VacancyID)
// 		if err != nil {
// 			ctx.JSON(200, gin.H{
// 				"status": "Err",
// 				"info":   fmt.Sprintf("Error in vacancy part. VacancyID: %d", body.VacancyID),
// 				"error":  err.Error(),
// 			})
// 			return
// 		}

// 		err = storage.CheckUserExist(body.UID)
// 		if err != nil {
// 			ctx.JSON(200, gin.H{
// 				"status": "Err",
// 				"info":   fmt.Sprintf("Error in user part. UID: %d", body.UID),
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		respID, err := storage.MakeResponse(body.UID, body.VacancyID)
// 		if err != nil {
// 			ctx.JSON(200, gin.H{
// 				"status": "Err",
// 				"info":   fmt.Sprintf("Error in 'MakeResponse' part. UID: %d. VacancyID: %d", body.UID, body.VacancyID),
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"status":     "OK!",
// 			"responseID": respID,
// 		})

// 	}
// }

// func GetTimeToken(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader := ctx.GetHeader("Authorization")
// 		fmt.Println(authHeader)
// 		if authHeader == "" {
// 			ctx.JSON(401, gin.H{
// 				"status": "Err",
// 				"error":  "Authorization header is required"},
// 			)
// 			return
// 		}

// 		// Проверяем, что заголовок начинается с "Bearer "
// 		if !strings.HasPrefix(authHeader, "Bearer ") {
// 			ctx.JSON(401, gin.H{
// 				"status": "Err",
// 				"error":  "Invalid authorization format"},
// 			)
// 			return
// 		}

// 		// Извлекаем токен, удаляя "Bearer " из строки
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		// fmt.Println(tokenString)
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return []byte(secretKEY), nil
// 		})
// 		// fmt.Println(token)
// 		if err != nil {
// 			ctx.JSON(200, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			// fmt.Println("Expires at:", time.Unix(int64(claims["exp"].(float64)), 0))
// 			ctx.JSON(200, gin.H{
// 				"status":      "OK",
// 				"token":       "valid",
// 				"expiredTime": int64(claims["exp"].(float64)),
// 				"nowTime":     time.Now().Unix(),
// 			})
// 		} else {
// 			ctx.JSON(200, gin.H{
// 				"status": "Err",
// 				"error":  "something get wrong! Please write to nick-005",
// 			})
// 			return
// 		}

// 	}
// }

// func AuthMiddleWare() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader := ctx.GetHeader("Authorization")
// 		if authHeader == "" {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"status": "Err",
// 				"error":  "Authorization header is required"},
// 			)
// 			ctx.Abort()
// 			return
// 		}

// 		// Проверяем, что заголовок начинается с "Bearer "
// 		if !strings.HasPrefix(authHeader, "Bearer ") {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"status": "Err",
// 				"error":  "Invalid authorization format"},
// 			)
// 			ctx.Abort()
// 			return
// 		}

// 		// Извлекаем токен, удаляя "Bearer " из строки
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		// fmt.Println(tokenString)
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return []byte(secretKEY), nil
// 		})
// 		// fmt.Println(token)
// 		if err != nil || !token.Valid {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			ctx.Abort()
// 			return
// 		}
// 		ctx.Next()
// 	}
// }

// type TokenForUser struct {
// 	Status string
// 	Token  string
// }

// // @Summary Выдать новый токен пользователю
// // @Description Позволяет выдать новый токен пользователю, если у него нету актуального 'Bearer Token' или был, но он уже не действителен.
// // @Tags token
// // @Accept  json
// // @Produce  json
// // @Param UserEmailNPassword body RequestNewToken true "Актуальные логин (почта) и пароль пользователя"
// // @Success 200 {object} TokenForUser "Возвращает актуальный и новый токен для пользователя. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!"
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить body, который отвечает за данные пользователя!"
// // @Failure 401 {object} SimpleError "Возвращает ошибку, если не удалось найти пользователя в БД, который соответствовал бы данным, которые были получены сервером в результате этого запроса!"
// // @Router /auth/user [post]
// func GetTokenForUser(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var body RequestNewToken
// 		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse body! Please check our body in request!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}

// 		data, err := storage.CheckPasswordNEmail(body.Email, body.Password)
// 		if err != nil {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}

// 		token, err := storage.CreateAccessToken(data.Email, "user")
// 		if err != nil {
// 			ctx.JSON(200, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"status": "OK",
// 			"token":  token,
// 		})
// 	}
// }

// type RequestEmployer struct {
// 	INN      string `json:"inn"`
// 	Password string `json:"password"`
// }

// // @Summary Выдать новый токен работодателю
// // @Description Позволяет выдать новый токен работодателю, если у него нету актуального 'Bearer Token' или был, но он уже не действителен.
// // @Tags token
// // @Accept  json
// // @Produce  json
// // @Param ДанныеРаботодателя body RequestNewToken true "Актуальные логин (ИНН?) и пароль работодателя"
// // @Success 200 {object} TokenForUser "Возвращает актуальный и новый токен для работодателя. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!"
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить body, который отвечает за данные работодателя!"
// // @Failure 401 {object} SimpleError "Возвращает ошибку, если не удалось найти работодателя в БД, который соответствовал бы данным, которые были получены сервером в результате этого запроса!"
// // @Router /auth/user [post]
// func GetTokenForEmployer(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var body RequestEmployer
// 		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse body! Please check our body in request!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}

// 		data, err := storage.CheckPasswordNEmail(body.INN, body.Password)
// 		if err != nil {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}

// 		token, err := storage.CreateAccessToken(data.Email, "user")
// 		if err != nil {
// 			ctx.JSON(200, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"status": "OK",
// 			"token":  token,
// 		})
// 	}
// }

// type AddNewUser struct {
// 	TokenForUser
// 	UID int64
// }

// // @Summary Создать нового пользователя
// // @Description Позволяет добавить нового пользователя в систему, если пользователя с такими данными не существовало!
// // @Tags user
// // @Accept  json
// // @Produce  json
// // @Param UserData body RequestAdd true "Данные пользователя. А именно: Почта (email), пароль (password), name (имя), номер телефона (phoneNumber)"
// // @Success 200 {object} AddNewUser "Возвращает актуальный токен для пользователя, а также ID пользователя. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!"
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить body, который отвечает за данные пользователя!"
// // @Failure 401 {object} SimpleError "Возвращает ошибку, если не удалось добавить пользователя в БД, который соответствовал бы данным, которые были получены сервером в результате этого запроса или не удалось создать для него токен! Конкретная ошибка будет в результате запроса!"
// // @Router /user [post]
// func PostUser(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var body RequestAdd
// 		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse body! Please check our body in request!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		uid, err := storage.AddUser(body.Email, body.Password, body.Name, body.PhoneNumber)
// 		if err != nil {
// 			ctx.JSON(401, gin.H{
// 				"status": "Er",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		token, err := storage.CreateAccessToken(body.Email, "user")
// 		if err != nil {
// 			ctx.JSON(401, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"status": "OK",
// 			"token":  token,
// 			"UID":    uid,
// 		})
// 	}
// }

// // TODO: Переделать так, чтобы не было видно другие вакансий, у которых is_visible == false

// // @Summary Получить все вакансии работодателя
// // @Description Позволяет получить массив данных о всех вакансиях, которые есть у работодателя. Для этого нужно передать ID работодателя!
// // @Tags employer
// // @Produce  json
// // @Param id query int true "ID работодателя"
// // @Success 200 {object} []sqlite.ResponseVac "Возвращает массив актуальных вакансий от одного работодателя."
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить ID"
// // @Failure 401 {object} SimpleError "Возвращает ошибку, если не удалось получить список всех вакансий! Конкретная ошибка будет в результате запроса!"
// // @Router /emp/vacs [get]
// func GetVacancyByEmployer(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		id, err := strconv.Atoi(ctx.Query("id"))
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse ID in path! Please check our id in request! He must be an integer type!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		result, err := storage.GetAllVacsForEmployee(id)
// 		if err != nil {
// 			ctx.JSON(401, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, result)
// 	}
// }

// type NewEmployer struct {
// 	Status string
// 	Emp_id int64
// }

// // TODO: Сделать так, чтобы этот endpoint выдавал еще и токен

// // @Summary Создать работодателя
// // @Description Позволяет создать работодателя в системе. Будет возвращен ID и токен для работодателя!
// // @Tags employer
// // @Accept json
// // @Produce  json
// // @Param EmpData body RequestEmployee true "Данные работодателя"
// // @Success 200 {object} NewEmployer "Возвращает ID (И попозже будет Token) работодателя."
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить body-request!"
// // @Failure 401 {object} SimpleError "Возвращает ошибку, если не добавить работодателя с корректными данными. Конкретная ошибка будет в результате запроса!"
// // @Router /emp [post]
// func PostEmployer(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var req RequestEmployee
// 		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse body in request! Please check your body in request!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		// status = 1  == Заблокирован
// 		// status = 2  == Активен																  |
// 		// status = 3  == Требует активации.													  |
// 		// В данном случае, у нас по умолчанию будет работодателю требоваться активация аккаунта  V
// 		id, err := storage.AddEmployee(req.NameOrganization, req.PhoneNumber, req.Email, req.INN, 3)
// 		if err != nil {
// 			ctx.JSON(200, fmt.Errorf("error in add employer. Error is: %w", err).Error())
// 			return
// 		}
// 		token, err := storage.CreateAccessToken(req.Email, "emp")
// 		if err != nil {
// 			ctx.JSON(401, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"emp_id": id,
// 			"token":  token,
// 			"status": "OK",
// 		})

// 	}

// }

// type NewVacancy struct {
// 	Status    string
// 	VacancyID int64
// }

// // @Summary Создать вакансию
// // @Description Позволяет создать новую вакансию в системе. Будет возвращен ID вакансии!
// // @Tags vacancy
// // @Accept json
// // @Produce json
// // @Param VacData body Vacancy_Body true "Данные вакансии"
// // @Success 200 {object} NewVacancy "Возвращает ID вакансии."
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить body-request!"
// // @Failure 401 {object} SimpleError "Возвращает ошибку, если не удалось добавить вакансию с переданными данными. Конкретная ошибка будет в результате запроса!"
// // @Router /vac [post]
// func PostVacancy(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var body Vacancy_Body
// 		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse body in request! Please check your body in request!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}

// 		_, err := storage.GetEmployee(body.Emp_ID)
// 		if err != nil {
// 			ctx.JSON(401, gin.H{
// 				"status": "Error",
// 				"info":   "That employer doesn't exist! Please check ur request!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}

// 		vac_id, err := storage.AddVacancy(body.Emp_ID, body.Vac_Name, body.Price, body.Email, body.PhoneNumber, body.Location, body.Experience, body.About, body.Is_visible)
// 		if err != nil {
// 			ctx.JSON(401, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"vacancyID": vac_id,
// 			"status":    "Success",
// 		})
// 	}

// }

// // @Summary Получить данные работодателя по его ID
// // @Description Позволяет получить данные работодателя по его ID.
// // @Tags employer
// // @Produce json
// // @Param id query int true "ID работодателя"
// // @Success 200 {object} sqlite.RequestEmployee "Возвращает ID вакансии."
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить ID работодателя из path!"
// // @Failure 401 {object} SimpleError "Возвращает ошибку, если не удалось получить данные работодателя, который соответствует переданному ID. Конкретная ошибка будет в результате запроса!"
// // @Router /empID [get]
// func GetEmployerByID(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		id, err := strconv.Atoi(ctx.Query("id"))
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse ID in path! Please check our id in request! He must be an integer type!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}

// 		res, err := storage.GetEmployee(id)
// 		if err != nil {
// 			ctx.JSON(401, gin.H{
// 				"status": "Error",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"status":  "OK!",
// 			"EmpData": res,
// 		})

// 	}

// }

// type TakeVacancyByID struct {
// 	Status string
// 	sqlite.ResponseVac
// }

// // @Summary Получить данные о вакансии по её ID
// // @Description Позволяет получить данные о вакансии по её ID.
// // @Tags vacancy
// // @Produce json
// // @Param id query int true "ID вакансии"
// // @Success 200 {object} TakeVacancyByID "Возвращает данные вакансии"
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить ID вакансии из строки запроса!"
// // @Failure 401 {object} SimpleError "Возвращает ошибку, если не удалось получить данные работодателя, который соответствует переданному ID. Конкретная ошибка будет в результате запроса!"
// // @Router /vacID [get]
// func GetVacancyByID(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		ab := ctx.Query("id")
// 		fmt.Println(ab)
// 		id, err := strconv.Atoi(ab)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse ID in path! Please check our id in request! He must be an integer type!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		response, err := storage.VacancyByID(id)
// 		if err != nil {
// 			ctx.JSON(401, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"status":      "OK!",
// 			"vacancyData": response,
// 		})

// 	}

// }

// // @Summary не использовать! УДАЛИТЬ!
// // @Tags delete
// // @Success 200 {string} GetAllVacancy
// // @Router /all/vac [get]
// func GetAllVacancy(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		response, err := storage.GetAllVacancy()
// 		if err != nil {
// 			ctx.JSON(200, fmt.Errorf("ERROR IN GET ALL VACANCY in SQLITE. %w", err).Error())
// 			return
// 		}
// 		ctx.JSON(200, response)
// 	}

// }

// type ListOfVacancies struct {
// 	Status   string
// 	Response []sqlite.VacancyTake
// }

// // @Summary Получить список вакансий
// // @Description Позволяет получить определенное кол-во вакансий.
// // @Tags vacancy
// // @Accept json
// // @Produce json
// // @Param limit query int true "Лимит сколько вакансий"
// // @Param lastID query int true "С какого ID надо показывать вакансии"
// // @Success 200 {object} ListOfVacancies "Возвращает список данных вакансий"
// // @Failure 400 {object} InfoError "Возвращает ошибку, если не удалось распарсить body вакансий!"
// // @Failure 401 {object} SimpleError "Возвращает ошибку, если не удалось получить данные вакансий. Конкретная ошибка будет в результате запроса!"
// // @Router /vac [get]
// func GetVacancy(storage *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var body RequestVac

// 		limit, err := strconv.Atoi(ctx.Query("limit"))
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse body in request! Please check your body in request!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		lastID, err := strconv.Atoi(ctx.Query("lastID"))

// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"status": "Err",
// 				"info":   "Error in parse body in request! Please check your body in request!",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		body.Last_id = lastID
// 		body.Limit = limit
// 		// if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
// 		// ctx.JSON(http.StatusBadRequest, gin.H{
// 		// 	"status": "Err",
// 		// 	"info":   "Error in parse body in request! Please check your body in request!",
// 		// 	"error":  err.Error(),
// 		// })
// 		// return
// 		// }
// 		response, err := storage.VacancyByLimit(body.Limit, body.Last_id)
// 		if err != nil {
// 			ctx.JSON(401, gin.H{
// 				"status": "Err",
// 				"error":  err.Error(),
// 			})
// 			return
// 		}
// 		ctx.JSON(200, gin.H{
// 			"status":    "OK!",
// 			"vacancies": response,
// 		})
// 	}

// }
