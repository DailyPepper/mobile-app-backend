package main

import (
	"net/http"

	_ "github.com/DailyPepper/mobile-app-backend/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type UserData struct {
	Values []int `json:"values"` // Массив из 10 значений
}

type ResponseOK struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type AverageResult struct {
	Result float64 `json:"result"`
}

var storedNumbers []int

// postUserData добавляет значения пользователя
// @Summary Добавить массив чисел
// @Description Сохраняет массив из 10 чисел
// @Tags Пользовательские данные
// @Accept json
// @Produce json
// @Param userData body UserData true "Массив из 10 чисел"
// @Success 200 {object} ResponseOK "Успешное добавление"
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

	storedNumbers = input.Values

	c.JSON(http.StatusOK, gin.H{"message": "Successfully added numbers"})
}

// getResultUser возвращает среднее значение чисел
// @Summary Получить среднее значение
// @Description Вычисляет среднее значение сохраненных чисел
// @Tags Пользовательские данные
// @Produce json
// @Success 200 {object} AverageResult "Среднее значение"
// @Failure 404 {object} ResponseError "Ошибка: Числа не найдены"
// @Router /resultUser [get]
func getResultUser(c *gin.Context) {
	if len(storedNumbers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No numbers found"})
		return
	}

	var sum int
	for _, num := range storedNumbers {
		sum += num
	}
	average := float64(sum) / float64(len(storedNumbers))

	c.JSON(http.StatusOK, gin.H{"result": average})
}

// @title Swagger Example API
// @version 1.0
// @description Это пример API для работы с пользовательскими данными.
// @host localhost:8081
// @BasePath /
func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/userData", postUserData)
	r.GET("/resultUser", getResultUser)

	r.Run(":8081")
}
