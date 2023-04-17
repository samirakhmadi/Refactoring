package service

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	UserNotFound = errors.New("user_not_found")
)

type errResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func errInvalidRequest(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func errInternalError(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal error.",
		ErrorText:      err.Error(),
	}
}

func writeCors(resp http.ResponseWriter) {
	resp.Header().Set("Access-Control-Allow-Credentials", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Origin", "Accept, Content-Type, Content-Length, Accept-Encoding")
}

func writeResponse(resp http.ResponseWriter, msg []byte, statusCode int) {
	resp.Write(msg)
	resp.WriteHeader(statusCode)
}
