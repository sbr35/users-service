package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sbr35/wallets-users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (handler *Login) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	var user models.LoginParams
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		return NewHTTPError(err, "Invalid Login Payload", 400)
	}
	var result models.User

	err = handler.collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&result)

	if err != nil {
		return NewHTTPError(err, "Invaid Email", 400)
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		return NewHTTPError(err, "Wrong Password", 400)
	}
	fmt.Printf("%+v", result)
	userid := result.ID
	useridstr := userid.Hex()
	fmt.Printf("%v", useridstr)
	token, err := TokenCreator(useridstr)

	if err != nil {
		return NewHTTPError(err, "Error while generation token, Try again", 500)
	}

	var loginResponse models.LoginResponse
	loginResponse.Token = token
	loginResponse.ID = result.ID
	json.NewEncoder(w).Encode(loginResponse)
	return nil
}
