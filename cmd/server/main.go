package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx" // swagger embed files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"main.go/docs"
	s "main.go/internal/api/Struct"
	"main.go/internal/api/employee"
	"main.go/internal/api/get"
	"main.go/internal/api/response"
	candid "main.go/internal/api/user"
	"main.go/internal/api/vacancy"
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

	// Только для деплоя
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	apiV1 := router.Group("/api/v1")
	{
		// ~ АДМИН ФУНКЦИОНАЛ

		// & Статус - 	получение и добавление

		// & Опыт - 	получение и добавление

		// ! Удаление соискателей
		apiV1.DELETE("/adm/user", AuthMiddleWare(), MakeTransaction(storage), candid.DeleteUser(storage))

		// ! Удаление работодателей
		apiV1.DELETE("/adm/emp", AuthMiddleWare(), MakeTransaction(storage), employee.DeleteUser(storage))

		// & Статус

		// * ----------------------- Все записи -----------------------
		apiV1.GET("/status", AuthMiddleWare(), MakeTransaction(storage), GetAllStatus(storage))

		// ^ ----------------------- Добавить запись -----------------------
		apiV1.POST("/status", AuthMiddleWare(), MakeTransaction(storage), AddNewStatus(storage))

		// & Работодатель

		// * ----------------------- Получить список всех работодателей -----------------------
		apiV1.GET("/emp", AuthMiddleWare(), MakeTransaction(storage), employee.GetAllEmployee(storage))

		// * ----------------------- Авторизовать работодателя (выдать новый токен) -----------------------
		apiV1.GET("/emp/auth", MakeTransaction(storage), employee.AuthorizationMethodEmp(storage))

		// ^ ----------------------- Добавить/зарегестрировать работодателя -----------------------
		apiV1.POST("/emp", MakeTransaction(storage), employee.PostNewEmployer(storage))

		// ? ----------------------- Обновить данные работодателя -----------------------
		apiV1.PUT("/emp", AuthMiddleWare(), MakeTransaction(storage), employee.PutEmployeeInfo(storage))

		// ? ----------------------- Обновить статус отклика на вакансию -----------------------
		apiV1.PATCH("/vac/response", AuthMiddleWare(), MakeTransaction(storage), response.PatchResponseStatus(storage))

		// & Опыт

		// * ----------------------- Все записи -----------------------
		apiV1.GET("/exp", AuthMiddleWare(), MakeTransaction(storage), GetAllExperience(storage))

		// ^ ----------------------- Добавить -----------------------
		apiV1.POST("/exp", AuthMiddleWare(), MakeTransaction(storage), PostNewExperience(storage))

		// & Соискатели

		// * ----------------------- Получить все данные пользователя -----------------------
		apiV1.GET("/user", AuthMiddleWare(), MakeTransaction(storage), candid.GetCandidateInfo(storage))

		// * ----------------------- Зачем то получение всех пользователей -----------------------
		apiV1.GET("/user/all", MakeTransaction(storage), candid.GetAllCandidates(storage))

		// * ----------------------- Авторизация пользователя (обновить/получить токен пользователя) -----------------------
		apiV1.GET("/user/auth", MakeTransaction(storage), candid.AuthorizationMethod(storage))

		// * -----------------------  Все резюме пользователя -----------------------
		apiV1.GET("/user/resume", AuthMiddleWare(), MakeTransaction(storage), candid.GetResumeOfCandidates(storage))

		// * ----------------------- Все отклики пользователя -----------------------
		apiV1.GET("/user/response", AuthMiddleWare(), MakeTransaction(storage), candid.GetAllUserResponse(storage))

		// ^ ----------------------- Добавить/зарегестрировать нового пользователя -----------------------
		apiV1.POST("/user", MakeTransaction(storage), candid.PostNewCandidate(storage))

		// ^ ----------------------- Добавить резюме -----------------------
		apiV1.POST("/user/resume", AuthMiddleWare(), MakeTransaction(storage), candid.PostNewResume(storage))

		// ^ ----------------------- Добавить отклик на вакансии -----------------------
		apiV1.POST("/vac/response", AuthMiddleWare(), MakeTransaction(storage), response.PostNewRespone(storage))

		// ? ----------------------- Обновить данные пользователя -----------------------
		apiV1.PUT("/user", AuthMiddleWare(), MakeTransaction(storage), candid.PutCandidateInfo(storage))

		// ? ----------------------- Обновить данные резюме пользователя -----------------------
		apiV1.PUT("/user/resume", AuthMiddleWare(), MakeTransaction(storage), candid.PutCandidateResume(storage))

		// ! ----------------------- Удалить резюме -----------------------
		apiV1.DELETE("/user/resume", AuthMiddleWare(), MakeTransaction(storage), candid.DeleteResume(storage))

		// ! ----------------------- Удаление отклика на вакансию -----------------------
		apiV1.DELETE("/vac/response", AuthMiddleWare(), MakeTransaction(storage), response.DeleteResponse(storage))

		// & Вакансии

		// * ----------------------- Все вакансии работодателя -----------------------
		apiV1.GET("/vac/emp", AuthMiddleWare(), MakeTransaction(storage), vacancy.GetAllVacanciesByEmployee(storage))

		// * ----------------------- Все вакансии работодателя по 'странично' -----------------------
		apiV1.GET("/vac", AuthMiddleWare(), MakeTransaction(storage), vacancy.GetVacancyWithLimit(storage))

		// * ----------------------- Все отклики на вакансию -----------------------
		apiV1.GET("/vac/response", AuthMiddleWare(), MakeTransaction(storage), response.GetAllResponseByVacancy(storage))

		// ^ ----------------------- Добавить новую вакансию -----------------------
		apiV1.POST("/vac", AuthMiddleWare(), MakeTransaction(storage), vacancy.PostNewVacancy(storage))

		// ? ----------------------- Обновить вакансии -----------------------
		apiV1.PUT("/vac", AuthMiddleWare(), MakeTransaction(storage), vacancy.PutVacancy(storage))

		// ! ----------------------- Удаление вакансии -----------------------
		apiV1.DELETE("/vac", AuthMiddleWare(), MakeTransaction(storage), vacancy.DeleteVacancy(storage))

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
				// log.Printf("произошла ошибка при попытке закоммитить изменения. error: %v", err)
				// ctx.JSON(http.StatusInternalServerError, gin.H{
				// 	"status": "Err",
				// 	"info":   "Ошибка при попытке закоммитить изменения в БД. Обратитесь к backend разрабу!",
				// 	"error":  err.Error(),
				// })
				return
			}
		}
	}
}

