package request

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
)

type httpClient struct {
}

func NewHttpClient() *httpClient {
	return &httpClient{}
}

func (c *httpClient) Send(ctx context.Context, url string) ([]byte, error) {
	request, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	response, requestErr := http.DefaultClient.Do(request)
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	if requestErr != nil {
		log.Printf("url %s %s\n", url, requestErr)
		return nil, requestErr
	}

	data, readErr := ioutil.ReadAll(response.Body)
	return data, readErr
}
