package wraplogrus

import "github.com/sirupsen/logrus"

// Config is a configuration logger.
type Config struct {
	Channel       string
	BuildID       string
	Level         string
	TextFormatter *logrus.TextFormatter
	JSONFormatter *logrus.JSONFormatter
}

// New create a new logger.
func New(conf Config) *logrus.Entry {
	if conf.Level == "" {
		conf.Level = "error"
	}

	logLevel, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		logLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logLevel)

	if conf.TextFormatter != nil {
		logrus.SetFormatter(conf.TextFormatter)
	}

	if conf.JSONFormatter != nil {
		logrus.SetFormatter(conf.JSONFormatter)
	}

	return logrus.WithFields(
		logrus.Fields{
			"channel":  conf.Channel,
			"build_id": conf.BuildID,
		},
	)
}
