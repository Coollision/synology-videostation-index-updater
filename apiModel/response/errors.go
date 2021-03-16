package response

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse >>>
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string   `json:"status"`           // user-level status message
	ErrorText  string   `json:"error,omitempty"`  // application-level error message, for debugging
	ErrorList  []string `json:"errors,omitempty"` // application-level error messages, for debugging
}

// Error >>>
func (e ErrResponse) Error() string {
	return fmt.Sprintf("Got StatusCode: %d, with : %s; %s;", e.HTTPStatusCode, e.StatusText, e.ErrorText)
}

// Render >>>
func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// IsErrResponse >>>
func IsErrResponse(err error) (ErrResponse, bool) {
	eResp, ok := err.(ErrResponse)
	return eResp, ok
}

// ErrorToRenderer checks if a response is an Errresponse, if it is not it will convert the error to an internalServerError
func ErrorToRenderer(err error) render.Renderer {
	if errResp, ok := IsErrResponse(err); ok {
		return errResp
	}
	return ErrInternalServer(err)
}

// ErrInvalidRequest >>>
func ErrInvalidRequest(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrInvalidRequestN >>>
func ErrInvalidRequestN(errors []string) ErrResponse {
	return ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Bad request.",
		ErrorList:      errors,
	}
}

// ErrUnauthorized >>>
func ErrUnauthorized(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized.",
		ErrorText:      err.Error(),
	}
}

// ErrForbidden >>>
func ErrForbidden() ErrResponse {
	return ErrResponse{
		Err:            nil,
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     "Forbidden.",
		ErrorText:      "You have insufficient acces rights, please log in as the correct user",
	}
}

// ErrNotFound >>>
func ErrNotFound(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Resource not found.",
		ErrorText:      err.Error(),
	}
}

// ErrInternalServer >>>
func ErrInternalServer(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Something went wrong.",
		ErrorText:      err.Error(),
	}
}

// ErrNoContent >>>
func ErrNoContent() ErrResponse {
	return ErrResponse{
		HTTPStatusCode: http.StatusNoContent,
		StatusText:     "No Content To Display.", // will not be returned => a 204 makes queriers ignore the response body but lets put it in there anyway
	}
}

// ErrBadRequest >>>
func ErrBadRequest(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Bad Request",
		ErrorText:      err.Error(),
	}
}

// ErrNotImplemented >>>
func ErrNotImplemented() ErrResponse {
	return ErrResponse{
		Err:            nil,
		HTTPStatusCode: http.StatusNotImplemented,
		StatusText:     "The Dev is still working on it.",
		ErrorText:      "",
	}
}

// ErrTeapot >>>
func ErrTeapot() ErrResponse {
	err := fmt.Errorf("i'm a teapot, I refuse to brew coffie")
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusTeapot,
		StatusText:     "I am a teapot.",
		ErrorText:      err.Error(),
	}
}
