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
		apiV1.GET("/status", MakeTransaction(storage), GetAllStatus(storage))
		apiV1.POST("/status", MakeTransaction(storage), AddNewStatus(storage))

		// & Работодатель
		apiV1.POST("/emp", MakeTransaction(storage), PostNewEmployer(storage))
		apiV1.GET("/emp", MakeTransaction(storage), GetAllEmployee(storage))
		apiV1.GET("/emp/auth", MakeTransaction(storage), AuthorizationMethodEmp(storage))

		// & Опыт
		apiV1.POST("/exp", MakeTransaction(storage), PostNewExperience(storage))
		apiV1.GET("/exp", MakeTransaction(storage), GetAllExperience(storage))

		// & Соискатели
		apiV1.POST("/user", MakeTransaction(storage), PostNewCandidate(storage))
		apiV1.GET("/user", AuthMiddleWare(), MakeTransaction(storage), GetCandidateInfo(storage))
		apiV1.GET("/user/all", MakeTransaction(storage), GetAllCandidates(storage))
		apiV1.GET("/user/auth", MakeTransaction(storage), AuthorizationMethod(storage))
		apiV1.PUT("/user", AuthMiddleWare(), MakeTransaction(storage), PutCandidateInfo(storage))

		// & Резюме
		apiV1.POST("/resume", AuthMiddleWare(), MakeTransaction(storage), PostNewResume(storage))
		apiV1.GET("/resume", AuthMiddleWare(), MakeTransaction(storage), GetResumeOfCandidates(storage))

		// & Вакансии
		apiV1.POST("/vac", AuthMiddleWare(), MakeTransaction(storage), PostNewVacancy(storage))
		apiV1.GET("/vac", AuthMiddleWare(), MakeTransaction(storage), GetAllVacanciesByEmployee(storage))

		// apiV1.GET("/token/check", GetTimeToken(storage))
		// apiV1.GET("/emp/vacs", GetVacancyByEmployer(storage))
		// apiV1.POST("/user/otklik", AuthMiddleWare(), PostResponseOnVacancy(storage))

		// apiV1.GET("/user/otkliks/:id", AuthMiddleWare(), GetAllUserResponse(storage))

	}

	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8089")
}

func GetUserRoleFromContext(ctx *gin.Context) (string, bool) {
	roleGet, fjd := ctx.Get("role")
	if !fjd {
		return "", false
	}
	role, ok := roleGet.(string)
	if !ok {

		return "", false
	}
	return role, true
}

func GetUserIDFromContext(ctx *gin.Context) (int, bool) {

	id, isThere := ctx.Get("id")
	if !isThere {

		return -1, false
	}

	uid, ok := id.(int)
	if !ok {
		return -1, false
	}
	return uid, true
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

func AuthorizationMethod(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		var req s.Authorization
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}

		data, err := sqlp.GetCandidateByLogin(tx, req.Email, req.Password)
		if err == sql.ErrNoRows {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Такого пользователя нету! Проверьте логин и пароль",
				"error":  err.Error(),
			})
			return
		} else if err != nil {
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
		token, err := sqlp.CreateAccessToken(claim)
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
	}
}

func AuthorizationMethodEmp(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		var req s.Authorization
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}

		data, err := sqlp.GetEmployeeLogin(tx, req.Email, req.Password)
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
		token, err := sqlp.CreateAccessToken(claim)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка при создании токена аутентификации",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":        "Ok!",
			"Employee_Info": data,
			"token":         token,
		})

	}
}

func PutCandidateInfo(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		var req s.RequestCandidate
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}
		uid, ok := GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		err := sqlp.UpdateCandidateInfo(tx, req, uid)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле для обновления данных о соискателе",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status": "Ok!",
			"info":   "Данные успешно обновлены!",
		})
	}
}

func GetAllVacanciesByEmployee(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		emp_id, ok := GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		data, err := sqlp.GetAllVacanciesByEmployee(tx, emp_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле для получения данных о вакансиях работодателя",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":       "Ok!",
			"vacancy_info": data,
			"emp_id":       emp_id,
		})

	}
}

func PostNewVacancy(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		role, ok := GetUserRoleFromContext(ctx)
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
				"info":   "У вас нету прав добавлять вакансии!",
			})
			return
		}
		emp_id, ok := GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}

		var req s.ResponseVac
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}
		employee, err := sqlp.GetEmployeeByID(tx, emp_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле для получения данных о работодателе",
				"error":  err.Error(),
			})
			return
		}
		data, err := sqlp.PostNewVacancy(tx, req, emp_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле для получения данных о вакансиях работодателя",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":        "Ok!",
			"vacancy_info":  data,
			"employee_info": employee,
		})
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
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uid, ok := GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}

		data, err := sqlp.GetAllResumeByCandidate(tx, uid)
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
	}
}

func GetCandidateInfo(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uid, ok := GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}

		data, err := sqlp.GetCandidateById(tx, uid)
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
			"candidate_info": data,
		})
	}
}

func GetAllCandidates(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		data, err := sqlp.GetAllCandidates(tx)
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
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		var req s.RequestCandidate
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}

		data, err := sqlp.PostNewCandidate(tx, req)
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
		token, err := sqlp.CreateAccessToken(claim)
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
	}
}

func PostNewResume(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		var req s.RequestResume
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}
		uid, ok := GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}

		err := sqlp.PostNewResume(tx, req, uid)
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

func PostNewEmployer(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		var req s.RequestEmployee
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}
		data, err := sqlp.PostNewEmployer(tx, req)
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
		token, err := sqlp.CreateAccessToken(claim)
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
	}
}

func GetAllEmployee(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		data, err := sqlp.GetAllEmployee(tx)
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
