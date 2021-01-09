package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/techartificer/swiftex/logger/hooks"
)

// Log is user defined logrus log
var Log *logrus.Logger

//SetupLog initilize logrus logger
func SetupLog() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.AddHook(hooks.NewHook())
}
