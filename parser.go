package seeker

import (
	"strconv"
	"strings"
)

type RangeHeader struct {
	Unit   string
	Ranges []*Range
}

type Range struct {
	Start int64
	End   int64
}

// ParseRangeHeader parses a Content-Range header value.
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
			return nil, ErrInvalidRangeHeader
		}
		values := strings.SplitN(rangeString, "-", 2)
		if len(values) < 2 {
			return nil, ErrInvalidRangeHeader
		}
		start, err := strconv.ParseInt(values[0], 10, 0)
		if err != nil {
			return nil, ErrInvalidRangeHeader
		}
		var end int64
		if values[1] == "" {
			end = -1
		} else {
			end, err = strconv.ParseInt(values[1], 10, 0)
			if err != nil {
				return nil, ErrInvalidRangeHeader
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
