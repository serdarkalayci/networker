package handlers

import "github.com/serdarkalayci/networker/data"

// APIContext handler for getting and updating Ratings
type APIContext struct {
	LogRecord data.ConsolidatedLog
}

// NewAPIContext returns a new APIContext handler with the given logger
func NewAPIContext() *APIContext {
	// consolidatedLog := data.ConsolidatedLog{
	// 	TotalCalls: 1,
	// }
	return &APIContext{data.ConsolidatedLog{}}
}
