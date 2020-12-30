package multiplexer

import (
	"fmt"
	"testing"
	"time"

	"github.com/DmitryTelepnev/http-multiplexer/internal/infrastructure/config"
	"github.com/DmitryTelepnev/http-multiplexer/mock/service/request"
	"github.com/stretchr/testify/assert"
)

func TestSendRequests_ErrorOnFirstURL(t *testing.T) {
	cfg := config.Multiplexer{
		AllRequestsTimeOut:   1 * time.Second,
		MaxOneTimeRequests:   1,
		RequestTimeOut:       2 * time.Second,
		MaxUrlsInRequest:     1,
		MaxActiveConnections: 1,
	}

	err := fmt.Errorf("test")

	errs := map[string]error{"url1": err}

	inmemoryMultiplexer := NewMultiplexer(cfg, request.NewInmemory(errs, nil))

	_, multiplexErr := inmemoryMultiplexer.SendRequests([]string{"url1", "url2", "url3", "url4", "url5"})
	assert.Error(t, err, multiplexErr)
}

func TestSendRequests_ErrorOnURLInMiddle(t *testing.T) {
	cfg := config.Multiplexer{
		AllRequestsTimeOut:   1 * time.Second,
		MaxOneTimeRequests:   2,
		RequestTimeOut:       2 * time.Second,
		MaxUrlsInRequest:     1,
		MaxActiveConnections: 1,
	}

	err := fmt.Errorf("test")

	errs := map[string]error{"url4": err}

	inmemoryMultiplexer := NewMultiplexer(cfg, request.NewInmemory(errs, nil))

	_, multiplexErr := inmemoryMultiplexer.SendRequests([]string{"url1", "url2", "url3", "url4", "url5"})
	assert.Error(t, err, multiplexErr)
}

func TestSendRequests_ErrorOnURLInMiddleAndUrlsProcessedAtOneTime(t *testing.T) {
	cfg := config.Multiplexer{
		AllRequestsTimeOut:   1 * time.Second,
		MaxOneTimeRequests:   5,
		RequestTimeOut:       2 * time.Second,
		MaxUrlsInRequest:     1,
		MaxActiveConnections: 1,
	}

	err := fmt.Errorf("test")

	errs := map[string]error{"url4": err}

	inmemoryMultiplexer := NewMultiplexer(cfg, request.NewInmemory(errs, nil))

	_, multiplexErr := inmemoryMultiplexer.SendRequests([]string{"url1", "url2", "url3", "url4", "url5"})
	assert.Error(t, err, multiplexErr)
}

func TestSendRequests_WithoutErrWithData(t *testing.T) {
	cfg := config.Multiplexer{
		AllRequestsTimeOut:   1 * time.Second,
		MaxOneTimeRequests:   5,
		RequestTimeOut:       2 * time.Second,
		MaxUrlsInRequest:     1,
		MaxActiveConnections: 1,
	}

	errs := map[string]error{}

	testData := []byte("test")
	inmemoryMultiplexer := NewMultiplexer(cfg, request.NewInmemory(errs, testData))

	data, multiplexErr := inmemoryMultiplexer.SendRequests([]string{"url1", "url2", "url3", "url4", "url5"})
	assert.Nil(t, multiplexErr)
	assert.Equal(t, map[string]string{"url1": "test", "url2": "test", "url3": "test", "url4": "test", "url5": "test"}, data)
}
