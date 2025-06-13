package vacancy

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	s "main.go/internal/api/Struct"
	"main.go/internal/api/get"
	sqlp "main.go/internal/storage/postSQL"
)

// @Summary Изменить видимость вакансии
// @Description Позволяет изменить видимость вакансии. Доступно только пользователям группы employee и ADMIN
// @Tags vacancy
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param VacancyID query int true "ID вакансии, которую работодатель хочет скрыть или вернуть на всеобщее обозрение"
// @Success 200 {object} s.StatusInfo "Возвращает статус 'Ok!'"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/visible [patch]
func PatchVisibleVacancy(storag *sqlx.DB) gin.HandlerFunc {
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
				"Info":   "Ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		vacID, err := strconv.Atoi(ctx.Query("VacancyID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "Ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
			})
			return
		}
		data, err := sqlp.GetVacancyInfoByID(tx, vacID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Такой вакансии нету в системе! Перепроверьте данные и попробуйте снова",
				"Error":  err.Error(),
			})
			return
		}

		err = sqlp.PatchVisibilityVacancy(tx, vacID, emp_id, !data.IsVisible)
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

// @Summary Обновить информцию о вакансии
// @Description Позволяет обновить всю основную информацию о вакансии. Доступно только пользователям группы employee и ADMIN
// @Tags vacancy
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param VacancyInfo body s.VacancyPut  true "Данные о вакансии, на которые нужно обновить в системе"
// @Success 200 {object} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac [put]
func PutVacancy(storag *sqlx.DB) gin.HandlerFunc {
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
				"Info":   "Ошибка в парсинге запроса! Пожалуйста перепроверьте ваши данные в Body запроса и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		}
		if req.Email == "" || req.VacancyName == "" || req.ID <= 0 || req.PhoneNumber == "" || req.Price <= 0 || req.About == "" || req.ExperienceId <= 0 || req.Location == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Вы не передали все необходимые данные! Пожалуйста перепроверьте данные, которые вы передаете в Body запроса и попробуйте снова!",
				"Error":  fmt.Errorf("одно или несколько полей с данными у вас отсутствуют или имеют неверное значение").Error(),
			})
			return
		}
		uid, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в попытке получить ID пользователя из заголовка токена",
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

// @Summary Получить кол-во ВИДИМЫХ вакансий в системе
// @Description Позволяет получить количество ВИДИМЫХ вакансий в системе, доступных для получения. Доступно всем пользователям
// @Tags vacancy
// @Accept json
// @Produce json
// @Success 200 {object} s.NumberOfVacancies "Возвращает статус 'Ok!' и количество вакансий"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
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

// @Summary фильтр ВИДИМЫХ вакансий
// @Description Возвращает список всех вакансий, которые будут соответствовать передаваемым требованиям. Имееют доступ все.
// @Tags vacancy
// @Produce json
// @Param ExpID query int false "ID опыта"
// @Param Min query int false "Минимальная ЗП"
// @Param Max query int false "Максимальная ЗП"
// @Success 200 {object} s.VacanciesByLimitResponse "Возвращает статус 'Ok!' и массив всех данных вакансий"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если произошла при попытке получить передаваемые данные. Или если параметров для фильтрации вообще не будет"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если произошла на стороне сервера."
// @Router /vac/filter [get]
func FilterVacanciesByParams(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		queryParams := ctx.Request.URL.Query()

		isExp := queryParams.Has("ExpID")
		isMin := queryParams.Has("Min")
		isMax := queryParams.Has("Max")

		if !(isExp || isMin || isMax) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Вы передали ни одного фильтра!",
			})
			return
		}
		var err error
		var ExpID int
		var Min int
		var Max int

		if isExp {
			ExpID, err = strconv.Atoi(queryParams.Get("ExpID"))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"Status": "Err",
					"Error":  err.Error(),
					"Info":   "Ошибка при попытке получить кол-во ExpID! проверьте его и попробуйте снова",
				})
			}
		}
		if isMin {
			Min, err = strconv.Atoi(queryParams.Get("Min"))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"Status": "Err",
					"Error":  err.Error(),
					"Info":   "Ошибка при попытке получить кол-во Min! проверьте его и попробуйте снова",
				})
			}
		}
		if isMax {
			Max, err = strconv.Atoi(queryParams.Get("Max"))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"Status": "Err",
					"Error":  err.Error(),
					"Info":   "Ошибка при попытке получить кол-во Max! проверьте его и попробуйте снова",
				})
			}
		}
		data, err := sqlp.GetVacanciesByFilter(tx, ExpID, Max, Min, isExp, isMax, isMin)
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

// @Summary Поиск ВИДИМЫХ вакансий
// @Description Возвращает список всех вакансий, у которых название будет иметь предаваемую часть слова в наименовании вакансии. Имееют доступ все.
// @Tags vacancy
// @Produce json
// @Param Text query string true "Искомый текст в названии вакансии"
// @Success 200 {object} s.VacanciesByLimitResponse "Возвращает статус 'Ok!' и массив всех данных вакансий"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если произошла на стороне сервера."
// @Router /vac/search [get]
func SearchVacancies(storage *sqlx.DB) gin.HandlerFunc {
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

// @Summary Получение списка ВИДИМЫХ вакансий по 'странично'
// @Description Позволяет получить всю основную информацию про все ВИДИМЫЕ вакансии, которые у есть, но в ограниченном количестве. Limit - кол-во вакансий, которое нужно вернуть. LastID - после какого ID будет идти отсчёт limit.
// @Tags vacancy
// @Accept json
// @Produce json
// @Param Page query int true "Номер страницы, которую нужно отобразить"
// @Param PerPage query int true "Кол-во вакансий, в соответствии с которым нужно вернуть их"
// @Success 200 {object} s.VacanciesByLimitResponse "Возвращает статус 'Ok!' и массив всех данных вакансий"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac [get]
func GetVacancyWithLimit(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		page, err := strconv.Atoi(ctx.Query("Page"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "Ошибка при попытке получить номер страницы! проверьте его и попробуйте снова",
			})
			return
		}
		perpage, err := strconv.Atoi(ctx.Query("PerPage"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "Ошибка при попытке получить кол-во вакансий на страницу! проверьте его и попробуйте снова",
			})
			return
		}

		data, err := sqlp.GetVacancyLimit(tx, page, perpage)
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

