package candid

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	s "main.go/internal/api/Struct"
	"main.go/internal/api/get"
	sqlp "main.go/internal/storage/postSQL"
)

var expirationTime = time.Now().Add(24 * time.Hour)

// @Summary Удаление аккаунта соискателя
// @Description Позволяет удалить соискателя из системы. Доступ имеют только пользователи роли ADMIN
// @Security ApiKeyAuth
// @Tags ADMIN
// @Produce json
// @Param userID query int true "ID пользователя, которого нужно удалить"
// @Success 200 {array} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /adm/user [delete]
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
		user, err := strconv.Atoi(ctx.Query("userID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  err.Error(),
				"info":   "ошибка при попытке получить ID соискателя! проверьте его и попробуйте снова",
			})
			return
		}
		err = sqlp.DeleteCandidate(tx, user)
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

// @Summary Все отклики соискателя
// @Description Позволяет получить массив всех откликов соискателя. В результате клиент получит ID отклика, данные о всех вакансиях, на которые он откликнулся, а также статус этого отклика
// @Security ApiKeyAuth
// @Tags candidate
// @Accept json
// @Produce json
// @Success 200 {array} s.ResponseByVac "Возвращает ID отклика, данные об этой вакансии, на которую откликнулся пользователь и статус отклика "
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/response [get]
func GetAllUserResponse(storage *sqlx.DB) gin.HandlerFunc {
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
		data, err := sqlp.GetResponseByCandidate(tx, uid)
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
			"data":   data,
		})

	}
}

// @Summary Обновить данные об резюме соискателя
// @Description Позволяет обновить данные, которые касаются только резюме соискателя. Доступ имеют роли Candidate и ADMIN
// @Security ApiKeyAuth
// @Tags candidate
// @Accept json
// @Produce json
// @Param ResumaData body s.RequestResumeUpdate true "Данные, которые можно изменить. Это только опыт (стаж) и описание. НО также указываете ID резюме, которое необходимо изменить!"
// @Success 200 {array} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/resume [put]
func PutCandidateResume(storag *sqlx.DB) gin.HandlerFunc {
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
		var req s.RequestResumeUpdate
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
		err := sqlp.UpdateCandidateResume(tx, req, uid)
		if err != nil {
			ctx.JSON(200, gin.H{
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

// @Summary Удалить резюме соискателя
// @Description Позволяет удалить данные об резюме пользователя. Доступ имеют роли Candidate и ADMIN
// @Security ApiKeyAuth
// @Tags candidate
// @Accept json
// @Produce json
// @Param resume_id query int true "ID резюме пользователя, чтобы найти и удалить его"
// @Success 200 {array} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/resume [delete]
func DeleteResume(storag *sqlx.DB) gin.HandlerFunc {
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
		vac_id, err := strconv.Atoi(ctx.Query("resume_id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  err.Error(),
				"info":   "ошибка при попытке получить ID резюме! проверьте его и попробуйте снова",
			})
			return
		}
		err = sqlp.DeleteResume(tx, vac_id, uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
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

// @Summary Добавить нового соискателя
// @Description Позволяет добавлять нового соискателя в систему. В ответе клиент получит токен, с помощью которого сможет получить доступ к некоторому функционалу. Доступ имеют роли Candidate и ADMIN
// @Tags candidate
// @Accept json
// @Produce json
// @Param Candidate_info body s.RequestCandidate true "Основные данные для добавления соискателя. В поле статус указывайте ID, который уже есть в системе!"
// @Success 200 {array} s.ResponseCreateCandiate "Возвращает статус 'Ok!', данные нового пользователя и его персональный токен, который можно использовать в течении 24 часов!"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user [post]
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка при создании токена аутентификации",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":         "Ok!",
			"candidate_Info": data,
			"token":          token,
		})
	}
}

// @Summary Получить информцию о соискателе
// @Description Позволяет получить всю основную информацию о соискателе при помощи его персонального токена
// @Tags candidate
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} s.InfoCandidate "Возвращает статус 'Ok!' и данные пользователя"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user [get]
func GetCandidateInfo(storag *sqlx.DB) gin.HandlerFunc {
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
			"candidate_Info": data,
		})
	}
}

// @Summary Обновить информцию о соискателе
// @Description Позволяет обновить всю основную информацию о соискателе при помощи его персонального токена и тела запроса. Доступно только пользователям группы Candidate и ADMIN
// @Tags candidate
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param CandidateInfo body s.RequestCandidate true "Данные о соискателе, на которые нужно обновить в системе"
// @Success 200 {array} s.InfoCandidate "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user [put]
func PutCandidateInfo(storag *sqlx.DB) gin.HandlerFunc {
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
		var req s.RequestCandidate
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
		err := sqlp.UpdateCandidateInfo(tx, req, uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
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

// @Summary Получить информцию про всех соискателях
// @Description Позволяет получить всю основную информацию про всех соискателях. Доступно только пользователям с ролью ADMIN
// @Tags candidate
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} s.InfoCandidate "Возвращает статус 'Ok!' и массив всех данных о соискателях"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/all [get]
func GetAllCandidates(storag *sqlx.DB) gin.HandlerFunc {
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
		data, err := sqlp.GetAllCandidates(tx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":          "Ok!",
			"candidates_Info": data,
		})
	}
}

// @Summary Добавить новое резюме для соискателя
// @Description Позволяет добавить к соискателю новое резюме.
// @Tags candidate
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param InfoResume body s.RequestResume true "Основные данные для резюме. В поле experience_id указывайте ID, который уже есть в системе!"
// @Success 200 {array} s.Ok "Возвращает статус 'Ok!"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/resume [post]
func PostNewResume(storag *sqlx.DB) gin.HandlerFunc {
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
		var req s.RequestResume
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

		err := sqlp.PostNewResume(tx, req, uid)
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
		})

	}
}

// @Summary Информация про все резюме
// @Description Позволяет получить всю основную информацию про все резюме пользователя, которые у него есть в системе. Доступно для всех пользователей, но токен обязательный!
// @Tags candidate
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param candidate_id query int true "ID соискателя для получения его всех резюмешек"
// @Success 200 {array} s.ResumeResult "Возвращает статус 'Ok!' и массив всех данных резюме соискателя"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {array} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/resume [get]
func GetResumeOfCandidates(storag *sqlx.DB) gin.HandlerFunc {
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
		uid, err := strconv.Atoi(ctx.Query("candidate_id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"error":  err.Error(),
				"info":   "ошибка при попытке получить ID пользователя! проверьте его и попробуйте снова",
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

// @Summary Авторизовать соискателя
// @Description Позволяет получить новый токен для соискателя, чтобы у него сохранился доступ к функционалу
// @Tags candidate
// @Accept json
// @Produce json
// @Param email query string true "email соискателя"
// @Param password query string true "password соискателя"
// @Success 200 {array} s.ResponseCreateCandiate "Возвращает статус 'Ok!', данные соискателя и новый токен"
// @Failure 400 {array} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {array} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/auth [get]
func AuthorizationMethod(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uEmail := ctx.Query("email")
		uPassword := ctx.Query("password")

		data, err := sqlp.GetCandidateByLogin(tx, uEmail, uPassword)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Такого пользователя нету! Проверьте логин и пароль",
				"error":  err.Error(),
			})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Err",
				"info":   "Ошибка при создании токена аутентификации",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":         "Ok!",
			"candidate_Info": data,
			"token":          token,
		})
	}
}
