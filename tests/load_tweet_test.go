package tests

import (
	"os"
	"path/filepath"
	"testing"
	"tweet-audit/src/app/archive"
)

func TestLoadTweet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tweets.js")

	contents := []byte(`window.YTD.tweets.part0 = [
			{
				"tweet": {
				"id_str": "123",
				"full_text": "hello",
				"created_at": "Mon Jan 01 00:00:00 +0000 2024",
				"lang": "en"
				}
			}
	]`)

	if err := os.WriteFile(path, contents, 0o644); err != nil {
		t.Fatal(err)
	}

	tweets, err := archive.LoadTweetDetails(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(tweets) != 1 {
		t.Fatalf("len = %d, want 1", len(tweets))
	}

	if tweets[0].Tweet.ID != "123" {
		t.Fatalf("Expected 123 but got %s", tweets[0].Tweet.ID)
	}

	if tweets[0].Tweet.FullText != "hello" {
		t.Fatalf("Expected hello but got %s", tweets[0].Tweet.FullText)
	}
}
