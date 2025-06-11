package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	mailer "main.go/internal/email-sender"
	sqlp "main.go/internal/storage/postSQL"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	mailer := mailer.New(
		os.Getenv("SMTP_HOSTING"),
		465,
		os.Getenv("SMTP_DOMEN"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_DOMEN"),
		2, // Количество горутин-воркеров
	)

	// Graceful shutdown: закрываем mailer при завершении
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		mailer.Close()
		os.Exit(0)
	}()
	host := os.Getenv("DB_DOMEN")
	port := 5432
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_U_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	connstring := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s target_session_attrs=read-write",
		host, port, dbname, user, password)
	storage, err := sqlx.Connect("pgx", connstring)
	if err != nil {
		log.Fatalln("Произошла ошибка в инициализации бд: ", err.Error())
	}
	defer storage.Close()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.RedirectTrailingSlash = false
	docs.SwaggerInfo.BasePath = "/api/v1"
	apiV1 := router.Group("/api/v1")
	{
		// ~ АДМИН ФУНКЦИОНАЛ
		// ! Удаление соискателей
		apiV1.DELETE("/adm/user", AuthMiddleWare(), MakeTransaction(storage), candid.DeleteUser(storage))

		// ! Удаление работодателей
		apiV1.DELETE("/adm/emp", AuthMiddleWare(), MakeTransaction(storage), employee.DeleteUser(storage))

		// ! Удаление статуса
		apiV1.DELETE("/adm/status", AuthMiddleWare(), MakeTransaction(storage), DeleteStatus(storage))

		// ! Удаление опыта
		apiV1.DELETE("/adm/exp", AuthMiddleWare(), MakeTransaction(storage), DeleteExperience(storage))

		// * Проверка токена на валидность
		apiV1.GET("/adm/token", CheckToken())

		// * Авторизация всех пользователей, вне зависимости от роли: Соискатель или работодатель
		apiV1.GET("/auth", MakeTransaction(storage), candid.AuthorizationMethodForAnybody(storage))

		// & Статус
		// * ----------------------- Все записи -----------------------
		apiV1.GET("/status", MakeTransaction(storage), GetAllStatus(storage))

		// ^ ----------------------- Добавить запись -----------------------
		apiV1.POST("/status", AuthMiddleWare(), MakeTransaction(storage), AddNewStatus(storage))

		// & Работодатель
		// * ----------------------- Получить данные работодателя -----------------------
		apiV1.GET("/emp", AuthMiddleWare(), MakeTransaction(storage), employee.GetEmployeeInfo(storage))

		// * ----------------------- Получить список всех работодателей -----------------------
		apiV1.GET("/emp/all", AuthMiddleWare(), MakeTransaction(storage), employee.GetAllEmployee(storage))

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
		apiV1.GET("/exp", MakeTransaction(storage), GetAllExperience(storage))

		// ^ ----------------------- Добавить -----------------------
		apiV1.POST("/exp", AuthMiddleWare(), MakeTransaction(storage), PostNewExperience(storage))

		// & Соискатели

		apiV1.GET("/user/recover", MakeTransaction(storage)) // TODO доделать

		// * ----------------------- Получить все данные пользователя -----------------------
		apiV1.GET("/user", AuthMiddleWare(), MakeTransaction(storage), candid.GetCandidateInfo(storage))

		// * ----------------------- Зачем то получение всех пользователей -----------------------
		apiV1.GET("/user/all", AuthMiddleWare(), MakeTransaction(storage), candid.GetAllCandidates(storage))

		// * ----------------------- Авторизация пользователя (обновить/получить токен пользователя) -----------------------
		apiV1.GET("/user/auth", MakeTransaction(storage), candid.AuthorizationMethod(storage))

		// * -----------------------  Все резюме пользователя -----------------------
		apiV1.GET("/user/resume", AuthMiddleWare(), MakeTransaction(storage), candid.GetResumeOfCandidates(storage))

		// * ----------------------- Все отклики пользователя -----------------------
		apiV1.GET("/user/response", AuthMiddleWare(), MakeTransaction(storage), candid.GetAllUserResponse(storage))

		// ^ ----------------------- Добавить/зарегестрировать нового пользователя -----------------------
		apiV1.POST("/user", MakeTransaction(storage), candid.PostNewCandidate(storage, mailer))

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
		apiV1.GET("/vac", MakeTransaction(storage), vacancy.GetVacancyWithLimit(storage))

		// * ----------------------- Все вакансии работодателя по 'странично' по времени -----------------------
		apiV1.GET("/vac/time", MakeTransaction(storage), vacancy.GetVacancyWithLimitByTime(storage))

		// * ----------------------- Получить информацию о вакансии -----------------------
		apiV1.GET("/vac/info", MakeTransaction(storage), vacancy.GetVacancyInfoByID(storage))

		// * ----------------------- Количество вакансий в системе -----------------------
		apiV1.GET("/vac/num", MakeTransaction(storage), vacancy.GetVacanciesNumbers(storage))

		apiV1.GET("/vac/user", AuthMiddleWare(), MakeTransaction(storage), vacancy.GetAllResponseByVacancy(storage))

		// * ----------------------- Все отклики на вакансию -----------------------
		apiV1.GET("/vac/response", AuthMiddleWare(), MakeTransaction(storage), response.GetAllResponseByVacancy(storage))

		// ^ ----------------------- Добавить новую вакансию -----------------------
		apiV1.POST("/vac", AuthMiddleWare(), MakeTransaction(storage), vacancy.PostNewVacancy(storage))

		// ? ----------------------- Обновить вакансии -----------------------
		apiV1.PUT("/vac", AuthMiddleWare(), MakeTransaction(storage), vacancy.PutVacancy(storage))

		// ? ----------------------- Обновить видимость вакансии -----------------------
		apiV1.PATCH("/vac/visible", AuthMiddleWare(), MakeTransaction(storage), vacancy.PatchVisibleVacancy(storage))

		// ! ----------------------- Удаление вакансии -----------------------
		apiV1.DELETE("/vac", AuthMiddleWare(), MakeTransaction(storage), vacancy.DeleteVacancy(storage))

		apiV1.GET("/search", MakeTransaction(storage), SearchSystemXD(storage))
	}

	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}

