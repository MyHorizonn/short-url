package testdb

import "short-url/internal/storage/urls"

type TestDb struct {
}

func (r *TestDb) Create(url urls.Url) error {
	return nil
}

func (r *TestDb) Get(key uint64) (string, error) {
	return "github.com/myhorizonn/short-url", nil
}

func (r *TestDb) IsExists(key uint64) (bool, error) {
	return true, nil
}
