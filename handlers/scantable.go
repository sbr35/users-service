package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sbr35/wallets-users/db"
	"github.com/sbr35/wallets-users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (userhandler *UserHandler) ScanTable(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()
	offset, err := strconv.ParseInt(queryParams.Get("offset"), 10, 64)
	if err != nil {
		return NewHTTPError(err, "provide offset value as query parameter", 400)
	}
	limit, err := strconv.ParseInt(queryParams.Get("limit"), 10, 64)
	if err != nil {
		return NewHTTPError(err, "provide limit value as query parameter", 400)
	}
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)

	collection, err := db.UsersCollection()

	if err != nil {
		return NewHTTPError(err, "Database is not responding", 500)
	}

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return NewHTTPError(err, "Database query error", 500)
	}

	var users models.Users
	for cur.Next(context.TODO()) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return NewHTTPError(err, "Error in listing user", 400)
		}
		users = append(users, &user)
	}
	if err := cur.Err(); err != nil {
		return NewHTTPError(err, "Database Error occured", 400)
	}
	cur.Close(context.TODO())
	json.NewEncoder(w).Encode(users)
	return nil
}
