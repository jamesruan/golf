package processor

import (
	"github.com/jamesruan/golf/event"
	"time"
)

type loggerP struct {
	name   string
	logger event.Logger
}

// LoggerP returns a processor that handle event to a logger
func NewLoggerP(name string, logger event.Logger) P {
	process_exiting := exitSignal
	go func() {
		loggerWg.Add(1)
		defer loggerWg.Done()
		var timeout <-chan time.Time
		for {
			select {
			case e := <-logger.Queue():
				logger.Log(e)
			case <-process_exiting:
				timeout = time.After(100 * time.Millisecond)
				process_exiting = nil
			case <-timeout:
				return
			}
		}
	}()
	return &loggerP{
		name:   name,
		logger: logger,
	}
}

func (p loggerP) Name() string {
	return p.name
}

func (p loggerP) Process(e *event.Event) {
	p.logger.Enqueue(e)
	return
}

func (p loggerP) Flush() {
	p.logger.Flush()
}
