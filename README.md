# Go LOggering Framework

## Features and TODOs
   - [X] Entry -> Handler... -> Logger architecture where Logger runs in a separate goroutine and allows asynchronous logging.
   - [X] Flexible handlers that can broadcast and forward logging events to logger in different condition.
   - [X] Multi-goroutine safe handlers that can change event route condition online.
   - [X] Support console based oupput with colorized formatter.
   - [X] Support file based output and rotate with [lumberjack](https://github.com/natefinch/lumberjack), response to os.Signal for triggered rotate.
   - [ ] Support network output with file based spooling.

