package seeker

import (
	"strings"
	"errors"
	"fmt"
	"strconv"
)

type RangeHeader struct {
	Unit   string
	Ranges []*Range
}

type Range struct {
	Start int64
	End   int64
}

var (
	ErrInvalidRangeHeader = errors.New("invalid range header")
	ErrEndSmallerThanStart = errors.New("invalid range: end is smaller than start")
)

func errInvalidRange(str string) error {
	return errors.New(fmt.Sprintf(`invalid range: "%s"`, str))
}

func ParseRangeHeader(header string) (*RangeHeader, error) {
	spl := strings.SplitN(header, "=", 2)
	if len(spl) < 2 {
		return nil, ErrInvalidRangeHeader
	}

	rh := &RangeHeader{
		Unit: spl[0],
	}

	// parse ranges
	rangeStrings := strings.Split(spl[1], ",")
	for _, rangeString := range rangeStrings {
		rangeString = strings.TrimSpace(rangeString)
		if rangeString == "" {
			return nil, errInvalidRange(rangeString)
		}
		values := strings.SplitN(rangeString, "-", 2)
		if len(values) < 2 {
			return nil, errInvalidRange(rangeString)
		}
		start, err := strconv.ParseInt(values[0], 10, 0)
		if err != nil {
			return nil, errInvalidRange(rangeString)
		}
		var end int64
		if values[1] == "" {
			end = -1
		} else {
			end, err = strconv.ParseInt(values[1], 10, 0)
			if err != nil {
				return nil, errInvalidRange(rangeString)
			}
			if end < start {
				return nil, ErrEndSmallerThanStart
			}
		}

		r := &Range{
			Start: start,
			End:   end,
		}
		rh.Ranges = append(rh.Ranges, r)
	}

	return rh, nil
}
