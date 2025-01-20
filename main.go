package main

import (
	"net/http"

	_ "github.com/DailyPepper/mobile-app-backend/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Структура для POST запроса
type UserData struct {
	Values []int `json:"values"` // Массив из 10 значений
}

var storedNumbers []int // Глобальная переменная для хранения чисел

// Эндпоинт POST для добавления значений (сохраняем как userData)
func postUserData(c *gin.Context) {
	var input UserData

	// Пробуем распарсить тело запроса в структуру
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Проверяем, что в запросе ровно 10 значений
	if len(input.Values) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must provide exactly 10 numbers"})
		return
	}

	// Сохраняем значения в глобальную переменную
	storedNumbers = input.Values

	c.JSON(http.StatusOK, gin.H{"message": "Successfully added numbers"})
}

// Эндпоинт GET для получения среднего значения (результат пользователя)
func getResultUser(c *gin.Context) {
	if len(storedNumbers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No numbers found"})
		return
	}

	// Суммируем все числа и делим на 10 для вычисления среднего значения
	var sum int
	for _, num := range storedNumbers {
		sum += num
	}
	average := float64(sum) / float64(len(storedNumbers))

	c.JSON(http.StatusOK, gin.H{"result": average})
}

func main() {
	r := gin.Default()

	// Swagger UI маршрут
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Ваши эндпоинты
	r.POST("/userData", postUserData)
	r.GET("/resultUser", getResultUser)

	r.Run(":8081")
}
