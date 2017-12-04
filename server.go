package sse

import (
	"io"
	"net"
	"net/http"
	"os"
	"syscall"
)

// ResponseWriter is an interface used by a HandlerFunc to send
// server-sent events (SSE) to the remove client.
type ResponseWriter interface {
	// WriteEvent sends a given event and data stream to the remote client.
	// If the remote has closed the stream, WriteEvent returns io.EOF.
	WriteEvent(*Event, io.Reader) error
}

type responseWriter struct {
	remote  io.Writer
	flusher http.Flusher
	buf     []byte
}

const (
	msgID    = "id: "
	msgEvent = "event: "
	msgData  = "data: "
	msgEnd   = "\n"
)

func (rw *responseWriter) WriteEvent(e *Event, r io.Reader) error {
	if e == nil && r == nil {
		return nil
	}

	err := rw.writeEvent(e, r)

	if isBrokenPipe(err) {
		return io.EOF
	}

	return err
}

func (rw *responseWriter) writeEvent(e *Event, r io.Reader) error {
	if e != nil {
		if e.ID != "" {
			_, err := writeStrings(rw.remote, msgID, e.ID, msgEnd)
			if err != nil {
				return err
			}
		}

		if e.Name != "" {
			_, err := writeStrings(rw.remote, msgEvent, e.Name, msgEnd)
			if err != nil {
				return err
			}
		}
	}

	if r != nil {
		_, err := paddedCopyBuffer(rw.remote, r, msgData, msgEnd, rw.buf)
		if err != nil {
			return err
		}
	}

	_, err := io.WriteString(rw.remote, msgEnd)
	if err != nil {
		return err
	}

	rw.flusher.Flush()

	return nil
}

// writeStrings takes multiple strings and writes them to w
// using io.WriteString.
func writeStrings(w io.Writer, v ...string) (n int, err error) {
	var x int

	for _, s := range v {
		x, err = io.WriteString(w, s)
		if err != nil {
			return
		}

		n += x
	}

	return
}

// HandlerFunc is invoked for every request for a server-sent events (SSE)
// stream.
type HandlerFunc func(ResponseWriter, *http.Request)

// ConnectionGreeting is the initial string sent to the client to signal
// a successful server-sent events (SSE) stream has been established.
const ConnectionGreeting = ": connected\n\n"

func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	rw := &responseWriter{remote: w, flusher: flusher, buf: make([]byte, 4*1024)}

	_, err := io.WriteString(w, ConnectionGreeting)
	if err != nil {
		return
	}

	flusher.Flush()

	hf(rw, r)
}

// isBrokenPipe returns true if given error is a broken TCP connection.
func isBrokenPipe(err error) bool {
	if err == nil {
		return false
	}

	if opErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
			return sysErr.Err == syscall.EPIPE
		}
	}

	return false
}
