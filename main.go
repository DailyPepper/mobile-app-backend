package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/DailyPepper/mobile-app-backend/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type UserData struct {
	Values []int `json:"values"`
}

type ResponseOK struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type AverageResult struct {
	Result  float64 `json:"result"`
	Comment string  `json:"comment"`
	Date    string  `json:"date"`
}

var storedNumbers []int
var allCalculations []AverageResult

// postUserData добавляет значения пользователя
// @Summary Добавление данных пользователя BASFI
// @Description Сохраняет данные пользователя и возвращает расчет BASFI
// @Tags Пользовательские данные
// @Accept json
// @Produce json
// @Param userData body UserData true "Массив из 10 чисел"
// @Success 200 {object} AverageResult "Успешное добавление с расчетом BASFI"
// @Failure 400 {object} ResponseError "Ошибка: Неверные данные"
// @Router /userData [post]
func postUserData(c *gin.Context) {
	var input UserData

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if len(input.Values) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must provide exactly 10 numbers"})
		return
	}

	var sum int
	for _, num := range input.Values {
		sum += num
	}
	average := float64(sum) / float64(len(input.Values))

	var basfiComment string
	if average <= 3 {
		basfiComment = fmt.Sprintf("Ваш индекс BASFI: %.1f. Отсутствие ограничений.", average)
	} else if average >= 4 && average <= 6 {
		basfiComment = fmt.Sprintf("Ваш индекс BASFI: %.1f. Умеренные ограничения.", average)
	} else {
		basfiComment = fmt.Sprintf("Ваш индекс BASFI: %.1f. Невозможность выполнить определенное действие. Рекомендуется консультация врача.", average)
	}

	storedNumbers = input.Values
	calculation := AverageResult{
		Result:  average,
		Comment: basfiComment,
		Date:    time.Now().Format(time.RFC3339),
	}

	allCalculations = append(allCalculations, calculation)

	fmt.Printf("POST data saved: %v\n", calculation)

	c.JSON(http.StatusOK, calculation)
}

// getResultUser возвращает все сохраненные расчеты
// @Summary Получение всех расчетов BASFI
// @Description Возвращает все расчеты BASFI, если данные были добавлены
// @Tags Пользовательские данные
// @Produce json
// @Success 200 {array} AverageResult "Все расчеты BASFI"
// @Failure 404 {object} ResponseError "Ошибка: Данные не найдены"
// @Router /resultUser [get]
func getResultUser(c *gin.Context) {
	if len(allCalculations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}

	fmt.Printf("GET data sent: %v\n", allCalculations)

	c.JSON(http.StatusOK, allCalculations)
}

// @title Swagger index BASFI API
// @version 1.0
// @description Это API для работы с пользовательскими данными BASFI.
// @host localhost:8081
// @BasePath /
func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/userData", postUserData)
	r.GET("/resultUser", getResultUser)

	r.Run(":8081")
}
