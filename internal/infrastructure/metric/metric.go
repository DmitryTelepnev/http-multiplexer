package metric

import "net/http"

type Metric interface {
	CollectHTTP(handler http.Handler) http.Handler
}
