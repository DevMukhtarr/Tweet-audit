package archive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"tweet-audit/src/app/types"
)

func LoadTweetDetails(filepath string) ([]types.TweetWrapper, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// find json array
	start := bytes.Index(data, []byte("["))
	end := bytes.LastIndex(data, []byte("]"))

	if start == -1 || end == -1 || end < start {
		return nil, errors.New("invalid tweet array")
	}

	jsonData := data[start : end+1]

	var tweets []types.TweetWrapper

	if err = json.Unmarshal(jsonData, &tweets); err != nil {
		return nil, err
	}

	return tweets, err
}

func SendToChannel(channel chan types.TweetWrapper) {
	var tweets []types.TweetWrapper
	tweets, err := LoadTweetDetails("../twitter-2026-11-06-12828b4e5188433b120ab869aa0b626f2ead885d7106f1c24b3698088e0fdf14/data/tweets.js")
	if err != nil {
		fmt.Println(err)
	}
	// push to channel
	for _, tweet := range tweets {
		channel <- tweet
	}
	close(channel)
}
