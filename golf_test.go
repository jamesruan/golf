package golf

import (
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	SetFilter("", WithLevel(DEBUG))
	Debugf("%v", DEBUG)
	Infof("%v", INFO)
	Logf("%v", LOG)
	Warnf("%v", WARN)
	Errorf("%v", ERROR)
	time.Sleep(time.Second)
}
