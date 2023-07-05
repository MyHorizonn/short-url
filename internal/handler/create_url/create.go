package create_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"short-url/internal/encdec"
	handler_types "short-url/internal/handler/types"
	"short-url/internal/storage/urls"
	"strconv"
	"time"
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
			expire_days, expireErr := strconv.Atoi(os.Getenv("EXPIRE_DAYS"))
			if err != nil {
				http.Error(w, fmt.Sprintf("error creating, %v", expireErr), http.StatusInternalServerError)
				return
			}
			time_now := time.Now()
			time_now = time_now.AddDate(0, 0, expire_days)
			url := urls.Url{OriginalUrl: req.Url, ShortUrl: shortUrl, Key: key, Expire: time_now}
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
