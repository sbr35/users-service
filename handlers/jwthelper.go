package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sbr35/users-service/models"
	"github.com/twinj/uuid"
)

func TokenCreator(userid string) (*models.LoginToken, error) {
	token := &models.LoginToken{}
	token.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	token.AccessUuid = uuid.NewV4().String()
	var err error

	accessTokenClaim := jwt.MapClaims{}
	accessTokenClaim["access_uuid"] = token.AccessUuid
	accessTokenClaim["userid"] = userid
	accessTokenClaim["expires"] = token.AtExpires

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)
	fmt.Println("Access Secret Token", os.Getenv("ACCESS_TOKEN_SECRET"))
	token.AccessToken, err = newAccessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractToken(r *http.Request) string {
	theToken := r.Header.Get("Authorization")
	return theToken
	// fmt.Println(theToken)
	// splitToken := strings.Split(theToken, " ")
	// if len(splitToken) == 2 {
	// 	return splitToken[1]
	// }
	// return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	getToken := ExtractToken(r)
	token, err := jwt.Parse(getToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method confirm to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) (jwt.MapClaims, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
