package handlers

import (
	"net/http"

	"github.com/serdarkalayci/networker/data"
)

// swagger:route GET /log Log log
// Returns OK if there's no problem
// responses:
//	200: OK

// Logs returns OK handles GET requests
func (p *APIContext) Logs(rw http.ResponseWriter, r *http.Request) {
	respondWithJSON(rw, r, 200, p.LogRecord)
}

// Reset resets the log record history
func (p *APIContext) Reset(rw http.ResponseWriter, r *http.Request) {
	p.LogRecord = data.ConsolidatedLog{}
	respondEmpty(rw, r, 200)
}
