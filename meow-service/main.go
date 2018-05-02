package main

import (
	"github.com/gin-gonic/gin"
)

func newRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/meows", insertMeowEndpoint)
	}
	return
}

func main() {
	router := newRouter()
	router.Run(":3000")
}
