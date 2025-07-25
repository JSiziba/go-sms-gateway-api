package middleware

import (
	"encoding/json"
	"go-sms-gateway-api/config"
	"net/http"
)

func AuthMiddleware(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			isAuthRequired := cfg.XRequireWhiskAuth
			if !isAuthRequired {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("X-Require-Whisk-Auth")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Authorization header required"})
				return
			}

			authSecret := cfg.XRequireWhiskAuthSecret

			if authHeader != authSecret {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid or missing authorization secret"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
