package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"sync"
	"time"
	"tweet-audit/src/app/archive"
	"tweet-audit/src/app/config"
	"tweet-audit/src/app/pipeline"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"google.golang.org/genai"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// loader channel
	ch := pipeline.Channel()
	go archive.SendToChannel(ch)

	// processed result channel
	resultChannel := pipeline.ProcessedResultCollector()

	ctx := context.Background()

	// csv file
	file, err := os.Create("tweets.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(file)

	writer.Write([]string{
		"tweet_url",
		"deleted",
	})

	go pipeline.CSVWriter(resultChannel, writer, config.ReadConfig())

	// gemini client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		log.Fatal(err)
	}

	// rate limiter
	geminiRateLimiter := rate.NewLimiter(
		rate.Every(12*time.Second),
		1,
	)

	// Workers
	numberOfWorkers := 2

	var wg sync.WaitGroup

	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			pipeline.Worker(i, ch, client, geminiRateLimiter, resultChannel, config.ReadConfig())
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
