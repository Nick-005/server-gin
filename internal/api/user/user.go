package candid

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	get "main.go/internal/api/Get"
	s "main.go/internal/api/Struct"
	sqlp "main.go/internal/storage/postSQL"
)

var expirationTime = time.Now().Add(24 * time.Hour)

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
			ctx.JSON(200, gin.H{
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
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка при создании токена аутентификации",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":         "Ok!",
			"condidate_Info": data,
			"token":          token,
		})
	}
}

func GetCandidateInfo(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

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
			"candidate_info": data,
		})
	}
}

func PutCandidateInfo(storag *sqlx.DB) gin.HandlerFunc {
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
			ctx.JSON(200, gin.H{
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

func GetAllCandidates(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		data, err := sqlp.GetAllCandidates(tx)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":          "Ok!",
			"Candidates_Info": data,
		})
	}
}

func PostNewResume(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

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
			ctx.JSON(200, gin.H{
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

func GetResumeOfCandidates(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		uid, ok := get.GetUserIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "ошибка в попытке получить ID пользователя из заголовка токена",
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

func AuthorizationMethod(storag *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		var req s.Authorization
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}

		data, err := sqlp.GetCandidateByLogin(tx, req.Email, req.Password)
		if err == sql.ErrNoRows {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Такого пользователя нету! Проверьте логин и пароль",
				"error":  err.Error(),
			})
			return
		} else if err != nil {
			ctx.JSON(200, gin.H{
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
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка при создании токена аутентификации",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":         "Ok!",
			"condidate_Info": data,
			"token":          token,
		})
	}
}
