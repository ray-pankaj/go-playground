package countingrequests

import (
	"fmt"
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

func (c *countingHandler) updateMyKey() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}

func (c *countingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.updateMyKey()
	fmt.Fprintf(w, "%d\n", c.n)
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
