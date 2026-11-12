package config

import (
	"encoding/json"
	"os"
	"tweet-audit/src/app/types"
)

func ReadConfig() types.AlignmentCriteria {
	data, err := os.ReadFile("../config.example.json")

	if err != nil {
		panic(err)
	}

	var config types.AlignmentCriteria
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}