func MakeTransaction(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx, err := storage.Beginx()
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в создании транзакции для БД",
				"Error":  err.Error(),
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
				return
			}
		}
	}
}

// @Summary Поиск вакансий
// @Description Возвращает список всех вакансий, у которых название будет иметь предаваемую часть слова в наименовании вакансии. Имееют доступ все.
// @Tags vacancy
// @Produce json
// @Param Text query string true "Искомый текст в названии вакансии"
// @Success 200 {object} s.VacanciesByLimitResponse "Возвращает статус 'Ok!' и массив всех данных вакансий"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если произошла на стороне сервера."
// @Router /search [get]
func SearchSystemXD(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		name := ctx.Query("Text")
		data, err := sqlp.GetVacanciesBySearchingSubstring(tx, name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для получения данных о вакансиях",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status":      "Ok!",
			"VacancyInfo": data,
		})
	}
}

// @Summary Получение списка опыта
// @Description Возвращает список всех опыта, который будет использоваться в дальнейшем. Имееют доступ все.
// @Tags ADMIN
// @Produce json
// @Success 200 {object} s.GetAllStatuses "Возвращает массив всех значений опыта. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки!"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если произошла на стороне сервера."
// @Router /exp [get]
func GetAllExperience(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		data, err := sqlp.GetAllExperience(tx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Status": "Ok!",
			"Data":   data,
		})

	}
}

