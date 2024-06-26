package handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	RequestIDKey = "requestId"
	UserIDKey    = "userID"
)

func loggingMiddleware(lg *logrus.Entry) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String() // generate a new UUID if requestID is not provided
			}

			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				userID = "anonymous" // default userID if not provided
			}

			// Attach the requestId and userID to the context
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
			ctx = context.WithValue(ctx, UserIDKey, userID)
			r = r.WithContext(ctx)

			// Create a new logger entry with the requestId and userID
			entry := lg.WithFields(logrus.Fields{
				"requestId": requestID,
				"userID":    userID,
			})

			// Pass the entry down the request chain
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "logger", entry)))
		})
	}
}
