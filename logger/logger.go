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

// Errorln user defined logrus log function
func Errorln(args ...interface{}) {
	Log.Errorln(args...)
}

// Infoln user defined logrus log function
func Infoln(args ...interface{}) {
	Log.Infoln(args...)
}

// Println user defined logrus log function
func Println(args ...interface{}) {
	Log.Println(args...)
}

// Printf user defined logrus log function
func Printf(format string, args ...interface{}) {
	Log.Printf(format, args...)
}
