package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/DmitryTelepnev/http-multiplexer/internal/infrastructure/config"
	"github.com/DmitryTelepnev/http-multiplexer/internal/service/multiplexer"
)

type MultiplexHandler struct {
	cfg            *config.Multiplexer
	wg             *sync.WaitGroup
	requestLimiter chan struct{}
	multiplexer    *multiplexer.Multiplexer
}

func NewMultiplexHandler(cfg *config.Multiplexer, wg *sync.WaitGroup, requestLimiter chan struct{}, multiplex *multiplexer.Multiplexer) *MultiplexHandler {
	return &MultiplexHandler{
		cfg:            cfg,
		wg:             wg,
		requestLimiter: requestLimiter,
		multiplexer:    multiplex,
	}
}

type Response struct {
	Msg  string            `json:"msg"`
	Data map[string]string `json:"data"`
}

func (h *MultiplexHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.requestLimiter <- struct{}{}
	defer func() {
		<-h.requestLimiter
	}()
	h.wg.Add(1)
	defer h.wg.Done()

	writer.Header().Add("Content-type", "application/json")

	defer func() {
		closeBodyErr := request.Body.Close()
		if closeBodyErr != nil {
			log.Printf("can't close request body: %v", closeBodyErr)
		}
	}()

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	var urls []string
	unmarshalErr := json.NewDecoder(request.Body).Decode(&urls)
	if unmarshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(writer).Encode(Response{
			Msg: "can't unmarshal request",
		})
		return
	}

	if len(urls) > h.cfg.MaxUrlsInRequest {
		writer.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(writer).Encode(Response{
			Msg: "too many urls",
		})
		return
	}

	urlsData, multiplexErr := h.multiplexer.SendRequests(urls)
	msg := "Ok"
	if multiplexErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		msg = multiplexErr.Error()
		urlsData = map[string]string{}
	}

	writeErr := json.NewEncoder(writer).Encode(Response{
		Msg:  msg,
		Data: urlsData,
	})
	if writeErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
