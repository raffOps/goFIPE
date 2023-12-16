package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/heath-check", controllers.GetHealthCheck)
	return r
}
