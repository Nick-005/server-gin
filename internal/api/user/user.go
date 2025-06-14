package candid

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
	"main.go/internal/api/get"
	mailer "main.go/internal/email-sender"
	sqlp "main.go/internal/storage/postSQL"
)

var expirationTime = time.Now().Add(24 * time.Hour)

// @Summary Удаление аккаунта соискателя
// @Description Позволяет удалить соискателя из системы. Доступ имеют только пользователи роли ADMIN
// @Security ApiKeyAuth
// @Tags ADMIN
// @Produce json
// @Param UserID query int true "ID пользователя, которого нужно удалить"
// @Success 200 {object} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /adm/user [delete]
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
		user, err := strconv.Atoi(ctx.Query("UserID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить ID соискателя! проверьте его и попробуйте снова",
			})
			return
		}
		err = sqlp.DeleteCandidate(tx, user)
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

// @Summary Все отклики соискателя
// @Description Позволяет получить массив всех откликов соискателя. В результате клиент получит ID отклика, данные о всех вакансиях, на которые он откликнулся, а также статус этого отклика
// @Security ApiKeyAuth
// @Tags candidate
// @Accept json
// @Produce json
// @Success 200 {object} s.ResponsesByVac "Возвращает ID отклика, данные об этой вакансии, на которую откликнулся пользователь и статус отклика "
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/response [get]
func GetAllUserResponse(storage *sqlx.DB) gin.HandlerFunc {
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
		if role != "candidate" && role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
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
		data, err := sqlp.GetResponseByCandidate(tx, uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status":    "Ok!",
			"Responses": data,
		})

	}
}

