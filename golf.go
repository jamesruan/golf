package golf

import (
	"github.com/jamesruan/golf/formatter/discard"
	"github.com/jamesruan/golf/formatter/text"
	"io/ioutil"
	"os"
	"sync"
)

var sinkWg sync.WaitGroup
var sinkCloseSignal chan struct{} = make(chan struct{})

func closeAllSink() {
	close(sinkCloseSignal)
	sinkWg.Wait()
}

var (
	DefaultSink      = DefaultStreamSink(os.Stderr, text.Console)
	DefaultPlainSink = DefaultStreamSink(os.Stderr, text.Plain)
	DiscardSink      = DefaultStreamSink(ioutil.Discard, discard.Default)
)

var (
	DefaultEntry      = NewTopicEntry("", DefaultSink)
	DefaultPlainEntry = NewTopicEntry("", DefaultPlainSink)
)
