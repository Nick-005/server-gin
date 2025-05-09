package employee

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	s "main.go/internal/api/Struct"
	get "main.go/internal/api/get"
	sqlp "main.go/internal/storage/postSQL"
)

var expirationTime = time.Now().Add(24 * time.Hour)

func PostNewEmployer(storage *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("tx").(*sqlx.Tx)

		var req s.RequestEmployee
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "Err",
				"info":   "Error in parse body in request! Please check your body in request!",
				"error":  err.Error(),
			})
			return
		}
		data, err := sqlp.PostNewEmployer(tx, req)
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
			Role:  "employee",
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
			"status":    "OK!",
			"allStatus": data,
			"token":     token,
		})
	}
}

func GetAllEmployee(storag *sqlx.DB) gin.HandlerFunc {
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
		data, err := sqlp.GetAllEmployee(tx)
		if err != nil {
			ctx.JSON(200, gin.H{
				"status": "Err",
				"info":   "Ошибка в SQL файле",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":        "Ok!",
			"Info_Employes": data,
		})

	}
}

func AuthorizationMethodEmp(storag *sqlx.DB) gin.HandlerFunc {
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

		data, err := sqlp.GetEmployeeLogin(tx, req.Email, req.Password)
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
			Role:  "employee",
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
			"status":        "Ok!",
			"Employee_Info": data,
			"token":         token,
		})

	}
}
