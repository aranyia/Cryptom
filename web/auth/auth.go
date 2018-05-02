package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//SecurityMiddleware defines http.HandlerFunc to be chained before handlers of protected resources.
func SecurityMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		var validationError error

		var authHeader = request.Header.Get("Authorization")
		if len(authHeader) > 0 {
			var authToken = authHeader[len("Bearer "):]

			_, validationError = jwt.Parse(authToken, func(*jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SIGN_KEY")), nil
			})
		} else {
			validationError = errors.New("missing authentication")
		}

		if validationError == nil {
			next.ServeHTTP(respWriter, request)
		} else {
			http.Error(respWriter, validationError.Error(), http.StatusUnauthorized)
		}
	})
}

// GenerateAuthToken generates a JWT authentication token.
func GenerateAuthToken() (string, error) {
	var jwtSignKey = []byte(os.Getenv("JWT_SIGN_KEY"))

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()

	return token.SignedString(jwtSignKey)
}

// IsAuthenticationValid determines whether a provided API-key is valid.
func IsAuthenticationValid(apiKey string) bool {
	var validApiKey = os.Getenv("CRYPTOM_API_KEY")
	return apiKey == validApiKey
}
