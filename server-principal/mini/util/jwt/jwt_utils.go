package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mini/models"
	"net/http"
	"time"
)

var jwtkey = []byte("hdcuahchuhcuiahdaushduoikahcjkxnvlkjnfg")

func MakeJWTToken(curUser models.UserInfo) (string, error) {
	expireTime := time.Now().Add(24 * 365 * 50 * time.Hour) //	过期时间

	claims := &Claims{
		UserClaim: UserClaim{Id: curUser.Id, Phone: curUser.Phone, UUID: curUser.UUID},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, tokenErr := token.SignedString(jwtkey)

	if tokenErr != nil {
		return "", errors.New("0")
	}
	return tokenString, nil
}

const (
	TokenHeader = "X-auth"
)

func ParseToken(ctx *gin.Context, resp *models.Resp) (UserClaim, bool) {
	//	从请求头中获取 tokenString
	var tokenString string
	tokenString = ctx.Request.Header.Get(TokenHeader)
	//	如果请求头中不存在，报错
	if tokenString == "" {
		ctx.JSON(http.StatusOK, resp.Fail(401, "请先登录"))
		return UserClaim{}, false
	}
	//	存在的话，进行解析
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) { return []byte(jwtkey), nil })
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusOK, resp.Fail(401, "请先登录"))
		return UserClaim{}, false
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		ctx.JSON(http.StatusOK, resp.Fail(401, "请先登录"))
		return UserClaim{}, false
	}
	//	判断 身份信息 是否真实
	right, code, err1 := models.UserInfoDao.Identifying(claims.UserClaim.Id, claims.UserClaim.Phone, claims.UserClaim.UUID)
	if !right {
		ctx.JSON(http.StatusOK, resp.Fail(code, err1.Error()))
		return UserClaim{}, false
	}

	return UserClaim{claims.UserClaim.Id, claims.UserClaim.Phone, claims.UserClaim.UUID}, true
}
