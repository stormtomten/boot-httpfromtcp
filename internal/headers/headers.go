package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	} else if idx == 0 {
		return 2, true, nil
	}

	line := string(data[:idx])
	idc := strings.Index(line, ":")
	if idc == -1 {
		return 0, false, fmt.Errorf("headers: invalid field-name: %s", line)
	}
	if strings.HasSuffix(line[:idc], " ") {
		return 0, false, fmt.Errorf("headers: invalid field-name: %s", line)
	}

	name := strings.TrimSpace(line[:idc])
	for i := 0; i < len(name); i++ {
		if !isValidHeaderChar(name[i]) {
			return 0, false, fmt.Errorf("headers: invalid token found in field-name: %s", line)
		}
	}

	value := strings.TrimSpace(line[idc+1:])

	h.Set(name, value)
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)

	v, exists := h[key]
	if exists {
		value = strings.Join([]string{v, value}, ", ")
	}
	h[key] = value
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func isValidHeaderChar(b byte) bool {
	switch {
	case '0' <= b && b <= '9':
		return true
	case 'A' <= b && b <= 'Z':
		return true
	case 'a' <= b && b <= 'z':
		return true
	case slices.Contains(tokenChars, b):
		return true
	}

	return false
}
