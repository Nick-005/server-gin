package vacancy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	get "main.go/internal/api/Get"
	s "main.go/internal/api/Struct"
	sqlp "main.go/internal/storage/postSQL"
)

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
		if role != "employee" {
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
