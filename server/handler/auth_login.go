package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AuthPayload is expected payload.
type AuthLoginPayload struct {
	Role string `json:"role"`
	Key  string `json:"key"`
}

// TokenCreator is an interface for creating new JWT provided role and key.
type TokenCreator interface {
	NewToken(role, key string) (token string, expires time.Time, err error)
}

// AuthLogin is a handler that creates jwt auth token and stores it in cookie for
// future authorized requests.
func AuthLogin(tokenCreator TokenCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload AuthLoginPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			HandleError(w, fmt.Errorf("error decoding auth payload: %v", err), http.StatusBadRequest, true)
			return
		}

		token, expires, err := tokenCreator.NewToken(payload.Role, payload.Key)
		if err != nil {
			HandleError(w, fmt.Errorf("error creating new auth token: %v", err), http.StatusUnauthorized, true)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  expires,
			HttpOnly: true,
			Path:     "/",
		})

		w.WriteHeader(http.StatusNoContent)
	}
}
