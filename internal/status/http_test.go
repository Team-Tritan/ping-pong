package status

import (
	"net/http"
	"testing"
	"time"
)

func Test_Backoff(t *testing.T) {
	timeout := 100 * time.Millisecond

	backedoff_timeout := backoff(timeout, 5)
	t.Logf("%f", backedoff_timeout.Seconds())
}

type TestTransport struct {
	retriesCounted      int
	retriesUntilSuccess int
}

func (tt *TestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tt.retriesCounted++

	if tt.retriesCounted == tt.retriesUntilSuccess {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}

	return &http.Response{
		StatusCode: http.StatusBadGateway,
	}, nil
}

func Test_HttpRetryTransport(t *testing.T) {
	testTransport := &TestTransport{
		retriesCounted:      0,
		retriesUntilSuccess: 3,
	}

	transport := &retryableTransport{
		transport:    testTransport,
		maxRetries:   5,
		retryTimeout: 500 * time.Millisecond,
	}

	cli := &http.Client{
		Transport: transport,
	}

	_, err := cli.Get("example.com")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if testTransport.retriesCounted != testTransport.retriesUntilSuccess {
		t.Error("assert failed: testTransport.retriesCounted != testTransport.retriesUntilSuccess")
		t.Fail()
	}
}
