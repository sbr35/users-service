package handlers

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Login struct {
	logger     *log.Logger
	collection *mongo.Collection
}

func NewLogin(logger *log.Logger, collection *mongo.Collection) *Login {
	return &Login{logger, collection}
}

func (handler *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err = error(nil)
	if r.Method == http.MethodPost {
		err = handler.LoginHandler(w, r)
	}
	if err == nil {
		w.WriteHeader(200)
		return
	}
	handler.logger.Printf("Error Occured: %v", err)
	clintError, ok := err.(CustomError)
	if !ok {
		w.WriteHeader(500)
		return
	}
	body, err := clintError.RespBody()
	if err != nil {
		handler.logger.Printf("Error Occured: %v", err)
		w.WriteHeader(500)
		return
	}
	status, headers := clintError.RespHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write(body)
}
