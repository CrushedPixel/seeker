package seeker

import (
	"io"
	"net/http"
	"os"
)

// Send handles a request to access a resource, supporting partial content requests.
func Send(w http.ResponseWriter, req *http.Request, in io.ReadSeeker, length int64) {
	if err := send(w, req, in, length); err != nil {
		err.Send(w)
	}
}

// Send handles a request to access a file, supporting partial content requests.
func SendFile(w http.ResponseWriter, req *http.Request, file *os.File) {
	if err := sendFile(w, req, file); err != nil {
		err.Send(w)
	}
}
