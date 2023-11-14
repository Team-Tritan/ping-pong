package status

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"
)

type HTTPStatistics struct {
	dnsStart          time.Time
	DnsDone           time.Duration
	TlsHandshakeStart time.Duration
	TlsHandshakeDone  time.Duration
	GotFirstByte      time.Duration
	TcpStart          time.Duration
	TcpDone           time.Duration

	started time.Time
}

func withHttpTrace(req *http.Request, stats *HTTPStatistics) *http.Request {
	return req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		// If we are reusing the connection, we want to reset the start time
		GotConn: func(connInfo httptrace.GotConnInfo) {
			stats.started = time.Now()
		},
		DNSStart: func(di httptrace.DNSStartInfo) {
			stats.dnsStart = time.Now()
		},
		DNSDone: func(di httptrace.DNSDoneInfo) {
			stats.DnsDone = time.Since(stats.dnsStart)
		},
		TLSHandshakeStart: func() {
			stats.TlsHandshakeStart = time.Since(stats.started)
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			stats.TlsHandshakeDone = time.Since(stats.started)
		},
		GotFirstResponseByte: func() {
			stats.GotFirstByte = time.Since(stats.started)
		},
		ConnectStart: func(network, addr string) {
			stats.started = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			stats.TcpDone = time.Since(stats.started)
		},
	}))
}

func newHttpTraceRequest(method, url string, body io.ReadCloser) (*http.Request, *HTTPStatistics, error) {
	stats := &HTTPStatistics{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	return withHttpTrace(req, stats), stats, nil
}
