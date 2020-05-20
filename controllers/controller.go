package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"ohmytech.io/picture-rekognition/middlewares"
	"ohmytech.io/picture-rekognition/models"
)

const (
	// MAXLIMIT :
	MAXLIMIT = int64(100)
)

// ControllerControllerer :
type ControllerControllerer interface {
	List(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// JSONResponseList :
type JSONResponseList struct {
	Entities interface{} `json:"entities"`
	Counter  int64       `json:"counter"`
}

// JSONResponse :
type JSONResponse struct {
	Entity interface{} `json:"entity"`
}

// JSONResponseError :
type JSONResponseError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// ErrorResponse :
func ErrorResponse(w http.ResponseWriter, status int, message string, err error) {
	re := JSONResponseError{
		Status: status,
		Error:  message,
	}

	mconv, _ := json.Marshal(re)

	log.Printf(">> Message|Error : %s >> %s", message, err)

	w.WriteHeader(status)
	w.Write(mconv)
}

func extractFilter(filter *models.QueryFilter, queryValues url.Values, paginationMaxLimit int64) {
	if 0 < len(queryValues["offset"]) {
		n, err := strconv.ParseInt(queryValues["offset"][0], 10, 64)
		if nil == err {
			filter.Offset = n
		}
	}

	if 0 < len(queryValues["limit"]) {
		n, err := strconv.ParseInt(queryValues["limit"][0], 10, 64)
		if nil == err {
			if paginationMaxLimit >= n {
				filter.Limit = n
			}
		}
	}

	if 0 < len(queryValues["order"]) {
		filter.Order = queryValues["order"][0]
	}

	if 0 < len(queryValues["sort"]) {
		filter.Sort = queryValues["sort"][0]
	}

	if 0 < len(queryValues["search"]) {
		filter.Search = queryValues["search"][0]
	}

	log.Println("extractFilter", filter)
}

func extractIDFromMuxVars(matchVars map[string]string) (int64, error) {
	ID, ok := matchVars["id"]
	if !ok {
		return 0, errors.New("ID cannot be empty")
	}

	entityID, err := strconv.ParseInt(ID, 10, 64)
	if nil != err {
		return 0, fmt.Errorf("Conversion error: %s", err.Error())
	}

	return entityID, nil
}

var fs = http.FileServer(http.Dir("./static"))

// Handlers :
func Handlers() *mux.Router {
	// Main router
	r := mux.NewRouter()

	home := HomeController{}

	// For all
	r.HandleFunc("/", home.Home).Methods(http.MethodGet, http.MethodHead, http.MethodOptions)
	r.HandleFunc("/healthz", home.Healtz).Methods(http.MethodGet)
	r.PathPrefix("/favicon.ico").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/static/{alll}").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.PathPrefix("/static/{name}.{extension}").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/.well-known/dnt-policy.txt", home.DNTpolicy).Methods(http.MethodGet)

	picture := PictureController{}

	r.HandleFunc("/pictures", picture.List).Methods(http.MethodGet)
	r.HandleFunc("/picture/upload", picture.Upload).Methods(http.MethodPost)
	r.HandleFunc("/picture/{id}", picture.Edit).Methods(http.MethodGet)

	r.Use(middlewares.Cors)
	r.Use(middlewares.Logger)

	r.NotFoundHandler = http.HandlerFunc(home.NotFound)

	return r
}