// @Summary Получение списка опыта
// @Description Возвращает список всех опыта, который будет использоваться в дальнейшем. Имееют доступ только пользователи роли ADMIN.
// @Security ApiKeyAuth
// @Tags ADMIN
// @Produce json
// @Success 200 {array} s.GetStatus "Возвращает массив всех значений опыта. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки!"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить из токена ID."
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Router /exp [get]
func GetAllExperience(storage *sqlx.DB) gin.HandlerFunc {
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
		if role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
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
			"status": "Ok!",
			"data":   data,
		})

	}
}

// @Summary Добавление новой записи в таблицу с опытом
// @Description Добавляет новую запись в таблицу, которая отвечает за хранение "констант опыта"
// @Security ApiKeyAuth
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param name query string true "Наименование нового опыта"
// @Success 200 {array} s.Ok "Добавляет новое значение в таблицу"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить из токена ID (авторизовать пользователя)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Router /exp [post]
func PostNewExperience(storage *sqlx.DB) gin.HandlerFunc {
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
		if role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
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

// @Summary Добавление новой записи в таблицу с статусом
// @Description Добавляет новую запись в таблицу, которая отвечает за хранение "констант статуса"
// @Security ApiKeyAuth
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param name query string true "Наименование нового статуса"
// @Success 200 {array} s.Ok "Добавляет новое значение в таблицу и просто возвращает статус 'Ok!'"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить из токена ID (авторизовать пользователя)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Router /status [post]
func AddNewStatus(storage *sqlx.DB) gin.HandlerFunc {
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
		if role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
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

// @Summary Получение списка статусов
// @Description Возвращает список всех значений статусов, который будет использоваться в дальнейшем. Имееют доступ только пользователи роли ADMIN.
// @Security ApiKeyAuth
// @Tags ADMIN
// @Accept json
// @Produce json
// @Success 200 {array} s.GetStatus "Возвращает массив всех значений статусов. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки!"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить из токена ID (авторизовать пользователя)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Router /status [get]
func GetAllStatus(storage *sqlx.DB) gin.HandlerFunc {
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
		if role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
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
