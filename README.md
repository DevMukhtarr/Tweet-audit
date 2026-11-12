# Tweet audit

This project reads an X archive and asks Gemini which tweets no longer fit a set of personal rules. Tweets that come back marked for deletion are written to a CSV so they can be removed by hand.

## How a run moves

The loader reads `tweets.js` and pushes one tweet at a time onto a buffered channel. If that channel is full, the loader waits. It does not pile the whole archive into the queue.

A worker pulls tweets off the channel and groups them into a batch. Each tweet sent to Gemini includes its id and full text. The archive also carries the created date and language, but the model only needs the id and the text to make a decision.

```text
loader (waits if the channel is full)
  -> channel
  -> worker (batch tweets; if Gemini fails, the same batch is retried)
  -> results channel
  -> CSV
```

If a batch fails, those tweets are already out of the jobs channel, but they are still in the worker's batch. The worker sends that same batch again instead of skipping them. The pause starts at 15 seconds and doubles until it reaches 2 minutes.

Results go onto a second channel. One writer turns each result into a row:

```text
tweet_url,deleted
https://x.com/<username>/status/<tweet id>,true
```

`deleted` is the model's decision: `true` means the tweet does not match the rules.

## Alignment

The rules live in `config.example.json`. The shape looks like this:

```json
{
  "criteria": {
    "forbidden_words": ["crypto", "NFT", "hustlegrindset"],
    "professional_check": true,
    "tone": "respectful and thoughtful",
    "exclude_politics": true
  },
  "username": "your-handle"
}
```

`username` is only used to build the profile URL. The Gemini API key stays in `.env` as `GEMINI_API_KEY` and is not part of this file.

## Run

From the repository root, with `.env` and `config.example.json` in place:

```bash
go run ./src/cmd/server
```

`GET /health` returns `{"status":"ok"}`. The audit writes `tweets.csv` while the server is running.

Tests, which do not call Gemini:

```bash
go test ./tests/
```
