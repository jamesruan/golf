package processor

import (
	"github.com/jamesruan/golf/event"
	"sync"
)

type loggerP struct {
	name       string
	logger     event.Logger
	ch_stopped chan struct{}
}

var loggerWg = new(sync.WaitGroup)

func WaitAllLoggerStop() {
	loggerWg.Wait()
}

// LoggerP returns a processor that handle event to a logger
func NewLoggerP(name string, logger event.Logger) ResidentP {
	return &loggerP{
		name:       name,
		logger:     logger,
		ch_stopped: make(chan struct{}),
	}
}

func (p loggerP) Name() string {
	return p.name
}

func (p loggerP) Stopped() <-chan struct{} {
	return p.ch_stopped
}

func (p *loggerP) Start(stop <-chan struct{}) P {
	go func() {
		loggerWg.Add(1)
		defer loggerWg.Done()
		for {
			select {
			case e := <-p.logger.Queue():
				p.logger.Log(e)
			case <-stop:
				for {
					select {
					case e := <-p.logger.Queue():
						p.logger.Log(e)
					default:
						p.logger.Flush()
						return
					}
				}
			}
		}
	}()
	return p
}

func (p loggerP) Process(e *event.Event) {
	p.logger.Enqueue(e)
	return
}

func (p loggerP) Flush() {
	p.logger.Flush()
}
