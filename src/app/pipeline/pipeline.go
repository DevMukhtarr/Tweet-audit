package pipeline

import (
	"context"
	"encoding/csv"
	"log"
	"strconv"
	"time"
	"tweet-audit/src/app/ai"
	"tweet-audit/src/app/types"

	"golang.org/x/time/rate"
	"google.golang.org/genai"
)

func Channel() chan types.TweetWrapper {
	channel := make(chan types.TweetWrapper, 100)
	return channel
}

func ProcessedResultCollector() chan types.AuditResult {
	channel := make(chan types.AuditResult, 100)
	return channel
}

func Worker(id int, job <-chan types.TweetWrapper, client *genai.Client, geminiRateLimiter *rate.Limiter, resultCh chan<- types.AuditResult, alignmentCriteria types.AlignmentCriteria) {

	const batchSize = 10
	batch := make([]types.TweetWrapper, 0, batchSize)
	// var batchTweet
	for tweet := range job {
		batch = append(batch, tweet)
		// rate limiter
		if len(batch) == batchSize {
			results, err := sendBatch(batch, alignmentCriteria, client, geminiRateLimiter)
			if err != nil {
				return
			}
			for _, result := range results {
				resultCh <- result
			}
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		results, err := sendBatch(batch, alignmentCriteria, client, geminiRateLimiter)
		if err != nil {
			return
		}
		for _, result := range results {
			resultCh <- result
		}
	}
}

func sendBatch(batch []types.TweetWrapper, alignmentCriteria types.AlignmentCriteria, client *genai.Client, geminiRateLimiter *rate.Limiter) ([]types.AuditResult, error) {
	delay := 15 * time.Second
	for attempt := 1; ; attempt++ {
		if err := geminiRateLimiter.Wait(context.Background()); err != nil {
			return nil, err
		}

		results, err := ai.SendTweetToGemini(context.Background(), batch, alignmentCriteria, client)
		if err == nil {
			return results, nil
		}

		log.Printf("gemini batch failed (attempt %d), retrying same tweets: %v", attempt, err)
		time.Sleep(delay)
		if delay < 2*time.Minute {
			delay *= 2
		}
	}
}

func CSVWriter(resultjob <-chan types.AuditResult, writer *csv.Writer, config types.AlignmentCriteria) error {
	for result := range resultjob {
		err := writer.Write([]string{
			"https://x.com/" + config.Username + "/status/" + result.TweetID,
			strconv.FormatBool(result.Delete),
		})

		if err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}