// @Summary Добавление новой записи в таблицу с опытом
// @Description Добавляет новую запись в таблицу, которая отвечает за хранение "констант опыта"
// @Security ApiKeyAuth
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param Name query string true "Наименование нового опыта"
// @Success 200 {object} s.Ok "Добавляет новое значение в таблицу"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить из токена ID (авторизовать пользователя)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /exp [post]
func PostNewExperience(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		role, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		if role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		name := ctx.Query("Name")
		err := sqlp.PostNewExperience(tx, name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Status": "Ok!",
		})

	}
}

// @Summary Добавление новой записи в таблицу с статусом
// @Description Добавляет новую запись в таблицу, которая отвечает за хранение "констант статуса"
// @Security ApiKeyAuth
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param Name query string true "Наименование нового статуса"
// @Success 200 {object} s.Ok "Добавляет новое значение в таблицу и просто возвращает статус 'Ok!'"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить из токена ID (авторизовать пользователя)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /status [post]
func AddNewStatus(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		tx := ctx.MustGet("tx").(*sqlx.Tx)
		role, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		if role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		name := ctx.Query("Name")
		err := sqlp.PostNewStatus(tx, name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status": "Ok!",
		})

	}
}

// @Summary Получение списка статусов
// @Description Возвращает список всех значений статусов, который будет использоваться в дальнейшем. Имееют доступ все.
// @Tags ADMIN
// @Accept json
// @Produce json
// @Success 200 {object} s.GetAllStatuses "Возвращает массив всех значений статусов. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки!"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить из токена ID (авторизовать пользователя)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /status [get]
func GetAllStatus(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		data, err := sqlp.GetAllStatus(tx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Status": "OK!",
			"Data":   data,
		})
	}
}

// @Summary Удаление статуса
// @Description Позволяет удалить запись из системы. Доступ имеют только пользователи роли ADMIN
// @Security ApiKeyAuth
// @Tags ADMIN
// @Produce json
// @Param Name query string true "наименование записи, которую нужно удалить"
// @Success 200 {object} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /adm/status [delete]
func DeleteStatus(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		role, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		if role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		name := ctx.Query("name")
		err := sqlp.DeleteStatusByName(tx, name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status": "OK!",
			"Info":   "данные успешно удалены!",
		})
	}
}

// @Summary Удаление опыта
// @Description Позволяет удалить запись из системы. Доступ имеют только пользователи роли ADMIN
// @Security ApiKeyAuth
// @Tags ADMIN
// @Produce json
// @Param Name query string true "наименование записи, которую нужно удалить"
// @Success 200 {object} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /adm/exp [delete]
func DeleteExperience(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		role, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		if role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		name := ctx.Query("name")
		err := sqlp.DeleteExperienceByName(tx, name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status": "OK!",
			"Info":   "данные успешно удалены!",
		})

	}
}

// @Summary Проверка токена
// @Description Позволяет проверить токен пользователя на актуальность
// @Tags ADMIN
// @Produce json
// @Param Token query string true "токен, который надо проверить"
// @Success 200 {object} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /adm/token [get]
func CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		tokenString := ctx.Query("Token")
		// fmt.Println(tokenString)
		claim := &s.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_TOKEN_EMP")), nil
		})
		if err != nil {

			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Error":  fmt.Sprintf("Ошибка при дешифровке токена! error: %v", err),
			},
			)
			return

		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Error":  "Невалидный токен! Пожалуйста перепроверьте его",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Status": "Ok!",
			"Info":   "Токен валидный, всё ок",
		})
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Error":  "Authorization заголовок обязательый, а его нету! Пожалуйста перепроверьте его наличие и попробуйте снова!"},
			)
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Error":  "Не верный формат авторизации. Добавить или перепроверить правильность написания Bearer перед токеном"},
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
				"Status": "Err",
				"Error":  fmt.Sprintf("Ошибка при дешифровке токена! error: %v", err),
			},
			)
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Error":  "Невалидный токен! Пожалуйста перепроверьте его",
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
