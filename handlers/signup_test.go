package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sbr35/users-service/db"
	"github.com/sbr35/wallets-users/handlers"
)

func TestSignUp(t *testing.T) {
	var payload = []byte(`{"email":"test@gmail.com","firstname":"shohidul","lastname":"bari","password":"14035"}`)
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	logger := log.New(os.Stdout, "users-api ", log.LstdFlags)
	collection, err := db.UsersCollection()
	if err != nil {
		log.Fatal(err)
	}
	handler := handlers.NewUserHandler(logger, collection)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler Return Code: got %v wat %v", status, http.StatusOK)
	}
}
