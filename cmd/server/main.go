package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx" // swagger embed files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"main.go/docs"
	s "main.go/internal/api/Struct"
	"main.go/internal/config"
	sqlp "main.go/internal/storage/postSQL"
)

var expirationTime = time.Now().Add(24 * time.Hour)

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

		apiV1.POST("/resume", AuthMiddleWare(), PostNewResume(storage))
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

			"Info":   data,
			"status": "Ok!",
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
		claim := &s.Claims{
			ID:    data.ID,
			Role:  "candidate",
			Email: data.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token, err := sqlp.CreateToken_Second(claim)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка при создании токена аутентификации",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":         "Ok!",
			"condidate_Info": data,
			"token":          token,
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
		id, isThere := ctx.Get("id")
		if !isThere {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "User ID not found in context",
			})
			return
		}

		uid, ok := id.(int)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Invalid user ID type in context",
			})
			return
		}

		err = sqlp.PostNewResume(storag, req, uid)
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
		claim := &s.Claims{
			ID:    data.ID,
			Role:  "employee",
			Email: data.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token, err := sqlp.CreateToken_Second(claim)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка при создании токена аутентификации",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":    "OK!",
			"allStatus": data,
			"token":     token,
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

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// проверяем, что запрос содержит заголовок "Authorization"
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  "Authorization заголовок обязательый, а его нету! Переделывай"},
			)
			ctx.Abort()
			return
		}

		// // Проверяем, что заголовок начинается с "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  "не верный формат авторизации. Добавить или перепроверить правильность написания Bearer перед токеном"},
			)
			ctx.Abort()
			return
		}

		// // Извлекаем токен, удаляя "Bearer " из строки
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claim := &s.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_TOKEN_EMP")), nil
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  fmt.Sprintf("ошибка при дешифровке токена! error: %v", err),
			},
			)
			ctx.Abort()
			return
		}

		// // fmt.Println(tokenString)
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("неожиданный метод подписи токена: %v", token.Header["alg"])
		// 	}
		// 	return []byte(os.Getenv("JWT_SECRET_TOKEN_EMP")), nil
		// })
		// // fmt.Println(token)
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  "невалидный токен! Пожалуйста перепроверьте его",
			})
			ctx.Abort()
			return
		}
		ctx.Set("id", claim.ID)
		ctx.Set("email", claim.Email)
		fmt.Println(claim.ID)
		ctx.Next()
	}
}
