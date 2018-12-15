// Package golf is a go logging framework
//
// Examples
//
// Just use DefaultEntry for console based use.
//
//    import (
//    	"github.com/jamesruan/golf"
//    	"github.com/jamesruan/golf/event"
//    )
//
//    func main() {
//    	logger := golf.DefaultEntry
//    	logger.Infof("Hello, World!")
//
//    	topic_logger := golf.NewTopicEntry("mytopic", golf.DefaultSink)
//    	topic_logger.Infof("log with topic")
//
//    	field_logger := topic_logger.WithFields(event.Field{
//    		Name: "field",
//    		Value: "field_value",
//    	})
//
//    	field_logger.Infof("log with field")
//    	field_logger.Fatalf("make sure log is flushed")
//    }
//
//
// Filter logging level
//
//    import (
//    	"github.com/jamesruan/golf"
//    	"github.com/jamesruan/golf/event"
//    	"github.com/jamesruan/golf/handlers"
//    )
//
//    func main() {
//    	levelFilterSink := handlers.NewLevel(event.WARN, golf.DefaultSink, nil)
//    	logger := golf.NewTopicEntry("mytopic", levelFilterSink)
//    	logger.Infof("info") // dropped
//    	logger.Warnf("warn")
//
//    	levelFilterSink.SetLevel(event.INFO) // can be goroutine-safely called
//    	logger.Infof("info") // displayed
//    	logger.Fatalf("fatal")
//    }
package golf

import (
	"github.com/jamesruan/golf/formatter/discard"
	"github.com/jamesruan/golf/formatter/text"
	"github.com/jamesruan/golf/logger"
)

var (
	DefaultSink      = DefaultStreamSink(logger.Stderr, text.Console)
	DefaultPlainSink = DefaultStreamSink(logger.Stderr, text.Plain)
	DiscardSink      = DefaultStreamSink(logger.Discard, discard.Default)
)

var (
	DefaultEntry      = NewTopicEntry("", DefaultSink)
	DefaultPlainEntry = NewTopicEntry("", DefaultPlainSink)
)
