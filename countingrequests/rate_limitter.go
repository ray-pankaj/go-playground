package countingrequests

import (
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimitter struct {
	Rate          int
	Period        time.Duration
	LimitDuration time.Duration
}

var incrementScript = redis.NewScript(`
local given_key = KEYS[1]
local time_to_live = ARGV[1]
local limit = tonumber(ARGV[2])
local block_key = given_key .. "_blocked"
local time_to_block = ARGV[3]
local saved_key = redis.call("GET", given_key)
local current_value = 0
if not saved_key then
	redis.call("SETEX", given_key, time_to_live, 1)
	current_value = 1
else
	current_value = redis.call("INCR", given_key)
end
if current_value >= limit then
	redis.call("SETEX", block_key, time_to_block, 1)
end
return {}
`)

func (rl *RateLimitter) isMitigatedClient(user_ip string) bool {
	blockedKey := user_ip + "_blocked"
	_, err := rdb.Get(ctx, blockedKey).Result()
	if err != nil {
		if redis.Nil.Error() == err.Error() {
			return false
		}
		log.Println(err)
		return false
	}
	return true
}

func (rl *RateLimitter) updateClientRequestCount(user_ip string) {
	_, err := incrementScript.Run(ctx, rdb, []string{user_ip}, rl.Period.Seconds(), rl.Rate, rl.LimitDuration.Seconds()).Result()
	if err != nil {
		log.Println(err)
	}
}

func (rl *RateLimitter) rateLimitHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.isMitigatedClient(r.Header.Get("X-Forwarded-For")) {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		go rl.updateClientRequestCount(r.Header.Get("X-Forwarded-For"))
		next.ServeHTTP(w, r)
	})
}
