package controllers

import (
	"net/http"
)

// PictureController :
type PictureController struct {
}

// List :
func (h PictureController) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// Edit :
func (h PictureController) Edit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return

}

// Update :
func (h PictureController) Update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// Upload :
func (h PictureController) Upload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}
