package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/luisthieme/GoMotion/internal"
)

type contextKey string
const ClientProfileKey = contextKey("clientProfile")

func TokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var clientId = r.URL.Query().Get("clientId")
		clientProfile, ok := internal.Database[clientId]

		if !ok || clientId == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		token := r.Header.Get("Authorization")

		if !isValidToken(clientProfile, token) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), ClientProfileKey, clientProfile)
		r = r.WithContext(ctx)

		next.ServeHTTP(w,r)

	}
}

func isValidToken(profile internal.ClientProfile, token string) bool{
	if strings.HasPrefix(token, "Bearer ") {
		return strings.TrimPrefix(token, "Bearer ") == profile.Token
	}

	return false
}
