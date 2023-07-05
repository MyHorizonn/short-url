package urls

import "time"

type Url struct {
	OriginalUrl string    `json:"originalurl" db:"original"`
	ShortUrl    string    `json:"shorturl" db:"short"`
	Key         uint64    `json:"key" db:"urlkey"`
	Expire      time.Time `json:"expire" db:"expire"`
}

type Storage interface {
	IsExists(key uint64) (bool, error)
	Create(url Url) error
	Get(key uint64) (string, error)
	DelExpire() error
}
