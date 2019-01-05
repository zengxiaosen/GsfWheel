package main

import (
	"github.com/gin-gonic/gin"
	"gsf/src/jwt/router"
)

func main()  {

	r := gin.Default()

	router.InitRouter(r)
	r.Run("127.0.0.1:8088") // listen and server on 0.0.0.0:8080


}
