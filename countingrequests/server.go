package countingrequests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
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
	var counter_from_redis int64

	defer func() {
		fmt.Println(*r)

		if r := recover(); r != nil {
			// log.Print("panic in updating counter\n", r, string(debug.Stack()))
			// TODO: change fmt's to log.* and figure out levels
			fmt.Println("panic in updating counter", r, string(debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%d\n", counter_from_redis)
	}()

	counter_from_redis = c.updateRedisKey()
}

func Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// http.HandleFunc("/counter", ch.ServeHTTP)
	http.Handle("/counter", &countingHandler{rdb: rdb})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
