package grpcserver

import (
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

const verbosityLevel = 99

func newVerbosityLogger(log *logrus.Entry) {
	if log == nil || log.Logger.Level != logrus.DebugLevel {
		return
	}

	vlog := &verbosityFormatLogger{log: log}

	grpclog.SetLoggerV2(
		grpclog.NewLoggerV2WithVerbosity(vlog, vlog, vlog, verbosityLevel),
	)
}

type verbosityFormatLogger struct {
	log *logrus.Entry
}

func (l *verbosityFormatLogger) Write(p []byte) (n int, err error) {
	s := strings.TrimRight(string(p), "\n")
	s = strings.ReplaceAll(s, `"`, "")
	l.log.Debug(s)

	return len(p), nil
}
