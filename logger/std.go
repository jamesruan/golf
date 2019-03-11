package logger

import (
	"io"
	"io/ioutil"
	"os"
)

type RotateLogger interface {
	io.Writer
	Rotate() error
}

var (
	Stderr  = os.Stderr
	Stdout  = os.Stdout
	Discard = ioutil.Discard
)
