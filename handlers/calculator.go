package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/AnatolyKoltun/go-calculator-service/config"
	"github.com/AnatolyKoltun/go-calculator-service/models"
	"github.com/gin-gonic/gin"
)

func Calculate(c *gin.Context) {
	var req models.RequestBody

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных: " + err.Error()})
		return
	}

	var result float64
	switch req.Operator {
	case "+":
		result = req.Argument1 + req.Argument2
	case "-":
		result = req.Argument1 - req.Argument2
	case "*":
		result = req.Argument1 * req.Argument2
	case "/":
		if req.Argument2 == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Деление на ноль запрещено"})
			return
		}
		result = req.Argument1 / req.Argument2
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неподдерживаемая операция"})
		return
	}

	calculation := models.Calculation{
		Argument1: req.Argument1,
		Argument2: req.Argument2,
		Operator:  req.Operator,
		Result:    result,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO calculations (argument1, argument2, operator, result, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := config.DB.QueryRow(context.Background(), query,
		calculation.Argument1, calculation.Argument2, calculation.Operator,
		calculation.Result, calculation.CreatedAt).Scan(&calculation.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения в БД: " + err.Error()})
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

	query := `SELECT id, argument1, argument2, operator, result, created_at 
              FROM calculations WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if filter.DateFrom != "" {
		query += ` AND created_at >= $` + string(rune('0'+argCount))
		args = append(args, filter.DateFrom+" 00:00:00")
		argCount++
	}

	if filter.DateTo != "" {
		query += ` AND created_at <= $` + string(rune('0'+argCount))
		args = append(args, filter.DateTo+" 23:59:59")
		argCount++
	}

	query += ` ORDER BY created_at DESC`

	rows, err := config.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка запроса к БД: " + err.Error()})
		return
	}
	defer rows.Close()

	var calculations []models.Calculation
	for rows.Next() {
		var calc models.Calculation
		err := rows.Scan(&calc.ID, &calc.Argument1, &calc.Argument2,
			&calc.Operator, &calc.Result, &calc.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных: " + err.Error()})
			return
		}
		calculations = append(calculations, calc)
	}

	if rows.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке результатов: " + rows.Err().Error()})
		return
	}

	if calculations == nil {
		calculations = []models.Calculation{}
	}

	c.JSON(http.StatusOK, calculations)
}
