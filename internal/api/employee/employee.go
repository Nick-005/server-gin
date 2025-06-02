package employee

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	s "main.go/internal/api/Struct"
	get "main.go/internal/api/get"
	sqlp "main.go/internal/storage/postSQL"
)

var expirationTime = time.Now().Add(24 * time.Hour)

// @Summary Удаление аккаунта работодателя
// @Description Позволяет удалить работодателя из системы. Доступ имеют только пользователи роли ADMIN
// @Security ApiKeyAuth
// @Tags ADMIN
// @Produce json
// @Param empID query int true "ID работодателя, которого нужно удалить"
// @Success 200 {array} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /adm/emp [delete]
func DeleteUser(storage *sqlx.DB) gin.HandlerFunc {
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
		user, err := strconv.Atoi(ctx.Query("empID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  err.Error(),
				"info":   "ошибка при попытке получить ID работодателя! проверьте его и попробуйте снова",
			})
			return
		}
		err = sqlp.DeleteEmployee(tx, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status": "Ok!",
			"info":   "Данные успешно удалены!",
		})
	}
}

// @Summary Обновить информцию о работодателе
// @Description Позволяет обновить всю основную информацию о работодателе при помощи его персонального токена и тела запроса. Доступно только пользователям группы employee и ADMIN
// @Tags employer
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Employee_info body s.RequestEmployee true "Данные о работодателе, на которые нужно обновить в системе"
// @Success 200 {array} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp [put]
func PutEmployeeInfo(storag *sqlx.DB) gin.HandlerFunc {
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
		if role != "employee" && role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		var req s.RequestEmployee
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
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
		err := sqlp.UpdateEmployeeInfo(tx, req, uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле для обновления данных резюме пользователя",
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

// @Summary Добавить нового работодателя
// @Description Позволяет добавлять нового работодателя в систему. В ответе клиент получит токен, с помощью которого сможет получить доступ к некоторому функционалу.
// @Tags employer
// @Accept json
// @Produce json
// @Param Employee_info body s.RequestEmployee true "Основные данные для добавления работодателя. В поле статус указывайте ID, который уже есть в системе!"
// @Success 200 {array} s.ResponseCreateEmployee "Возвращает статус 'Ok!', данные работодателя и новый токен"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp [post]
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
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

// @Summary Получить информцию про всех работодателей
// @Description Позволяет получить всю основную информацию про всех работодатлей. Доступно только пользователям с ролью ADMIN
// @Tags employer
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} s.SuccessEmployer "Возвращает статус 'Ok!' и массив всех данных о работодателях"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp/all [get]
func GetAllEmployee(storag *sqlx.DB) gin.HandlerFunc {
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
		data, err := sqlp.GetAllEmployee(tx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
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

// @Summary Получить информцию про работодателя
// @Description Позволяет получить всю основную информацию про работодатля. Доступно всем авторизованным пользователям, но токен обязателен!
// @Tags employer
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param employerID query int true "ID работодателя"
// @Success 200 {array} s.ResponseEmployeeInfo "Возвращает статус 'Ok!' и данные о работодателе"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp [get]
func GetEmployeeInfo(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		_, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		emp_id, err := strconv.Atoi(ctx.Query("employeeID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  err.Error(),
				"info":   "ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
			})
			return
		}
		data, err := sqlp.GetEmployeeByID(tx, emp_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":   "Ok!",
			"Employee": data,
		})

	}
}

// @Summary Авторизовать работодателя
// @Description Позволяет получить новый токен для работодателя, чтобы у него сохранился доступ к функционалу
// @Tags employer
// @Accept json
// @Produce json
// @Param email query string true "email работодателя"
// @Param password query string true "password работодателя"
// @Success 200 {array} s.ResponseCreateEmployee "Возвращает статус 'Ok!', данные работодателя и новый токен"
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp/auth [get]
func AuthorizationMethodEmp(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uEmail := ctx.Query("email")
		uPassword := ctx.Query("password")

		data, err := sqlp.GetEmployeeLogin(tx, uEmail, uPassword)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
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