// @Summary Получение списка вакансий по 'странично' по ВРЕМЕНИ
// @Description Позволяет получить всю основную информацию про все вакансии, которые у есть, но в ограниченном количестве. Limit - кол-во вакансий, которое нужно вернуть. CreatedAt - время, после которого будет идти отсчёт limit.
// @Tags vacancy
// @Accept json
// @Produce json
// @Param Limit query int true "Кол-во вакансий, в соответствии с которым нужно вернуть их"
// @Param CreatedAt query string true "время, после которого будет идти отсчёт limit. Сюда указываем время создания последней отображаемой вакансии. Работает, только если использовать время в формате, как в примере: '2025-06-06T22:40:44Z' или '2006-01-02T15:04:05Z'"
// @Success 200 {object} s.ResponseInfoByVacancyByTimes "Возвращает статус 'Ok!' и массив всех данных вакансий"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/time [get]
func GetVacancyWithLimitByTime(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		var cursor struct {
			CreatedAt time.Time `form:"CreatedAt" time_format:"2006-01-02T15:04:05Z"`
			Limit     int       `form:"Limit,default=5"`
		}

		if err := ctx.ShouldBindQuery(&cursor); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid query parameters"})
			return
		}

		data, err := sqlp.GetVacancyLimitByTimes(tx, cursor.Limit, cursor.CreatedAt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для получения данных о вакансиях",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status":        "Ok!",
			"VacanciesInfo": data,
		})
	}
}

// @Summary Удаление вакансии
// @Description Позволяет удалить вакансию из системы. Доступ имеют только пользователи роли employee и ADMIN
// @Security ApiKeyAuth
// @Tags vacancy
// @Produce json
// @Param VacancyID query int true "ID вакансии, которую нужно удалить"
// @Success 200 {object} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac [delete]
func DeleteVacancy(storage *sqlx.DB) gin.HandlerFunc {
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
				"Info":   "Ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		vac_id, err := strconv.Atoi(ctx.Query("VacancyID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "Ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
			})
			return
		}

		err = sqlp.DeleteVacancy(tx, emp_id, vac_id, role)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "Произошла ошибка при попытке удалить резюме",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status": "Ok!",
			"Info":   "Успешно удалили данные!",
		})
	}
}

// @Summary Добавить новую вакансию
// @Description Позволяет добавлять новую вакансию в систему. В ответе клиент получит данные вакансии и работодателя. Доступ имеют роли Employee и ADMIN
// @Security ApiKeyAuth
// @Tags vacancy
// @Accept json
// @Produce json
// @Param VacancyInfo body s.ResponseVac true "Основные данные для добавления вакансии. В поле exp_id указывайте ID, который уже есть в системе!"
// @Success 200 {object} s.ResponseCreateNewVacancy "Возвращает статус 'Ok!', данные новой вакансии и работодателя"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac [post]
func PostNewVacancy(storag *sqlx.DB) gin.HandlerFunc {
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
				"Info":   "Ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}

		var req s.ResponseVac
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в парсинге запроса! Пожалуйста перепроверьте ваши данные в Body запроса и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		}
		if req.Email == "" || req.VacancyName == "" || req.PhoneNumber == "" || req.Price <= 0 || req.About == "" || req.ExperienceId <= 0 || req.Location == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Вы не передали все необходимые данные! Пожалуйста перепроверьте данные, которые вы передаете в Body запроса и попробуйте снова!",
				"Error":  fmt.Errorf("одно или несколько полей с данными у вас отсутствуют или имеют неверное значение").Error(),
			})
			return
		}
		employee, err := sqlp.GetEmployeeByID(tx, emp_id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для получения данных о работодателе",
				"Error":  err.Error(),
			})
			return
		}
		data, err := sqlp.PostNewVacancy(tx, req, emp_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
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

// @Summary Данные вакансии по ID
// @Description Позволяет получить все данные вакансии по её ID. Если да, то какой у неё статус.
// @Tags vacancy
// @Accept json
// @Produce json
// @Param VacancyID query int true "ID вакансии, о которой хотите получить данные"
// @Success 200 {object} s.ResponseInfoByVacancy "Возвращает информацию о вакансии"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/info [get]
func GetVacancyInfoByID(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		vac_id, err := strconv.Atoi(ctx.Query("VacancyID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "Ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
			})
			return
		}
		data, err := sqlp.GetVacancyInfoByID(tx, vac_id)
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
// @Param VacancyID query int true "ID вакансии, на которую надо посмотреть отклик"
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
				"Info":   "Ошибка в попытке получить роль пользователя из заголовка токена",
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
				"Info":   "Ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		vac_id, err := strconv.Atoi(ctx.Query("VacancyID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "Ошибка при попытке получить ID вакансии! проверьте его и попробуйте снова",
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
		emp_id, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в попытке получить ID пользователя из заголовка токена",
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
