package multiplexer

import (
	"github.com/DmitryTelepnev/http-multiplexer/internal/infrastructure/config"
	"github.com/DmitryTelepnev/http-multiplexer/internal/service/request"
	"context"
	"log"
	"sync"
)

type Multiplexer struct {
	cfg    config.Multiplexer
	client request.Client
}

func NewMultiplexer(cfg config.Multiplexer, client request.Client) *Multiplexer {
	return &Multiplexer{cfg: cfg, client: client}
}

func (m *Multiplexer) SendRequests(urls []string) (map[string]string, error) {
	urlsData := make(map[string]string, len(urls))
	limiter := make(chan struct{}, m.cfg.MaxOneTimeRequests)

	var mutex sync.RWMutex
	var requestErr error

	multiplexContext, multiplexCancel := context.WithTimeout(context.Background(), m.cfg.AllRequestsTimeOut)
	defer multiplexCancel()

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()

			limiter <- struct{}{}
			defer func() {
				<- limiter
			}()

			log.Printf("process url %s", url)

			ctx, cancel := context.WithTimeout(multiplexContext, m.cfg.RequestTimeOut)
			data, err := m.client.Send(ctx, url)

			if err != nil {
				mutex.Lock()
				if requestErr == nil {
					requestErr = err
				}
				mutex.Unlock()

				cancel()
				multiplexCancel()
				return
			}

			log.Printf("successfully processed url %s", url)

			mutex.Lock()
			urlsData[url] = string(data)
			mutex.Unlock()
		}(url)
	}

	wg.Wait()
	close(limiter)

	return urlsData, requestErr
}