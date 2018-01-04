package seeker

import (
	"os"
	"github.com/gin-gonic/gin"
	"net/http"
)

// implements margo.Response
type SeekableFileResponse struct {
	file *os.File
}

func NewMargoResponse(file *os.File) *SeekableFileResponse {
	return &SeekableFileResponse{
		file: file,
	}
}

func (s *SeekableFileResponse) Send(c *gin.Context) error {
	err := SendSeekableFile(c, s.file)
	if err != nil {
		// check if error is user input error
		if err == ErrInvalidRangeHeader ||
			err == ErrEndSmallerThanStart ||
			err == ErrRangeUnitUnsupported ||
			err == ErrMultipleRangesUnsupported ||
			err == ErrRangeOutOfBounds {
			// set bad request status
			c.Status(http.StatusBadRequest)
			return nil
		}

		// otherwise, let margo error handler decide
		// what to do with the error
		return err
	}

	return nil
}
