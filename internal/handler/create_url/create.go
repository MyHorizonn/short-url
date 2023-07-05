package create_handler

import (
	"encoding/json"
	"fmt"
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
			http.Error(w, fmt.Sprintf("error reading request body, %v", err), http.StatusBadRequest)
			return
		}
		if req == (handler_types.Req{}) {
			http.Error(w, "error empty data", http.StatusBadRequest)
			return
		}
		var shortUrl string
		var key uint64
		shortUrl, key = encdec.Encode(req.Url)
		exists, exErr := db.IsExists(key)
		if exErr != nil {
			http.Error(w, fmt.Sprintf("error reading data, %v", exErr), http.StatusBadRequest)
			return
		}
		if !exists {
			url := urls.Url{OriginalUrl: req.Url, ShortUrl: shortUrl, Key: key}
			createErr := db.Create(url)
			if createErr != nil {
				http.Error(w, fmt.Sprintf("error creating, %v", createErr), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		response := handler_types.Resp{Url: shortUrl}
		jsonErr := json.NewEncoder(w).Encode(response)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error %v", jsonErr), http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
