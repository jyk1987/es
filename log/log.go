package log

import (
	"bytes"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
	Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.999",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			fs := strings.Split(f.Function, ".")
			b := bytes.Buffer{}
			b.WriteString(filename)
			b.WriteString(":")
			b.WriteString(strconv.Itoa(f.Line))
			b.WriteString("->")
			b.WriteString(fs[len(fs)-1])
			b.WriteString("():")
			return "", b.String()
		},
	})
	// time.Now().Format()
}
