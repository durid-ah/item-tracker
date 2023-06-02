package userendpoints

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/durid-ah/item-tracker/dto"
	"github.com/durid-ah/item-tracker/helpers"
	"github.com/durid-ah/item-tracker/services"
	"github.com/golang-jwt/jwt/v4"
)

// Create the Signin handler
func Signin(db *sql.DB) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	   var creds dto.Credentials
	   userSvc := services.UserService{ Db: db }

	   err := json.NewDecoder(r.Body).Decode(&creds)
	   if err != nil {
		   w.WriteHeader(http.StatusBadRequest)
		   return
	   }

	   user, err := userSvc.GetByUsername(creds.Username)
	   if err != nil || !helpers.ValidatePassword(creds.Password, user.Password) {
		   w.WriteHeader(http.StatusUnauthorized)
		   return
	   }

	   expirationTime := time.Now().Add(5 * time.Minute)
	   claims := &dto.Claims{
		   Username: creds.Username,
		   RegisteredClaims: jwt.RegisteredClaims{
			   ExpiresAt: jwt.NewNumericDate(expirationTime),
		   },
	   }

	   // Declare the token with the algorithm used for signing, and the claims
	   token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	   // Create the JWT string
	   tokenString, err := token.SignedString(dto.Key)
	   if err != nil {
		   w.WriteHeader(http.StatusInternalServerError)
		   return
	   }

	   // Finally, we set the client cookie for "token" as the JWT we just generated
	   // we also set an expiry time which is the same as the token itself
	   http.SetCookie(w, &http.Cookie{
		   Name:    "token",
		   Value:   tokenString,
		   Expires: expirationTime,
	   })
   })
}
