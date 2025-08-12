package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

var CommentRateLimits = make(map[int]*RateLimit)

func CheckRateLimitComment(ratelimit *RateLimit, window time.Duration) bool {
	if time.Now().Before(ratelimit.BlockedUntil) {
		return false
	}
	if time.Now().After(ratelimit.BlockedUntil) && ratelimit.count > 50 {
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

func UserInfosComments(r *http.Request) (*RateLimit, bool) {
	rateLimit := &RateLimit{
		count:        0,
		FirstTime:    time.Now(),
		BlockedUntil: time.Time{},
		UserId:       -1,
	}
	_, userID = IsLoggedIn(r)
	rateLimit.UserId = userID
	return rateLimit, true
}

func CommentsRatelimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRateLimit, ok := UserInfosComments(r)
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

		ratelimit, exists := CommentRateLimits[userRateLimit.UserId]
		if !exists {
			AddUserToTheMap_comment(userRateLimit)
			ratelimit = userRateLimit
		}

		if !CheckRateLimitComment(ratelimit, 1*time.Minute) {
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

func AddUserToTheMap_comment(ratelimit *RateLimit) {
	CommentRateLimits[ratelimit.UserId] = ratelimit
}
