package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"short-url/internal/handler"
	"short-url/internal/storage/urls"
	"short-url/internal/storage/urls/postgres"
	"short-url/internal/storage/urls/redisdb"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	var choose_db string
	if len(os.Args) > 1 {
		choose_db = os.Args[1]
	}
	var db urls.Storage
	err := godotenv.Load("./././.env")
	if err != nil {
		log.Fatalln(err)
	}
	switch choose_db {
	case "postgres":
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
		dbOp, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalln(err)
		}
		dbOp.SetMaxOpenConns(10)
		db = &postgres.Postgres{Client: dbOp}
		go func() {
			fmt.Println("Deleting expired urls...")
			err := db.DelExpire()
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Duration(time.Hour))
		}()
	default:
		db = &redisdb.Redis{Client: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: "",
			DB:       0,
		})}
	}
	handler.StartServer(db)
}
