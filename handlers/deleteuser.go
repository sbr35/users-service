package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sbr35/users-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (handler *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	claims, err := TokenValid(r)
	if err != nil {
		return NewHTTPError(err, "Authorizaton Failed", 400)
	}
	if int64(claims["expires"].(float64)) < time.Now().Unix() {
		return NewHTTPError(nil, "Session Timeout. Log in again", 400)
	}
	userid := claims["userid"].(string)
	userObjectID, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return NewHTTPError(err, "User Id parsing error. Authorization Issue", 400)
	}
	_, err = handler.collection.DeleteOne(context.Background(), bson.M{"_id": userObjectID})
	if err != nil {
		return NewHTTPError(err, "Delete Unsuccessful", 400)
	}
	var resp models.Response
	resp.Result = "Delete Operation Successful"
	json.NewEncoder(w).Encode(resp)
	return nil
}
