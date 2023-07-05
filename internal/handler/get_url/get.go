package get_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"short-url/internal/encdec"
	handler_types "short-url/internal/handler/types"
	"short-url/internal/storage/urls"
)

func GetOriginalURL(w http.ResponseWriter, r *http.Request, db urls.Storage) {
	switch r.Method {
	case "GET":
		var req handler_types.Req
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body, %v", err), http.StatusBadRequest)
			return
		}

		if req == (handler_types.Req{}) {
			http.Error(w, fmt.Sprint("error empty data"), http.StatusBadRequest)
			return
		}

		res := encdec.Decode(req.Url)
		exists, exErr := db.IsExists(res)
		if exErr != nil {
			http.Error(w, fmt.Sprintf("error reading db row, %v", exErr), http.StatusBadRequest)
			return
		}
		if exists {
			origUrl, getErr := db.Get(res)
			if getErr != nil {
				http.Error(w, fmt.Sprintf("error reading db row, %v", getErr), http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			response := handler_types.Resp{Url: origUrl}
			jsonErr := json.NewEncoder(w).Encode(response)
			if jsonErr != nil {
				http.Error(w, fmt.Sprintf("error %v", jsonErr), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, fmt.Sprintf("error url does not exist"), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
