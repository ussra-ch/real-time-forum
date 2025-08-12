package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type RateLimitPosts struct {
	count        int
	FirstTime    time.Time
	BlockedUntil time.Time
	UserId       int
}
type ErrorStruct struct {
	Type string
	Text string
}

var (
	PostRateLimits = make(map[int]*RateLimitPosts)
	userID         int
)

func CheckRateLimitPost(ratelimit *RateLimitPosts, window time.Duration) bool {
	if time.Now().Before(ratelimit.BlockedUntil) {
		return false
	}
	if time.Now().After(ratelimit.BlockedUntil) && ratelimit.count > 10 {
		ratelimit.FirstTime = time.Now()
		ratelimit.BlockedUntil = time.Time{}
		ratelimit.count = 0
	}
	ratelimit.count++
	if ratelimit.count > 2 {
		ratelimit.BlockedUntil = time.Now().Add(window)
		return false
	}
	return true
}

func UserInfosPosts(r *http.Request) (*RateLimitPosts, bool) {
	rateLimit := &RateLimitPosts{
		count:        0,
		FirstTime:    time.Now(),
		BlockedUntil: time.Time{},
		UserId:       -1,
	}
	_, userID = IsLoggedIn(r)

	rateLimit.UserId = userID
	return rateLimit, true
}

func RateLimitPostsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRateLimit, ok := UserInfosPosts(r)
		_, userID = IsLoggedIn(r)
		if !ok {
			errorr := ErrorStruct{
				Type: "error",
				Text: "Unauthorized",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorr)
			return
		}

		ratelimit, exists := PostRateLimits[userRateLimit.UserId]
		if !exists {
			AddUserToTheMap_Post(userRateLimit)
			ratelimit = userRateLimit
		}

		if !CheckRateLimitPost(ratelimit, 1*time.Hour) {
			errorr := ErrorStruct{
				Type: "error",
				Text: "Too many requests",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(errorr)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func AddUserToTheMap_Post(ratelimit *RateLimitPosts) {
	PostRateLimits[ratelimit.UserId] = ratelimit
}
