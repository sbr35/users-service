package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Response struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type LoginToken struct {
	AccessToken string `json:"accesstoken"`
	AccessUuid  string `json:"accessuuid"`
	AtExpires   int64  `json:"atexpires"`
}

type LoginResponse struct {
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Token *LoginToken        `json:"token"`
}

type RegistrationResponse struct {
	ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Result string             `json:"result"`
}

type UpdateResponse struct {
	Result string `json:"result"`
}

type RegistrationParams struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateParams struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Email     string             `json:"email"`
	FirstName string             `json:"firstname"`
	LastName  string             `json:"lastname"`
	Password  string             `json:"-"`
}

type Users []*User
