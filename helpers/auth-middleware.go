package helpers

import (
	"net/http"
	"context"

	"github.com/durid-ah/item-tracker/dto"

	"github.com/golang-jwt/jwt"
)


const UserNameKey = "appUserName"

func WithAuth(handler http.Handler) http.Handler {
	
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			// We can obtain the session token from the requests cookies, which come with every request
			c, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					// If the cookie is not set, return an unauthorized status
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				// For any other type of error, return a bad request status
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Get the JWT string from the cookie
			tknStr := c.Value

			// Initialize a new instance of `Claims`
			claims := dto.Claims{}

			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			tkn, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
				return dto.Key, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			newCtx := context.WithValue(ctx, UserNameKey, claims.Username)

			handler.ServeHTTP(w, r.WithContext(newCtx))
		})
}