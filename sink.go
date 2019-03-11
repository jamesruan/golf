package golf

import (
	"os"
	"sync"
)

var sinkWg sync.WaitGroup
var sinkCloseSignal chan struct{} = make(chan struct{})

func closeAllSink() {
	close(sinkCloseSignal)
	sinkWg.Wait()
}

// used for signalling the logger to flush and rotate when the logger implements logger.RotateLogger
var SinkRotateSignal chan os.Signal = make(chan os.Signal, 1)

type sinkBase struct{}

// Register mark the sink started and waiting for close signal
func (b sinkBase) register() {
	sinkWg.Add(1)
}

// Closing return the channel for closing signal
func (b sinkBase) closing() <-chan struct{} {
	return sinkCloseSignal
}

// Deregister mark the sink stopped
func (b sinkBase) deregister() {
	sinkWg.Done()
}
