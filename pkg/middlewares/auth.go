package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

const (
	Auhorization = "Authorization"
)

type UserToken struct {
	ID   string `json:"id"`
	Uuid string `json:"uuid"`
}

var sessions map[string]UserToken = map[string]UserToken{"Token1": {ID: "ID-1", Uuid: uuid.NewString()}}

func AuthMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := ExtractToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, present := sessions[token]

		if !present {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		ctx = context.WithValue(ctx, "user_info", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ExtractToken(r *http.Request) (token string, err error) {
	authHeader := r.Header.Get(Auhorization)
	if authHeader == "" || len(authHeader) < 7 {
		return "", errors.New("Please authenticate to use API")
	}
	headerType := authHeader[0:7]
	if headerType != "Bearer " {
		return "", errors.New("Bad auth header")
	}
	token = authHeader[7:]
	return
}
