package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Response struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type LoginToken struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
	AccessUuid   string `json:"accessuuid"`
	RefreshUuid  string `json:"refreshuuid"`
	AtExpires    int64  `json:"atexpires"`
	RtExpires    int64  `json:"rtexpires"`
}

type LoginResponse struct {
	Token *LoginToken `json:"token"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Email     string             `json:"email"`
	FirstName string             `json:"firstname"`
	LastName  string             `json:"lastname"`
	Password  string             `json:"-"`
}

type Users []*User
