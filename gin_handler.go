package seeker

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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

var (
	ErrRangeUnitUnsupported      = errors.New(`seeker only supports range unit "bytes"`)
	ErrMultipleRangesUnsupported = errors.New("seeker only supports single ranges")
	ErrRangeOutOfBounds          = errors.New("range is out of bounds")
)

func SendSeekable(c *gin.Context, in io.ReadSeeker, length int64) error {
	// indicate this resource can be partially requested
	c.Header(acceptRangesHeader, unitBytes)

	if r := c.GetHeader(rangeHeader); r != "" {
		// parse range header if set
		parsed, err := ParseRangeHeader(r)
		if err != nil {
			return err
		}

		if parsed.Unit != unitBytes {
			return ErrRangeUnitUnsupported
		}
		if len(parsed.Ranges) != 1 {
			return ErrMultipleRangesUnsupported
		}

		rng := parsed.Ranges[0]
		if rng.End == -1 {
			rng.End = length - 1
		}
		if rng.End >= length {
			return ErrRangeOutOfBounds
		}

		rangeLen := rng.End - rng.Start + 1

		c.Status(http.StatusPartialContent)
		c.Header(contentLengthHeader, strconv.FormatInt(rangeLen, 10))
		c.Header(contentRangeHeader, fmt.Sprintf("%s %d-%d/%d", unitBytes, rng.Start, rng.End, length))

		// seek resource to start of range
		_, err = in.Seek(rng.Start, 0)
		if err != nil {
			return err
		}
		// write range to response,
		// ignoring any errors, as they are most likely
		// caused by the client closing the pipe.
		io.CopyN(c.Writer, in, rangeLen)
	} else {
		c.Status(http.StatusOK)
		c.Header(contentLengthHeader, strconv.FormatInt(length, 10))

		// write entire resource to response,
		// ignoring any errors, as they are most likely
		// caused by the client closing the pipe.
		io.Copy(c.Writer, in)
	}

	return nil
}

func SendSeekableFile(c *gin.Context, file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	return SendSeekable(c, file, stat.Size())
}
