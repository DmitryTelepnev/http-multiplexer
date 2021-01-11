package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/DmitryTelepnev/http-multiplexer/internal/infrastructure/config"
	"github.com/DmitryTelepnev/http-multiplexer/internal/infrastructure/handler"
	"github.com/DmitryTelepnev/http-multiplexer/internal/infrastructure/metric"
	"github.com/DmitryTelepnev/http-multiplexer/internal/infrastructure/signal"
	"github.com/DmitryTelepnev/http-multiplexer/internal/service/multiplexer"
	"github.com/DmitryTelepnev/http-multiplexer/internal/service/request"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.MustConfigure()
	metrics := metric.NewPrometheus()

	http.Handle(cfg.K8S.MetricEndpoint, promhttp.Handler())
	http.Handle(cfg.K8S.HealthEndpoint, &handler.Health{})

	httpClient := request.NewHttpClient()
	httpMultiplexer := multiplexer.NewMultiplexer(cfg.Multiplexer, httpClient)

	var wg sync.WaitGroup
	multiplexHandler := handler.NewMultiplexHandler(&cfg.Multiplexer, &wg, httpMultiplexer)
	handlerWithRequestLimiter := handler.WithRequestLimiter(multiplexHandler, cfg.Multiplexer.MaxActiveConnections)
	http.Handle("/", metrics.CollectHTTP(handlerWithRequestLimiter))

	go func() {
		listenErr := http.ListenAndServe(fmt.Sprintf(":%d", cfg.K8S.Port), nil)
		if listenErr != nil {
			panic(listenErr)
		}
	}()

	signalHandler := signal.NewHandler(func() error {
		// graceful stop with resource closing
		wg.Wait()
		log.Println("Gracefully stopped")
		return nil
	})

	signalHandler.Poll()
}
