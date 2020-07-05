package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sbr35/wallets-users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (handler *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	claims, err := TokenValid(r)
	if err != nil {
		return NewHTTPError(err, "Authorization Failed", 400)
	}

	var params models.UpdateParams
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NewHTTPError(err, "Error in reading request body", 400)
	}
	err = json.Unmarshal(body, &params)
	if err != nil {
		return NewHTTPError(err, "Error in parsing request body", 400)
	}
	if int64(claims["expires"].(float64)) < time.Now().Unix() {
		return NewHTTPError(nil, "Session Timeout. Log In again", 400)
	}
	userid, err := primitive.ObjectIDFromHex(claims["userid"].(string))
	updateParams := bson.M{}
	if params.FirstName != "" {
		updateParams["firstname"] = params.FirstName
	}
	if params.LastName != "" {
		updateParams["lastname"] = params.LastName
	}
	if params.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 5)
		if err != nil {
			return NewHTTPError(err, "Error in new password hashing", 500)
		}
		updateParams["password"] = string(hash)
	}

	updateBson := bson.M{
		"$set": updateParams,
	}

	_, err = handler.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": userid},
		updateBson,
	)
	if err != nil {
		return NewHTTPError(err, "Error in updating user. Try Again", 500)
	}
	var resp models.UpdateResponse
	resp.Result = "Update Successful"
	json.NewEncoder(w).Encode(resp)
	return nil
}
