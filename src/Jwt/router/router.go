package router

import (
	"github.com/gin-gonic/gin"
	"gsf/src/jwt/jwtAuthService"
	"gsf/src/jwt/router/api"
)

func InitRouter(r *gin.Engine) *gin.Engine {



	apiv1 := r.Group("/")
	apiv1.GET("auth", api.GetAuth)

	apiv1.Use(jwtAuthService.JWT())
	{
		apiv1.GET("hello", api.GetHello)
	}

	return r
}