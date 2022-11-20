package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/go-post-api/database"
	"example.com/go-post-api/helpers"
	"example.com/go-post-api/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	db, ctx := database.ConnectDB()
	defer database.CloseDB()
	cursor, err := db.Collection(helpers.USER_COLLECTION).Find(ctx, bson.M{})
	if err != nil {
		helpers.SeverError(err, res)
		return
	}

	defer cursor.Close(ctx)

	var users = make([]models.SanitizedUser, 0)

	for cursor.Next(ctx) {
		var user models.SanitizedUser
		cursor.Decode(&user)
		users = append(users, user)
	}
	fmt.Println(users)
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(users)
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var user models.User
	json.NewDecoder(req.Body).Decode(&user)
	fmt.Println(user)
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}

	db, ctx := database.ConnectDB()
	defer database.CloseDB()
	result, err := db.Collection(helpers.USER_COLLECTION).InsertOne(ctx, user)
	if err != nil {
		helpers.SeverError(err, res)
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(result)
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	// convert id to object id that monogo understands
	id, _ := primitive.ObjectIDFromHex(params["id"])
	db, ctx := database.ConnectDB()
	defer database.CloseDB()
	var user models.SanitizedUser
	if err := db.Collection(helpers.USER_COLLECTION).FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&user); err != nil {
		helpers.SeverError(err, res)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(user)
}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	db, ctx := database.ConnectDB()
	defer database.CloseDB()
	result, err := db.Collection(helpers.USER_COLLECTION).DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		helpers.SeverError(err, res)
		return
	}
	helpers.SeverError(err, res)
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(result)

}

func UpdateUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user models.User
	json.NewDecoder(req.Body).Decode(&user)
	db, ctx := database.ConnectDB()
	defer database.CloseDB()

	// validate := validator.New()
	// err := validate.Struct(user)
	// if err != nil {
	// 	helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
	// 	return
	// }

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: user}}
	result, err := db.Collection(helpers.USER_COLLECTION).UpdateOne(ctx, filter, update)
	if err != nil {
		helpers.SeverError(err, res)
		return
	}
	json.NewEncoder(res).Encode(result)
}

// Sign Up Handler
func SignUp(res http.ResponseWriter, req *http.Request) {
	db, ctx := database.ConnectDB()
	defer database.CloseDB()

	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}

	//Validate data from client
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}

	// Check if email already exists
	result := db.Collection(helpers.USER_COLLECTION).FindOne(ctx, bson.D{{Key: "email", Value: user.Email}})
	//if there is no err, it means user exists
	if result.Err() == nil {
		helpers.ErrorResonse("Email already exist", http.StatusBadRequest, res)
		return
	}

	// Hash password
	user.Password, err = helpers.GeneratehashPassword(user.Password)
	if err != nil {
		helpers.SeverError(err, res)
		return
	}

	//insert user details in database
	createRes, err := db.Collection(helpers.USER_COLLECTION).InsertOne(ctx, user)
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(createRes)
}

// Sign In Handler
func SignIn(res http.ResponseWriter, req *http.Request) {
	db, ctx := database.ConnectDB()
	defer database.CloseDB()

	// Extract client data from request
	var authdetails models.Authentication
	err := json.NewDecoder(req.Body).Decode(&authdetails)
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}

	var authuser models.User
	err = db.Collection(helpers.USER_COLLECTION).FindOne(ctx, bson.D{{Key: "email", Value: authdetails.Email}}).Decode(&authuser)
	if err != nil {
		helpers.ErrorResonse("Invalid email or password", http.StatusBadRequest, res)
		return
	}

	// Check password
	check := helpers.CheckPasswordHash(authuser.Password, authdetails.Password)
	if !check {
		helpers.ErrorResonse("Invalid email or password", http.StatusBadRequest, res)
		return
	}

	// Generate JWT Token
	validToken, err := helpers.GenerateJWT(authuser.ID.Hex())
	if err != nil {
		helpers.SeverError(err, res)
		return
	}

	var token models.Token
	token.Token = validToken
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(token)
}
