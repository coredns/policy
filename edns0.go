package policy

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
	typeEDNS0Bytes = iota
	typeEDNS0Hex
	typeEDNS0IP
)

var edns0Types = map[string]uint16{
	"bytes":   typeEDNS0Bytes,
	"hex":     typeEDNS0Hex,
	"address": typeEDNS0IP,
}

type edns0Opt struct {
	name     string
	dataType uint16
	size     int
	start    int
	end      int
	metrics  bool
}

func newEdns0Opt(sCode, name, sType, sSize, sStart, sEnd string) (uint16, *edns0Opt, error) {
	code, err := strconv.ParseUint(sCode, 0, 16)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not parse EDNS0 code: %s", err)
	}

	dataType, ok := edns0Types[strings.ToLower(sType)]
	if !ok {
		return 0, nil, fmt.Errorf("Unknown EDNS0 data type %q", sType)
	}

	size, err := strconv.ParseInt(sSize, 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not parse EDNS0 data size: %s", err)
	}

	start, err := strconv.ParseInt(sStart, 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not parse EDNS0 start index: %s", err)
	}

	end, err := strconv.ParseInt(sEnd, 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not parse EDNS0 end index: %s", err)
	}

	if start > 0 && end > 0 && end <= start {
		return 0, nil, fmt.Errorf("End index should be > start index (actual %d <= %d)", end, start)
	}

	if size > 0 {
		if start > 0 && start >= size {
			return 0, nil, fmt.Errorf("Start index should be < size (actual %d >= %d)", start, size)
		}

		if end > 0 && end > size {
			return 0, nil, fmt.Errorf("End index should be <= size (actual %d > %d)", end, size)
		}
	}

	return uint16(code), &edns0Opt{
		name:     name,
		dataType: dataType,
		size:     int(size),
		start:    int(start),
		end:      int(end),
	}, nil
}

func (o *edns0Opt) makeHexString(b []byte) string {
	if o.size > 0 && o.size != len(b) {
		return ""
	}

	start := 0
	if o.start > 0 {
		if o.start >= len(b) {
			return ""
		}

		start = o.start
	}

	end := len(b)
	if o.end > 0 {
		if o.end > len(b) {
			return ""
		}

		end = o.end
	}

	return hex.EncodeToString(b[start:end])
}
