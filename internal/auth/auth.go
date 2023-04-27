package auth

import (
	"bcraftTestTask/internal/logging"
	"bcraftTestTask/internal/models"
	"bcraftTestTask/internal/properties"
	u "bcraftTestTask/internal/utils"
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/user/new", "/user/login",
			"/probes/liveness", "/probes/readiness"}
		requestPath := r.URL.Path

		config, err := properties.GetConfig()
		if err != nil {
			logger := logging.GetLogger()
			logger.Fatal(err)
		}

		for _, value := range notAuth {
			if value == requestPath {
				ctx := context.WithValue(r.Context(), "user", uint64(0))
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			u.Respond(w, u.Message(false, "Missing auth token",
				"REST API JwtAuthentication"), http.StatusUnauthorized)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			u.Respond(w, u.Message(false, "Invalid/Malformed auth token",
				"REST API JwtAuthentication"), http.StatusUnauthorized)
			return
		}

		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.ProgramSettings.TokenPassword), nil
		})

		if err != nil {
			u.Respond(w, u.Message(false, "Malformed authentication token",
				"REST API JwtAuthentication"), http.StatusUnauthorized)
			return
		}

		if !token.Valid {

			u.Respond(w, u.Message(false, "Token is not valid",
				"REST API JwtAuthentication"), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
