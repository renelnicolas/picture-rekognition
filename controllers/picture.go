package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"ohmytech.io/picture-rekognition/awsbuckets3"

	"ohmytech.io/picture-rekognition/awsdynamodb"
)

// PictureController :
type PictureController struct {
}

// List :
func (h PictureController) List(w http.ResponseWriter, r *http.Request) {
	result, awserr := awsdynamodb.FindAll(nil, "img-rekognition")
	if awserr != nil {
		http.Error(w, awserr.Err, http.StatusInternalServerError)
		return
	}

	mconv, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mconv)
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
	err := r.ParseMultipartForm(200000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formdata := r.MultipartForm

	files := formdata.File["multiplefiles"]

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileName := files[i].Filename
		filePath := "/tmp/" + fileName

		out, err := os.Create(filePath)

		defer out.Close()
		if err != nil {
			http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusInternalServerError)
			return
		}

		awsbuckets3.UploadObject(sess, filePath, "dev-img-rekognition", `waiting/`+fileName)
	}

	w.WriteHeader(http.StatusOK)
	return
}
