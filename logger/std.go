package logger

import (
	"io"
	"io/ioutil"
	"os"
)

type StreamLogger = io.Writer

type RotateLogger interface {
	StreamLogger
	Rotate() error
}

var (
	Stderr  = os.Stderr
	Stdout  = os.Stdout
	Discard = ioutil.Discard
)
