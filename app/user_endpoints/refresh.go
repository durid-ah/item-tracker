package userendpoints

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/durid-ah/item-tracker/dto"
	"github.com/golang-jwt/jwt/v4"
)

func Refresh(db *sql.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			tknStr := c.Value
			claims := &dto.Claims{}
			tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
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

			// We ensure that a new token is not issued until enough time has elapsed
			// In this case, a new token will only be issued if the old token is within
			// 30 seconds of expiry. Otherwise, return a bad request status
			if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Now, create a new token for the current use, with a renewed expiration time
			expirationTime := time.Now().Add(5 * time.Minute)
			claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(dto.Key)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Set the new token as the users `token` cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})
		})
}