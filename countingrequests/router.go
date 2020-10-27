package countingrequests

import (
	"time"

	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	rl := RateLimiter{60, time.Minute, time.Minute * 5}
	router.Use(loggingHandler)
	router.Use(rl.rateLimitHandler)
	router.HandleFunc("/counter", getCountHandler).Methods("GET")
	router.HandleFunc("/counter", countHandler).Methods("POST")
	return router
}
