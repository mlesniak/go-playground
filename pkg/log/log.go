package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var commitWarned = false

func init() {
	// On local environment ignore all file logging.
	env := os.Getenv("ENVIRONMENT")
	if env == "local" {
		logrus.Info("Local environment, using solely console output")
		return
	}

	// Set to JSON.
	formatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(formatter)
}

// New returns a usable log entry with a set of default fields.
func New() *logrus.Entry {
	commit := os.Getenv("COMMIT")
	if commit == "" {
		if !commitWarned {
			commitWarned = true
			logrus.Warn("Environment variable 'commit' not set, unable to record hash in every log")
		}
		return logrus.NewEntry(logrus.New())
	}

	return logrus.WithFields(logrus.Fields{"commit": commit})
}
