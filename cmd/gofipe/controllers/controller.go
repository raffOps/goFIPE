package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/database"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/models"
	"log"
	"net/http"
)

func GetHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "api up"})
}

func CreateVehicle(c *gin.Context) {
	var vehicles models.Vehicle
	if err := c.ShouldBindJSON(&vehicles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := database.DB.Create(&vehicles)
	log.Println("Error:", result.Error, "Rows affected:", result.RowsAffected)
	c.JSON(http.StatusCreated, gin.H{"id": vehicles.ID})
}
