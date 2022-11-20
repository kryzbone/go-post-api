package helpers

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
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

func Validate(obj interface{}, res http.ResponseWriter) error {
	validate := validator.New()
	err := validate.Struct(obj)
	if err != nil {
		ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return err
	}
	return nil
}
