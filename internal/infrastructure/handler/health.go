package handler

import "net/http"

type Health struct {
}

func (h *Health) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
