package controllers

import "net/http"

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(`{"mesage": "Welcome to Go Post Api"}`))
}
