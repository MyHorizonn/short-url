package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"short-url/internal/encdec"
	"short-url/internal/storage/urls"
)

type resp struct {
	Url string `json:"url"`
}

type req struct {
	Url string `json:"url"`
}

func createShortURL(w http.ResponseWriter, r *http.Request, db urls.Storage) {
	switch r.Method {
	case "POST":
		var req req
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
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func getOriginalURL(w http.ResponseWriter, r *http.Request, db urls.Storage) {
	switch r.Method {
	case "GET":
		var req req
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body, %v", err), http.StatusInternalServerError)
		}

		res := encdec.Decode(req.Url)
		exists, exErr := db.IsExists(res)
		if exErr != nil {
			http.Error(w, fmt.Sprintf("error reading db row, %v", exErr), http.StatusBadRequest)
		}
		if exists {
			origUrl, getErr := db.Get(res)
			if getErr != nil {
				http.Error(w, fmt.Sprintf("error reading db row, %v", getErr), http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			response := resp{Url: origUrl}
			jsonErr := json.NewEncoder(w).Encode(response)
			if jsonErr != nil {
				log.Fatalf(jsonErr.Error())
			}
		} else {
			http.Error(w, fmt.Sprintf("error row does not exist"), http.StatusBadRequest)
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func StartServer(db urls.Storage) {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		createShortURL(w, r, db)
	})
	http.HandleFunc("/get_original", func(w http.ResponseWriter, r *http.Request) {
		getOriginalURL(w, r, db)
	})

	log.Fatalln(http.ListenAndServe(":9000", nil))
}
