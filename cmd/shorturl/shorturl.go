package main

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"short-url/internal/server"
	"short-url/internal/storage/urls"
	"short-url/internal/storage/urls/postgres"
	"short-url/internal/storage/urls/redisdb"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	var choose_db string
	if len(os.Args) > 1 {
		choose_db = os.Args[1]
	}
	var db urls.Storage
	switch choose_db {
	case "postgres":
		connStr := "host=localhost port=5432 user=postgres password=1234 dbname=shorturl sslmode=disable"
		dbOp, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalln(err)
			return
		}
		dbOp.SetMaxOpenConns(10)
		db = &postgres.Postgres{Client: dbOp}
	default:
		db = &redisdb.Redis{Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})}
	}
	rand.Seed(time.Now().UnixNano())
	server.StartServer(db)
}
