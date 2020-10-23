package countingrequests

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gorilla/handlers"
)

func updateRedisKey() (int64, error) {
	counter, err := rdb.Incr(ctx, "mycounter").Result()
	if err != nil {
		return -1, err
	}
	return counter, nil

}

func loggingHandler(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func getCountHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {

		if r := recover(); r != nil {
			log.Print("panic in updating counter\n", r, string(debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	counter_from_redis, err := rdb.Get(ctx, "mycounter").Result()
	if err != nil {
		//respondWithError(w, err)
	} else {
		//respondWithSuccess(w, counter_from_redis)
		fmt.Fprintf(w, "%v\n", counter_from_redis)
	}
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {

		if r := recover(); r != nil {
			log.Print("panic in updating counter\n", r, string(debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	counter_from_redis, err := updateRedisKey()
	if err != nil {
		//respondWithError(w, err)
	} else {
		//respondWithSuccess(w, counter_from_redis)
		fmt.Fprintf(w, "%v\n", counter_from_redis)
	}
}
