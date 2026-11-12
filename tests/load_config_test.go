package tests

import (
	"slices"
	"testing"
	"tweet-audit/src/app/config"
	"tweet-audit/src/app/types"
)

func TestLoadConfig(t *testing.T) {

	values := types.AlignmentCriteria{
		Criteria: types.Criteria{
			ForbiddenWords:    []string{"crypto", "NFT"},
			ProfessionalCheck: true,
			Tone:              "friendly",
			ExcludePolitics:   true,
		},
		Username: "devmukhtar",
	}

	config := config.ReadConfig()

	if config.Criteria.ExcludePolitics != values.Criteria.ExcludePolitics ||
		!slices.Equal(values.Criteria.ForbiddenWords, config.Criteria.ForbiddenWords) ||
		config.Criteria.Tone != values.Criteria.Tone ||
		config.Criteria.ProfessionalCheck != values.Criteria.ProfessionalCheck ||
		config.Username != values.Username {
		t.Errorf("Values don't match")
	}
}
