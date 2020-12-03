package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// swagger:route GET / Index index
// Returns OK if there's no problem
// responses:
//	200: OK

// Index returns OK handles GET requests
func (p *APIContext) Index(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
}

// swagger:route GET /version Index version
// Returns version information
// responses:
//	200: OK

// Version returns the version info for the service by reading from a static file
func (p *APIContext) Version(rw http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("./static/version.txt")
	if err != nil {
		dat = append(dat, '0')
	}
	fmt.Fprintf(rw, "Welcome to Networker API! Version:%s", dat)
}
