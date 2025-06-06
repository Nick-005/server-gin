package vacancy

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	s "main.go/internal/api/Struct"
	"main.go/internal/api/get"
	sqlp "main.go/internal/storage/postSQL"
)

// @Summary Обновить информцию о вакансии
// @Description Позволяет обновить всю основную информацию о вакансии. Доступно только пользователям группы employee и ADMIN
// @Tags vacancy
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param VacancyInfo body s.VacancyPut  true "Данные о вакансии, на которые нужно обновить в системе"
// @Success 200 {array} s.InfoCandidate "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac [put]
func PutVacancy(storag *sqlx.DB) gin.HandlerFunc {
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
		var req s.VacancyPut
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
		err := sqlp.UpdateVacancyInfo(tx, req, uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для обновления данных вакансии",
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

// @Summary Получить кол-во вакансий в системе
// @Description Позволяет получить количество вакансий в системе, доступных для получения. Доступно всем пользователям
// @Tags vacancy
// @Accept json
// @Produce json
// @Success 200 {array} s.NumberOfVacancies "Возвращает статус 'Ok!' и количество вакансий"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/num [get]
func GetVacanciesNumbers(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		number, err := sqlp.GetNumberOfVacancies(tx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для получения данных о вакансиях",
				"Error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Status":   "Ok!",
			"Quantity": number,
		})
	}
}

// @Summary Получение списка вакансий по 'странично'
// @Description Позволяет получить всю основную информацию про все вакансии, которые у есть, но в ограниченном количестве. Limit - кол-во вакансий, которое нужно вернуть. LastID - после какого ID будет идти отсчёт limit.
// @Tags vacancy
// @Accept json
// @Produce json
// @Param limit query int true "Кол-во вакансий, в соответствии с которым нужно вернуть их"
// @Param last_id query int true "После какого ID будет идти отсчёт limit"
// @Success 200 {array} s.VacancyData_Limit "Возвращает статус 'Ok!' и массив всех данных вакансий"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac [get]
func GetVacancyWithLimit(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить кол-во limit! проверьте его и попробуйте снова",
			})
			return
		}
		last_id, err := strconv.Atoi(ctx.Query("last_id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить кол-во last_id! проверьте его и попробуйте снова",
			})
			return
		}
		data, err := sqlp.GetVacancyLimit(tx, limit, last_id)
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

// @Summary Удаление вакансии
// @Description Позволяет удалить вакансию из системы. Доступ имеют только пользователи роли employee и ADMIN
// @Security ApiKeyAuth
// @Tags vacancy
// @Produce json
// @Param vacancyID query int true "ID вакансии, которую нужно удалить"
// @Success 200 {array} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac [delete]
func DeleteVacancy(storage *sqlx.DB) gin.HandlerFunc {
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
				"Info":   "У вас нету прав добавлять вакансии!",
			})
			return
		}

		emp_id, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		vac_id, err := strconv.Atoi(ctx.Query("vacancyID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
			})
			return
		}

		err = sqlp.DeleteVacancy(tx, emp_id, vac_id, role)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "произошла ошибка при попытке удалить резюме",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status": "Ok!",
			"Info":   "успешно удалили данные!",
		})
	}
}

// @Summary Добавить новую вакансию
// @Description Позволяет добавлять новую вакансию в систему. В ответе клиент получит данные вакансии и работодателя. Доступ имеют роли Employee и ADMIN
// @Security ApiKeyAuth
// @Tags vacancy
// @Accept json
// @Produce json
// @Param Vacancy_Info body s.ResponseVac true "Основные данные для добавления вакансии. В поле exp_id указывайте ID, который уже есть в системе!"
// @Success 200 {array} s.ResponseCreateNewVacancy "Возвращает статус 'Ok!', данные новой вакансии и работодателя"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac [post]
func PostNewVacancy(storag *sqlx.DB) gin.HandlerFunc {
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
				"Info":   "У вас нету прав добавлять вакансии!",
			})
			return
		}
		emp_id, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}

		var req s.ResponseVac
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Error in parse body in request! Please check your body in request!",
				"Error":  err.Error(),
			})
			return
		}
		employee, err := sqlp.GetEmployeeByID(tx, emp_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для получения данных о работодателе",
				"Error":  err.Error(),
			})
			return
		}
		data, err := sqlp.PostNewVacancy(tx, req, emp_id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для получения данных о вакансиях работодателя",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status":       "Ok!",
			"VacancyInfo":  data,
			"EmployerInfo": employee,
		})
	}
}

// @Summary Проверка отклика
// @Description Позволяет узнать, откликнулся ли ранее пользователь на эту вакансию. Если да, то какой у неё статус.
// @Tags vacancy
// @Accept json
// @Produce json
// @Param vacancyID query int true "ID вакансии, о которой хотите получить данные"
// @Success 200 {array} s.ResponseInfoByVacancy "Возвращает информацию о вакансии"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/info [get]
func GetVacancyInfoByID(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		vac_id, err := strconv.Atoi(ctx.Query("vacancyID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
			})
			return
		}
		data, err := sqlp.GetVacancyByID(tx, vac_id)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Такой вакансии нету в системе! Перепроверьте данные и попробуйте снова",
				"Error":  err.Error(),
			})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле откликов",
				"Error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"VacancyInfo": data,
			"Status":      "Ok!",
		})

	}
}

// @Summary Проверка отклика
// @Description Позволяет узнать, откликнулся ли ранее пользователь на эту вакансию. Если да, то какой у неё статус.
// @Security ApiKeyAuth
// @Tags vacancy
// @Accept json
// @Produce json
// @Param vacancyID query int true "ID вакансии, на которую надо посмотреть отклик"
// @Success 200 {array} s.ResponseOnVacancy "Возвращает откликнулся ли уже пользователь на эту вакансию и если это правда, то возвращает статус отклика"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/user [get]
func GetAllResponseByVacancy(storage *sqlx.DB) gin.HandlerFunc {
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
		if role == "employee" || role == "ADMIN" {
			ctx.JSON(200, gin.H{
				"Status": "Ok!",
				"Info":   "Вы не можете откликаться на вакансии",
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
		vac_id, err := strconv.Atoi(ctx.Query("vacancyID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
			})
			return
		}
		data, err := sqlp.GetResponseOnVacancy(tx, uid, vac_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле откликов",
				"Error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"IsResponsed":    data.IsResponsed,
			"StatusResponse": data.Status,
			"Status":         "Ok!",
		})

	}
}

// @Summary Все вакансии одного работодателя
// @Description Позволяет получить массив всех вакансий работодателя. В результате клиент получит ID работодателя и массив всех его вакансий.
// @Security ApiKeyAuth
// @Tags vacancy
// @Accept json
// @Produce json
// @Success 200 {array} s.ResponseAllVacancyByEmployee "Возвращает ID отклика, данные об этой вакансии, на которую откликнулся пользователь и статус отклика "
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/emp [get]
func GetAllVacanciesByEmployee(storag *sqlx.DB) gin.HandlerFunc {
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
		emp_id, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		data, err := sqlp.GetAllVacanciesByEmployee(tx, emp_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для получения данных о вакансиях работодателя",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status":        "Ok!",
			"VacanciesInfo": data,
			"EmployerID":    emp_id,
		})

	}
}
