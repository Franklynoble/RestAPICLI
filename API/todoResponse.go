package main

import (
	"encoding/json"
	"time"

	"github.com/Franklynoble/todocli"
)

type todoResponse struct {
	Results todocli.List `json:"results"`
}

func (r *todoResponse) MarshalJSON() ([]byte, error) {

	resp := struct {
		Results      todocli.List `json:"results"`
		Date         int64        `json:"date"`
		TotalResults int          `json:"total_results"`
	}{
		Results:      r.Results,
		Date:         time.Now().Unix(),
		TotalResults: len(r.Results),
	}
	return json.Marshal(resp)
}
