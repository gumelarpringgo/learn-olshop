package handler

import (
	"context"
	"errors"
	"learn/config"
	"net/http"
	"strings"
)

var errNoAuthHeaderIncluded = errors.New("no authorization header included")

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")

		if accessToken == "" {
			WriteErrorResponse(w, http.StatusUnauthorized, errNoAuthHeaderIncluded)
			return
		}

		token := ""
		arrayToken := strings.Split(accessToken, " ")
		if len(arrayToken) == 2 {
			token = arrayToken[1]
		}

		claims, err := config.Parse(token)
		if err != nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), "userInfo", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
