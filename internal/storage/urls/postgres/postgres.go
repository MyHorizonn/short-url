package postgres

import (
	"database/sql"
	"short-url/internal/storage/urls"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Client *sql.DB
}

func (p *Postgres) Create(url urls.Url) error {
	_, err := p.Client.Exec("INSERT INTO urls(original, short, key) Values($1, $2, $3)", url.OriginalUrl, url.ShortUrl, url.Key)
	return err
}

func (p *Postgres) Get(key uint64) (string, error) {
	var url string = ""
	err := p.Client.QueryRow("select original from urls where key = $1", key).Scan(&url)
	if err != nil || url == "" {
		return "", err
	}
	return url, nil
}

func (p *Postgres) IsExists(key uint64) (bool, error) {
	res, err := p.Client.Exec("select * from urls where key = $1", key)
	if err != nil {
		return false, nil
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
