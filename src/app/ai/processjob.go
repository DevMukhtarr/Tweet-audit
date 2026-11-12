package ai

import (
	"context"
	"encoding/json"
	"time"
	"tweet-audit/src/app/types"

	"google.golang.org/genai"
)

func SendTweetToGemini(ctx context.Context, batch []types.TweetWrapper, alignment types.AlignmentCriteria, client *genai.Client) ([]types.AuditResult, error) {
	callContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	rules, err := json.Marshal(alignment)

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{
				Text: `You are auditing a batch of tweets against the alignment rules below.

				Input is a JSON array of tweets: [{"id":"...","full_text":"..."}, ...]

				For every input tweet, decide whether to delete it.
				- delete=true means delete (does not align with the rules)
				- delete=false means keep (aligns with the rules)

				Return ONLY a valid JSON array. No markdown. No extra text.
				Include every input id once, same order as the input:

				[{"id":"123","delete":true},{"id":"456","delete":false}]

				Alignment rules:
				` +
					string(rules),
			}},
		},
	}

	var parsed []struct {
		ID     string `json:"id"`
		Delete bool   `json:"delete"`
	}

	items := make([]types.BatchItem, 0, len(batch))

	for _, t := range batch {
		items = append(items, types.BatchItem{
			ID:   t.Tweet.ID,
			Text: t.Tweet.FullText,
		})
	}

	payload, err := json.Marshal(items)
	if err != nil {
		return []types.AuditResult{}, err
	}

	response, err := client.Models.GenerateContent(
		callContext,
		"gemini-3.5-flash-lite",
		genai.Text(
			string(payload), //send text in batch
		),
		config,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(response.Text()), &parsed); err != nil {
		return []types.AuditResult{}, err
	}

	results := make([]types.AuditResult, 0, len(parsed))
	for _, result := range parsed {
		results = append(results, types.AuditResult{
			TweetID: result.ID,
			Delete:  result.Delete,
		})
	}

	return results, nil
}
