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
// @Tags Admin
// @Produce json
// @Param EmployerID query int true "ID работодателя, которого нужно удалить"
// @Success 200 {object} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /adm/emp [delete]
func DeleteUser() gin.HandlerFunc {
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
		user, err := strconv.Atoi(ctx.Query("EmployerID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "Ошибка при попытке получить ID работодателя! проверьте его и попробуйте снова",
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

// @Summary Изменить статус работодателя
// @Description Позволяет изменить статус работодателя. Доступно только пользователям группы ADMIN
// @Tags Admin
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param EmployerID query int true "ID работодателя, статус которого нужно обновить"
// @Param StatusID query int true "ID статуса, на который нужно поменять"
// @Success 200 {object} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /adm/emp [patch]
func PatchEmployerStatus(storag *sqlx.DB) gin.HandlerFunc {
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
		queryParams := ctx.Request.URL.Query()
		EmpID, err := strconv.Atoi(queryParams.Get("EmployerID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в парсинге запроса! Пожалуйста перепроверьте ваши передаваемые данные и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		}
		StatusID, err := strconv.Atoi(queryParams.Get("StatusID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в парсинге запроса! Пожалуйста перепроверьте ваши передаваемые данные и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		}
		err = sqlp.PatchStatusEmployer(tx, StatusID, EmpID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для обновления данных отклика на вакансию",
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

// @Summary Обновить информцию о работодателе
// @Description Позволяет обновить всю основную информацию о работодателе при помощи его персонального токена и тела запроса. Доступно только пользователям группы employee и ADMIN
// @Tags Employer
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param EmployerInfo body s.RequestEmployer true "Данные о работодателе, на которые нужно обновить в системе"
// @Success 200 {object} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp [put]
func PutEmployeeInfo(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		email, ok := get.GetUserEmailFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в попытке получить почту пользователя из заголовка токена",
			})
			return
		}
		role, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в попытке получить роль пользователя из заголовка токена",
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
		var req s.RequestEmployer
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в парсинге запроса! Пожалуйста перепроверьте ваши данные в Body запроса и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		}
		if req.Email == "" || req.NameOrganization == "" || req.PhoneNumber == "" || req.Status_id <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Вы не передали все необходимые данные! Пожалуйста перепроверьте данные, которые вы передаете в Body запроса и попробуйте снова!",
				"Error":  fmt.Errorf("Одно или несколько полей с данными у вас отсутствуют или имеют неверное значение").Error(),
			})
			return
		}
		if req.Email != email {
			ok, err := sqlp.CheckEmailIsValid(tx, req.Email)
			if err != nil || !ok {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"Status": "Err",
					"Error":  err.Error(),
				})
				return
			}
		}

		uid, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в попытке получить ID пользователя из заголовка токена",
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
// @Tags Employer
// @Accept json
// @Produce json
// @Param EmployerInfo body s.RequestEmployee true "Основные данные для добавления работодателя. В поле статус указывайте ID, который уже есть в системе!"
// @Success 200 {object} s.ResponseCreateEmployer "Возвращает статус 'Ok!', данные работодателя и новый токен"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp [post]
func PostNewEmployer(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		var req s.RequestEmployee
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в парсинге запроса! Пожалуйста перепроверьте ваши данные в Body запроса и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		}

		ok, err := sqlp.CheckEmailIsValid(tx, req.Email)
		if err != nil || !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
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
		claim := &s.Claims{}
		if data.Status.ID == 2 {
			claim = &s.Claims{
				ID:    data.ID,
				Role:  "ADMIN",
				Email: data.Email,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 5 * 12)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}
		} else {
			claim = &s.Claims{
				ID:    data.ID,
				Role:  "employee",
				Email: data.Email,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(expirationTime),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}
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
// @Tags Admin
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} s.SuccessAllEmployers "Возвращает статус 'Ok!' и массив всех данных о работодателях"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /adm/emp [get]
func GetAllEmployee(storag *sqlx.DB) gin.HandlerFunc {
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
			"Status":        "Ok!",
			"EmployersInfo": data,
		})

	}
}

// @Summary Получить информцию про работодателя
// @Description Позволяет получить всю основную информацию про работодатля. Доступно всем авторизованным пользователям, поэтому токен обязателен!
// @Tags Employer
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param EmployerID query int true "ID работодателя"
// @Success 200 {object} s.ResponseEmployerInfo "Возвращает статус 'Ok!' и данные о работодателе"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp [get]
func GetEmployeeInfo(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		emp_id, err := strconv.Atoi(ctx.Query("EmployerID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "Ошибка при попытке получить ID работодателя! Проверьте его и попробуйте снова",
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
		// if emp_id == uid || role == "ADMIN" {
		// 	fmt.Println("Работодатель получил свои данные или админом")
		// } else {
		// 	fmt.Println("Получены данные не админом и не собственником данных")
		// }
		ctx.JSON(200, gin.H{
			"Status":       "Ok!",
			"EmployerInfo": data,
		})

	}
}

// @Summary Авторизовать работодателя
// @Description Позволяет получить новый токен для работодателя, чтобы у него сохранился доступ к функционалу
// @Tags Employer
// @Accept json
// @Produce json
// @Param Email query string true "email работодателя"
// @Param Password query string true "password работодателя"
// @Success 200 {object} s.ResponseCreateEmployer "Возвращает статус 'Ok!', данные работодателя и новый токен"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если такого пользователя в системе нету."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /emp/auth [get]
func AuthorizationMethodEmp(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uEmail := ctx.Query("Email")
		uPassword := ctx.Query("Password")

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
