package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/go-post-api/controllers"
	"example.com/go-post-api/helpers"
	"example.com/go-post-api/routers"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Post API Project in Golang")
	helpers.LoadEnv()

	// Get port value from env variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Create mux Router
	router := mux.NewRouter().StrictSlash(true)

	// Setup routes
	router.HandleFunc("/", controllers.IndexHandler)
	router.HandleFunc("/signup/", controllers.SignUp).Methods("POST")
	router.HandleFunc("/login/", controllers.SignIn).Methods("POST")
	routers.UserRouter(router)
	routers.PostRouter(router)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
