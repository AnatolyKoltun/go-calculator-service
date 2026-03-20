package handlers

import (
	"net/http"

	"github.com/AnatolyKoltun/go-calculator-service/models"
	"github.com/AnatolyKoltun/go-calculator-service/services"
	"github.com/gin-gonic/gin"
)

var calcService services.CalculationService

func Calculate(c *gin.Context) {
	var req models.RequestBody

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных: " + err.Error()})
		return
	}

	calculation, errCount := calcService.Count(req)

	if errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errCount.Error()})
		return
	}

	c.JSON(http.StatusOK, calculation)
}

func GetCalculations(c *gin.Context) {
	var filter models.FilterRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметров: " + err.Error()})
		return
	}

	calculations, errList := calcService.GetCalculationsList(filter)

	if errList != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errList.Error()})
		return
	}

	c.JSON(http.StatusOK, calculations)
}
