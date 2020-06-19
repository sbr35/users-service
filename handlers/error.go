package handlers

import (
	"encoding/json"
	"fmt"
)

type CustomError interface {
	Error() string
	RespBody() ([]byte, error)
	RespHeaders() (int, map[string]string)
}

type HTTPError struct {
	MainError error  `json:"-"`
	Detail    string `json:"detail"`
	Status    int    `json:"-"`
}

func (he *HTTPError) Error() string {
	if he.MainError == nil {
		return he.Detail
	}
	return he.Detail + " : " + he.MainError.Error()
}

func (he *HTTPError) RespBody() ([]byte, error) {
	body, err := json.Marshal(he)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}
	return body, nil
}

func (he *HTTPError) RespHeaders() (int, map[string]string) {
	return he.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

func NewHTTPError(err error, detail string, status int) error {
	return &HTTPError{
		MainError: err,
		Detail:    detail,
		Status:    status,
	}
}
