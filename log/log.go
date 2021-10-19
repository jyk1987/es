package log

import (
	_ "github.com/gogf/gf"
	"github.com/gogf/gf/os/glog"
)

var Log *glog.Logger

func init() {
	Log = glog.New()
	Log.SetStdoutPrint(true)
	Log.SetFlags(glog.F_FILE_SHORT | glog.F_TIME_STD)
}
