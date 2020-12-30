package signal

import (
	"os"
	OSSignal "os/signal"
	"syscall"
)

type handler struct {
	signalChan chan os.Signal
	onShutdown func() error
}

func NewHandler(onShutdown func() error) *handler {
	signalChan := make(chan os.Signal)
	OSSignal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	return &handler{
		signalChan: signalChan,
		onShutdown: onShutdown,
	}
}

func (h *handler) Poll() {
	<-h.signalChan
	shutdownErr := h.onShutdown()
	if shutdownErr != nil {
		panic(shutdownErr)
	}
}
