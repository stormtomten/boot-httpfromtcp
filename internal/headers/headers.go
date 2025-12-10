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

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)

	v, exists := h[key]
	if exists {
		value = strings.Join([]string{v, value}, ", ")
	}
	h[key] = value
}

func (h Headers) Override(key, value string) {
	delete(h, key)
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	val, ok := h[key]
	return val, ok
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
