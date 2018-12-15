# Go LOggering Framework

See documents [![GoDoc](https://godoc.org/github.com/jamesruan/golf?status.svg)](https://godoc.org/github.com/jamesruan/golf).

## Features and TODOs
   - [X] Entry -> Handler... -> Logger architecture where Logger runs in a separate goroutine and allows asynchronous logging.
   - [X] Flexible handlers that can broadcast and forward logging events to logger in different condition.
   - [X] Multi-goroutine safe handlers that can change event route condition online.
   - [X] Support file based output and rotate with [lumberjack](https://github.com/natefinch/lumberjack), response to os.Signal for triggered rotate.
   - [ ] Support network output with file based spooling.

   - Console formatter
     - [X] Colorized.
     - [X] Call stack print.
     - [X] Calling position (file:line).
     - [X] Line-bufferred output.

## Example

```go
import (
	"github.com/jamesruan/golf"
	"github.com/jamesruan/golf/event"
)

func main() {
	logger := golf.DefaultEntry
	logger.Infof("Hello, World!")

	topic_logger := golf.NewTopicEntry("mytopic", golf.DefaultSink)
	topic_logger.Infof("log with topic")

	field_logger := topic_logger.WithFields(event.Field{
		Name: "field",
		Value: "field_value",
	})

	field_logger.Infof("log with field")
	field_logger.Fatalf("make sure log is flushed")
}
```
