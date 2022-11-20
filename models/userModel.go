package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" validate:"email,required" bson:",omitempty"`
	FirstName string             `json:"firstName" validate:"required" bson:",omitempty"`
	LastName  string             `json:"lastName" validate:"required" bson:",omitempty"`
	Password  string             `json:"password,omitempty" validate:"required" bson:",omitempty"`
}

type SanitizedUser struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email"`
	FirstName string             `json:"firstName" `
	LastName  string             `json:"lastName"`
}
