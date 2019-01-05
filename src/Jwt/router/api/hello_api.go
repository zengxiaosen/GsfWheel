package api

import (
	"github.com/gin-gonic/gin"
	"gsf/src/jwt/e"
	"net/http"
)

func GetHello(c *gin.Context) {
	res := Gin{C: c}
	name := c.Query("name")
	res.Response(http.StatusOK, e.SUCCESS, "response: " + name + " after checking the token")
	return
}

