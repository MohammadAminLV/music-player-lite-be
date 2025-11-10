package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Track struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Album  *string `json:"album"`
	Poster string  `json:"poster"`
	URL    string  `json:"url"`
}

type Payload struct {
	Data []Track `json:"data"`
}

func readPayload(fileName string) (Payload, error) {
	var p Payload

	path, err := filepath.Abs(fileName)
	if err != nil {
		return p, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}
	if err := json.Unmarshal(b, &p); err != nil {
		return p, err
	}
	return p, nil
}

func main() {
	// Defaults
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "data.json"
	}

	r := gin.Default()

	// Permissive CORS for development; lock down in production.
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	api := r.Group("/api")
	{
		// GET /api/tracks -> returns the entire JSON payload
		api.GET("/tracks", func(c *gin.Context) {
			payload, err := readPayload(dataFile)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read data file", "detail": err.Error()})
				return
			}
			c.JSON(http.StatusOK, payload)
		})

		// GET /api/tracks/:index -> returns single track
		api.GET("/tracks/:index", func(c *gin.Context) {
			payload, err := readPayload(dataFile)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read data file", "detail": err.Error()})
				return
			}
			idxStr := c.Param("index")
			idx, err := strconv.Atoi(idxStr)
			if err != nil || idx < 0 || idx >= len(payload.Data) {
				c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
				return
			}
			c.JSON(http.StatusOK, payload.Data[idx])
		})
	}

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
