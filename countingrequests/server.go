package countingrequests

import (
	"context"
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

var ctx = context.Background()

func (c *countingHandler) updateRedisKey() int64 {
	counter, err := c.rdb.Incr(ctx, "mycounter").Result()
	if err != nil {
		panic(err)
	}
	return counter
}

func (c *countingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.updateMyKey()
	var counter_from_redis int64
	counter_from_redis = c.updateRedisKey()
	fmt.Fprintf(w, "%d\n", counter_from_redis)
}

func Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// http.HandleFunc("/counter", ch.ServeHTTP)
	http.Handle("/counter", &countingHandler{rdb: rdb})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
