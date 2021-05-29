package worker

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dmitryrn/resource-md5/config"
	"github.com/dmitryrn/resource-md5/transport"
	"io"
	"net/http"
	"sync"
)

type (
	Worker struct {
		config     config.Config
		httpClient transport.HTTPClient
	}
)

func NewWorker(config config.Config, httpClient transport.HTTPClient) *Worker {
	return &Worker{
		config:     config,
		httpClient: httpClient,
	}
}

// Do method will run multiple workers in parallel according to config.Config
func (w *Worker) Do() (<-chan [2]string, <-chan error) {
	ch := make(chan [2]string)
	errCh := make(chan error)

	ctx := context.Background()

	wg := sync.WaitGroup{}

	maxWorkers := int(w.config.Parallel)

	jobs := make(chan string, len(w.config.URLs))
	for _, url := range w.config.URLs {
		jobs <- url
	}
	close(jobs)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for url := range jobs {
				hash, err := w.do(ctx, url)
				if err != nil {
					errCh <- fmt.Errorf("worker failed for url %s: %w", url, err)
					return
				}
				ch <- [2]string{url, hash}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
		close(errCh)
	}()

	return ch, errCh
}

func (w *Worker) do(ctx context.Context, url string) (string, error) {
	resp, err := w.httpClient.Get(ctx, url)
	if err != nil {
		return "", fmt.Errorf("fetch error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("got non-200 status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body error: %w", err)
	}

	hash := md5.Sum(body)
	md5Hash := hex.EncodeToString(hash[:])

	return md5Hash, nil
}
