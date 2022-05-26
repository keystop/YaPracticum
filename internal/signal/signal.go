package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func HandleQuit(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	cancel()
}