// @Summary Обновить данные об резюме соискателя
// @Description Позволяет обновить данные, которые касаются только резюме соискателя. Доступ имеют роли Candidate и ADMIN
// @Security ApiKeyAuth
// @Tags candidate
// @Accept json
// @Produce json
// @Param ResumeData body s.RequestResumeUpdate true "Данные, которые можно изменить. Это только опыт (стаж) и описание. НО также указываете ID резюме, которое необходимо изменить!"
// @Success 200 {object} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/resume [put]
func PutCandidateResume(storag *sqlx.DB) gin.HandlerFunc {
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
		if role != "candidate" && role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		var req s.RequestResumeUpdate
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Error in parse body in request! Please check your body in request!",
				"Error":  err.Error(),
			})
			return
		}
		if req.Experience <= 0 || req.Description == "" || req.Resume_id <= 0 {
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
				"Info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}
		err := sqlp.UpdateCandidateResume(tx, req, uid)
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

// @Summary Удалить резюме соискателя
// @Description Позволяет удалить данные об резюме пользователя. Доступ имеют роли Candidate и ADMIN
// @Security ApiKeyAuth
// @Tags candidate
// @Accept json
// @Produce json
// @Param ResumeID query int true "ID резюме пользователя, чтобы найти и удалить его"
// @Success 200 {object} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/resume [delete]
func DeleteResume(storag *sqlx.DB) gin.HandlerFunc {
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
		if role != "candidate" && role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
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
		vac_id, err := strconv.Atoi(ctx.Query("ResumeID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить ID резюме! проверьте его и попробуйте снова",
			})
			return
		}
		err = sqlp.DeleteResume(tx, vac_id, uid)
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

// @Summary Подтвердить email
// @Description Позволяет изменить статус подтверждения email пользователя.
// @Tags ADMIN
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Token query string true "Токен для подтверждения почты, который приходит на почту пользователю"
// @Success 200 {object} s.StatusInfo "Возвращает статус 'Ok!'"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /verify [patch]
// func PatchVerifyStatus(storag *sqlx.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		tx := ctx.MustGet("tx").(*sqlx.Tx)

// 		tokenString := ctx.Query("Token")
// 		claim := &s.ClaimsToVerify{}
// 		token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
// 			return []byte(os.Getenv("JWT_SECRET_TOKEN_EMP")), nil
// 		})
// 		if err != nil {

// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"Status": "Err",
// 				"Error":  fmt.Sprintf("Ошибка при дешифровке токена! error: %v", err),
// 			},
// 			)
// 			return

// 		}
// 		if !token.Valid {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"Status": "Err",
// 				"Error":  "Невалидный токен! Пожалуйста перепроверьте его",
// 			})
// 			return
// 		}
// 		err = sqlp.ConfirmUserEmail(tx, claim.Email)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"Status": "Err",
// 				"Info":   "Произошла ошибка в ",
// 				"Error":  err.Error(),
// 			})
// 			return
// 		}

// 		ctx.JSON(200, gin.H{
// 			"Status": "Ok!",
// 			"Info":   "Данные успешно обновлены!",
// 		})
// 	}
// }

// @Summary Добавить нового соискателя
// @Description Позволяет добавлять нового соискателя в систему. В ответе клиент получит токен, с помощью которого сможет получить доступ к некоторому функционалу. Доступ имеют роли Candidate и ADMIN
// @Tags candidate
// @Accept json
// @Produce json
// @Param CandidateInfo body s.RequestCandidate true "Основные данные для добавления соискателя. В поле статус указывайте ID, который уже есть в системе!"
// @Success 200 {object} s.ResponseCreateCandidate "Возвращает статус 'Ok!', данные нового пользователя и его персональный токен, который можно использовать в течении 24 часов!"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user [post]
func PostNewCandidate(storag *sqlx.DB, mailer *mailer.Mailer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		var req s.RequestCandidate
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Error in parse body in request! Please check your body in request!",
				"Error":  err.Error(),
			})
			return
		}
		if req.Email == "" || req.Name == "" || req.Password == "" || req.PhoneNumber == "" || req.Status_id <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Вы не передали все необходимые данные! Пожалуйста перепроверьте данные, которые вы передаете в Body запроса и попробуйте снова!",
				"Error":  fmt.Errorf("Одно или несколько полей с данными у вас отсутствуют или имеют неверное значение").Error(),
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
		data, err := sqlp.PostNewCandidate(tx, req)
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
				Role:  "candidate",
				Email: data.Email,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(expirationTime),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}
		}
		tokenVerify, err := sqlp.GetGenerateTokenToVerify(data.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка при создании токена для подтверждения почты пользователя. Обратитесь в поддержку!",
				"Error":  err.Error(),
			})
			return
		}
		link := fmt.Sprint("https://isp-workall.online/api/v1/user/confirm-email?Token=")
		textToSend := fmt.Sprintf("Здравствуйте, %s!\n\nБлагодарим вас за регистрацию на нашем сервисе!\n\nДля подтверждения почты, пожалуйста, перейдите по ссылке ниже:\n%s%s", data.Name, link, tokenVerify)
		mailer.SendAsync(data.Email, "Подтверждения почты!", textToSend)
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
			"Status":        "Ok!",
			"CandidateInfo": data,
			"Token":         token,
		})
	}
}

// @Summary Подтверждение почты пользователя
// @Description Позволяет подтвердить почту пользователя
// @Tags ADMIN
// @Produce json
// @Param Token query string true "токен, который надо проверить"
// @Success 200 {object} s.StatusInfo "Возвращает статус и краткую информацию "
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/confirm-email [get]
func CheckToken(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		tokenString := ctx.Query("Token")
		email, err := sqlp.ParseVerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Error":  fmt.Sprintf("Ошибка при дешифровке токена! error: %v", err),
			})
			return
		}
		err = sqlp.ConfirmUserEmail(tx, email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Error":  fmt.Sprintf("Ошибка в SQL файле! error: %v", err),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Status": "Ok!",
			"Info":   "Почта успешно подтверждена!",
		})
	}
}

