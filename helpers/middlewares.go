package helpers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Authorization Middleware
func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// fmt.Println(req.URL.Path)
		// // Non Auth routets
		// noAuthRoutes := []string{"/", "/login/", "/create/"}
		// // Skip middleware
		// if ok := contains(noAuthRoutes, req.URL.Path); ok {
		// 	next.ServeHTTP(res, req)
		// 	return
		// }

		if req.Header["Authorization"] == nil {
			ErrorResonse("Authentication credentials were not provided.", http.StatusUnauthorized, res)
			return
		}

		secretkey := os.Getenv("SECRET")
		var mySigningKey = []byte(secretkey)

		rawToken := req.Header["Authorization"][0]
		rawToken = strings.Split(rawToken, " ")[1]

		token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			ErrorResonse("Invalid Token", http.StatusUnauthorized, res)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// TODO: retrieve user and store on req object
			if claims["user_id"] == "" {
				ErrorResonse("Authentication credentials invalid.", http.StatusUnauthorized, res)
				return
			}
			next.ServeHTTP(res, req)
			return
		}

		ErrorResonse("Not Authorized", http.StatusUnauthorized, res)
	})
}

// func LoggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Do stuff here
// 		fmt.Println(r.RequestURI)
// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(w, r)
// 	})
// }
