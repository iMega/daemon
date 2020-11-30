package redis

import (
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func newLogger(l *logrus.Entry) {
	var wlog io.Writer

	wlog = &logger{log: l}

	if l == nil || l.Logger.Level != logrus.DebugLevel {
		wlog = ioutil.Discard
	}

	redis.SetLogger(
		log.New(wlog, "", 0),
	)
}

type logger struct {
	log *logrus.Entry
}

func (l *logger) Write(p []byte) (n int, err error) {
	s := strings.TrimRight(string(p), "\n")
	s = strings.ReplaceAll(s, `"`, "")
	l.log.Debug(s)

	return len(p), nil
}
