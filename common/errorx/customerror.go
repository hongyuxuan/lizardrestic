package errorx

import (
	"fmt"
	"net/http"
)

type LizardresticError struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HttpErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(code int, message string, data interface{}) error {
	return &LizardresticError{Code: code, Message: message, Data: data}
}

func NewDefaultError(message string, a ...any) error {
	return &LizardresticError{Code: http.StatusInternalServerError, Message: fmt.Sprintf(message, a...)}
}

func (e *LizardresticError) Error() string {
	return e.Message
}

func (e *LizardresticError) GetData() *HttpErrorResponse {
	return &HttpErrorResponse{
		Code:    e.Code,
		Message: e.Message,
		Data:    e.Data,
	}
}
