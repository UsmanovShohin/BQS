package pkg

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}
