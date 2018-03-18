package seeker

import (
	"errors"
	"net/http"
)

// errorResponse is a custom error type
// that can be written to a ResponseWriter.
type errorResponse struct {
	error
	status int
}

func (r *errorResponse) Send(w http.ResponseWriter) {
	w.WriteHeader(r.status)
	w.Write([]byte(r.Error()))
}

func newErrorResponse(status int, text string) *errorResponse {
	return &errorResponse{errors.New(text), status}
}

// ErrInvalidRangeHeader indicates an invalid range header.
var ErrInvalidRangeHeader = newErrorResponse(http.StatusBadRequest, "invalid range header")

// ErrEndSmallerThanStart indicates the end of requested byte range being a smaller value than the start.
var ErrEndSmallerThanStart = newErrorResponse(http.StatusBadRequest, "invalid range: end is smaller than start")

// ErrUnsupportedRangeUnit indicates an unsupported range unit in the Content-Range header.
// The only supported unit is "bytes".
var ErrUnsupportedRangeUnit = newErrorResponse(http.StatusBadRequest, `seeker only supports range unit "bytes"`)

// ErrMultipleRanges indicates multiple byte ranges were requested.
// Seeker only supports single-range requests.
var ErrMultipleRanges = newErrorResponse(http.StatusBadRequest, "seeker only supports single ranges")

// ErrRangeOutOfBounds indicates a byte range that was out of bounds was requested.
var ErrRangeOutOfBounds = newErrorResponse(http.StatusBadRequest, "range is out of bounds")

// ErrIO indicates an error accessing the resource.
var ErrIO = newErrorResponse(http.StatusInternalServerError, "internal server error")
