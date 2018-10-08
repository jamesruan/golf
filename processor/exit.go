package processor

import (
	"os"
	"sync"
)

var exitSignal = make(chan struct{})
var loggerWg = new(sync.WaitGroup)

// Exit wait for all LoggerP to do their cleanup and call os.Exit(-1)
func Exit() {
	select {
	case <-exitSignal:
	default:
		close(exitSignal)
	}
	// wait for all logger to flush its queue
	loggerWg.Wait()
	os.Exit(-1)
}
