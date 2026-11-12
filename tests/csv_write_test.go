package tests

import (
	"encoding/csv"
	"os"
	"testing"
	"tweet-audit/src/app/config"
	"tweet-audit/src/app/pipeline"
	"tweet-audit/src/app/types"
)

func TestCSVWrite(t *testing.T) {
	file, err := os.CreateTemp("", "tweets-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	// csv writer
	writer := csv.NewWriter(file)
	writer.Write([]string{
		"tweet_url",
		"deleted",
	})

	// auditResult sample
	result := types.AuditResult{
		TweetID: "123456",
		Delete:  true,
	}

	resultChannel := pipeline.ProcessedResultCollector()
	resultChannel <- result
	close(resultChannel)

	if err := pipeline.CSVWriter(resultChannel, writer, config.ReadConfig()); err != nil {
		t.Fatal(err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 2 {
		t.Fatalf("rows = %d, want 2", len(records))
	}

	wantURL := "https://x.com/devmukhtar/status/123456"

	if records[1][0] != wantURL {
		t.Fatalf("url = %q, want %q", records[1][0], wantURL)
	}

	if records[1][1] != "true" {
		t.Fatalf("deleted = %q, want true", records[1][1])
	}
	// fmt.Println(records[1][0])
	// fmt.Println(records[1][1])
	// fmt.Println(records[1])
}
