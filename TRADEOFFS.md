
# Architecture choice
I went with this architecture because each stage only deals with the tweets in front of it. The loader reads the archive and sends tweets onto a bounded channel. A full channel blocks the loader, so the workers are not flooded. Each worker groups tweets into batches of 10 and sends one batch to Gemini. If that call fails, the same batch is retried, because those tweets have left the channel but are still in the worker's slice. Results go onto a second bounded channel, and the CSV writer appends them under the header written when the file was created.

# Batching
I batch 10 tweets per call. A shared limiter already starts calls about 12 seconds apart, so one tweet per request would not blow the per-minute cap by itself. It would spend Gemini daily allowance on a single tweet, and thousands of tweets would take thousands of calls. The cost is a larger prompt, and a failure retries all 10. The retry sends that same prompt again.

# Workers
There are two workers, and both share that one limiter. They do not each keep their own clock. The limiter only spaces when a call may start, so a slow call can still overlap the next one. The jobs channel holds at most 100 tweets. A batch of 10 is one Gemini call, not ten.

# Error handling
On an error or a timeout, the worker keeps the batch and tries again. The pause starts at 15 seconds and doubles until it reaches 2 minutes. Those tweets are not dropped while the process is running. They are lost only if the process is stopped before the call succeeds.

# The model i used
I use `gemini-3.5-flash-lite` instead of a thinking model. A yes-or-no decision on a short batch does not need a long reasoning pass, and the thinking model was timing out before it returned JSON. Each call is allowed 180 seconds.

# The loader
The loader still decodes the whole archive into one slice first. The channel limits the queue, not parse memory, so a much larger archive would sit in memory up front, tweets move in small groups.

# CSV Writer
The CSV columns are `tweet_url` and `deleted`. The URL is `https://x.com/`, the username from the config, `/status/`, and the tweet id. The model is not asked for the username. `deleted` is the boolean Gemini returned for that id. Rows are written only after a batch succeeds. The results channel also holds at most 100 items.

# Tests
The tests load a small archive fixture, load the alignment config, and write one CSV row. They do not call Gemini.

# Why i chose go
I chose Go because the concurrency tools are part of the language. A worker is a goroutine, and a bounded channel is the queue. A full channel blocks the sender without another library. A failed Gemini call returns an error, and the worker retries that batch instead of losing it.