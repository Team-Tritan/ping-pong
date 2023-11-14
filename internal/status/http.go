package status

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	ErrMaxRetries = errors.New("reached max retries")
)

type HTTPStatusClient struct {
	RetryTimeout time.Duration
	MaxRetries   int
	Backoff      bool
}

type retryableTransport struct {
	transport    http.RoundTripper
	maxRetries   int
	retryTimeout time.Duration
}

func newRetryableTransit(maxRetries int, retryTimeout time.Duration) *http.Client {
	transport := &retryableTransport{
		transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		maxRetries:   maxRetries,
		retryTimeout: retryTimeout,
	}

	return &http.Client{
		Transport: transport,
	}
}

func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}

	return resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout
}

func backoff(timeout time.Duration, retries int) time.Duration {
	return time.Duration(float64(timeout.Nanoseconds()) * float64(retries))
}

func drain(resp *http.Response) {
	if resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.transport.RoundTrip(req)

	// Copy the body bytes so that we can send them in the next request
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	retries := 0
	for shouldRetry(err, resp) && retries < t.maxRetries {
		time.Sleep(backoff(t.retryTimeout, retries))

		if resp != nil {
			drain(resp)
		}

		if req.Body != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		resp, err = t.transport.RoundTrip(req)

		retries++
	}

	return resp, err
}

func NewHTTPStatus() *HTTPStatusClient {
	return &HTTPStatusClient{
		RetryTimeout: time.Millisecond * 500,
		MaxRetries:   5,
		Backoff:      true,
	}
}

func (status *HTTPStatusClient) CheckStatus(method string, url string, retry bool) (*HTTPStatus, error) {
	return status.checkStatus(method, url, retry)
}

type HTTPStatus struct {
	Body       string
	Retries    int
	Statistics *HTTPStatistics
}

func (status *HTTPStatusClient) checkStatus(method string, url string, retry bool) (*HTTPStatus, error) {
	var client *http.Client
	if retry {
		client = newRetryableTransit(status.MaxRetries, status.RetryTimeout)
	} else {
		client = &http.Client{}
	}

	req, stats, err := newHttpTraceRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPStatus{
		Body:       string(body),
		Statistics: stats,
		Retries:    0,
	}, nil
}
