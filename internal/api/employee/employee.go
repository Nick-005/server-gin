package employee

import (
	"database/sql"
	"fmt"
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
				"Status": "Err",
				"Info":   "ошибка в попытке получить роль пользователя из заголовка токена",
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
		user, err := strconv.Atoi(ctx.Query("empID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить ID работодателя! проверьте его и попробуйте снова",
			})
			return
		}
		err = sqlp.DeleteEmployee(tx, user)
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
			"Info":   "Данные успешно удалены!",
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
				"Status": "Err",
				"Info":   "ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		if role != "employee" && role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		var req s.RequestEmployee
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Error in parse body in request! Please check your body in request!",
				"Error":  err.Error(),
			})
			return
		}
		uid, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		err := sqlp.UpdateEmployeeInfo(tx, req, uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для обновления данных резюме пользователя",
				"Error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Status": "Ok!",
			"Info":   "Данные успешно обновлены!",
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
				"Status": "Err",
				"Info":   "Error in parse body in request! Please check your body in request!",
				"Error":  err.Error(),
			})
			return
		}
		data, err := sqlp.PostNewEmployer(tx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
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
				"Status": "Err",
				"Info":   "Ошибка при создании токена аутентификации",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status":       "OK!",
			"EmployerInfo": data,
			"Token":        token,
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
				"Status": "Err",
				"Info":   "ошибка в попытке получить роль пользователя из заголовка токена",
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
		data, err := sqlp.GetAllEmployee(tx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Status":       "Ok!",
			"EmployerInfo": data,
		})

	}
}

// @Summary Получить информцию про работодателя
// @Description Позволяет получить всю основную информацию про работодатля. Доступно всем авторизованным пользователям, поэтому токен обязателен!
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
		role, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		emp_id, err := strconv.Atoi(ctx.Query("employerID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить ID работодателя! Проверьте его и попробуйте снова",
			})
			return
		}

		uid, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}

		data, err := sqlp.GetEmployeeByID(tx, emp_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}
		if emp_id == uid || role == "ADMIN" {
			fmt.Println("Работодатель получил свои данные или админом")
		} else {
			fmt.Println("Получены данные не админом и не собственником данных")
		}
		ctx.JSON(200, gin.H{
			"Status":       "Ok!",
			"EmployerInfo": data,
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
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если такого пользователя в системе нету."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp/auth [get]
func AuthorizationMethodEmp(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uEmail := ctx.Query("email")
		uPassword := ctx.Query("password")

		data, err := sqlp.GetEmployeeLogin(tx, uEmail, uPassword)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "Такого работодателя не было найдено в системе! Перепроверьте данные и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Произошла ошибка на стороне сервера. Пишите этому горе разрабу. Ошибка в SQL файле",
				"Error":  err.Error(),
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
				"Status": "Err",
				"Info":   "Ошибка при создании токена аутентификации",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status":       "Ok!",
			"EmployerInfo": data,
			"Token":        token,
		})

	}
}
