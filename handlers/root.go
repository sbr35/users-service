package handlers

import (
	"log"
	"net/http"
)

type UserHandler struct {
	logger *log.Logger
}

func NewUserHandler(logger *log.Logger) *UserHandler {
	return &UserHandler{logger}
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
		err = handler.ScanTable(w, r)
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
