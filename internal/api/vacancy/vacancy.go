package vacancy

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	get "main.go/internal/api/Get"
	s "main.go/internal/api/Struct"
	sqlp "main.go/internal/storage/postSQL"
)

func DeleteVacancy(storage *sqlx.DB) gin.HandlerFunc {
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
				"info":   "У вас нету прав добавлять вакансии!",
			})
			return
		}
		emp_id, ok := get.GetUserIDFromContext(ctx)
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

		err = sqlp.DeleteVacancy(tx, emp_id, vac_id)
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

func PostNewVacancy(storag *sqlx.DB) gin.HandlerFunc {
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
				"info":   "У вас нету прав добавлять вакансии!",
			})
			return
		}
		emp_id, ok := get.GetUserIDFromContext(ctx)
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

func GetAllVacanciesByEmployee(storag *sqlx.DB) gin.HandlerFunc {
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
		emp_id, ok := get.GetUserIDFromContext(ctx)
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
