package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetOutput(os.Stderr)
	Log.SetLevel(logrus.DebugLevel)
	// Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.999",
	})
	// time.Now().Format()
}
