package types

type AlignmentCriteria struct {
	Criteria Criteria `json:"criteria"`
	Username string   `json:"username"`
}
type Criteria struct {
	ForbiddenWords    []string `json:"forbidden_words"`
	ProfessionalCheck bool     `json:"professional_check"`
	Tone              string   `json:"tone"`
	ExcludePolitics   bool     `json:"exclude_politics"`
}
