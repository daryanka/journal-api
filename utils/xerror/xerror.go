package xerror

import "net/http"

type XerrorT struct {
	StatusCode int `json:"status_code"`
	Message string `json:"message,omitempty"`
	Error bool `json:"error"`
	Type string `json:"type,omitempty"`
}

func NewBadRequest(err string, errType...string) *XerrorT {
	e := &XerrorT{
		StatusCode: http.StatusBadRequest,
		Message: err,
		Error: true,
	}

	if len(errType) > 0 {
		e.Type = errType[0]
	}
	return e
}

func NewInternalError(err string, errType ...string) *XerrorT {
	e := &XerrorT{
		StatusCode: http.StatusInternalServerError,
		Message: err,
		Error: true,
	}

	if len(errType) > 0 {
		e.Type = errType[0]
	}
	return e
}
