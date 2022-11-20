package routers

import (
	"example.com/go-post-api/controllers"
	"example.com/go-post-api/helpers"
	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router) {
	// Create a subrouter for users '/users/*'
	userRouter := router.PathPrefix("/users").Subrouter()

	userRouter.HandleFunc("/", controllers.GetUsers).Methods("GET")
	userRouter.HandleFunc("/", controllers.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}/", controllers.UpdateUser).Methods("PATCH")
	userRouter.HandleFunc("/{id}/", controllers.GetUser).Methods("GET")
	userRouter.HandleFunc("/{id}/", controllers.DeleteUser).Methods("DELETE")

	userRouter.Use(helpers.IsAuthorized)
}
