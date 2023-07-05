package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "short-url/grpc/protob"
	"short-url/internal/encdec"
	"short-url/internal/storage/urls"
	"short-url/internal/storage/urls/postgres"
	"short-url/internal/storage/urls/redisdb"
	"time"
	"os"
	"strconv"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"database/sql"
)

var (
	port = flag.Int("port", 50051, "gRPC server port")
	Db urls.Storage
 )

type server struct {
	pb.UnimplementedShortUrlerServer
}

func init() {
	var choose_db string
	if len(os.Args) > 1 {
		choose_db = os.Args[1]
	}
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
		Db = &postgres.Postgres{Client: dbOp}
		go func() {
			fmt.Println("Deleting expired urls...")
			err := Db.DelExpire()
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Duration(time.Hour))
		}()
	default:
		Db = &redisdb.Redis{Client: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: "",
			DB:       0,
		})}
	}
}

func (s *server) GetShortUrl(ctx context.Context, in *pb.GetShortUrlRequest) (*pb.GetShortUrlResponse, error) {
	log.Printf("Received: %v", in.GetOrigUrl())
	var shortUrl string
	var key uint64
	shortUrl, key = encdec.Encode(in.GetOrigUrl())
	exists, exErr := Db.IsExists(key)
	if exErr != nil {
		return &pb.GetShortUrlResponse{ShortUrl: ""}, nil
	}
	if !exists {
		expire_days, expireErr := strconv.Atoi(os.Getenv("EXPIRE_DAYS"))
		if expireErr != nil {
			return &pb.GetShortUrlResponse{ShortUrl: ""}, nil
		}
		time_now := time.Now()
		time_now = time_now.AddDate(0, 0, expire_days)
		url := urls.Url{OriginalUrl: in.GetOrigUrl(), ShortUrl: shortUrl, Key: key, Expire: time_now}
		createErr := Db.Create(url)
		if createErr != nil {
			return &pb.GetShortUrlResponse{ShortUrl: ""}, nil
		}
	}
	return &pb.GetShortUrlResponse{ShortUrl: shortUrl}, nil
}

func (s *server) GetOriginalUrl(ctx context.Context, in *pb.GetOriginalUrlRequest) (*pb.GetOriginalUrlResponse, error) {
	log.Printf("Received: %v", in.GetOrigUrl())
	res := encdec.Decode(in.GetOrigUrl())
	exists, exErr := Db.IsExists(res)
	if exErr != nil {
		return &pb.GetOriginalUrlResponse{ShortUrl: ""}, nil
	}
	if exists {
		origUrl, getErr := Db.Get(res)
		if getErr != nil {
			return &pb.GetOriginalUrlResponse{ShortUrl: origUrl}, nil
		}
	} else {
		return &pb.GetOriginalUrlResponse{ShortUrl: ""}, nil
	}
	return &pb.GetOriginalUrlResponse{ShortUrl: ""}, nil
}


func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterShortUrlerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}