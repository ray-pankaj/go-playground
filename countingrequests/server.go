package countingrequests

import (
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
)

type countingHandler struct {
	mu  sync.Mutex
	n   int
	rdb *redis.Client
}

func (c *countingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// ch := new(countingHandler)
	// ch.rdb = rdb
	// http.HandleFunc("/counter", ch.ServeHTTP)
	http.Handle("/counter", &countingHandler{rdb: rdb})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
