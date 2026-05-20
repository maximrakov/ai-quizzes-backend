package handler

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	pkgjwt "github.com/maximrakov/ai-quizzes-backend/pkg/jwt"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

func AuthMiddleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "отсутствует токен авторизации", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := pkgjwt.Parse(tokenStr, secret)
		if err != nil {
			http.Error(w, "недействительный токен", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		var bodyStr string
		if r.Body != nil && r.ContentLength != 0 {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				bodyStr = string(bodyBytes)
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		args := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"ip", ip,
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
		}
		if bodyStr != "" {
			args = append(args, "body", bodyStr)
		}
		logger.Info("incoming request", args...)
	})
}
