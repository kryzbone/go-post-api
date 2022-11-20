package helpers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

const DB_NAME = "go-api"
const USER_COLLECTION = "users"
const POST_COLLECTIOM = "posts"

func LoadEnv() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// client request validator function
func Validate(obj interface{}, res http.ResponseWriter) error {
	validate := validator.New()
	err := validate.Struct(obj)
	if err != nil {
		ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return err
	}
	return nil
}

// Hash password function
func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check hashed password function
func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate token function
func GenerateJWT(user_id string) (string, error) {
	//get secret from env variables
	secretkey := os.Getenv("secret")

	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	//Add user id and expiry to token
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		CheckErr(err)
		return "", err
	}
	return tokenString, nil
}

// Fnction to check if value is in slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
