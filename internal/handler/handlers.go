package handler

import (
	"log"
	"net/http"
	create_url_handler "short-url/internal/handler/create_url"
	get_url_handler "short-url/internal/handler/get_url"
	"short-url/internal/storage/urls"
)

func StartServer(db urls.Storage) {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		create_url_handler.CreateShortURL(w, r, db)
	})
	http.HandleFunc("/get_original", func(w http.ResponseWriter, r *http.Request) {
		get_url_handler.GetOriginalURL(w, r, db)
	})

	log.Fatalln(http.ListenAndServe(":9000", nil))
}
