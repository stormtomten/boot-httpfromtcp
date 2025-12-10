package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid header no colon/field-name
	headers = NewHeaders()
	data = []byte("       Host  localhost42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Single header extra whitespace
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 37, n)
	assert.False(t, done)

	// Test: Field without value
	headers = NewHeaders()
	data = []byte("     Keep-alive:             \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "", headers["keep-alive"])
	assert.Equal(t, 31, n)
	assert.False(t, done)

	// Test: Done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.True(t, done)

	// Test: Valid 2 headers with existing headers + done
	headers = NewHeaders()
	data = []byte(" Something: or other\r\n  Else: there is this\r\n\r\n")
	n, done, err = headers.Parse(data)
	sum := n
	n2, done, err := headers.Parse(data[n:])
	sum += n2
	n3, done, err := headers.Parse(data[sum:])
	sum += n3
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "or other", headers["something"])
	assert.Equal(t, "there is this", headers["else"])
	assert.Equal(t, 47, sum)
	assert.True(t, done)

	// Test: Data without CRLF
	headers = NewHeaders()
	data = []byte("     Some: Data    ")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Empty header
	headers = NewHeaders()
	data = []byte("Valid-name:\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 13, n)
	assert.Equal(t, "", headers["valid-name"])
	assert.False(t, done)

	// Test: Invalid char in field-name
	headers = NewHeaders()
	data = []byte("Inv@lid-name:\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Append to existing key
	headers = NewHeaders()
	data = []byte("Users: Murvel\r\n   Users:Ylle\r\n  Users:  Harald\r\n\r\n")
	sum = 0
	n, done, err = headers.Parse(data[sum:])
	require.NoError(t, err)
	require.Equal(t, "Murvel", headers["users"])
	require.Equal(t, 15, n)
	sum += n
	n, done, err = headers.Parse(data[sum:])
	require.NoError(t, err)
	require.Equal(t, "Murvel, Ylle", headers["users"])
	require.Equal(t, 15, n)
	sum += n
	require.Equal(t, 30, sum)
	n, done, err = headers.Parse(data[sum:])
	require.NoError(t, err)
	require.Equal(t, "Murvel, Ylle, Harald", headers["users"])
	require.Equal(t, 18, n)
	sum += n
	require.Equal(t, 48, sum)
	n, done, err = headers.Parse(data[sum:])
	require.NoError(t, err)
	require.Equal(t, 2, n)
	sum += n
	require.Equal(t, 50, sum)
	require.True(t, done)
}
