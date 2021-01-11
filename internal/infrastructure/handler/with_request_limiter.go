package handler

import "net/http"

type withRequestLimiter struct {
	requestLimiter   chan struct{}
	decoratedHandler http.Handler
}

func (h *withRequestLimiter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	select {
	case h.requestLimiter <- struct{}{}:
		defer func() {
			<-h.requestLimiter
		}()
		h.decoratedHandler.ServeHTTP(writer, request)
	default:
		writer.WriteHeader(http.StatusTooManyRequests)
	}
}

func WithRequestLimiter(handler http.Handler, maxRequestsNum int) http.Handler {
	return &withRequestLimiter{
		requestLimiter:   make(chan struct{}, maxRequestsNum),
		decoratedHandler: handler,
	}
}
