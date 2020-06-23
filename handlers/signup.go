package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sbr35/wallets-users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (handler *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	var user models.RegistrationParams
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res models.RegistrationResponse

	if err != nil {
		return NewHTTPError(err, "Error in body parsing", 400)
	}

	var result models.User
	err = handler.collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

			if err != nil {
				return NewHTTPError(err, "Error while hashing the password", 500)
			}

			user.Password = string(hash)
			resp, err := handler.collection.InsertOne(context.TODO(), user)
			if err != nil {
				return NewHTTPError(err, "Error in Saving user information", 400)
			}

			oid, ok := resp.InsertedID.(primitive.ObjectID)
			if !ok {
				return NewHTTPError(nil, "User created but failed to get the user id", 404)
			}
			res.ID = oid
			res.Result = "Registration Successful"
			json.NewEncoder(w).Encode(res)
			return nil
		}
		return NewHTTPError(nil, "Error occured", 500)
	}
	return NewHTTPError(nil, "User already exists!!", 400)
}