// @Summary Получить информцию о соискателе
// @Description Позволяет получить всю основную информацию о соискателе при помощи его ID. Доступно всем авторизованным пользователям, поэтому токен обязателен!
// @Tags candidate
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param CandidateID query int true "ID соискателя"
// @Success 200 {object} s.GetAllFromCandidates "Возвращает статус 'Ok!' и данные пользователя"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user [get]
func GetCandidateInfo(storag *sqlx.DB) gin.HandlerFunc {
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
		candidId, err := strconv.Atoi(ctx.Query("CandidateID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "ошибка в попытке получить ID пользователя из параметра запроса. Перепроверьте данные и попробуйте снова!",
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

		data, err := sqlp.GetCandidateById(tx, candidId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}
		if candidId == uid || role == "ADMIN" {
			fmt.Println("Соискатель получил свои данные или админом")
		} else {
			fmt.Println("Получены данные не админом и не собственником данных")
		}
		ctx.JSON(200, gin.H{
			"Status":        "Ok!",
			"CandidateInfo": data,
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
// @Success 200 {object} s.StatusInfo "Возвращает статус 'Ok!' и небольшую информацию"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user [put]
func PutCandidateInfo(storag *sqlx.DB) gin.HandlerFunc {
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
		if role != "candidate" && role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		var req s.RequestCandidate
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Ошибка при парсинге данных! Пожалуйста перепроверьте данные, которые вы передаете в Body запроса и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		}
		if req.Email == "" || req.Name == "" || req.PhoneNumber == "" || req.Status_id == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "Вы не передали все необходимые данные! Пожалуйста перепроверьте данные, которые вы передаете в Body запроса и попробуйте снова!",
				"Error":  fmt.Errorf("одно или несколько полей с данными у вас отсутствуют или имеют неверное значение").Error(),
			})
			return
		}
		uEmail, ok := get.GetUserEmailFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		if uEmail != req.Email {
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
				"Info":   "ошибка в попытке получить ID пользователя из заголовка токена",
			})
			return
		}

		err := sqlp.UpdateCandidateInfo(tx, req, uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле для обновления данных о соискателе",
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

// @Summary Получить информцию про всех соискателях
// @Description Позволяет получить всю основную информацию про всех соискателях. Доступно только пользователям с ролью ADMIN
// @Tags candidate
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} s.InfoAboutAllCandidates "Возвращает статус 'Ok!' и массив всех данных о соискателях"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/all [get]
func GetAllCandidates(storag *sqlx.DB) gin.HandlerFunc {
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
		data, err := sqlp.GetAllCandidates(tx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Status":         "Ok!",
			"CandidatesInfo": data,
		})
	}
}

// @Summary Добавить новое резюме для соискателя
// @Description Позволяет добавить к соискателю новое резюме.
// @Tags candidate
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param ResumeInfo body s.RequestResume true "Основные данные для резюме. В поле experience_id указывайте ID, который уже есть в системе!"
// @Success 200 {object} s.Ok "Возвращает статус 'Ok!"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/resume [post]
func PostNewResume(storag *sqlx.DB) gin.HandlerFunc {
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
		if role != "candidate" && role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "У вас нету прав к этому функционалу!",
			})
			return
		}
		var req s.RequestResume
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

		err := sqlp.PostNewResume(tx, req, uid)
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
		})

	}
}

// @Summary Информация про все резюме
// @Description Позволяет получить всю основную информацию про все резюме пользователя, которые у него есть в системе. Доступно для всех пользователей, но токен обязательный!
// @Tags candidate
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param CandidateID query int true "ID соискателя для получения его всех резюмешек"
// @Success 200 {object} s.ResumeResult "Возвращает статус 'Ok!' и массив всех данных резюме соискателя"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 401 {object} s.InfoError "Возвращает ошибку, если у пользователя нету доступа к этому функционалу."
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/resume [get]
func GetResumeOfCandidates(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)
		_, ok := get.GetUserRoleFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Info":   "ошибка в попытке получить роль пользователя из заголовка токена",
			})
			return
		}
		uid, err := strconv.Atoi(ctx.Query("CandidateID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "Err",
				"Error":  err.Error(),
				"Info":   "ошибка при попытке получить ID пользователя! проверьте его и попробуйте снова",
			})
			return
		}
		data, err := sqlp.GetAllResumeByCandidate(tx, uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Произошла ошибка на стороне сервера. Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{

			"ResumesInfo":   data.Resumes,
			"CandidateInfo": data.Candidate,
			"Status":        "Ok!",
		})
	}
}

