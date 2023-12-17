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

func CreateVeiculo(c *gin.Context) {
	var veiculo models.Veiculo
	if err := c.ShouldBindJSON(&veiculo); err != nil {
		log.Println("error binding json", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	result := database.DB.Create(&veiculo)
	log.Println("Error:", result.Error, "Rows affected:", result.RowsAffected)
	c.JSON(http.StatusCreated, gin.H{"id": veiculo.ID})
}
