package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	User  primitive.ObjectID `json:"user" validate:"required" bson:",omitempty"`
	Title string             `json:"title" validate:"required" bson:",omitempty"`
	Body  string             `json:"body" bson:",omitempty"`
}
