package handlers

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	logger     *log.Logger
	collection *mongo.Collection
}

func NewUserHandler(logger *log.Logger, collection *mongo.Collection) *UserHandler {
	return &UserHandler{logger, collection}
}

func (handler *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err = error(nil)
	if r.Method == http.MethodPost {
		err = handler.SignUp(w, r)
	}
	if r.Method == http.MethodPut {
		err = handler.UpdateUser(w, r)
	}

	if r.Method == http.MethodGet {
		err = handler.GetInfo(w, r)
	}

	if r.Method == http.MethodDelete {
		err = handler.DeleteUser(w, r)
	}

	if err == nil {
		w.WriteHeader(200)
		return
	}
	handler.logger.Printf("Error Occured: %v", err)
	clientError, ok := err.(CustomError)
	if !ok {
		w.WriteHeader(500)
		return
	}
	body, err := clientError.RespBody()
	if err != nil {
		handler.logger.Printf("An Error Occured: %v", err)
		w.WriteHeader(500)
		return
	}
	status, headers := clientError.RespHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write(body)
}
