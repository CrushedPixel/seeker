package seeker

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	acceptRangesHeader = "Accept-Ranges"
	unitBytes          = "bytes"

	rangeHeader         = "Range"
	contentLengthHeader = "Content-Length"
	contentRangeHeader  = "Content-Range"
)

// send writes a resource to the response, supporting partial content requests.
// If an error is returned, the response headers have not been written yet.
func send(w http.ResponseWriter, req *http.Request, in io.ReadSeeker, length int64) *errorResponse {
	// indicate this resource can be partially requested
	req.Header.Set(acceptRangesHeader, unitBytes)

	if r := req.Header.Get(rangeHeader); r != "" {
		// parse range header if set
		parsed, err := ParseRangeHeader(r)
		if err != nil {
			return err.(*errorResponse)
		}

		if parsed.Unit != unitBytes {
			return ErrUnsupportedRangeUnit
		}
		if len(parsed.Ranges) != 1 {
			return ErrMultipleRanges
		}

		rng := parsed.Ranges[0]
		if rng.End == -1 {
			rng.End = length - 1
		}
		if rng.End >= length {
			return ErrRangeOutOfBounds
		}

		// seek resource to start of range
		if _, err := in.Seek(rng.Start, 0); err != nil {
			return ErrIO
		}

		// calculate content length of range
		rangeLen := rng.End - rng.Start + 1
		w.Header().Set(contentLengthHeader, strconv.FormatInt(rangeLen, 10))
		w.Header().Set(contentRangeHeader, fmt.Sprintf("%s %d-%d/%d", unitBytes, rng.Start, rng.End, length))
		w.WriteHeader(http.StatusPartialContent)

		// write range to response,
		// ignoring any errors, as they are most likely
		// caused by the client closing the pipe.
		io.CopyN(w, in, rangeLen)
	} else {
		w.Header().Set(contentLengthHeader, strconv.FormatInt(length, 10))

		// write entire resource to response,
		// ignoring any errors, as they are most likely
		// caused by the client closing the pipe.
		io.Copy(w, in)
	}

	return nil
}

func sendFile(w http.ResponseWriter, req *http.Request, file *os.File) *errorResponse {
	stat, err := file.Stat()
	if err != nil {
		return ErrIO
	}

	return send(w, req, file, stat.Size())
}
