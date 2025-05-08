package main

/*

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

*/
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx" // swagger embed files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"main.go/docs"
	get "main.go/internal/api/Get"
	s "main.go/internal/api/Struct"
	"main.go/internal/api/employee"
	candid "main.go/internal/api/user"
	"main.go/internal/api/vacancy"
	"main.go/internal/config"
	sqlp "main.go/internal/storage/postSQL"
)

func main() {
	cfg := config.MustLoad()
	storage, err := sqlx.Connect("pgx", cfg.StoragePath)
	if err != nil {
		log.Fatalln("Произошла ошибка в инициализации бд: ", err.Error())
	}
	defer storage.Close()

	// Только для деплоя
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	apiV1 := router.Group("/api/v1")
	{
		// & Статус
		// * Все записи
		apiV1.GET("/status", MakeTransaction(storage), GetAllStatus(storage))

		// ^ Добавить запись
		apiV1.POST("/status", MakeTransaction(storage), AddNewStatus(storage))

		// & Работодатель
		// ^ Добавить/зарегестрировать работодателя
		apiV1.POST("/emp", MakeTransaction(storage), employee.PostNewEmployer(storage))

		// * Получить список всех работодателей
		apiV1.GET("/emp", MakeTransaction(storage), employee.GetAllEmployee(storage))

		// * Авторизовать работодателя (выдать новый токен)
		apiV1.GET("/emp/auth", MakeTransaction(storage), employee.AuthorizationMethodEmp(storage))

		// & Опыт
		// ^ Добавить
		apiV1.POST("/exp", MakeTransaction(storage), PostNewExperience(storage))
		// * Все записи
		apiV1.GET("/exp", MakeTransaction(storage), GetAllExperience(storage))

		// & Соискатели
		// ^ Добавить/зарегестрировать нового пользователя
		apiV1.POST("/user", MakeTransaction(storage), candid.PostNewCandidate(storage))

		// * Получить все данные пользователя
		apiV1.GET("/user", AuthMiddleWare(), MakeTransaction(storage), candid.GetCandidateInfo(storage))

		// * Зачем то получение всех пользователей
		apiV1.GET("/user/all", MakeTransaction(storage), candid.GetAllCandidates(storage))

		// * Авторизация пользователя (обновить/получить токен пользователя)
		apiV1.GET("/user/auth", MakeTransaction(storage), candid.AuthorizationMethod(storage))

		// ? Обновить данные пользователя
		apiV1.PUT("/user", AuthMiddleWare(), MakeTransaction(storage), candid.PutCandidateInfo(storage))

		// ^ Добавить резюме
		apiV1.POST("/user/resume", AuthMiddleWare(), MakeTransaction(storage), candid.PostNewResume(storage))

		// ! Удалить резюме
		apiV1.DELETE("/user/resume", AuthMiddleWare(), MakeTransaction(storage), candid.DeleteResume(storage))

		// * Все резюме пользователя
		apiV1.GET("/user/resume", AuthMiddleWare(), MakeTransaction(storage), candid.GetResumeOfCandidates(storage))

		// ^ Добавить отклик на вакансии
		apiV1.POST("/vac/response", AuthMiddleWare(), MakeTransaction(storage), PostNewRespone(storage))

		// * Все отклики пользователя
		apiV1.GET("/user/response", AuthMiddleWare(), MakeTransaction(storage), candid.GetAllUserResponse(storage))

		// ! Удаление отклика на вакансию
		apiV1.DELETE("/vac/response", AuthMiddleWare(), MakeTransaction(storage), DeleteResponse(storage))

		// & Вакансии
		// ^ Добавить новую вакансию
		apiV1.POST("/vac", AuthMiddleWare(), MakeTransaction(storage), vacancy.PostNewVacancy(storage))

		// * Все вакансии работодателя
		apiV1.GET("/vac", AuthMiddleWare(), MakeTransaction(storage), vacancy.GetAllVacanciesByEmployee(storage))

		// * Все отклики на вакансию
		apiV1.GET("/vac/response", AuthMiddleWare(), MakeTransaction(storage), GetAllResponseByVacancy(storage))

		// apiV1.GET("/token/check", GetTimeToken(storage))
		// apiV1.GET("/emp/vacs", GetVacancyByEmployer(storage))
		// apiV1.POST("/user/otklik", AuthMiddleWare(), PostResponseOnVacancy(storage))

	}

	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8089")
}

func MakeTransaction(storage *sqlx.DB) gin.HandlerFunc {
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
		ctx.Set("tx", tx)
		ctx.Next()

		if ctx.Writer.Status() < http.StatusBadRequest {
			if err := tx.Commit(); err != nil {
				log.Printf("произошла ошибка при попытке закоммитить изменения. error: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status": "Err",
					"info":   "Ошибка при попытке закоммитить изменения в БД. Обратитесь к backend разрабу!",
					"error":  err.Error(),
				})
			}
		}
	}
}

func GetAllResponseByVacancy(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		role, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		if role != "employee" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		vac_id, err := strconv.Atoi(ctx.Query("vacancy"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  err.Error(),
				"info":   "ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
			})
			return
		}
		data, err := sqlp.GetResponseByVacancy(tx, vac_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле откликов",
				"error":  err.Error(),
			})
			return
		}

		data.Vacancy, err = sqlp.GetVacancyByID(tx, vac_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле вакансии",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"vacancy":   data.Vacancy,
			"responses": data.Responses,
			"status":    "Ok!",
		})

	}
}

func DeleteResponse(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		role, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		if role != "candidate" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		uid, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		vac_id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  err.Error(),
				"info":   "ошибка при попытке получить ID резюме! проверьте его и попробуйте снова",
			})
			return
		}
		err = sqlp.DeleteResponse(tx, vac_id, uid)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": "Err",
				"error":  err.Error(),
				"info":   "произошла ошибка при попытке удалить резюме",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status": "Ok!",
			"info":   "успешно удалили данные!",
		})
	}

}

func PostNewRespone(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		role, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		if role != "candidate" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		uid, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		vac_id, err := strconv.Atoi(ctx.Query("vacancy"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  err.Error(),
				"info":   "ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
			})
			return
		}
		resp_id, err := sqlp.PostResponse(tx, uid, vac_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле добавления данных",
				"error":  err.Error(),
			})
			return
		}
		vac_data, err := sqlp.GetVacancyByID(tx, vac_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле вакансий",
				"error":  err.Error(),
			})
			return
		}
		status_info, err := sqlp.GetStatusByID(tx, 7)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле статуса",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"responseID":     resp_id,
			"vacancy":        vac_data,
			"responseStatus": status_info,
			"status":         "Ok!",
		})

	}
}

func GetAllExperience(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		data, err := sqlp.GetAllExperience(tx)
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

	}
}

func PostNewExperience(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		name := ctx.Query("name")
		err := sqlp.PostNewExperience(tx, name)
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

	}
}

func AddNewStatus(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		name := ctx.Query("name")
		err := sqlp.PostNewStatus(tx, name)
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

	}
}

func GetAllStatus(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		data, err := sqlp.GetAllStatus(tx)
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
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  "Authorization заголовок обязательый, а его нету! Переделывай"},
			)
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"error":  "не верный формат авторизации. Добавить или перепроверить правильность написания Bearer перед токеном"},
			)
			ctx.Abort()
			return
		}

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
		ctx.Set("role", claim.Role)
		fmt.Println(claim.ID)
		ctx.Next()
	}
}
