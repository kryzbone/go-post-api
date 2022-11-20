package routers

import (
	"example.com/go-post-api/controllers"
	"github.com/gorilla/mux"
)

func PostRouter(router *mux.Router) {
	postRouter := router.PathPrefix("/posts").Subrouter()

	postRouter.HandleFunc("/", controllers.GetPosts).Methods("GET")
	postRouter.HandleFunc("/", controllers.CreatePost).Methods("POST")
	postRouter.HandleFunc("/{id}/", controllers.GetPost).Methods("GET")
	postRouter.HandleFunc("/{id}/", controllers.UpdatePost).Methods("PATCH")
	postRouter.HandleFunc("/{id}/", controllers.DeletePost).Methods("DELETE")
}
