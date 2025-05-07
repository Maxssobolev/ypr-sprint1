package main

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type TemperatureResponse struct {
	Location    string    `json:"location"`
	Timestamp   time.Time `json:"timestamp"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Value       float64   `json:"value"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"`
}

func main() {
	router := gin.Default()
	apiRoutes := router.Group("/temperature")
	apiRoutes.GET("/:id", temperatureHandler)

	srv := &http.Server{
		Addr:    ":" + getEnv("PORT", "8081"),
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func temperatureHandler(c *gin.Context) {
	sensorID := c.Param("id")
	location := c.Query("location")

	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	temp := rand.Float64()*30 + 10 // от 10°C до 40°C

	response := TemperatureResponse{
		Location:  location,
		SensorID:  sensorID,
		Value:     math.Round(temp),
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}
