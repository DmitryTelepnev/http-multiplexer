package metric

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/DmitryTelepnev/http-multiplexer/internal/infrastructure/config"
	"github.com/prometheus/client_golang/prometheus"
)

type prom struct {
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewPrometheus() *prom {
	namespace := strings.NewReplacer(".", "_", "-", "_").
		Replace(config.Namespace + "_" + config.ApplicationName)

	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "http_requests",
		Help:      "The http requests counter",
	}, []string{"uri", "code"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "http_requests_duration",
		Help:      "The http request duration",
		Buckets:   []float64{0.1, 0.2, 0.3, 0.5, 1, 2, 5, 10, math.Inf(1)},
	}, []string{"uri", "code"})

	prometheus.MustRegister(
		requestCounter,
		requestDuration,
	)
	return &prom{
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

type statusReminder struct {
	http.ResponseWriter
	status int
}

func (sr *statusReminder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (m *prom) CollectHTTP(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		timeStart := time.Now()

		writerWithStatus := &statusReminder{writer, http.StatusOK}
		handler.ServeHTTP(writerWithStatus, request)

		m.requestCounter.
			WithLabelValues(request.RequestURI, fmt.Sprint(writerWithStatus.status)).
			Inc()
		m.requestDuration.
			WithLabelValues(request.RequestURI, fmt.Sprint(writerWithStatus.status)).
			Observe(time.Since(timeStart).Seconds())
	})
}
