package routes

import "github.com/gin-gonic/gin"
import "github.com/rjribeiro/goFIPE/cmd/controllers"

func HandleRequests() {
	r := gin.Default()
	r.GET("/heath-check", controllers.GetHealthCheck)
	_ = r.Run(":8080")

}
