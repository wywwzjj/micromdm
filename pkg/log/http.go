package log

import (
	"context"
	"net"
	"net/http"

	"github.com/felixge/httpsnoop"
)

// HTTP returns an HTTP logging middleware using the provided base logger.
func HTTP(l Logger) func(http.Handler) http.Handler { return handler{logger: l}.decorate }

type handler struct{ logger Logger }

func (h handler) decorate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := newTraceContext(r.Context())
		ctx = NewContext(ctx, h.logger)

		// https://github.com/felixge/httpsnoop#why-this-package-exists
		// https://github.com/golang/go/issues/18997
		var metrics httpsnoop.Metrics

		defer func() {
			logRequest(ctx, metrics.Code, r)
		}()

		metrics = httpsnoop.CaptureMetrics(next, w, r.WithContext(ctx))
	})
}

func logRequest(ctx context.Context, code int, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	url := *r.URL
	uri := r.RequestURI

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if r.ProtoMajor == 2 && r.Method == "CONNECT" {
		uri = r.Host
	}

	if uri == "" {
		uri = url.RequestURI()
	}

	keyvals := []interface{}{
		"method", r.Method,
		"status", code,
		"proto", r.Proto,
		"host", host,
		"user_agent", r.UserAgent(),
		"path", uri,
	}

	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		keyvals = append(keyvals, "x_forwarded_for", fwd)
	}

	if referer := r.Referer(); referer != "" {
		keyvals = append(keyvals, "referer", referer)
	}

	if code >= 500 {
		Info(FromContext(ctx)).Log(keyvals...)
	} else {
		Debug(FromContext(ctx)).Log(keyvals...)
	}
}
