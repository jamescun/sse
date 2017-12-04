# Server Sent Events

SSE implements a server-sent events (event source) client and server in Go.


## Server Example

```go
package main

import (
	"net/http"
	"time"
	"strings"

	"github.com/jamescun/sse"
)

func Ticker(w sse.ResponseWriter, r *http.Request) {
	for {
		err := w.WriteEvent(
			&sse.Event{Name: "tick"},
			strings.NewReader(time.Now().Format(time.RFC3339)),
		)
		if err != nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	http.ListenAndServe("127.0.0.1:8080", sse.HandlerFunc(Ticker))
}
```
