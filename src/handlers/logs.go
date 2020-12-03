package handlers

import "net/http"

// swagger:route GET /log Log log
// Returns OK if there's no problem
// responses:
//	200: OK

// Logs returns OK handles GET requests
func (p *APIContext) Logs(rw http.ResponseWriter, r *http.Request) {
	respondWithJSON(rw, r, 200, p.LogRecord)
}
