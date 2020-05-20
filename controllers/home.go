package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// HomeController :
type HomeController struct {
}

// Home :
func (h HomeController) Home(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusOK)
	return
}

// Healtz :
func (h HomeController) Healtz(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusOK)
	return
}

// DNTpolicy :
func (h HomeController) DNTpolicy(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusOK)
	return
}

// NotFound :
func (h HomeController) NotFound(w http.ResponseWriter, r *http.Request) {
	log.Println(strings.Repeat(">", 15) + " NotFound " + strings.Repeat("<", 15))
	log.Printf("%s || %s\n", r.Host, r.URL.Path)

	for k, v := range r.Header {
		fmt.Printf("%v : %v\n", strings.ToLower(k), string(v[0]))
	}

	ErrorResponse(w, http.StatusNotFound, `Not Found`, nil)
	return
}
