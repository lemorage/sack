package main

import (
	"fmt"
	"net/http"
	"strings"
)

// ResponseRecorder to capture the response body
type responseRecorder struct {
	http.ResponseWriter
	body       strings.Builder
	statusCode int
}

// the middleware to inject the WebSocket script into HTML responses
func injectWebSocketScriptMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the response
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		// If the content type is HTML, inject the WebSocket script
		if rr.Header().Get("Content-Type") == "text/html; charset=utf-8" {
			injectedContent := strings.ReplaceAll(rr.body.String(), "</body>", `<script>
				const ws = new WebSocket("ws://localhost:7536/ws");
				ws.onmessage = function(event) {
					if (event.data === "reload") {
						window.location.reload();
					}
				};
			</script></body>`)

			w.Header().Set("Content-Length", fmt.Sprint(len(injectedContent)))
			w.WriteHeader(rr.statusCode)
			w.Write([]byte(injectedContent))
		} else {
			// Otherwise, write the original response
			w.WriteHeader(rr.statusCode)
			w.Write([]byte(rr.body.String()))
		}
	})
}

func (rr *responseRecorder) Write(p []byte) (int, error) {
	return rr.body.Write(p)
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
}
