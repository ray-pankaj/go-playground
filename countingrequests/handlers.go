package countingrequests

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gorilla/handlers"
)

func loggingHandler(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func getCountHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {

		if r := recover(); r != nil {
			log.Print("panic in getting counter\n", r, string(debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	hitCount, err := getIntegerKeyFromRedis("mycounter")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		countRes := CountResponse{hitCount}
		respondWithJSON(w, countRes)
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
	hitCount, err := rdb.Incr(ctx, "mycounter").Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		countRes := CountResponse{hitCount}
		respondWithJSON(w, countRes)
	}
}
