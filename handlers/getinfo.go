package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sbr35/wallets-users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (handler *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	claims, err := TokenValid(r)
	if err != nil {
		return NewHTTPError(err, "Authorization Failed", 400)
	}
	userid := claims["userid"].(string)
	userObjectID, err := primitive.ObjectIDFromHex(userid)
	fmt.Println(userObjectID)
	if err != nil {
		return NewHTTPError(err, "User Id parsing error. Authorization Issue", 400)
	}

	// collection, err := db.UsersCollection()
	// if err != nil {
	// 	return NewHTTPError(err, "Database Connection Error", 500)
	// }
	var user models.User
	err = handler.collection.FindOne(context.Background(), bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		return NewHTTPError(err, "Information didn't found", 404)
	}
	json.NewEncoder(w).Encode(user)
	return nil
}
