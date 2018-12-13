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
