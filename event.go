package sse

import (
	"io"
)

// Event represents a single Event in a server-sent events stream.
// It MAY contain an ID or an event Name.
type Event struct {
	ID   string
	Name string
}

func (e *Event) Read(b []byte) (int, error) {
	return 0, io.EOF
}

// ProtocolError represents SSE errors.
type ProtocolError struct {
	ErrorString string
}

func (pe *ProtocolError) Error() string { return pe.ErrorString }