// TODO доделать
// @Summary Восстановить пароль
// @Description Позволяет восстановить пароль пользователю, если он забыл его
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param Password query string true "новый пароль пользователя"
// @Param Token query string true "новый пароль пользователя"
// @Success 200 {object} s.ResponseCreateCandidate "Возвращает статус 'Ok!', данные соискателя и новый токен"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/recover [get]
func RecoverPassword(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uEmail := ctx.Query("Token")
		uPassword := ctx.Query("Password")

		data, err := sqlp.GetCandidateByLogin(tx, uEmail, uPassword)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "Такого пользователя не было найдено в системе! Перепроверьте данные и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Произошла ошибка на стороне сервера. Ошибка в SQL файле",
				"Error":  err.Error(),
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
				"Status": "Err",
				"Info":   "Произошла ошибка на стороне сервера. Ошибка при создании токена аутентификации",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status":        "Ok!",
			"CandidateInfo": data,
			"Token":         token,
		})
	}
}

// @Summary Авторизовать пользователя
// @Description Позволяет получить новый токен для пользователя, чтобы у него сохранился доступ к функционалу
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param Email query string true "email пользователя"
// @Param Password query string true "password пользователя"
// @Success 200 {object} s.ResponseAuthorization "Возвращает статус 'Ok!', данные пользователя и его новый токен. Если он авторизовался как соискатель, то будут возвращены его данные. А если как работодатель, то тоже только его"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /auth [get]
func AuthorizationMethodForAnybody(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uEmail := ctx.Query("Email")
		uPassword := ctx.Query("Password")

		isEmp, err := sqlp.CheckUserByEmailOnEmployer(tx, uEmail)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Произошла ошибка при попытке проверить пользователя в таблице Работодателей. Ошибка в SQL файле",
				"Error":  err.Error(),
			})
			return
		}
		// fmt.Println(isEmp)
		if isEmp {
			data, err := sqlp.GetEmployeeLogin(tx, uEmail, uPassword)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"Status": "Err",
					"Info":   "Произошла ошибка при попытке получить данные работодателя. Ошибка в SQL файле",
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
				"Status":       "Ok!",
				"EmployerInfo": data,
				"Token":        token,
			})
		} else {
			data, err := sqlp.GetCandidateByLogin(tx, uEmail, uPassword)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"Status": "Err",
					"Info":   "Произошла ошибка при попытке получить данные соискателя. Ошибка в SQL файле",
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
					Role:  "candidate",
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
					"Info":   "Произошла ошибка на стороне сервера. Ошибка при создании токена аутентификации",
					"Error":  err.Error(),
				})
				return
			}
			ctx.JSON(200, gin.H{
				"Status":        "Ok!",
				"CandidateInfo": data,
				"Token":         token,
			})
		}

	}
}

// @Summary Авторизовать соискателя
// @Description Позволяет получить новый токен для соискателя, чтобы у него сохранился доступ к функционалу
// @Tags candidate
// @Accept json
// @Produce json
// @Param Email query string true "email соискателя"
// @Param Password query string true "password соискателя"
// @Success 200 {object} s.ResponseCreateCandidate "Возвращает статус 'Ok!', данные соискателя и новый токен"
// @Failure 400 {object} s.InfoError "Возвращает ошибку, если не удалось получить данные из запроса (токен или передача каких-либо других данных)"
// @Failure 500 {object} s.InfoError "Возвращает ошибку, если на сервере произошла непредвиденная ошибка."
// @Router /user/auth [get]
func AuthorizationMethod(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uEmail := ctx.Query("Email")
		uPassword := ctx.Query("Password")

		data, err := sqlp.GetCandidateByLogin(tx, uEmail, uPassword)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Status": "Err",
				"Info":   "Такого пользователя не было найдено в системе! Перепроверьте данные и попробуйте снова!",
				"Error":  err.Error(),
			})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Err",
				"Info":   "Произошла ошибка на стороне сервера. Ошибка в SQL файле",
				"Error":  err.Error(),
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
				"Status": "Err",
				"Info":   "Произошла ошибка на стороне сервера. Ошибка при создании токена аутентификации",
				"Error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Status":        "Ok!",
			"CandidateInfo": data,
			"Token":         token,
		})
	}
}
