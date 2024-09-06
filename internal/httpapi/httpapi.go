package httpapi

import (
	"encoding/json"
	"net/http"
)

const (
	RESPONSE_CODE__OK                = 0
	RESPONSE_CODE__UKNOWN            = 1
	RESPONSE_CODE__DB                = 2
	RESPONSE_CODE__UNAUTHORIZED      = 3
	RESPONSE_CODE__PERMISSION_DENIED = 4

	RESPONSE_CODE__INVALID_REQUEST = 10
	RESPONSE_CODE__NOT_FOUND       = 11
)

type JsonRepsonse struct {
	Status  int         `json:"-"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ListBody struct {
	Total    int64       `json:"total"`
	Records  interface{} `json:"records"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func (r *JsonRepsonse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	_ = json.NewEncoder(w).Encode(r)
}

func NewBadRequestResponse() *JsonRepsonse {
	return &JsonRepsonse{
		Status:  http.StatusBadRequest,
		Code:    RESPONSE_CODE__INVALID_REQUEST,
		Message: "bad request",
	}
}

func NewNotFoundErrorResponse() *JsonRepsonse {
	return &JsonRepsonse{
		Status:  http.StatusNotFound,
		Code:    RESPONSE_CODE__NOT_FOUND,
		Message: "not found",
	}
}

func NewUnknownErrorResponse() *JsonRepsonse {
	return &JsonRepsonse{
		Status:  http.StatusInternalServerError,
		Code:    RESPONSE_CODE__UKNOWN,
		Message: "unknown error",
	}
}

func NewInvalidResponse(msg string) *JsonRepsonse {
	return &JsonRepsonse{
		Status:  http.StatusBadRequest,
		Code:    RESPONSE_CODE__INVALID_REQUEST,
		Message: msg,
	}
}

// NewDataResponse creates a new JSON response object with the provided data.
//
// It takes data of type interface{} as a parameter and returns a pointer to JsonRepsonse.
func NewDataResponse(data interface{}, msg string) *JsonRepsonse {
	return &JsonRepsonse{
		Status:  http.StatusOK,
		Code:    RESPONSE_CODE__OK,
		Data:    data,
		Message: msg,
	}
}
