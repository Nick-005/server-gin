package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	s "main.go/internal/api/Struct"
	"main.go/internal/api/get"
	sqlp "main.go/internal/storage/postSQL"
)

// @Summary Все отклики соискателей на вакансию
// @Description Позволяет получить массив всех откликов соискателей на одну определенную вакансию.
// @Security ApiKeyAuth
// @Tags vacancy
// @Accept json
// @Produce json
// @Param vacancyID query int true "ID вакансии, на которую надо посмотреть все отклики"
// @Success 200 {array} s.ResponseAllResponsesOnVacancy "Возвращает данные вакансии и все её отклики"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/response [get]
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
		if role != "employee" && role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "Err",
				"info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		vac_id, err := strconv.Atoi(ctx.Query("vacancyID"))
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле откликов",
				"error":  err.Error(),
			})
			return
		}

		data.Vacancy, err = sqlp.GetVacancyByID(tx, vac_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
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

// @Summary Удалить отклик соискателя
// @Description Позволяет удалить данные об отклике пользователя на вакансию. Доступ имеют роли Candidate и ADMIN
// @Security ApiKeyAuth
// @Tags vacancy
// @Accept json
// @Produce json
// @Param vacancy_id query int true "ID вакансии"
// @Success 200 {array} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/response [delete]
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
		if role != "candidate" && role != "ADMIN" {
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
		vac_id, err := strconv.Atoi(ctx.Query("vacancy_id"))
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

// @Summary Изменить статус отклика
// @Description Позволяет изменить статус отклика на вакансию. Доступно только пользователям группы employee и ADMIN
// @Tags vacancy
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Help_data body s.ResponsePatch true "ID отклика, статус которого нужно обновить, а также ID статуса, на который нужно поменять"
// @Success 200 {array} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/response [patch]
func PatchResponseStatus(storag *sqlx.DB) gin.HandlerFunc {
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
		var req s.ResponsePatch
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}
		err := sqlp.PatchResponse(tx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле для обновления данных отклика на вакансию",
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

// @Summary Добавить новый отклик на вакансию
// @Description Позволяет создать в системе новый отклик соискателя на вакансию.
// @Tags vacancy
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param vacancyID query int true "ID вакансии, на которую нужно сделать отклик!"
// @Success 200 {array} s.ResponseCreateNewResponse "Возвращает статус 'Ok!, ID отклика, данные вакансии, на которую откликнулись и статус отклика"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /vac/response [post]
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
		if role != "candidate" && role != "ADMIN" {
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
		vac_id, err := strconv.Atoi(ctx.Query("vacancyID"))
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле добавления данных",
				"error":  err.Error(),
			})
			return
		}
		vac_data, err := sqlp.GetVacancyByID(tx, vac_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле вакансий",
				"error":  err.Error(),
			})
			return
		}
		status_info, err := sqlp.GetStatusByID(tx, 7)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
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
