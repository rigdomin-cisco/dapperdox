package static

import (
	"github.com/sirupsen/logrus"

	"github.com/kenjones-cisco/dapperdox/logger"
)

func log() logrus.Ext1FieldLogger {
	return logger.Logger().WithField("pkg", "handlers.static")
}
