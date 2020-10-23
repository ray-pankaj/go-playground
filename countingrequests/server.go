package countingrequests

import (
	"context"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

var ctx = context.Background()

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	log.Fatal(http.ListenAndServe(":8080", InitRouter()))
}
