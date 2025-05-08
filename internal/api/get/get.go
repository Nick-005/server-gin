package get

import "github.com/gin-gonic/gin"

func GetUserRoleFromContext(ctx *gin.Context) (string, bool) {
	roleGet, fjd := ctx.Get("role")
	if !fjd {
		return "", false
	}
	role, ok := roleGet.(string)
	if !ok {

		return "", false
	}
	return role, true
}

func GetUserIDFromContext(ctx *gin.Context) (int, bool) {

	id, isThere := ctx.Get("id")
	if !isThere {

		return -1, false
	}

	uid, ok := id.(int)
	if !ok {
		return -1, false
	}
	return uid, true
}
