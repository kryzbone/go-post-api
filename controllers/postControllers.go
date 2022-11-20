package controllers

import (
	"encoding/json"
	"net/http"

	"example.com/go-post-api/database"
	"example.com/go-post-api/helpers"
	"example.com/go-post-api/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPosts(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	db, ctx := database.ConnectDB()
	defer database.CloseDB()

	cursor, err := db.Collection(helpers.POST_COLLECTIOM).Find(ctx, bson.M{})
	if err != nil {
		helpers.SeverError(err, res)
		return
	}
	defer cursor.Close(ctx)

	var posts = make([]models.Post, 0)

	for cursor.Next(ctx) {
		var post models.Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(posts)
}

func CreatePost(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var post models.Post
	json.NewDecoder(req.Body).Decode(&post)

	db, ctx := database.ConnectDB()
	defer database.CloseDB()

	validate := validator.New()
	err := validate.Struct(post)
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}

	result, err := db.Collection(helpers.POST_COLLECTIOM).InsertOne(ctx, post)
	if err != nil {
		helpers.SeverError(err, res)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(result)
}

func GetPost(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}

	db, ctx := database.ConnectDB()
	defer database.CloseDB()

	var post models.Post
	filter := bson.D{{Key: "_id", Value: id}}
	err = db.Collection(helpers.POST_COLLECTIOM).FindOne(ctx, filter).Decode(&post)

	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusNotFound, res)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(post)
}

func UpdatePost(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}

	db, ctx := database.ConnectDB()
	defer database.CloseDB()

	var post models.Post
	json.NewDecoder(req.Body).Decode(&post)

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: post}}

	result, err := db.Collection(helpers.POST_COLLECTIOM).UpdateOne(ctx, filter, update)
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(result)
}

func DeletePost(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}

	db, ctx := database.ConnectDB()
	defer database.CloseDB()

	filter := bson.D{{Key: "_id", Value: id}}
	result, err := db.Collection(helpers.POST_COLLECTIOM).DeleteOne(ctx, filter)
	if err != nil {
		helpers.ErrorResonse(err.Error(), http.StatusBadRequest, res)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(result)
}
