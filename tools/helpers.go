package tools

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMessage interface {
	Error() string
	GetCode() int
	GetAll() errorMessage
}

type errorMessage struct {
	Err       string `json:"err"`
	HumanText string `json:"humanText"`
	Code      int    `json:"code"`
}

func (e *errorMessage) Error() string {
	return e.Err
}

func (e *errorMessage) GetCode() int {
	return e.Code
}

func (e *errorMessage) GetAll() errorMessage {
	return *e
}

func EncodeIntoResponseWriter(w http.ResponseWriter, message ErrorMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(message.GetCode())
	err := json.NewEncoder(w).Encode(message.GetAll())
	if err != nil {
		log.Printf("error while encoding error to writer: %v", err)
	}
}

func NewErrorMessage(err error, humanReadable string, code int) ErrorMessage {
	return &errorMessage{
		Err:       err.Error(),
		HumanText: humanReadable,
		Code:      code,
	}
}

func NewErrorMessageEncodeIntoWriter(err error, humanReadable string, code int) ErrorMessage {
	if err == nil {
		return nil
	}
	return &errorMessage{
		Err:       err.Error(),
		HumanText: humanReadable,
		Code:      code,
	}
}

func DecodeNewErrorMessage(resp *http.Response) ErrorMessage {
	if resp == nil {
		return nil
	}
	var response errorMessage
	if resp.StatusCode != 200 {
		err := json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			if err != nil {
				return &errorMessage{
					Err: err.Error(),
				}
			}
			return &errorMessage{
				Err: buf.String(),
			}
		}
		return &response
	}
	return nil
}
