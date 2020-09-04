package es

import (
	"github.com/sirupsen/logrus"
	"os"
)

var eslog *logrus.Logger

func init() {
	eslog = logrus.New()
	eslog.Out = os.Stdout
}
