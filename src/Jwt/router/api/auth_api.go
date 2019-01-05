package api

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"gsf/src/jwt/dao"
	"gsf/src/jwt/e"
	"gsf/src/jwt/jwtAuthService"
	"net/http"
)

type Auth struct {
	Username string
	Password string
}

func GetAuth(c *gin.Context) {


	userName := c.Query("username")
	password := c.Query("password")

	authService := Auth{Username: userName, Password: password}
	isExist, err := authService.Check()

	if err != nil {
		c.JSON(http.StatusOK, responseForm(bson.M{}, e.ERROR_AUTH_CHECK_TOKEN_FAIL, e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL)))
		return
	}

	if !isExist {
		c.JSON(http.StatusOK, responseForm(bson.M{}, e.ERROR_AUTH, e.GetMsg(e.ERROR_AUTH)))
		return
	}

	token, err := jwtAuthService.GenerateToken(userName, password)

	if err != nil {
		c.JSON(http.StatusOK, responseForm(bson.M{}, e.ERROR_AUTH_TOKEN, e.GetMsg(e.ERROR_AUTH_TOKEN)))
		return
	}

	c.JSON(http.StatusOK, responseForm(bson.M{}, e.SUCCESS, "token: "+token))

}

func (a *Auth) Check() (bool, error) {
	// 查看db层，是否有该用户名密码，如果有，则鉴权成功，分配token，否则失败
	return dao.CheckAuth(a.Username, a.Password)
}

func responseForm(data interface{}, code int, msg string) map[string]interface{} {
	return bson.M{"data": data, "code": code, "msg": msg}
}