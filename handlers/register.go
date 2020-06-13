package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/sbr35/wallets-users/db"
	"github.com/sbr35/wallets-users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	logger *log.Logger
}

func NewUserHandler(logger *log.Logger) *UserHandler {
	return &UserHandler{logger}
}

func (reg *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		reg.addUser(w, r)
		return
	}
	if r.Method == http.MethodGet {
		reg.getUser(w, r)
		return
	}
}

func (userhandler *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()
	var res models.Response
	offset, err := strconv.ParseInt(queryParams.Get("offset"), 10, 64)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	limit, err := strconv.ParseInt(queryParams.Get("limit"), 10, 64)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)

	collection, err := db.UsersCollection()

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var users models.Users
	for cur.Next(context.TODO()) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		users = append(users, &user)
	}
	if err := cur.Err(); err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	cur.Close(context.TODO())
	json.NewEncoder(w).Encode(users)
}

func (reg *UserHandler) addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res models.Response

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	collection, err := db.UsersCollection()

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var result models.User
	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

			if err != nil {
				res.Error = "Error while Hashing Password, Try Again"
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(res)
				return
			}

			user.Password = string(hash)
			_, err = collection.InsertOne(context.TODO(), user)

			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(res)
				return
			}

			res.Result = "Registration Successful"
			json.NewEncoder(w).Encode(res)
			return
		}
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	reg.logger.Println(err)
	res.Result = "Username already Exists!!"
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(res)
	return
}
