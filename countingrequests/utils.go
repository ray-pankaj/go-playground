package countingrequests

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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

func respondWithJSON(w http.ResponseWriter, resObject interface{}) {
	jsonObject, err := json.Marshal(resObject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonObject)
}
