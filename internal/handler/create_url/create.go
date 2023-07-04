package create_handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"short-url/internal/encdec"
	handler_types "short-url/internal/handler/types"
	"short-url/internal/storage/urls"
)

func CreateShortURL(w http.ResponseWriter, r *http.Request, db urls.Storage) {
	switch r.Method {
	case "POST":
		var req handler_types.Req
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body, %v", err), http.StatusInternalServerError)
		}
		var shortUrl string
		var key uint64
		shortUrl, key = encdec.Encode(req.Url)
		exists, exErr := db.IsExists(key)
		if exErr != nil {
			http.Error(w, fmt.Sprintf("error reading db row, %v", exErr), http.StatusBadRequest)
		}
		if !exists {
			url := urls.Url{OriginalUrl: req.Url, ShortUrl: shortUrl, Key: key}
			createErr := db.Create(url)
			if createErr != nil {
				log.Fatalln(err)
				http.Error(w, fmt.Sprintf("error creating db row, %v", createErr), http.StatusInternalServerError)
			}
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		response := handler_types.Resp{Url: shortUrl}
		jsonErr := json.NewEncoder(w).Encode(response)
		if jsonErr != nil {
			log.Fatalf(jsonErr.Error())
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
