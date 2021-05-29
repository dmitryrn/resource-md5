package worker

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

type TestHTTPClient struct{}

func (TestHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	if url == "test1" {
		return nil, errors.New("test1")
	}
	panic("bad arguments")
}

func TestWorker(t *testing.T) {
	w := &Worker{
		httpClient: TestHTTPClient{},
	}
	_, err := w.do(nil, "test1")
	if err == nil {
		t.Errorf("error should be non nil")
		return
	}
	if err.Error() != "fetch error: test1" {
		t.Errorf("got wrong error %s", err.Error())
		return
	}
}
