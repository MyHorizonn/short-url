package create_handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"short-url/internal/storage/urls/testdb"
	"testing"
)

func TestCreateShortUrl(t *testing.T) {

	// orig: github.com/myhorizonn/short-url
	// short: jO0j5jyCyj
	// key: 251567266899663446

	//testData := urls.Url{OriginalUrl: "github.com/myhorizonn/short-url", ShortUrl: "jO0j5jyCyj", Key: 251567266899663446}

	cases := []struct {
		name      string
		jsonStr   []byte
		expStatus int
	}{
		{
			name:      "empty data",
			jsonStr:   []byte(`{}`),
			expStatus: http.StatusBadRequest,
		},
		{
			name:      "ok",
			jsonStr:   []byte(`{"url": "github.com/myhorizonn/short-url"}`),
			expStatus: http.StatusCreated,
		},
	}

	emptyDb := &testdb.TestDb{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateShortURL(w, r, emptyDb)
	})

	var jsonStr = []byte(`{"url": "github.com/myhorizonn/short-url"}`)
	req, err := http.NewRequest("GET", "/create", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err.Error())
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Fatalf("\nTest: %s, expected %d got %d", "method not allowed", http.StatusMethodNotAllowed, status)
	}

	for _, tCase := range cases {
		req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(tCase.jsonStr))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err.Error())
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != tCase.expStatus {
			t.Fatalf("\nTest: %s, expected %d got %d", tCase.name, tCase.expStatus, status)
		}
	}

}
