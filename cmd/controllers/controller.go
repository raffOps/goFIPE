package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "api up"})
}
