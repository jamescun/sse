package sse

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaddedCopyBuffer(t *testing.T) {
	tests := []struct {
		Name           string
		Prefix, Suffix string
		BufferSize     int
		Source         []byte
		Written        int64
		Expected       []byte
		Error          error
	}{
		{"Overflow", "foo", "bar\n", 8, []byte("BAZ"), 24, []byte("fooBbar\nfooAbar\nfooZbar\n"), nil},
		{"SSE", "data: ", "\n", 32 * 1024, []byte("foo"), 10, []byte("data: foo\n"), nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var dst bytes.Buffer
			src := bytes.NewReader(test.Source)

			n, err := paddedCopyBuffer(&dst, src, test.Prefix, test.Suffix, make([]byte, test.BufferSize))
			if test.Error == nil {
				if assert.NoError(t, err) {
					assert.Equal(t, test.Written, n)
					assert.Equal(t, test.Expected, dst.Bytes())
				}
			} else {
				assert.Equal(t, test.Error, err)
			}
		})
	}
}
