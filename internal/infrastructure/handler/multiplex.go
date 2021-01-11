package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/DmitryTelepnev/http-multiplexer/internal/infrastructure/config"
	"github.com/DmitryTelepnev/http-multiplexer/internal/service/multiplexer"
)

type MultiplexHandler struct {
	cfg         *config.Multiplexer
	wg          *sync.WaitGroup
	multiplexer *multiplexer.Multiplexer
}

func NewMultiplexHandler(cfg *config.Multiplexer, wg *sync.WaitGroup, multiplex *multiplexer.Multiplexer) *MultiplexHandler {
	return &MultiplexHandler{
		cfg:         cfg,
		wg:          wg,
		multiplexer: multiplex,
	}
}

type Response struct {
	Msg  string            `json:"msg"`
	Data map[string]string `json:"data"`
}

func (h *MultiplexHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.wg.Add(1)
	defer h.wg.Done()

	writer.Header().Add("Content-type", "application/json")

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
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
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(Response{
			Msg: "too many urls",
		})
		return
	}

	if len(urls) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(Response{
			Msg: "empty urls list",
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
