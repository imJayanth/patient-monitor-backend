package response

import (
	"patient-monitor-backend/internal/errors"
)

type Response struct {
	Data          interface{}         `json:"data"`
	ResponseError errors.RestAPIError `json:"error"`
}

func NewResponse(data interface{}, error errors.RestAPIError) *Response {
	return &Response{Data: data, ResponseError: error}
}
