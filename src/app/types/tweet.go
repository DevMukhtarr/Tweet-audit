package types

type TweetWrapper struct {
	Tweet Tweet `json:"tweet"`
}

type Tweet struct {
	ID        string `json:"id_str"`
	FullText  string `json:"full_text"`
	CreatedAt string `json:"created_at"`
	Language  string `json:"lang"`
}

type AuditResult struct {
	TweetID string
	Delete  bool
}

type BatchItem struct {
	ID   string `json:"id"`
	Text string `json:"full_text"`
}
