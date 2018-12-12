package golf

import (
	"sync"
)

var sinkWg sync.WaitGroup
var sinkCloseSignal chan struct{} = make(chan struct{})

func closeAllSink() {
	close(sinkCloseSignal)
	sinkWg.Wait()
}
