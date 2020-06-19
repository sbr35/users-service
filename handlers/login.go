package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sbr35/wallets-users/db"
	"github.com/sbr35/wallets-users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	logger *log.Logger
}

func NewLogin(logger *log.Logger) *Login {
	return &Login{logger}
}

func (login *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.LoginParams
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)

	if err != nil {
		login.logger.Fatal(err)
	}
	login.logger.Printf("%v", user)
	collection, err := db.UsersCollection()

	if err != nil {
		login.logger.Fatal(err)
	}

	var result models.User
	var res models.Response

	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&result)

	if err != nil {
		res.Error = "Invalid Email"
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid Password"
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(res)
		return
	}
	fmt.Printf("%+v", result)
	userid := result.ID
	useridstr := userid.Hex()
	fmt.Printf("%v", useridstr)
	token, err := TokenCreator(useridstr)

	if err != nil {
		res.Error = "Error while generation token, Try again"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	var loginResponse models.LoginResponse
	loginResponse.Token = token
	loginResponse.ID = result.ID
	json.NewEncoder(w).Encode(loginResponse)
	return
}
