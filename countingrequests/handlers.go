package countingrequests

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/gorilla/handlers"
)

func updateRedisKey() (int64, error) {
	counter, err := rdb.Incr(ctx, "mycounter").Result()
	if err != nil {
		return -1, err
	}
	return counter, nil

}

func getIntegerKeyFromRedis(keyName string) (int64, error) {
	stringKey, err := rdb.Get(ctx, keyName).Result()
	if err != nil {
		return -1, err
	}
	intKey, err := strconv.ParseInt(stringKey, 10, 64)
	if err != nil {
		return -1, err
	}
	return intKey, nil
}

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
	hitCount, err := updateRedisKey()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		countRes := CountResponse{hitCount}
		respondWithJSON(w, countRes)
	}
}
func respondWithJSON(w http.ResponseWriter, resObject interface{}) {
	jsonObject, err := json.Marshal(resObject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonObject)
}
