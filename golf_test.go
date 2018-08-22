package golf

import (
	"testing"
)

func TestDefault(t *testing.T) {
	Debugf("%d", 0)
	Infof("%d", 1)
	Logf("%d", 2)
	Warnf("%d", 3)
	Errorf("%d", 4)
}
