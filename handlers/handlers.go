package handlers

import (
	"time"
	"log"
	"os"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/sbr35/wallets/services/users/db"
	"github.com/sbr35/wallets/services/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/twinj/uuid"
	jwt "github.com/dgrijalva/jwt-go"

)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)

	if(err != nil){
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
	res.Result = "Username already Exists!!"
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(res)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user);

	if(err != nil){
		log.Fatal(err)
	}

	collection, err := db.UsersCollection();

	if err != nil {
		log.Fatal(err)
	}

	var result models.User
	var res models.Response

	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result);
	
	if err != nil {
		res.Error = "Invalid Email"
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password));

	if err != nil {
		res.Error = "Invalid Password"
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(res)
		return
	}

	token, err := TokenCreator(result.Email)

	if err != nil {
		res.Error = "Error while generation token, Try again"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	var loginResponse models.LoginResponse
	loginResponse.Token = token
	json.NewEncoder(w).Encode(loginResponse)
	return
}

func TokenCreator(user_email string) (*models.LoginToken, error){
	token := &models.LoginToken{}
	token.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	token.AccessUuid = uuid.NewV4().String()
	token.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	token.RefreshUuid = uuid.NewV4().String()
	var err error

	accessTokenClaim := jwt.MapClaims{}
	accessTokenClaim["access_uuid"] = token.AccessUuid
	accessTokenClaim["user_email"] =  user_email
	accessTokenClaim["expires"] = token.AtExpires

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)
	token.AccessToken, err = newAccessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}


	refreshTokenClaim := jwt.MapClaims{}
	refreshTokenClaim["refresh_uuid"] = token.RefreshUuid
	refreshTokenClaim["user_email"] = user_email
	refreshTokenClaim["expires"] = token.RtExpires

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaim)
	token.RefreshToken, err = newRefreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))

	if err != nil {
		return nil, err
	}

	return token, nil
}
